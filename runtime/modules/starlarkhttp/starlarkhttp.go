// The MIT License (MIT)
//
// Copyright (c) 2018 QRI, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package starlarkhttp

// This module uses qri-io/starlib's http module as a base, and adds Tidbyt
// specific headers to all requests.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	util "github.com/qri-io/starlib/util"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// AsString unquotes a starlark string value
func AsString(x starlark.Value) (string, error) {
	return strconv.Unquote(x.String())
}

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('http.star', 'http')
const ModuleName = "http.star"

var (
	// StarlarkHTTPClient is a factory method for creating the http client used to create the http module.
	// override with a custom function before calling LoadModule
	StarlarkHTTPClient = func() *http.Client {
		return http.DefaultClient
	}
	// StarlarkHTTPGuard is a global RequestGuard used in LoadModule. override with a custom
	// implementation before calling LoadModule
	StarlarkHTTPGuard RequestGuard
)

// Encodings for form data.
//
// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/POST
const (
	formEncodingMultipart = "multipart/form-data"
	formEncodingURL       = "application/x-www-form-urlencoded"
)

// LoadModule creates an http Module
func LoadModule() (starlark.StringDict, error) {
	var m = &Module{cli: StarlarkHTTPClient()}
	if StarlarkHTTPGuard != nil {
		m.rg = StarlarkHTTPGuard
	}
	ns := starlark.StringDict{
		"http": m.Struct(),
	}
	return ns, nil
}

// RequestGuard controls access to http by checking before making requests
// if Allowed returns an error the request will be denied
type RequestGuard interface {
	Allowed(thread *starlark.Thread, req *http.Request) (*http.Request, error)
}

// Module joins http tools to a dataset, allowing dataset
// to follow along with http requests
type Module struct {
	cli *http.Client
	rg  RequestGuard
}

// Struct returns this module's methods as a starlark Struct
func (m *Module) Struct() *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(starlarkstruct.Default, m.StringDict())
}

// StringDict returns all module methods in a starlark.StringDict
func (m *Module) StringDict() starlark.StringDict {
	return starlark.StringDict{
		"get":     starlark.NewBuiltin("get", m.reqMethod("get")),
		"put":     starlark.NewBuiltin("put", m.reqMethod("put")),
		"post":    starlark.NewBuiltin("post", m.reqMethod("post")),
		"delete":  starlark.NewBuiltin("delete", m.reqMethod("delete")),
		"patch":   starlark.NewBuiltin("patch", m.reqMethod("patch")),
		"options": starlark.NewBuiltin("options", m.reqMethod("options")),
	}
}

// reqMethod is a factory function for generating starlark builtin functions for different http request methods
func (m *Module) reqMethod(method string) func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var (
			urlv         starlark.String
			params       = &starlark.Dict{}
			headers      = &starlark.Dict{}
			formBody     = &starlark.Dict{}
			formEncoding starlark.String
			auth         starlark.Tuple
			body         starlark.String
			jsonBody     starlark.Value
			ttl          starlark.Int
		)

		if err := starlark.UnpackArgs(method, args, kwargs, "url", &urlv, "params?", &params, "headers", &headers, "body", &body, "form_body", &formBody, "form_encoding", &formEncoding, "json_body", &jsonBody, "auth", &auth, "ttl_seconds?", &ttl); err != nil {
			return nil, err
		}

		rawurl, err := AsString(urlv)
		if err != nil {
			return nil, err
		}
		if err = setQueryParams(&rawurl, params); err != nil {
			return nil, err
		}

		req, err := http.NewRequest(strings.ToUpper(method), rawurl, nil)
		if err != nil {
			return nil, err
		}
		if m.rg != nil {
			req, err = m.rg.Allowed(thread, req)
			if err != nil {
				return nil, err
			}
		}

		if err = setHeaders(req, headers); err != nil {
			return nil, err
		}
		if err = setStandardHeaders(req, thread, ttl); err != nil {
			return nil, err
		}
		if err = setAuth(req, auth); err != nil {
			return nil, err
		}
		if err = SetBody(req, body, formBody, formEncoding, jsonBody); err != nil {
			return nil, err
		}

		res, err := m.cli.Do(req)
		if err != nil {
			return nil, err
		}

		r := &Response{*res}
		return r.Struct(), nil
	}
}

func setQueryParams(rawurl *string, params *starlark.Dict) error {
	keys := params.Keys()
	if len(keys) == 0 {
		return nil
	}

	u, err := url.Parse(*rawurl)
	if err != nil {
		return err
	}

	q := u.Query()
	for _, key := range keys {
		keystr, err := AsString(key)
		if err != nil {
			return err
		}

		val, _, err := params.Get(key)
		if err != nil {
			return err
		}
		if val.Type() != "string" {
			return fmt.Errorf("expected param value for key '%s' to be a string. got: '%s'", key, val.Type())
		}
		valstr, err := AsString(val)
		if err != nil {
			return err
		}

		q.Set(keystr, valstr)
	}

	u.RawQuery = q.Encode()
	*rawurl = u.String()
	return nil
}

func setAuth(req *http.Request, auth starlark.Tuple) error {
	if len(auth) == 0 {
		return nil
	} else if len(auth) == 2 {
		username, err := AsString(auth[0])
		if err != nil {
			return fmt.Errorf("parsing auth username string: %s", err.Error())
		}
		password, err := AsString(auth[1])
		if err != nil {
			return fmt.Errorf("parsing auth password string: %s", err.Error())
		}
		req.SetBasicAuth(username, password)
		return nil
	}
	return fmt.Errorf("expected two values for auth params tuple")
}

func setStandardHeaders(req *http.Request, thread *starlark.Thread, ttl starlark.Int) error {
	// Set app identifier.
	req.Header.Set("X-Tidbyt-App", getAppIdentifier(thread))

	// Set ttl for caching client.
	ttl64, ok := ttl.Int64()
	if !ok {
		return fmt.Errorf("ttl_seconds must be valid integer (not %s)", ttl.String())
	}
	req.Header.Set("X-Tidbyt-Cache-Seconds", fmt.Sprintf("%d", ttl64))

	return nil
}

func getAppIdentifier(thread *starlark.Thread) string {
	parts := strings.Split(thread.Name, "/")
	if len(parts) > 0 {
		return parts[0]
	}

	return thread.Name
}

func setHeaders(req *http.Request, headers *starlark.Dict) error {
	keys := headers.Keys()
	if len(keys) == 0 {
		return nil
	}

	for _, key := range keys {
		keystr, err := AsString(key)
		if err != nil {
			return err
		}

		val, _, err := headers.Get(key)
		if err != nil {
			return err
		}
		if val.Type() != "string" {
			return fmt.Errorf("expected param value for key '%s' to be a string. got: '%s'", key, val.Type())
		}
		valstr, err := AsString(val)
		if err != nil {
			return err
		}

		req.Header.Add(keystr, valstr)
	}

	return nil
}

func SetBody(req *http.Request, body starlark.String, formData *starlark.Dict, formEncoding starlark.String, jsondata starlark.Value) error {
	if !util.IsEmptyString(body) {
		uq, err := strconv.Unquote(body.String())
		if err != nil {
			return err
		}
		req.Body = io.NopCloser(strings.NewReader(uq))
		// Specifying the Content-Length ensures that https://go.dev/src/net/http/transfer.go doesnt specify Transfer-Encoding: chunked which is not supported by some endpoints.
		// This is required when using io.NopCloser method for the request body (see ShouldSendChunkedRequestBody() in the library mentioned above).
		req.ContentLength = int64(len(uq))

		return nil
	}

	if jsondata != nil && jsondata.String() != "" {
		req.Header.Set("Content-Type", "application/json")

		v, err := util.Unmarshal(jsondata)
		if err != nil {
			return err
		}
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(data))
		req.ContentLength = int64(len(data))
	}

	if formData != nil && formData.Len() > 0 {
		form := url.Values{}
		for _, key := range formData.Keys() {
			keystr, err := AsString(key)
			if err != nil {
				return err
			}

			val, _, err := formData.Get(key)
			if err != nil {
				return err
			}
			if val.Type() != "string" {
				return fmt.Errorf("expected param value for key '%s' to be a string. got: '%s'", key, val.Type())
			}
			valstr, err := AsString(val)
			if err != nil {
				return err
			}

			form.Add(keystr, valstr)
		}

		var contentType string
		switch formEncoding {
		case formEncodingURL, "":
			contentType = formEncodingURL
			req.Body = io.NopCloser(strings.NewReader(form.Encode()))
			req.ContentLength = int64(len(form.Encode()))

		case formEncodingMultipart:
			var b bytes.Buffer
			mw := multipart.NewWriter(&b)
			defer mw.Close()

			contentType = mw.FormDataContentType()

			for k, values := range form {
				for _, v := range values {
					w, err := mw.CreateFormField(k)
					if err != nil {
						return err
					}
					if _, err := w.Write([]byte(v)); err != nil {
						return err
					}
				}
			}

			req.Body = io.NopCloser(&b)

		default:
			return fmt.Errorf("unknown form encoding: %s", formEncoding)
		}

		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", contentType)
		}
	}

	return nil
}

// Response represents an HTTP response, wrapping a go http.Response with
// starlark methods
type Response struct {
	http.Response
}

// Struct turns a response into a *starlark.Struct
func (r *Response) Struct() *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
		"url":         starlark.String(r.Request.URL.String()),
		"status_code": starlark.MakeInt(r.StatusCode),
		"headers":     r.HeadersDict(),
		"encoding":    starlark.String(strings.Join(r.TransferEncoding, ",")),

		"body": starlark.NewBuiltin("body", r.Text),
		"json": starlark.NewBuiltin("json", r.JSON),
	})
}

// HeadersDict flops
func (r *Response) HeadersDict() *starlark.Dict {
	d := new(starlark.Dict)
	for key, vals := range r.Header {
		if err := d.SetKey(starlark.String(key), starlark.String(strings.Join(vals, ","))); err != nil {
			panic(err)
		}
	}
	return d
}

// Text returns the raw data as a string
func (r *Response) Text(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body.Close()
	// reset reader to allow multiple calls
	r.Body = io.NopCloser(bytes.NewReader(data))

	return starlark.String(string(data)), nil
}

// JSON attempts to parse the response body as JSON
func (r *Response) JSON(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var data interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	r.Body.Close()
	// reset reader to allow multiple calls
	r.Body = io.NopCloser(bytes.NewReader(body))
	return util.Marshal(data)
}

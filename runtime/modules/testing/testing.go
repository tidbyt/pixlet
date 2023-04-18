package testing

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	starlibhttp "github.com/qri-io/starlib/http"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/starlarktest"
	"tidbyt.dev/pixlet/encode"
	"tidbyt.dev/pixlet/render"
	"tidbyt.dev/pixlet/runtime/modules/render_runtime"
)

const (
	ModuleName = "testing"
)

var (
	once                 sync.Once
	module               starlark.StringDict
	stubbedHttpTransport StubbedHttpTransport
)

type StubbedHttpTransport struct {
	realTransport http.RoundTripper
	stubs         []Stub
	allowRealHttp bool
}

type Stub struct {
	thread    *starlark.Thread
	matcher   *starlark.Function
	responder *starlark.Function
}

func InitHttpStub(allowRealHttp bool) *StubbedHttpTransport {
	stubbedHttpTransport = StubbedHttpTransport{
		realTransport: http.DefaultTransport,
		allowRealHttp: allowRealHttp,
	}
	starlibhttp.Client = &http.Client{Transport: &stubbedHttpTransport}
	return &stubbedHttpTransport
}

func (stub *Stub) Matches(req *http.Request) (bool, error) {
	resultValue, err := stub.matcher.CallInternal(stub.thread, starlark.Tuple{starlark.String(req.Method), starlark.String(req.URL.String())}, nil)
	if err != nil {
		return false, err
	}
	return bool(resultValue.Truth()), nil
}

func (stub *Stub) Respond(req *http.Request) (*http.Response, error) {
	responseValue, err := stub.responder.CallInternal(stub.thread, starlark.Tuple{starlark.String(req.Method), starlark.String(req.URL.String())}, nil)
	if err != nil {
		return nil, err
	}
	responseDict, ok := responseValue.(*starlark.Dict)
	if !ok {
		return nil, fmt.Errorf("Responder must return a dict")
	}

	response := http.Response{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Request:    req,
	}
	for _, keySlValue := range responseDict.Keys() {
		keySlString, ok := keySlValue.(starlark.String)
		if !ok {
			return nil, fmt.Errorf("keys in response dict must be string")
		}

		slVal, _, err := responseDict.Get(keySlValue)
		if err != nil {
			return nil, err
		}
		switch keySlString {
		case starlark.String("status_code"):
			status_code, err := statusFromStarlarkValue(slVal)
			if err != nil {
				return nil, err
			}
			response.StatusCode = int(status_code)
			response.Status = http.StatusText(int(status_code))
		case starlark.String("headers"):
			headers, err := headersFromStarlarkValue(slVal)
			if err != nil {
				return nil, err
			}
			response.Header = headers
		case starlark.String("body"):
			body, err := bodyFromStarlarkValue(slVal)
			if err != nil {
				return nil, err
			}
			response.Body = body
		default:
			return nil, fmt.Errorf("Unknown key in response dict: %s", keySlString)
		}
	}

	return &response, nil
}

func statusFromStarlarkValue(value starlark.Value) (int, error) {
	slInt, ok := value.(starlark.Int)
	if !ok {
		return -1, fmt.Errorf("Response status_code key must be an int")
	}
	status_code, ok := slInt.Int64()
	if !ok {
		return -1, fmt.Errorf("Response status_code key has invalid integer")
	}
	return int(status_code), nil
}

func headersFromStarlarkValue(value starlark.Value) (http.Header, error) {
	slDict, ok := value.(*starlark.Dict)
	if !ok {
		return nil, fmt.Errorf("Response headers must be a dict")
	}

	header := make(http.Header)

	for _, keySlValue := range slDict.Keys() {
		keySlString, ok := keySlValue.(starlark.String)
		if !ok {
			return nil, fmt.Errorf("keys in headers dict must be string")
		}
		valueSlValue, _, err := slDict.Get(keySlValue)
		if err != nil {
			return nil, err
		}
		valueSlString, ok := valueSlValue.(starlark.String)
		if !ok {
			return nil, fmt.Errorf("values in headers dict must be string")
		}
		header.Set(keySlString.GoString(), valueSlString.GoString())
	}

	return header, nil
}

func bodyFromStarlarkValue(value starlark.Value) (io.ReadCloser, error) {
	slString, ok := value.(starlark.String)
	if !ok {
		return nil, fmt.Errorf("Response body must be a string")
	}
	return ioutil.NopCloser(strings.NewReader(slString.GoString())), nil
}

func (trans *StubbedHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, stub := range trans.stubs {
		match, err := stub.Matches(req)
		if err != nil {
			return nil, err
		}
		if match {
			return stub.Respond(req)
		}
	}
	if trans.allowRealHttp {
		return trans.realTransport.RoundTrip(req)
	}
	return nil, fmt.Errorf("Unstubbed HTTP requests not allowed.")
}

func (trans *StubbedHttpTransport) AddStub(thread *starlark.Thread, matcher *starlark.Function, responder *starlark.Function) {
	trans.stubs = append(trans.stubs, Stub{thread, matcher, responder})
}

func (trans *StubbedHttpTransport) ClearStubs() {
	trans.stubs = []Stub{}
}

func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		module = starlark.StringDict{
			ModuleName: &starlarkstruct.Module{
				Name: ModuleName,
				Members: starlark.StringDict{
					"assert_render_webp": starlark.NewBuiltin("assert_render_webp", assertRendersWebp),
					"save_webp":          starlark.NewBuiltin("save_webp", saveWebp),
					"stub_http":          starlark.NewBuiltin("stub_http", stubHTTP),
				},
			},
		}
	})

	return module, nil
}

func unpackRoots(rootables starlark.Value) ([]render.Root, error) {
	if root, ok := rootables.(render_runtime.Rootable); ok {
		return []render.Root{root.AsRenderRoot()}, nil
	} else if returnList, ok := rootables.(*starlark.List); ok {
		roots := make([]render.Root, returnList.Len())
		iter := returnList.Iterate()
		defer iter.Done()
		i := 0
		var listVal starlark.Value
		for iter.Next(&listVal) {
			if listValRoot, ok := listVal.(render_runtime.Rootable); ok {
				roots[i] = listValRoot.AsRenderRoot()
			} else {
				return nil, fmt.Errorf(
					"expected Root(s) but found: %s (at index %d)",
					listVal.Type(),
					i,
				)
			}
			i++
		}
		return roots, nil
	} else {
		return nil, fmt.Errorf("expected Root(s) but found: %s", rootables.Type())
	}
}

func assertRendersWebp(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		rootables starlark.Value
		expected  starlark.String
	)

	if err := starlark.UnpackArgs(
		"assert_render_webp",
		args, kwargs,
		"root", &rootables,
		"expected", &expected,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for bytes: %s", err)
	}
	roots, err := unpackRoots(rootables)
	if err != nil {
		return nil, err
	}

	expectedFilepath := expected.GoString()
	if !path.IsAbs(expectedFilepath) {
		workingDir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		scriptPath := thread.Local("scriptPath")
		if scriptPath != nil {
			absoluteScriptDir := path.Dir(path.Clean(path.Join(workingDir, scriptPath.(string))))
			expectedFilepath = path.Join(absoluteScriptDir, expectedFilepath)
		}
	}

	expectedRawWebp, err := ioutil.ReadFile(expectedFilepath)
	if err != nil {
		return nil, err
	}

	screens := encode.ScreensFromRoots(roots)
	maxDuration := 0
	actualRawWebp, err := screens.EncodeWebP(maxDuration)
	if err != nil {
		return nil, err
	}

	if bytes.Compare(actualRawWebp, expectedRawWebp) != 0 {
		actualRawWebpPath, err := writeTempFile("actual.*.webp", actualRawWebp)

		var errorMsg string
		if err != nil {
			errorMsg = "actual webp did not match expected webp. Unable to write actual webp file."
		} else {
			errorMsg = fmt.Sprintf("actual webp did not match expected webp. Actual: %s", actualRawWebpPath)
		}
		reportError(thread, errorMsg)
	}

	return starlark.None, nil
}

func writeTempFile(pattern string, data []byte) (string, error) {
	tmpFile, err := ioutil.TempFile("", pattern)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	_, err = tmpFile.Write(data)
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func saveWebp(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		rootables starlark.Value
		output    starlark.String
	)

	if err := starlark.UnpackArgs(
		"save_webp",
		args, kwargs,
		"root", &rootables,
		"output", &output,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for bytes: %s", err)
	}
	roots, err := unpackRoots(rootables)
	if err != nil {
		return nil, err
	}

	outFilePath := output.GoString()

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	scriptDir := path.Dir(path.Clean(path.Join(workingDir, thread.Local("scriptPath").(string))))
	absoluteOutFilePath := path.Join(scriptDir, outFilePath)

	screens := encode.ScreensFromRoots(roots)
	maxDuration := 0
	buf, err := screens.EncodeWebP(maxDuration)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(absoluteOutFilePath, buf, 0644)
	if err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func reportError(thread *starlark.Thread, msg string) {
	stk := thread.CallStack()
	stk.Pop()
	errorStr := fmt.Sprintf("%sError: %s", stk, msg)
	starlarktest.GetReporter(thread).Error(errorStr)
}

func stubHTTP(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		matcher   *starlark.Function
		responder *starlark.Function
	)

	if err := starlark.UnpackArgs(
		"assert_render_webp",
		args, kwargs,
		"matcher", &matcher,
		"responder", &responder,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for bytes: %s", err)
	}

	stubbedHttpTransport.AddStub(thread, matcher, responder)

	return starlark.None, nil
}

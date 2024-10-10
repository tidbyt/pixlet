package runtime

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

func TestInitHTTP(t *testing.T) {
	c := NewInMemoryCache()
	InitHTTP(c)

	b, err := os.ReadFile("testdata/httpcache.star")
	assert.NoError(t, err)

	app, err := NewApplet("httpcache.star", b)
	assert.NoError(t, err)
	assert.NotNil(t, app)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}

// TestDetermineTTL tests the DetermineTTL function.
func TestDetermineTTL(t *testing.T) {
	type test struct {
		ttl         int
		retryAfter  int
		resHeader   string
		statusCode  int
		method      string
		expectedTTL time.Duration
	}

	tests := map[string]test{
		"test request cache control headers": {
			ttl:         3600,
			resHeader:   "",
			statusCode:  200,
			method:      "GET",
			expectedTTL: 3600 * time.Second,
		},
		"test response cache control headers": {
			ttl:         0,
			resHeader:   "public, max-age=3600, s-maxage=7200, no-transform",
			statusCode:  200,
			method:      "GET",
			expectedTTL: 3600 * time.Second,
		},
		"test too long response cache control headers": {
			ttl:         0,
			resHeader:   "max-age=604800",
			statusCode:  200,
			method:      "GET",
			expectedTTL: 3600 * time.Second,
		},
		"test max-age of zero": {
			ttl:         0,
			resHeader:   "max-age=0",
			statusCode:  200,
			method:      "GET",
			expectedTTL: 5 * time.Second,
		},
		"test both request and response cache control headers": {
			ttl:         3600,
			resHeader:   "public, max-age=60, s-maxage=7200, no-transform",
			statusCode:  200,
			method:      "GET",
			expectedTTL: 3600 * time.Second,
		},
		"test 500 response code": {
			ttl:         3600,
			resHeader:   "",
			statusCode:  500,
			method:      "GET",
			expectedTTL: 5 * time.Second,
		},
		"test too low ttl": {
			ttl:         3,
			resHeader:   "",
			statusCode:  200,
			method:      "GET",
			expectedTTL: 5 * time.Second,
		},
		"test DELETE request": {
			ttl:         0,
			resHeader:   "",
			statusCode:  200,
			method:      "DELETE",
			expectedTTL: 5 * time.Second,
		},
		"test POST request configured with TTL": {
			ttl:         30,
			resHeader:   "",
			statusCode:  200,
			method:      "POST",
			expectedTTL: 30 * time.Second,
		},
		"test POST request without configured TTL": {
			ttl:         0,
			resHeader:   "",
			statusCode:  200,
			method:      "POST",
			expectedTTL: 5 * time.Second,
		},
		"test 429 request": {
			ttl:         30,
			retryAfter:  60,
			resHeader:   "",
			statusCode:  429,
			method:      "GET",
			expectedTTL: 60 * time.Second,
		},
		"test 429 request below minimum": {
			ttl:         30,
			retryAfter:  3,
			resHeader:   "",
			statusCode:  429,
			method:      "GET",
			expectedTTL: 5 * time.Second,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := &http.Request{
				Header: map[string][]string{
					"X-Tidbyt-Cache-Seconds": {fmt.Sprintf("%d", tc.ttl)},
				},
				Method: tc.method,
			}

			res := &http.Response{
				Header: map[string][]string{
					"Cache-Control": {tc.resHeader},
				},
				StatusCode: tc.statusCode,
			}

			if tc.retryAfter > 0 {
				res.Header.Set("Retry-After", fmt.Sprintf("%d", tc.retryAfter))
			}

			ttl := determineTTL(req, res)
			assert.Equal(t, tc.expectedTTL, ttl)
		})
	}
}

func TestDetermineTTLJitter(t *testing.T) {
	req := &http.Request{
		Header: map[string][]string{
			"X-Tidbyt-Cache-Seconds": {"60"},
		},
		Method: "GET",
	}

	res := &http.Response{
		StatusCode: 200,
	}

	rand.Seed(50)
	ttl := DetermineTTL(req, res)
	assert.Equal(t, 64, int(ttl.Seconds()))
}

func TestDetermineTTLNoHeaders(t *testing.T) {
	req := &http.Request{
		Method: "GET",
	}

	res := &http.Response{
		StatusCode: 200,
	}

	ttl := DetermineTTL(req, res)
	assert.Equal(t, MinRequestTTL, ttl)
}

func TestSetCookieOnRedirect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Requests to "/login" set a cookie and redirect to /destination
		if strings.HasSuffix(r.URL.Path, "/login") {
			if len(r.Cookies()) != 0 {
				t.Errorf("Cookie was already set on initial call")
			}
			w.Header().Set("Set-Cookie", "doodad=foobar; path=/; HttpOnly")
			w.Header().Set("Location", "/destination")
			w.WriteHeader(302)
			return
		}
		// Requests to /destination must have cookie set
		if strings.HasSuffix(r.URL.Path, "/destination") {
			c, err := r.Cookie("doodad")
			if err != nil {
				t.Errorf("Expected cookie `doodad` not present") // Occurs if client has no cookie jar
			}
			if c.Value != "foobar" {
				t.Errorf("Cookie `doodad` value mismatch. Expected foobar, got %s", c.Value)
			}
			if _, err := w.Write([]byte(`{"hello":"world"}`)); err != nil {
				t.Fatal(err)
			}
			return
		}
		t.Errorf("Unexpected path requested: %s", r.URL.Path)
	}))

	starlark.Universe["test_server_url"] = starlark.String(ts.URL)
	c := NewInMemoryCache()
	InitHTTP(c)

	b, err := os.ReadFile("testdata/httpredirect.star")
	assert.NoError(t, err)

	app, err := NewApplet("httpredirect.star", b)
	assert.NoError(t, err)

	_, err = app.Run(context.Background())
	assert.NoError(t, err)

	// Run it again and check that we're not using the same cookie jar
	// across executions. If we were, the first request would error out.
	app2, err := NewApplet("httpredirect.star", b)
	assert.NoError(t, err)

	_, err = app2.Run(context.Background())
	assert.NoError(t, err)
}

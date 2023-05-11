package runtime

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitHTTP(t *testing.T) {
	c := NewInMemoryCache()
	InitHTTP(c)

	b, err := os.ReadFile("testdata/httpcache.star")
	assert.NoError(t, err)

	app := &Applet{}
	err = app.Load("httpcache.star", b, nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
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
			ttl:         30,
			resHeader:   "",
			statusCode:  200,
			method:      "DELETE",
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

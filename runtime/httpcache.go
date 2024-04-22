package runtime

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"tidbyt.dev/pixlet/runtime/modules/starlarkhttp"
)

const (
	MinRequestTTL    = 5 * time.Second
	MaxResponseTTL   = 1 * time.Hour
	HTTPTimeout      = 5 * time.Second
	MaxResponseBytes = 20 * 1024 * 1024 // 20MB
	HTTPCachePrefix  = "httpcache"
	TTLHeader        = "X-Tidbyt-Cache-Seconds"
)

// Status codes that are cacheable as defined here:
// https://developer.mozilla.org/en-US/docs/Glossary/Cacheable
var cacheableStatusCodes = map[int]bool{
	200: true,
	201: true,
	202: true,
	203: true,
	204: true,
	206: true,
	300: true,
	301: true,
	404: true,
	405: true,
	410: true,
	414: true,
	501: true,
}

type cacheClient struct {
	cache     Cache
	transport http.RoundTripper
}

func InitHTTP(cache Cache) {
	cc := &cacheClient{
		cache:     cache,
		transport: http.DefaultTransport,
	}

	httpClient := &http.Client{
		Transport: cc,
		Timeout:   HTTPTimeout * 2,
	}
	starlarkhttp.StarlarkHTTPClient = httpClient
}

// RoundTrip is an approximation of what our internal HTTP proxy does. It should
// behave the same way, and any discrepancy should be considered a bug.
func (c *cacheClient) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	ctx, cancel := context.WithTimeout(ctx, HTTPTimeout)
	defer cancel() // need to do this to not leak a goroutine

	key, err := cacheKey(req, false)
	if err != nil {
		return nil, fmt.Errorf("failed to generate cache key: %w", err)
	}
	// TODO: remove once old cache entries expire
	keyWithTTL, err := cacheKey(req, true)
	if err != nil {
		return nil, fmt.Errorf("failed to generate cache key: %w", err)
	}

	if req.Method == "GET" || req.Method == "HEAD" || req.Method == "POST" {
		b, exists, err := c.cache.Get(nil, key)
		if exists && err == nil {
			if res, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(b)), req); err == nil {
				res.Header.Set("tidbyt-cache-status", "HIT")
				return res, nil
			}
		}
		// TODO: remove once old entries expire
		b, exists, err = c.cache.Get(nil, keyWithTTL)
		if exists && err == nil {
			if res, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(b)), req); err == nil {
				res.Header.Set("tidbyt-cache-status", "HIT")
				return res, nil
			}
		}
	}

	resp, err := c.transport.RoundTrip(req.WithContext(ctx))
	if err == nil {
		resp.Body = http.MaxBytesReader(nil, resp.Body, MaxResponseBytes)
	}

	if err == nil && (req.Method == "GET" || req.Method == "HEAD" || req.Method == "POST") {
		ser, err := httputil.DumpResponse(resp, true)
		if err != nil {
			// if httputil.DumpResponse fails, it leaves the response body in an
			// undefined state, so we cannot continue
			return nil, fmt.Errorf("failed to serialize response for cache: %s", resp.Status)
		}

		ttl := DetermineTTL(req, resp)
		c.cache.Set(nil, key, ser, int64(ttl.Seconds()))
		resp.Header.Set("tidbyt-cache-status", "MISS")
	}

	return resp, err
}

func cacheKey(req *http.Request, keep_ttl bool) (string, error) {
	// TODO: remove keep_ttl param and make this always happen.
	ttl := req.Header.Get(TTLHeader)
	if !keep_ttl {
		req.Header.Del(TTLHeader)
	}
	r, err := httputil.DumpRequest(req, true)
	if err != nil {
		return "", fmt.Errorf("%s: %w", "failed to serialize request", err)
	}
	if ttl != "" {
		req.Header.Set(TTLHeader, ttl)
	}

	h := sha256.Sum256(r)
	key := hex.EncodeToString(h[:])

	app := req.Header.Get("X-Tidbyt-App")
	if app == "" {
		return key, nil
	}

	return fmt.Sprintf("%s:%s:%s", HTTPCachePrefix, app, key), nil
}

// DetermineTTL determines the TTL for a request based on the request and
// response. We first check request method / response status code to determine
// if we should actually cache the response. Then we check the headers passed in
// from starlark to see if the user configured a TTL. Finally, if the response
// is cachable but the developer didn't configure a TTL, we check the response
// to get a hint at what the TTL should be.
func DetermineTTL(req *http.Request, resp *http.Response) time.Duration {
	ttl := determineTTL(req, resp)

	// Jitter the TTL by 10% and double check that it's still greater than the
	// minimum TTL. If it's not, return the minimum TTL. The main thing we want
	// to avoid is a TTL of 0 given it will be cached forever.
	ttl = jitterDuration(ttl)
	if ttl < MinRequestTTL {
		return MinRequestTTL
	}

	return ttl
}

func determineTTL(req *http.Request, resp *http.Response) time.Duration {
	// If the response is a 429, we want to cache the response for the duration
	// the remote server told us to wait before retrying.
	if resp.StatusCode == 429 {
		retry := MinRequestTTL
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			if intValue, err := strconv.Atoi(retryAfter); err == nil {
				retry = time.Duration(intValue) * time.Second
			}
		}

		if retry < MinRequestTTL {
			return MinRequestTTL
		}

		return retry
	}

	// Check the status code to determine if the response is cacheable.
	_, ok := cacheableStatusCodes[resp.StatusCode]
	if !ok {
		return MinRequestTTL
	}

	// Determine the TTL based on the developer's configuration.
	ttl := determineDeveloperTTL(req)

	// We don't want to cache POST requests unless the developer explicitly
	// requests it.
	if ttl == 0 && !(req.Method == "GET" || req.Method == "HEAD") {
		return MinRequestTTL
	}

	// If the developer didn't configure a TTL, determine the TTL based on the
	// response.
	if ttl == 0 {
		ttl = determineResponseTTL(resp)
	}

	if ttl < MinRequestTTL {
		return MinRequestTTL
	}

	return ttl
}

func jitterDuration(duration time.Duration) time.Duration {
	jitter := int64(float64(duration) * 0.1)
	randomJitter := rand.Int63n(2*jitter+1) - jitter
	return time.Duration(duration + time.Duration(randomJitter))
}

func determineResponseTTL(resp *http.Response) time.Duration {
	resHeader := parseCacheControl(resp.Header.Get("Cache-Control"))
	value, ok := resHeader["max-age"]
	if ok {
		if intValue, ok := value.(int); ok {
			ttl := time.Duration(intValue) * time.Second

			// If we're using a response TTL, we're making the assumption that
			// the remote server is providing a reasonable TTL that a developer
			// didn't configure. In the case of weathermap, the TTL is 1 week,
			// but the developer is requesting a new map ever hour. So while the
			// old map _is_ valid for a week, we the app only cares about it for
			// one hour.
			if ttl > MaxResponseTTL {
				return MaxResponseTTL
			}
			return ttl
		}
	}

	return 0
}

func determineDeveloperTTL(req *http.Request) time.Duration {
	ttlHeader := req.Header.Get("X-Tidbyt-Cache-Seconds")
	if ttlHeader != "" {
		if intValue, err := strconv.Atoi(ttlHeader); err == nil {
			return time.Duration(intValue) * time.Second
		}
	}

	return 0
}

func parseCacheControl(header string) map[string]interface{} {
	directives := make(map[string]interface{})

	for _, directive := range strings.Split(header, ",") {
		parts := strings.SplitN(strings.TrimSpace(directive), "=", 2)

		key := strings.ToLower(parts[0])
		var value interface{} = true

		if len(parts) > 1 {
			value = parts[1]
			if intValue, err := strconv.Atoi(parts[1]); err == nil {
				value = intValue
			}
		}

		directives[key] = value
	}

	return directives
}

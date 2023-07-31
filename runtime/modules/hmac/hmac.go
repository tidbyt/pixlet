package hmac

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"sync"
	"time"

	godfe "github.com/newm4n/go-dfe"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	ModuleName = "hmac"
)

var (
	once        sync.Once
	module      starlark.StringDict
	empty       time.Time
	translation *godfe.PatternTranslation
)

func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		translation = godfe.NewPatternTranslation()
		module = starlark.StringDict{
			ModuleName: &starlarkstruct.Module{
				Name: ModuleName,
				Members: starlark.StringDict{
					"md5":    starlark.NewBuiltin("md5", fnHmac(md5.New)),
					"sha1":   starlark.NewBuiltin("sha1", fnHmac(sha1.New)),
					"sha256": starlark.NewBuiltin("sha256", fnHmac(sha256.New)),
				},
			},
		}
	})

	return module, nil
}

func fnHmac(hashFunc func() hash.Hash) func(*starlark.Thread, *starlark.Builtin, starlark.Tuple, []starlark.Tuple) (starlark.Value, error) {
	return func(t *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var (
			key starlark.Value
			s   starlark.String
			binary bool = false
		)
		if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "key", &key, "s", &s, "binary?", &binary); err != nil {
			return nil, err
		}

		var byteKey []byte
		switch key := key.(type) {
		case starlark.String:
			byteKey = []byte(string(key))
		case starlark.Bytes:
			byteKey = []byte(key)
		default:
			return nil, fmt.Errorf("%s: for parameter 1 got %s, want string or bytes", fn.Name(), key.Type())
		}

		h := hmac.New(hashFunc, byteKey)

		if _, err := h.Write([]byte(string(s))); err != nil {
			return starlark.None, err
		}

		digest := h.Sum(nil)
		if binary {
			return starlark.Bytes(digest), nil
		}
		return starlark.String(fmt.Sprintf("%x", digest)), nil
	}
}
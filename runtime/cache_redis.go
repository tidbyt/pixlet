package runtime

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var rdb *redis.Client

var (
	cacheRedisOnce   sync.Once
	cacheRedisModule starlark.StringDict
)

func LoadCacheRedisModule() (starlark.StringDict, error) {
	cacheRedisOnce.Do(func() {
		cacheRedisModule = starlark.StringDict{
			"cache_redis": &starlarkstruct.Module{
				Name: "cache_redis",
				Members: starlark.StringDict{
					"connect": starlark.NewBuiltin("connect", cacheRedisConnect),
					"get":     starlark.NewBuiltin("get", cacheRedisGet),
					"set":     starlark.NewBuiltin("set", cacheRedisSet),
				},
			},
		}
	})

	return cacheRedisModule, nil
}

func cacheRedisConnect(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		host     starlark.String
		username starlark.String
		password starlark.String
	)

	if err := starlark.UnpackArgs(
		"connect",
		args, kwargs,
		"host", &host,
		"username", &username,
		"password", &password,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for cache_redis.connect: %v", err)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     host.GoString(),
		Username: username.GoString(),
		Password: password.GoString(), // no password set
		DB:       0,                   // use default DB
	})

	return starlark.String("Connected"), nil
}

func cacheRedisGet(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var key starlark.String

	if err := starlark.UnpackArgs(
		"get",
		args, kwargs,
		"key", &key,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for cache.get: %v", err)
	}

	cacheKey := scopedCacheKey(thread, key)

	ctx := context.TODO()
	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return starlark.None, nil
	} else if err != nil {
		log.Printf("getting %s from cache: %v", cacheKey, err)
		return starlark.None, nil
	} else {
		return starlark.String(val), nil
	}
}

func cacheRedisSet(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		key starlark.String
		val starlark.String
		ttl starlark.Int
	)

	if err := starlark.UnpackArgs(
		"set",
		args, kwargs,
		"key", &key,
		"value", &val,
		"ttl_seconds?", &ttl,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for cache.set: %v", err)
	}

	cacheKey := scopedCacheKey(thread, key)

	ttl64, ok := ttl.Int64()
	if !ok {
		return nil, fmt.Errorf("ttl_seconds must be valid integer (not %s)", ttl.String())
	}

	if ttl64 < 0 {
		return nil, fmt.Errorf("ttl_seconds cannot be negative")
	}

	if ttl64 == 0 {
		ttl64 = DefaultExpirationSeconds
	}

	ctx := context.TODO()

	err := rdb.Set(ctx, cacheKey, val.GoString(), time.Duration(ttl.BigInt().Int64()*int64(time.Second))).Err()
	if err != nil {
		log.Printf("setting %s in cache: %v", cacheKey, err)
	}

	return starlark.None, nil
}

package random

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	ModuleName = "random"
)

var (
	once   sync.Once
	module starlark.StringDict
)

func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		rand.Seed(time.Now().UnixNano())
		module = starlark.StringDict{
			ModuleName: &starlarkstruct.Module{
				Name: ModuleName,
				Members: starlark.StringDict{
					"number": starlark.NewBuiltin("number", randomNumber),
				},
			},
		}
	})

	return module, nil
}

func randomNumber(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starMin starlark.Int
		starMax starlark.Int
	)

	if err := starlark.UnpackArgs(
		"number",
		args, kwargs,
		"min", &starMin,
		"max", &starMax,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for random number: %w", err)
	}

	min, ok := starMin.Int64()
	if !ok {
		return nil, fmt.Errorf("casting min to an int64")

	}

	max, ok := starMax.Int64()
	if !ok {
		return nil, fmt.Errorf("casting max to an int64")

	}

	if min < 0 {
		return nil, fmt.Errorf("min has to be 0 or greater")
	}

	if max < min {
		return nil, fmt.Errorf("max is less then min")
	}

	number := rand.Int63n(max-min+1) + min

	return starlark.MakeInt64(number), nil
}

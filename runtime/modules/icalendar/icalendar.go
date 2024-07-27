package icalendar

import (
	"sync"
	"time"

	godfe "github.com/newm4n/go-dfe"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	ModuleName = "icalendar"
)

var (
	once        sync.Once
	module      starlark.StringDict
	empty       time.Time
	translation *godfe.PatternTranslation
)

func LoadModule() (starlark.StringDict, error) {
	translation = godfe.NewPatternTranslation()
	once.Do(func() {
		module = starlark.StringDict{
			ModuleName: &starlarkstruct.Module{
				Name: ModuleName,
				Members: starlark.StringDict{
					"time":          starlark.NewBuiltin("time", times),
					"findNextEvent": starlark.NewBuiltin("findNextEvent", findNextEvent),
				},
			},
		}
	})

	return module, nil
}

func times(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.False, nil
}

func findNextEvent(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.False, nil
}

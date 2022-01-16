package sunrise

import (
	"fmt"
	"sync"
	"time"

	gosunrise "github.com/nathan-osman/go-sunrise"
	startime "go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	ModuleName = "sunrise"
)

var (
	once   sync.Once
	module starlark.StringDict
)

func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		module = starlark.StringDict{
			ModuleName: &starlarkstruct.Module{
				Name: ModuleName,
				Members: starlark.StringDict{
					"sunrise": starlark.NewBuiltin("sunrise", sunrise),
					"sunset":  starlark.NewBuiltin("sunset", sunset),
				},
			},
		}
	})

	return module, nil
}

func sunrise(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starLat  starlark.Float
		starLng  starlark.Float
		starDate startime.Time
	)

	if err := starlark.UnpackArgs(
		"sunrise",
		args, kwargs,
		"lat", &starLat,
		"lng", &starLng,
		"date", &starDate,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for sunrise: %s", err)
	}

	lat := float64(starLat)
	lng := float64(starLng)
	date := time.Time(starDate)
	rise, _ := gosunrise.SunriseSunset(lat, lng, date.Year(), date.Month(), date.Day())

	return startime.Time(rise), nil
}

func sunset(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starLat  starlark.Float
		starLng  starlark.Float
		starDate startime.Time
	)

	if err := starlark.UnpackArgs(
		"sunset",
		args, kwargs,
		"lat", &starLat,
		"lng", &starLng,
		"date", &starDate,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for sunset: %s", err)
	}

	lat := float64(starLat)
	lng := float64(starLng)
	date := time.Time(starDate)
	_, set := gosunrise.SunriseSunset(lat, lng, date.Year(), date.Month(), date.Day())

	return startime.Time(set), nil
}

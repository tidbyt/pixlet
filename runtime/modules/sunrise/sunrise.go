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
	empty  time.Time
)

func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		module = starlark.StringDict{
			ModuleName: &starlarkstruct.Module{
				Name: ModuleName,
				Members: starlark.StringDict{
					"sunrise":        starlark.NewBuiltin("sunrise", sunrise),
					"sunset":         starlark.NewBuiltin("sunset", sunset),
					"elevation":      starlark.NewBuiltin("elevation", elevation),
					"elevation_time": starlark.NewBuiltin("elevation_time", elevation_time),
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
	if rise == empty {
		return starlark.None, nil
	}

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
	if set == empty {
		return starlark.None, nil
	}

	return startime.Time(set), nil
}

func elevation(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starLat  starlark.Float
		starLng  starlark.Float
		starTime startime.Time
	)

	if err := starlark.UnpackArgs(
		"elevation",
		args, kwargs,
		"lat", &starLat,
		"lng", &starLng,
		"time", &starTime,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for elevation: %s", err)
	}

	lat := float64(starLat)
	lng := float64(starLng)
	when := time.Time(starTime)

	elev := gosunrise.Elevation(lat, lng, when)
	return starlark.Float(elev), nil
}

func elevation_time(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starLat  starlark.Float
		starLng  starlark.Float
		starElev starlark.Float
		starDate startime.Time
	)

	if err := starlark.UnpackArgs(
		"elevation_time",
		args, kwargs,
		"lat", &starLat,
		"lng", &starLng,
		"elev", &starElev,
		"date", &starDate,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for elevation: %s", err)
	}

	lat := float64(starLat)
	lng := float64(starLng)
	elev := float64(starElev)
	date := time.Time(starDate)

	morning, evening := gosunrise.TimeOfElevation(lat, lng, elev, date.Year(), date.Month(), date.Day())
	if morning == empty || evening == empty {
		return starlark.None, nil
	}
	starMorning := startime.Time(morning)
	starEvening := startime.Time(evening)

	return starlark.Tuple([]starlark.Value{starMorning, starEvening}), nil
}

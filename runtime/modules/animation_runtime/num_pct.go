package animation_runtime

import (
	"fmt"
	"reflect"

	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/render/animation"
)

func NumberOrPercentageFromStarlark(value starlark.Value, min, max float64, mapping map[string]float64) (animation.NumberOrPercentage, error) {
	if val, ok := starlark.AsFloat(value); ok {
		if min <= val && val <= max {
			return animation.Number{val}, nil
		} else {
			return nil, fmt.Errorf("invalid range for number: %f (expected number in range [0.0, 1.0])", val)
		}
	} else if str, ok := starlark.AsString(value); ok {
		return animation.ParsePercentage(str, mapping)
	}

	return nil, fmt.Errorf("invalid type for number or percentage: %s (expected number or string)", value.Type())
}

func PercentageFromStarlark(value starlark.Value, mapping map[string]float64) (animation.Percentage, error) {
	if val, err := NumberOrPercentageFromStarlark(value, 0.0, 1.0, mapping); err == nil {
		if n, ok := val.(animation.Number); ok {
			return animation.Percentage{n.Value}, nil
		} else if p, ok := val.(animation.Percentage); ok {
			return p, nil
		}

		return animation.Percentage{}, fmt.Errorf("invalid type returned: %v (unreachable code)", reflect.TypeOf(val))
	} else {
		return animation.Percentage{}, err
	}
}

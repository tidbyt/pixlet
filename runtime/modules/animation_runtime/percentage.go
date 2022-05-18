package animation_runtime

import (
	"fmt"

	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/render/animation"
)

func PercentageFromStarlark(value starlark.Value) (animation.Percentage, error) {
	if val, ok := starlark.AsFloat(value); ok {
		if 0.0 <= val && val <= 1.0 {
			return animation.Percentage{val}, nil
		} else {
			return animation.Percentage{}, fmt.Errorf("invalid range for percentage: %f (expected number in range [0.0, 1.0])", val)
		}
	}

	return animation.Percentage{}, fmt.Errorf("invalid type for percentage: %s (expected number in range [0.0, 1.0])", value.Type())
}

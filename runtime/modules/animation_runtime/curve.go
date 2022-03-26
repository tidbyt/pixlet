package animation_runtime

import (
	"fmt"

	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/render/animation"
)

func CurveFromStarlark(value starlark.Value) (animation.Curve, error) {
	if str, ok := value.(starlark.String); ok {
		if str.Len() == 0 {
			return animation.LinearCurve{}, nil
		} else if curve, err := animation.ParseCurve(str.GoString()); err == nil {
			return curve, nil
		} else {
			return animation.LinearCurve{}, fmt.Errorf("curve is not a valid curve string: %s", str.GoString())
		}
	}

	if fn, ok := value.(*starlark.Function); ok {
		if fn.NumParams() != 1 || fn.NumKwonlyParams() != 0 {
			return animation.LinearCurve{}, fmt.Errorf("invalid number of parameters to curve function: %s", fn.String())
		}

		return animation.CustomCurve{Function: fn}, nil
	}

	return animation.LinearCurve{}, nil
}

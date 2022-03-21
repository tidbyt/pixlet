package animation_runtime

import (
	"fmt"

	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/render/animation"
)

func CurveFromStarlark(value starlark.Value) (animation.Curve, error) {
	if str, ok := starlark.AsString(value); ok {
		return animation.ParseCurve(str)
	}

	return nil, fmt.Errorf("invalid type for curve: %s (expected string)", value.Type())
}

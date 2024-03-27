package animation

import (
	"tidbyt.dev/pixlet/render/canvas"
)

type Transform interface {
	Apply(ctx canvas.Canvas, origin Vec2f, rounding Rounding)
	Interpolate(other Transform, progress float64) (result Transform, ok bool)
}

func ExtendTransforms(lhs []Transform, rhs []Transform) []Transform {
	for i, transform := range rhs {
		if i >= len(lhs) {
			switch transform.(type) {
			case Translate:
				lhs = append(lhs, TranslateDefault)
			case Scale:
				lhs = append(lhs, ScaleDefault)
			case Rotate:
				lhs = append(lhs, RotateDefault)
			}
		}
	}

	return lhs
}

// See: https://www.w3.org/TR/css-transforms-1/#interpolation-of-transforms
func InterpolateTransforms(lhs, rhs []Transform, progress float64) (result []Transform, ok bool) {
	if len(lhs) == 0 && len(rhs) == 0 {
		return make([]Transform, 0), true
	}

	if len(lhs) < len(rhs) {
		lhs = ExtendTransforms(lhs, rhs)
	} else if len(lhs) > len(rhs) {
		rhs = ExtendTransforms(rhs, lhs)
	}

	result = make([]Transform, 0)

	for i := 0; i < len(lhs); i++ {
		if t, ok := lhs[i].Interpolate(rhs[i], progress); ok {
			result = append(result, t)
		} else {
			// This is the point where remaining transforms would be composed into matrices
			// and interpolated on the matrix level, but for simplicity is not supported.
			return make([]Transform, 0), false
		}
	}

	return result, true
}

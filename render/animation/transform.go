package animation

import (
	"github.com/fogleman/gg"
)

type Transform interface {
	Apply(ctx *gg.Context, origin Vec2f, rounding Rounding)
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
func InterpolateTransforms(lhs, rhs []Transform, progress float64) (result []Transform) {
	if len(lhs) == 0 && len(rhs) == 0 {
		return make([]Transform, 0)
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
			m0, m1 :=
				Matrix{matrix: ComposeMatrix(lhs[i:])},
				Matrix{matrix: ComposeMatrix(rhs[i:])}
			m2, _ := m0.Interpolate(m1, progress)
			result = append(result, m2)
			break
		}
	}

	return result
}

package animation

import (
	"math"

	"github.com/fogleman/gg"
)

type Matrix struct {
	matrix gg.Matrix
}

func (self Matrix) Apply(ctx *gg.Context, origin Vec2f, rounding Rounding) {
	translate, scale, angle := DecomposeMatrix(self.matrix)

	ctx.RotateAbout(gg.Radians(angle), origin.X, origin.Y)
	ctx.ScaleAbout(scale.X, scale.Y, origin.X, origin.Y)
	ctx.Translate(math.Round(translate.X), math.Round(translate.Y))
}

func (self Matrix) Interpolate(other Transform, progress float64) (result Transform, ok bool) {
	if other, ok := other.(Matrix); ok {
		t0, s0, a0 := DecomposeMatrix(self.matrix)
		t1, s1, a1 := DecomposeMatrix(other.matrix)
		t2, s2, a2 := interpolateMatrix(t0, s0, a0, t1, s1, a1, progress)
		result := gg.Identity().Rotate(gg.Radians(a2)).Scale(s2.X, s2.Y).Translate(t2.X, t2.Y)
		return Matrix{matrix: result}, true
	}

	return MatrixDefault, false
}

var MatrixDefault = Matrix{matrix: gg.Identity()}

func ComposeMatrix(transforms []Transform) gg.Matrix {
	result := gg.Identity()

	for _, transform := range transforms {
		if t, ok := transform.(Translate); ok {
			result = result.Translate(t.X, t.Y)
		} else if s, ok := transform.(Scale); ok {
			result = result.Scale(s.X, s.Y)
		} else if r, ok := transform.(Rotate); ok {
			result = result.Rotate(gg.Radians(r.Angle))
		}
	}

	return result
}

// See: https://www.w3.org/TR/css-transforms-1/#decomposing-a-2d-matrix
func DecomposeMatrix(m gg.Matrix) (t, s Vec2f, a float64) {
	t = Vec2f{m.X0, m.Y0}
	s = Vec2f{
		math.Sqrt(math.Pow(m.XX, 2.0) + math.Pow(m.YX, 2.0)),
		math.Sqrt(math.Pow(m.XY, 2.0) + math.Pow(m.YY, 2.0)),
	}

	// If determinant is negative, one axis was flipped.
	determinant := m.XX*m.YY - m.YX*m.XY
	if determinant < 0 {
		// Flip axis with minimum unit vector dot product.
		if m.XX < m.YY {
			s.X = -s.X
		} else {
			s.Y = -s.Y
		}
	}

	a = gg.Degrees(math.Atan2(m.YX, m.XX))

	return
}

// See: https://www.w3.org/TR/css-transforms-1/#interpolation-of-decomposed-2d-matrix-values
func interpolateMatrix(
	t0, s0 Vec2f, a0 float64,
	t1, s1 Vec2f, a1 float64,
	progress float64,
) (t2, s2 Vec2f, a2 float64) {
	// If x-axis of one is flipped, and y-axis of the other, convert to an unflipped rotation.
	if (s0.X < 0.0 && s1.Y < 0.0) || (s0.Y < 0.0 && s1.X < 0.0) {
		s0.X = -s0.X
		s0.Y = -s0.Y

		if a0 < 0.0 {
			a0 += 180.0
		} else {
			a0 -= 180.0
		}
	}

	// Don’t rotate the long way around.
	if a0 == 0.0 {
		a0 = 360.0
	}

	if a1 == 0.0 {
		a1 = 360.0
	}

	if math.Abs(a0-a1) > 180 {
		if a0 > a1 {
			a0 -= 360
		} else {
			a1 -= 360
		}
	}

	t2 = t0.Lerp(t1, progress)
	s2 = s0.Lerp(s1, progress)
	a2 = Lerp(a0, a1, progress)

	return
}

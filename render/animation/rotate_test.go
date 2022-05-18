package animation

import (
	"testing"
)

func assertInterpolateRotate(
	t *testing.T,
	expected float64,
	from float64,
	to float64,
	progress float64,
) {
	AssertInterpolate(t, Rotate{Angle: expected}, Rotate{Angle: from}, Rotate{Angle: to}, progress)
}

func TestInterpolateRotate(t *testing.T) {
	from := 0.0
	to := 360.0

	assertInterpolateRotate(t, 0.0, from, to, 0.0)
	assertInterpolateRotate(t, 36.0, from, to, 0.1)
	assertInterpolateRotate(t, 72.0, from, to, 0.2)
	assertInterpolateRotate(t, 90.0, from, to, 0.25)
	assertInterpolateRotate(t, 180.0, from, to, 0.5)
	assertInterpolateRotate(t, 360.0, from, to, 1.0)
	assertInterpolateRotate(t, 720.0, from, to, 2.0)
	assertInterpolateRotate(t, -360.0, from, to, -1.0)
}

package animation

import (
	"testing"
)

func assertInterpolateTranslate(
	t *testing.T,
	expected Vec2f,
	from Vec2f,
	to Vec2f,
	progress float64,
) {
	AssertInterpolate(t, Translate{expected}, Translate{from}, Translate{to}, progress)
}

func TestInterpolateTranslate(t *testing.T) {
	from := Vec2f{X: 0.0, Y: 0.0}
	to := Vec2f{X: 100.0, Y: 200.0}

	assertInterpolateTranslate(t, Vec2f{X: 0.0, Y: 0.0}, from, to, 0.0)
	assertInterpolateTranslate(t, Vec2f{X: 10.0, Y: 20.0}, from, to, 0.1)
	assertInterpolateTranslate(t, Vec2f{X: 33.0, Y: 66.0}, from, to, 0.33)
	assertInterpolateTranslate(t, Vec2f{X: 100.0, Y: 200.0}, from, to, 1.0)
	assertInterpolateTranslate(t, Vec2f{X: 1337.0, Y: 2674.0}, from, to, 13.37)
}

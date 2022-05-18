package animation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertInterpolate(
	t *testing.T,
	expected Transform,
	from Transform,
	to Transform,
	progress float64,
) {
	actual, ok := from.Interpolate(to, progress)
	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestExtendTransforms(t *testing.T) {
	assert.Equal(t,
		[]Transform{TranslateDefault, ScaleDefault, RotateDefault},
		ExtendTransforms(
			[]Transform{},
			[]Transform{TranslateDefault, ScaleDefault, RotateDefault}))

	assert.Equal(t,
		[]Transform{Translate{Vec2f{X: 5.0, Y: 5.0}}, ScaleDefault, RotateDefault},
		ExtendTransforms(
			[]Transform{Translate{Vec2f{X: 5.0, Y: 5.0}}},
			[]Transform{TranslateDefault, ScaleDefault, RotateDefault}))

	assert.Equal(t,
		[]Transform{Scale{Vec2f{X: 2.5, Y: 2.5}}, RotateDefault},
		ExtendTransforms(
			[]Transform{Scale{Vec2f{X: 2.5, Y: 2.5}}},
			[]Transform{ScaleDefault, RotateDefault}))
}

func TestInterpolateTransforms(t *testing.T) {
	// Matching transforms are interpolated.
	result, ok := InterpolateTransforms(
		[]Transform{Translate{Vec2f{X: 0.0, Y: 0.0}}, Scale{Vec2f{X: 1.0, Y: 1.0}}, Rotate{Angle: 0.0}},
		[]Transform{Translate{Vec2f{X: 5.0, Y: 5.0}}, Scale{Vec2f{X: 2.0, Y: 2.0}}, Rotate{Angle: 180.0}},
		0.5)

	assert.True(t, ok)
	assert.Equal(t, []Transform{Translate{Vec2f{X: 2.5, Y: 2.5}}, Scale{Vec2f{X: 1.5, Y: 1.5}}, Rotate{Angle: 90.0}}, result)

	// Missing transforms are extended by default transforms of that type.
	result, ok = InterpolateTransforms(
		[]Transform{Translate{Vec2f{X: 0.0, Y: 0.0}}, Scale{Vec2f{X: 2.0, Y: 2.0}}},
		[]Transform{Translate{Vec2f{X: 5.0, Y: 5.0}}},
		0.5)

	assert.True(t, ok)
	assert.Equal(t, []Transform{Translate{Vec2f{X: 2.5, Y: 2.5}}, Scale{Vec2f{X: 1.5, Y: 1.5}}}, result)

	// Mismatched transforms are not supported for simplicity reasons.
	result, ok = InterpolateTransforms(
		[]Transform{Scale{Vec2f{X: 1.0, Y: 1.0}}, Translate{Vec2f{X: 100, Y: 100}}},
		[]Transform{Rotate{Angle: 90.0}, Scale{Vec2f{X: 2.0, Y: 2.0}}},
		0.5)

	assert.False(t, ok)
}

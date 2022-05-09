package animation

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrigin(t *testing.T) {
	r := image.Rect(0, 0, 100, 100)

	assert.Equal(t, Vec2f{X: 0.0, Y: 0.0}, Origin{X: Percentage{0.0}, Y: Percentage{0.0}}.Transform(r))
	assert.Equal(t, Vec2f{X: 33.0, Y: 33.0}, Origin{X: Percentage{0.33}, Y: Percentage{0.33}}.Transform(r))
	assert.Equal(t, Vec2f{X: 50.0, Y: 50.0}, Origin{X: Percentage{0.5}, Y: Percentage{0.5}}.Transform(r))

	v := Origin{X: Percentage{0.666}, Y: Percentage{0.666}}.Transform(r)
	assert.InDelta(t, 66.6, v.X, 0.00001)
	assert.InDelta(t, 66.6, v.Y, 0.00001)

	assert.Equal(t, Vec2f{X: 100.0, Y: 100.0}, Origin{X: Percentage{1.0}, Y: Percentage{1.0}}.Transform(r))
}

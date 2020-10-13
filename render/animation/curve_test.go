package animation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinearCurve(t *testing.T) {
	lc := LinearCurve{}

	assert.InDelta(t, 0, lc.Transform(0), 0.00001)
	assert.InDelta(t, 1, lc.Transform(1), 0.00001)
	assert.InDelta(t, 0.31416, lc.Transform(0.31416), 0.00001)
	assert.InDelta(t, 0.5, lc.Transform(0.5), 0.00001)
}

func TestEaseCurves(t *testing.T) {
	assert.InDelta(t, 0.000000, EaseIn.Transform(0.000000), 0.01)
	assert.InDelta(t, 0.104000, EaseIn.Transform(0.219200), 0.01)
	assert.InDelta(t, 0.352000, EaseIn.Transform(0.481600), 0.01)
	assert.InDelta(t, 0.648000, EaseIn.Transform(0.734400), 0.01)
	assert.InDelta(t, 0.896000, EaseIn.Transform(0.924800), 0.01)
	assert.InDelta(t, 1.000000, EaseIn.Transform(1.000000), 0.01)

	/* Temporarily using a different curve
	assert.InDelta(t, 0.000000, EaseInOut.Transform(0.000000), 0.01)
	assert.InDelta(t, 0.104000, EaseInOut.Transform(0.123200), 0.01)
	assert.InDelta(t, 0.352000, EaseInOut.Transform(0.193600), 0.01)
	assert.InDelta(t, 0.648000, EaseInOut.Transform(0.302400), 0.01)
	assert.InDelta(t, 0.896000, EaseInOut.Transform(0.540800), 0.01)
	assert.InDelta(t, 1.000000, EaseInOut.Transform(1.000000), 0.01)
	*/

	assert.InDelta(t, 0.000000, EaseOut.Transform(0.000000), 0.01)
	assert.InDelta(t, 0.104000, EaseOut.Transform(0.008000), 0.01)
	assert.InDelta(t, 0.352000, EaseOut.Transform(0.064000), 0.01)
	assert.InDelta(t, 0.648000, EaseOut.Transform(0.216000), 0.01)
	assert.InDelta(t, 0.896000, EaseOut.Transform(0.512000), 0.01)
	assert.InDelta(t, 1.000000, EaseOut.Transform(1.000000), 0.01)
}

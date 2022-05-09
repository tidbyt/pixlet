package animation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRescale(t *testing.T) {
	// [0.0, 1.0] -> [0.0, 100.0]
	assert.Equal(t, 0.0, Rescale(0.0, 1.0, 0.0, 100.0, 0.0))
	assert.Equal(t, 50.0, Rescale(0.0, 1.0, 0.0, 100.0, 0.5))
	assert.Equal(t, 100.0, Rescale(0.0, 1.0, 0.0, 100.0, 1.0))

	// [0.0, 1.0] -> [-50.0, 100.0]
	assert.Equal(t, -50.0, Rescale(0.0, 1.0, -50.0, 100.0, 0.0))
	assert.Equal(t, 25.0, Rescale(0.0, 1.0, -50.0, 100.0, 0.5))
	assert.Equal(t, 100.0, Rescale(0.0, 1.0, -50.0, 100.0, 1.0))

	// [-33.2, 10.71] -> [0.0, 1.0]
	assert.Equal(t, 0.5283534502391255, Rescale(-33.2, 10.71, 0.0, 1.0, -10))
	assert.Equal(t, 0.8699612844454566, Rescale(-33.2, 10.71, 0.0, 1.0, 5))
	assert.Equal(t, 0.9838305625142336, Rescale(-33.2, 10.71, 0.0, 1.0, 10))
}

func TestLerp(t *testing.T) {
	// [0.0, 1.0]
	assert.Equal(t, 0.0, Lerp(0.0, 1.0, 0.0))
	assert.Equal(t, 0.1, Lerp(0.0, 1.0, 0.1))
	assert.Equal(t, 0.33, Lerp(0.0, 1.0, 0.33))
	assert.Equal(t, 0.5, Lerp(0.0, 1.0, 0.5))
	assert.Equal(t, 0.7533, Lerp(0.0, 1.0, 0.7533))
	assert.Equal(t, 1.0, Lerp(0.0, 1.0, 1.0))

	// [-1.0, 1.0]
	assert.Equal(t, -1.0, Lerp(-1.0, 1.0, 0.0))
	assert.Equal(t, -0.8, Lerp(-1.0, 1.0, 0.1))
	assert.Equal(t, -0.33999999999999997, Lerp(-1.0, 1.0, 0.33))
	assert.Equal(t, 0.0, Lerp(-1.0, 1.0, 0.5))
	assert.Equal(t, 0.5065999999999999, Lerp(-1.0, 1.0, 0.7533))
	assert.Equal(t, 1.0, Lerp(-1.0, 1.0, 1.0))

	// [0.0, 42.1337]
	assert.Equal(t, 0.0, Lerp(0.0, 42.1337, 0.0))
	assert.Equal(t, 4.21337, Lerp(0.0, 42.1337, 0.1))
	assert.Equal(t, 13.904121, Lerp(0.0, 42.1337, 0.33))
	assert.Equal(t, 21.06685, Lerp(0.0, 42.1337, 0.5))
	assert.Equal(t, 31.73931621, Lerp(0.0, 42.1337, 0.7533))
	assert.Equal(t, 42.1337, Lerp(0.0, 42.1337, 1.0))
}

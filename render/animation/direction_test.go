package animation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectionNormal(t *testing.T) {
	d := DirectionNormal

	assert.Equal(t, 7, d.FrameCount(1, 5))

	// Delay
	assert.Equal(t, 0.0, d.Progress(1, 5, 1.0, 0))

	// Progress
	assert.Equal(t, 0.0, d.Progress(1, 5, 1.0, 1))
	assert.Equal(t, 0.25, d.Progress(1, 5, 1.0, 2))
	assert.Equal(t, 0.5, d.Progress(1, 5, 1.0, 3))
	assert.Equal(t, 0.75, d.Progress(1, 5, 1.0, 4))
	assert.Equal(t, 1.0, d.Progress(1, 5, 1.0, 5))

	// Delay
	assert.Equal(t, 1.0, d.Progress(1, 5, 1.0, 6))

	// Fill
	assert.Equal(t, 0.5, d.Progress(1, 5, 0.5, 7))
}

func TestDirectionReverse(t *testing.T) {
	d := DirectionReverse

	assert.Equal(t, 7, d.FrameCount(1, 5))

	// Delay
	assert.Equal(t, 1.0, d.Progress(1, 5, 0.0, 0))

	// Progress
	assert.Equal(t, 1.0, d.Progress(1, 5, 0.0, 1))
	assert.Equal(t, 0.75, d.Progress(1, 5, 0.0, 2))
	assert.Equal(t, 0.5, d.Progress(1, 5, 0.0, 3))
	assert.Equal(t, 0.25, d.Progress(1, 5, 0.0, 4))
	assert.Equal(t, 0.0, d.Progress(1, 5, 0.0, 5))

	// Delay
	assert.Equal(t, 0.0, d.Progress(1, 5, 0.0, 6))

	// Fill
	assert.Equal(t, 0.5, d.Progress(1, 5, 0.5, 7))
}

func TestDirectionAlternate(t *testing.T) {
	d := DirectionAlternate

	assert.Equal(t, 12, d.FrameCount(1, 5))

	// Delay
	assert.Equal(t, 0.0, d.Progress(1, 5, 1.0, 0))

	// Progress
	assert.Equal(t, 0.0, d.Progress(1, 5, 1.0, 1))
	assert.Equal(t, 0.25, d.Progress(1, 5, 1.0, 2))
	assert.Equal(t, 0.5, d.Progress(1, 5, 1.0, 3))
	assert.Equal(t, 0.75, d.Progress(1, 5, 1.0, 4))
	assert.Equal(t, 1.0, d.Progress(1, 5, 1.0, 5))

	// Delay
	assert.Equal(t, 1.0, d.Progress(1, 5, 1.0, 6))

	// Progress
	assert.Equal(t, 1.0, d.Progress(1, 5, 1.0, 7))
	assert.Equal(t, 0.75, d.Progress(1, 5, 1.0, 8))
	assert.Equal(t, 0.5, d.Progress(1, 5, 1.0, 9))
	assert.Equal(t, 0.25, d.Progress(1, 5, 1.0, 10))
	assert.Equal(t, 0.0, d.Progress(1, 5, 1.0, 11))

	// Fill
	assert.Equal(t, 0.5, d.Progress(1, 5, 0.5, 12))
}

func TestDirectionAlternateReverse(t *testing.T) {
	d := DirectionAlternateReverse

	assert.Equal(t, 12, d.FrameCount(1, 5))

	// Delay
	assert.Equal(t, 1.0, d.Progress(1, 5, 1.0, 0))

	// Progress
	assert.Equal(t, 1.0, d.Progress(1, 5, 0.0, 1))
	assert.Equal(t, 0.75, d.Progress(1, 5, 0.0, 2))
	assert.Equal(t, 0.5, d.Progress(1, 5, 0.0, 3))
	assert.Equal(t, 0.25, d.Progress(1, 5, 0.0, 4))
	assert.Equal(t, 0.0, d.Progress(1, 5, 0.0, 5))

	// Delay
	assert.Equal(t, 0.0, d.Progress(1, 5, 1.0, 6))

	// Progress
	assert.Equal(t, 0.0, d.Progress(1, 5, 0.0, 7))
	assert.Equal(t, 0.25, d.Progress(1, 5, 0.0, 8))
	assert.Equal(t, 0.5, d.Progress(1, 5, 0.0, 9))
	assert.Equal(t, 0.75, d.Progress(1, 5, 0.0, 10))
	assert.Equal(t, 1.0, d.Progress(1, 5, 0.0, 11))

	// Fill
	assert.Equal(t, 0.5, d.Progress(1, 5, 0.5, 12))
}

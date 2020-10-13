package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTracerCircularPath(t *testing.T) {
	ic := ImageChecker{
		palette: map[string]color.RGBA{
			"1": color.RGBA{0xff, 0xff, 0xff, 0xff},
			".": color.RGBA{0, 0, 0, 0},
		},
	}

	tr := Tracer{
		Path:        &CircularPath{Radius: 4},
		TraceLength: 0,
	}

	// First quadrant
	assert.Equal(t, nil, ic.Check([]string{
		"........",
		"........",
		"........",
		"........",
		".......1",
		"........",
		"........",
		"........",
	}, tr.Paint(image.Rect(0, 0, 100, 100), 0)))

	assert.Equal(t, nil, ic.Check([]string{
		"........",
		"........",
		"........",
		"........",
		"........",
		".......1",
		"........",
		"........",
	}, tr.Paint(image.Rect(0, 0, 100, 100), 1)))

	assert.Equal(t, nil, ic.Check([]string{
		"........",
		"........",
		"........",
		"........",
		"........",
		"........",
		".......1",
		"........",
	}, tr.Paint(image.Rect(0, 0, 100, 100), 2)))

	assert.Equal(t, nil, ic.Check([]string{
		"........",
		"........",
		"........",
		"........",
		"........",
		"........",
		"........",
		"......1.",
	}, tr.Paint(image.Rect(0, 0, 100, 100), 3)))

	// Spot check third quadrant
	assert.Equal(t, nil, ic.Check([]string{
		"........",
		"1.......",
		"........",
		"........",
		"........",
		"........",
		"........",
		"........",
	}, tr.Paint(image.Rect(0, 0, 100, 100), 14)))

	// Last pixel and verify it loops
	assert.Equal(t, nil, ic.Check([]string{
		"........",
		"........",
		"........",
		".......1",
		"........",
		"........",
		"........",
		"........",
	}, tr.Paint(image.Rect(0, 0, 100, 100), 23)))

	assert.Equal(t, nil, ic.Check([]string{
		"........",
		"........",
		"........",
		"........",
		"........",
		".......1",
		"........",
		"........",
	}, tr.Paint(image.Rect(0, 0, 100, 100), 25)))

	// All in all, we should have 24 frames
	assert.Equal(t, 24, tr.FrameCount())
}

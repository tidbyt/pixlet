package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Column is just a Vector. See vector_test.go for full coverage.

func TestColumnPaint(t *testing.T) {
	c := Column{
		Expanded:   true,
		MainAlign:  "space_evenly",
		CrossAlign: "center",
		Children: []Widget{
			// A green  box
			Box{Width: 6, Height: 7, Color: color.RGBA{0, 0xff, 0, 0xff}},
			// A red  box
			Box{Width: 8, Height: 9, Color: color.RGBA{0xff, 0, 0, 0xff}},
		},
	}

	// On large canvas, height gets truncated to max of children,
	// while width expands to full size
	im := PaintWidget(c, image.Rect(0, 0, 25, 16+3), 0)
	assert.Equal(t, nil, checkImage([]string{
		"........",
		".gggggg.",
		".gggggg.",
		".gggggg.",
		".gggggg.",
		".gggggg.",
		".gggggg.",
		".gggggg.",
		"........",
		"rrrrrrrr",
		"rrrrrrrr",
		"rrrrrrrr",
		"rrrrrrrr",
		"rrrrrrrr",
		"rrrrrrrr",
		"rrrrrrrr",
		"rrrrrrrr",
		"rrrrrrrr",
		"........",
	}, im))
}

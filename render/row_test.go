package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Row is just a Vector. See vector_test.go for full coverage.

func TestRowPaint(t *testing.T) {
	r := Row{
		Expanded:   true,
		MainAlign:  "space_evenly",
		CrossAlign: "end",
		Children: []Widget{
			// A green  box
			Box{Width: 6, Height: 7, Color: color.RGBA{0, 0xff, 0, 0xff}},
			// A red  box
			Box{Width: 8, Height: 9, Color: color.RGBA{0xff, 0, 0, 0xff}},
		},
	}

	// On large canvas, height gets truncated to max of children,
	// while width expands to full size
	im := r.Paint(image.Rect(0, 0, 14+2, 17), 0)
	assert.Equal(t, nil, checkImage([]string{
		"........rrrrrrrr",
		"........rrrrrrrr",
		".gggggg.rrrrrrrr",
		".gggggg.rrrrrrrr",
		".gggggg.rrrrrrrr",
		".gggggg.rrrrrrrr",
		".gggggg.rrrrrrrr",
		".gggggg.rrrrrrrr",
		".gggggg.rrrrrrrr",
	}, im))
}

package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarqueeNoScroll(t *testing.T) {
	m := Marquee{
		Width: 6,
		Child: Row{
			Children: []Widget{
				Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
				Box{Width: 2, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
				Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
	}

	// Child fits so there's just 1 single frame
	assert.Equal(t, 1, m.FrameCount())
	im := m.Paint(image.Rect(0, 0, 100, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrggb",
		"rrrgg.",
		"rrr...",
	}, im))
}

func TestMarqueeScrolling(t *testing.T) {
	m := Marquee{
		Width: 6,
		Child: Row{
			Children: []Widget{
				Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
				Box{Width: 3, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
				Box{Width: 3, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
	}

	// The child's 9 pixels will be scrolled into view (6 frames),
	// scrolled out of view (9 frames) and then finally scrolled
	// back into view again (6 frames). 21 frames in total.
	assert.Equal(t, 21, m.FrameCount())

	// Scrolling into view
	assert.Equal(t, nil, checkImage([]string{
		"......",
		"......",
		"......",
	}, m.Paint(image.Rect(0, 0, 100, 100), 0)))

	assert.Equal(t, nil, checkImage([]string{
		"....rr",
		"....rr",
		"....rr",
	}, m.Paint(image.Rect(0, 0, 100, 100), 2)))

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, m.Paint(image.Rect(0, 0, 100, 100), 6)))

	// Scrolling out of view
	assert.Equal(t, nil, checkImage([]string{
		"rgggbb",
		"rggg..",
		"r.....",
	}, m.Paint(image.Rect(0, 0, 100, 100), 8)))

	assert.Equal(t, nil, checkImage([]string{
		"b.....",
		"......",
		"......",
	}, m.Paint(image.Rect(0, 0, 100, 100), 14)))

	assert.Equal(t, nil, checkImage([]string{
		"......",
		"......",
		"......",
	}, m.Paint(image.Rect(0, 0, 100, 100), 15)))

	// Scrolling back into view
	assert.Equal(t, nil, checkImage([]string{
		"...rrr",
		"...rrr",
		"...rrr",
	}, m.Paint(image.Rect(0, 0, 100, 100), 18)))

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, m.Paint(image.Rect(0, 0, 100, 100), 21)))

	// Later frames keep it fixed in the last frame. This makes
	// multiple simultaneous marquees look nice when they've
	// different length.

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, m.Paint(image.Rect(0, 0, 100, 100), 22)))

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, m.Paint(image.Rect(0, 0, 100, 100), 26)))

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, m.Paint(image.Rect(0, 0, 100, 100), 100000)))

}

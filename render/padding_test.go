package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

// By default, padding is added to child, regardless of bounds
func TestPadding(t *testing.T) {
	pad := Padding{
		Child: Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
		Pad: Insets{
			Left:   1,
			Top:    2,
			Right:  3,
			Bottom: 4,
		},
	}

	// Large bounds
	im := pad.Paint(image.Rect(0, 0, 20, 20), 0)
	assert.Equal(t, nil, checkImage([]string{
		".......",
		".......",
		".rrr...",
		".rrr...",
		".rrr...",
		".......",
		".......",
		".......",
		".......",
	}, im))

	// Small bounds
	im = pad.Paint(image.Rect(0, 0, 4, 4), 0)
	assert.Equal(t, nil, checkImage([]string{
		".......",
		".......",
		".rrr...",
		".rrr...",
		".rrr...",
		".......",
		".......",
		".......",
		".......",
	}, im))
}

// If expanded, the full bounds are used and child may be cropped
func TestPaddingExpanded(t *testing.T) {

	// Child fits, so it's placed in upper left corner
	pad := Padding{
		Child: Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
		Pad: Insets{
			Left:   1,
			Top:    1,
			Right:  1,
			Bottom: 1,
		},
		Expanded: true,
	}

	im := pad.Paint(image.Rect(0, 0, 7, 7), 0)
	assert.Equal(t, nil, checkImage([]string{
		".......",
		".rrr...",
		".rrr...",
		".rrr...",
		".......",
		".......",
		".......",
	}, im))

	// Child doesn't fit: crop
	im = pad.Paint(image.Rect(0, 0, 3, 3), 0)
	assert.Equal(t, nil, checkImage([]string{
		"...",
		".r.",
		"...",
	}, im))
}

// Same as TestPadding, with color
func TestColorPadding(t *testing.T) {
	pad := Padding{
		Child: Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
		Pad: Insets{
			Left:   1,
			Top:    2,
			Right:  3,
			Bottom: 4,
		},
		Color: color.RGBA{0, 0xff, 0, 0xff},
	}

	// Large bounds
	im := pad.Paint(image.Rect(0, 0, 20, 20), 0)
	assert.Equal(t, nil, checkImage([]string{
		"ggggggg",
		"ggggggg",
		"grrrggg",
		"grrrggg",
		"grrrggg",
		"ggggggg",
		"ggggggg",
		"ggggggg",
		"ggggggg",
	}, im))

	// Small bounds
	im = pad.Paint(image.Rect(0, 0, 4, 4), 0)
	assert.Equal(t, nil, checkImage([]string{
		"ggggggg",
		"ggggggg",
		"grrrggg",
		"grrrggg",
		"grrrggg",
		"ggggggg",
		"ggggggg",
		"ggggggg",
		"ggggggg",
	}, im))
}

// Same as TestPaddingExpanded, with color
func TestColorPaddingExpanded(t *testing.T) {

	// Child fits, so it's placed in upper left corner
	pad := Padding{
		Child: Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
		Pad: Insets{
			Left:   1,
			Top:    1,
			Right:  1,
			Bottom: 1,
		},
		Expanded: true,
		Color: color.RGBA{0, 0xff, 0, 0xff},
	}

	im := pad.Paint(image.Rect(0, 0, 7, 7), 0)
	assert.Equal(t, nil, checkImage([]string{
		"ggggggg",
		"grrrggg",
		"grrrggg",
		"grrrggg",
		"ggggggg",
		"ggggggg",
		"ggggggg",
	}, im))

	// Child doesn't fit: crop
	im = pad.Paint(image.Rect(0, 0, 3, 3), 0)
	assert.Equal(t, nil, checkImage([]string{
		"ggg",
		"grg",
		"ggg",
	}, im))
}

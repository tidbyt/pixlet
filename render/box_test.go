package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

// A box without a child is just a box
func TestBoxNoChild(t *testing.T) {

	// Transparent by default
	box := Box{}
	im := box.Paint(image.Rect(0, 0, 5, 5), 0)
	assert.Equal(t, nil, checkImage([]string{
		".....",
		".....",
		".....",
		".....",
		".....",
	}, im))

	// Color can be set
	box = Box{Color: color.RGBA{0xff, 0, 0, 0xff}}
	im = box.Paint(image.Rect(0, 0, 5, 5), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrr",
		"rrrrr",
		"rrrrr",
		"rrrrr",
		"rrrrr",
	}, im))

	// Specify Width and the box fills height bounds
	box = Box{Color: color.RGBA{0xff, 0, 0, 0xff}, Width: 3}
	im = box.Paint(image.Rect(0, 0, 5, 5), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrr",
		"rrr",
		"rrr",
		"rrr",
		"rrr",
	}, im))

	// Specify Height and it fills width bounds
	box = Box{Color: color.RGBA{0xff, 0, 0, 0xff}, Height: 3}
	im = box.Paint(image.Rect(0, 0, 5, 5), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrr",
		"rrrrr",
		"rrrrr",
	}, im))

	// Specify both and it ignores the bounds entirely
	box = Box{Color: color.RGBA{0xff, 0, 0, 0xff}, Width: 2, Height: 3}
	im = box.Paint(image.Rect(0, 0, 5, 5), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rr",
		"rr",
		"rr",
	}, im))
}

// Box with child centers the child
func TestBoxChildCenter(t *testing.T) {

	box := Box{
		Child: Box{
			Color:  color.RGBA{0xff, 0, 0, 0xff},
			Width:  2,
			Height: 2,
		},
	}
	im := box.Paint(image.Rect(0, 0, 4, 4), 0)
	assert.Equal(t, nil, checkImage([]string{
		"....",
		".rr.",
		".rr.",
		"....",
	}, im))

	// If perfect centering can be done, remaining pixels are on
	// the right and below the child.
	im = box.Paint(image.Rect(0, 0, 5, 5), 0)
	assert.Equal(t, nil, checkImage([]string{
		".....",
		".rr..",
		".rr..",
		".....",
		".....",
	}, im))

	// Centered horizontally here
	im = box.Paint(image.Rect(0, 0, 4, 2), 0)
	assert.Equal(t, nil, checkImage([]string{
		".rr.",
		".rr.",
	}, im))

	// Centered vertically here
	im = box.Paint(image.Rect(0, 0, 2, 4), 0)
	assert.Equal(t, nil, checkImage([]string{
		"..",
		"rr",
		"rr",
		"..",
	}, im))
}

// Box can place padding around child
func TestBoxPadding(t *testing.T) {
	// No padding and the child box will fill the bounds
	box := Box{
		Child: Box{
			Color: color.RGBA{0xff, 0, 0, 0xff},
		},
	}
	im := box.Paint(image.Rect(0, 0, 4, 4), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrr",
		"rrrr",
		"rrrr",
		"rrrr",
	}, im))

	// Add padding and bounds for child box shrinks
	box = Box{
		Child: Box{
			Color: color.RGBA{0xff, 0, 0, 0xff},
		},
		Padding: 1,
	}
	im = box.Paint(image.Rect(0, 0, 4, 4), 0)
	assert.Equal(t, nil, checkImage([]string{
		"....",
		".rr.",
		".rr.",
		"....",
	}, im))

	// More padding!
	box = Box{
		Child: Box{
			Color: color.RGBA{0xff, 0, 0, 0xff},
		},
		Padding: 3,
	}
	im = box.Paint(image.Rect(0, 0, 8, 8), 0)
	assert.Equal(t, nil, checkImage([]string{
		"........",
		"........",
		"........",
		"...rr...",
		"...rr...",
		"........",
		"........",
		"........",
	}, im))
}

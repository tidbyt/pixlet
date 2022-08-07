package render

import (
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"testing"
)

func TestVectorMainAlignStart(t *testing.T) {
	// A horizontally laid out vector
	v := Vector{
		Expanded:  true,
		Vertical:  false,
		MainAlign: "start",
		Children: []Widget{
			// A red  box
			Box{Width: 10, Height: 4, Color: color.RGBA{0xff, 0, 0, 0xff}},
			// A green  box
			Box{Width: 6, Height: 8, Color: color.RGBA{0, 0xff, 0, 0xff}},
			// A blue  box
			Box{Width: 2, Height: 12, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}

	// On large canvas, height gets truncated to max of children,
	// while width expands to full size
	im := PaintWidget(v, image.Rect(0, 0, 25, 20), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrrrrrrrggggggbb.......",
		"rrrrrrrrrrggggggbb.......",
		"rrrrrrrrrrggggggbb.......",
		"rrrrrrrrrrggggggbb.......",
		"..........ggggggbb.......",
		"..........ggggggbb.......",
		"..........ggggggbb.......",
		"..........ggggggbb.......",
		"................bb.......",
		"................bb.......",
		"................bb.......",
		"................bb.......",
	}, im))

	// Reduce height. Overflowing children are partially drawn.
	im = PaintWidget(v, image.Rect(0, 0, 25, 10), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrrrrrrrggggggbb.......",
		"rrrrrrrrrrggggggbb.......",
		"rrrrrrrrrrggggggbb.......",
		"rrrrrrrrrrggggggbb.......",
		"..........ggggggbb.......",
		"..........ggggggbb.......",
		"..........ggggggbb.......",
		"..........ggggggbb.......",
		"................bb.......",
		"................bb.......",
	}, im))

	// Reduce width further.
	im = PaintWidget(v, image.Rect(0, 0, 17, 10), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrrrrrrrggggggb",
		"rrrrrrrrrrggggggb",
		"rrrrrrrrrrggggggb",
		"rrrrrrrrrrggggggb",
		"..........ggggggb",
		"..........ggggggb",
		"..........ggggggb",
		"..........ggggggb",
		"................b",
		"................b",
	}, im))

	// Reduce so a child is completely cut out, and it's no longer
	// included in height.
	im = PaintWidget(v, image.Rect(0, 0, 16, 10), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrrrrrrrgggggg",
		"rrrrrrrrrrgggggg",
		"rrrrrrrrrrgggggg",
		"rrrrrrrrrrgggggg",
		"..........gggggg",
		"..........gggggg",
		"..........gggggg",
		"..........gggggg",
	}, im))

	// Perfect fit is a perfect fit.
	im = PaintWidget(v, image.Rect(0, 0, 18, 12), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrrrrrrrggggggbb",
		"rrrrrrrrrrggggggbb",
		"rrrrrrrrrrggggggbb",
		"rrrrrrrrrrggggggbb",
		"..........ggggggbb",
		"..........ggggggbb",
		"..........ggggggbb",
		"..........ggggggbb",
		"................bb",
		"................bb",
		"................bb",
		"................bb",
	}, im))

	// Flip it
	v.Vertical = true
	im = PaintWidget(v, image.Rect(0, 0, 10, 25), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrrrrrrr",
		"rrrrrrrrrr",
		"rrrrrrrrrr",
		"rrrrrrrrrr",
		"gggggg....",
		"gggggg....",
		"gggggg....",
		"gggggg....",
		"gggggg....",
		"gggggg....",
		"gggggg....",
		"gggggg....",
		"bb........",
		"bb........",
		"bb........",
		"bb........",
		"bb........",
		"bb........",
		"bb........",
		"bb........",
		"bb........",
		"bb........",
		"bb........",
		"bb........",
		"..........",
	}, im))
}

func TestVectorMainAlignStartVertical(t *testing.T) {
	// A vertically laid out vector
	v := Vector{
		Expanded:  true,
		Vertical:  true,
		MainAlign: "start",
		Children: []Widget{
			// A red  box
			Box{Width: 4, Height: 6, Color: color.RGBA{0xff, 0, 0, 0xff}},
			// A green  box
			Box{Width: 8, Height: 3, Color: color.RGBA{0, 0xff, 0, 0xff}},
			// A blue  box
			Box{Width: 12, Height: 2, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}

	// Width shrinks to fit
	im := PaintWidget(v, image.Rect(0, 0, 100, 11), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrr........",
		"rrrr........",
		"rrrr........",
		"rrrr........",
		"rrrr........",
		"rrrr........",
		"gggggggg....",
		"gggggggg....",
		"gggggggg....",
		"bbbbbbbbbbbb",
		"bbbbbbbbbbbb",
	}, im))

	// Height does not shrink
	im = PaintWidget(v, image.Rect(0, 0, 100, 13), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrr........",
		"rrrr........",
		"rrrr........",
		"rrrr........",
		"rrrr........",
		"rrrr........",
		"gggggggg....",
		"gggggggg....",
		"gggggggg....",
		"bbbbbbbbbbbb",
		"bbbbbbbbbbbb",
		"............",
		"............",
	}, im))

	// Children are partially drawn
	im = PaintWidget(v, image.Rect(0, 0, 10, 13), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrr......",
		"rrrr......",
		"rrrr......",
		"rrrr......",
		"rrrr......",
		"rrrr......",
		"gggggggg..",
		"gggggggg..",
		"gggggggg..",
		"bbbbbbbbbb",
		"bbbbbbbbbb",
		"..........",
		"..........",
	}, im))

	// Along both axes
	im = PaintWidget(v, image.Rect(0, 0, 10, 10), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrr......",
		"rrrr......",
		"rrrr......",
		"rrrr......",
		"rrrr......",
		"rrrr......",
		"gggggggg..",
		"gggggggg..",
		"gggggggg..",
		"bbbbbbbbbb",
	}, im))

	// And if a child is completely out of sight, it no longer
	// affects width
	im = PaintWidget(v, image.Rect(0, 0, 10, 9), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"gggggggg",
		"gggggggg",
		"gggggggg",
	}, im))
}

func TestVectorMainAlignEnd(t *testing.T) {

	v := Vector{
		Expanded:  true,
		Vertical:  true,
		MainAlign: "end",
		Children: []Widget{
			// A red  box
			Box{Width: 4, Height: 6, Color: color.RGBA{0xff, 0, 0, 0xff}},
			// A green  box
			Box{Width: 8, Height: 3, Color: color.RGBA{0, 0xff, 0, 0xff}},
			// A blue  box
			Box{Width: 5, Height: 2, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}

	// With MainAlignEnd, children are placed at the end along the
	// main axis
	im := PaintWidget(v, image.Rect(0, 0, 100, 13), 0)
	assert.Equal(t, nil, checkImage([]string{
		"........",
		"........",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"gggggggg",
		"gggggggg",
		"gggggggg",
		"bbbbb...",
		"bbbbb...",
	}, im))

	// Flip it!
	v.Vertical = false
	im = PaintWidget(v, image.Rect(0, 0, 19, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"..rrrrggggggggbbbbb",
		"..rrrrggggggggbbbbb",
		"..rrrrgggggggg.....",
		"..rrrr.............",
		"..rrrr.............",
		"..rrrr.............",
	}, im))

	// Reducing width/height
	im = PaintWidget(v, image.Rect(0, 0, 16, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrggggggggbbbb",
		"rrrrggggggggbbbb",
		"rrrrgggggggg....",
		"rrrr............",
		"rrrr............",
		"rrrr............",
	}, im))
	im = PaintWidget(v, image.Rect(0, 0, 16, 5), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrggggggggbbbb",
		"rrrrggggggggbbbb",
		"rrrrgggggggg....",
		"rrrr............",
		"rrrr............",
	}, im))
	im = PaintWidget(v, image.Rect(0, 0, 16, 5), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrggggggggbbbb",
		"rrrrggggggggbbbb",
		"rrrrgggggggg....",
		"rrrr............",
		"rrrr............",
	}, im))
}

func TestVectorMainAlignSpaceEvenly(t *testing.T) {

	v := Vector{
		Expanded:  true,
		Vertical:  false,
		MainAlign: "space_evenly",
		Children: []Widget{
			// A red  box
			Box{Width: 4, Height: 6, Color: color.RGBA{0xff, 0, 0, 0xff}},
			// A green  box
			Box{Width: 8, Height: 3, Color: color.RGBA{0, 0xff, 0, 0xff}},
			// A blue  box
			Box{Width: 5, Height: 2, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}

	// Canvas leaves 2 pixels spacing before and after each child
	im := PaintWidget(v, image.Rect(0, 0, 17+2*4, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"..rrrr..gggggggg..bbbbb..",
		"..rrrr..gggggggg..bbbbb..",
		"..rrrr..gggggggg.........",
		"..rrrr...................",
		"..rrrr...................",
		"..rrrr...................",
	}, im))

	// Adding 2 pixels width means canvas doesn't divide
	// evenly. The residual should be distributed start to end,
	// one pixel at a time.
	im = PaintWidget(v, image.Rect(0, 0, 17+2*4+2, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"...rrrr...gggggggg..bbbbb..",
		"...rrrr...gggggggg..bbbbb..",
		"...rrrr...gggggggg.........",
		"...rrrr....................",
		"...rrrr....................",
		"...rrrr....................",
	}, im))

	// If not expanded, this just shrink wraps to smallest size
	v.Expanded = false
	im = PaintWidget(v, image.Rect(0, 0, 17+2*4+2, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrggggggggbbbbb",
		"rrrrggggggggbbbbb",
		"rrrrgggggggg.....",
		"rrrr.............",
		"rrrr.............",
		"rrrr.............",
	}, im))

	// Flip it
	v.Vertical = true
	im = PaintWidget(v, image.Rect(0, 0, 100, 11+2*4+2), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"gggggggg",
		"gggggggg",
		"gggggggg",
		"bbbbb...",
		"bbbbb...",
	}, im))

	// And unexpand it
	v.Expanded = true
	im = PaintWidget(v, image.Rect(0, 0, 100, 11+2*4+2), 0)
	assert.Equal(t, nil, checkImage([]string{
		"........",
		"........",
		"........",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"........",
		"........",
		"........",
		"gggggggg",
		"gggggggg",
		"gggggggg",
		"........",
		"........",
		"bbbbb...",
		"bbbbb...",
		"........",
		"........",
	}, im))
}

func TestVectorCrossAlignCenter(t *testing.T) {

	v := Vector{
		Expanded:   true,
		Vertical:   false,
		MainAlign:  "space_evenly",
		CrossAlign: "center",
		Children: []Widget{
			// A red  box
			Box{Width: 4, Height: 6, Color: color.RGBA{0xff, 0, 0, 0xff}},
			// A green  box
			Box{Width: 8, Height: 3, Color: color.RGBA{0, 0xff, 0, 0xff}},
			// A blue  box
			Box{Width: 5, Height: 2, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}

	// Everything centered along the cross axis. Space isn't
	// evenly divisible for the green box.
	im := PaintWidget(v, image.Rect(0, 0, 17+2*4, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"..rrrr...................",
		"..rrrr..gggggggg.........",
		"..rrrr..gggggggg..bbbbb..",
		"..rrrr..gggggggg..bbbbb..",
		"..rrrr...................",
		"..rrrr...................",
	}, im))

	// Flip it!
	v.Vertical = true
	im = PaintWidget(v, image.Rect(0, 0, 10000, 11+2*4), 0)
	assert.Equal(t, nil, checkImage([]string{
		"........",
		"........",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"........",
		"........",
		"gggggggg",
		"gggggggg",
		"gggggggg",
		"........",
		"........",
		".bbbbb..",
		".bbbbb..",
		"........",
		"........",
	}, im))

	// Unexpand it
	v.Expanded = false
	im = PaintWidget(v, image.Rect(0, 0, 10000, 11+2*4), 0)
	assert.Equal(t, nil, checkImage([]string{
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"gggggggg",
		"gggggggg",
		"gggggggg",
		".bbbbb..",
		".bbbbb..",
	}, im))

	// Works with other main axis alignment
	v.Expanded = true
	v.MainAlign = "start"
	im = PaintWidget(v, image.Rect(0, 0, 10000, 11+2), 0)
	assert.Equal(t, nil, checkImage([]string{
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"gggggggg",
		"gggggggg",
		"gggggggg",
		".bbbbb..",
		".bbbbb..",
		"........",
		"........",
	}, im))
	v.MainAlign = "end"
	im = PaintWidget(v, image.Rect(0, 0, 10000, 11+2), 0)
	assert.Equal(t, nil, checkImage([]string{
		"........",
		"........",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"..rrrr..",
		"gggggggg",
		"gggggggg",
		"gggggggg",
		".bbbbb..",
		".bbbbb..",
	}, im))

}

func TestVectorCrossAlignEnd(t *testing.T) {

	v := Vector{
		Expanded:   true,
		Vertical:   false,
		MainAlign:  "space_evenly",
		CrossAlign: "end",
		Children: []Widget{
			// A red  box
			Box{Width: 4, Height: 6, Color: color.RGBA{0xff, 0, 0, 0xff}},
			// A green  box
			Box{Width: 8, Height: 3, Color: color.RGBA{0, 0xff, 0, 0xff}},
			// A blue  box
			Box{Width: 5, Height: 2, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}

	im := PaintWidget(v, image.Rect(0, 0, 17+2*4+1, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"...rrrr...................",
		"...rrrr...................",
		"...rrrr...................",
		"...rrrr..gggggggg.........",
		"...rrrr..gggggggg..bbbbb..",
		"...rrrr..gggggggg..bbbbb..",
	}, im))

	v.Vertical = true
	im = PaintWidget(v, image.Rect(0, 0, 10000, 11+2*4+1), 0)
	assert.Equal(t, nil, checkImage([]string{
		"........",
		"........",
		"........",
		"....rrrr",
		"....rrrr",
		"....rrrr",
		"....rrrr",
		"....rrrr",
		"....rrrr",
		"........",
		"........",
		"gggggggg",
		"gggggggg",
		"gggggggg",
		"........",
		"........",
		"...bbbbb",
		"...bbbbb",
		"........",
		"........",
	}, im))

}

func TestVectorCrossAlignStart(t *testing.T) {

	v := Vector{
		Expanded:   true,
		Vertical:   false,
		MainAlign:  "end",
		CrossAlign: "start",
		Children: []Widget{
			// A red  box
			Box{Width: 4, Height: 6, Color: color.RGBA{0xff, 0, 0, 0xff}},
			// A green  box
			Box{Width: 8, Height: 3, Color: color.RGBA{0, 0xff, 0, 0xff}},
			// A blue  box
			Box{Width: 5, Height: 2, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}

	im := PaintWidget(v, image.Rect(0, 0, 17+2*4+1, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		".........rrrrggggggggbbbbb",
		".........rrrrggggggggbbbbb",
		".........rrrrgggggggg.....",
		".........rrrr.............",
		".........rrrr.............",
		".........rrrr.............",
	}, im))

	v.Vertical = true
	im = PaintWidget(v, image.Rect(0, 0, 10000, 11+2*4+1), 0)
	assert.Equal(t, nil, checkImage([]string{
		"........",
		"........",
		"........",
		"........",
		"........",
		"........",
		"........",
		"........",
		"........",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"rrrr....",
		"gggggggg",
		"gggggggg",
		"gggggggg",
		"bbbbb...",
		"bbbbb...",
	}, im))

	v.Expanded = false
	im = PaintWidget(v, image.Rect(0, 0, 7, 11+2*4+1), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrr...",
		"rrrr...",
		"rrrr...",
		"rrrr...",
		"rrrr...",
		"rrrr...",
		"ggggggg",
		"ggggggg",
		"ggggggg",
		"bbbbb..",
		"bbbbb..",
	}, im))
}

func TestVectorMainAlignSpaceAround(t *testing.T) {
	v := Vector{
		Expanded:  true,
		Vertical:  false,
		MainAlign: "space_around",
		Children: []Widget{
			// A red  box
			Box{Width: 3, Height: 4, Color: color.RGBA{0xff, 0, 0, 0xff}},
			// A green  box
			Box{Width: 5, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
			// A blue  box
			Box{Width: 1, Height: 3, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}

	// MainAlignSpaceAround will put create equal space _between_
	// adjacent children, and half of that before and after the
	// first and last child. When the remaining space is not
	// evenly divisible, the residual gets distributed left to
	// right. When looking at individual pixels, this can look a
	// bit unintuitive, but I couldn't come up with anything
	// better. These tests illustrate (and verify) the behavior.

	im := PaintWidget(v, image.Rect(0, 0, 9, 20), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrgggggb",
		"rrrgggggb",
		"rrr.....b",
		"rrr......",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 10, 20), 0)
	assert.Equal(t, nil, checkImage([]string{
		".rrrgggggb",
		".rrrgggggb",
		".rrr.....b",
		".rrr......",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 11, 20), 0)
	assert.Equal(t, nil, checkImage([]string{
		".rrr.gggggb",
		".rrr.gggggb",
		".rrr......b",
		".rrr.......",
	}, im))

	// This is a bit unintuitive. With 3 empty pixels, we have 1
	// pixel between each child, 1/2 before the first and 1/2
	// after the last. We can't draw 1/2 pixel, so there's 0
	// padding before the red block, and 0 after the
	// blue. However, since that comes in 1 pixel short, the final
	// column of pixels is not painted, which comes out as
	// padding. So it makes sense, but it also doesn't.
	im = PaintWidget(v, image.Rect(0, 0, 12, 20), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrr.ggggg.b.",
		"rrr.ggggg.b.",
		"rrr.......b.",
		"rrr.........",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 13, 20), 0)
	assert.Equal(t, nil, checkImage([]string{
		".rrr.ggggg.b.",
		".rrr.ggggg.b.",
		".rrr.......b.",
		".rrr.........",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 14, 20), 0)
	assert.Equal(t, nil, checkImage([]string{
		".rrr..ggggg.b.",
		".rrr..ggggg.b.",
		".rrr........b.",
		".rrr..........",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 15, 20), 0)
	assert.Equal(t, nil, checkImage([]string{
		".rrr..ggggg..b.",
		".rrr..ggggg..b.",
		".rrr.........b.",
		".rrr...........",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 16, 20), 0)
	assert.Equal(t, nil, checkImage([]string{
		"..rrr..ggggg..b.",
		"..rrr..ggggg..b.",
		"..rrr.........b.",
		"..rrr...........",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 20, 20), 0)
	assert.Equal(t, nil, checkImage([]string{
		"..rrr....ggggg...b..",
		"..rrr....ggggg...b..",
		"..rrr............b..",
		"..rrr...............",
	}, im))

	// And flip it up a bit
	v.Vertical = true
	v.CrossAlign = "center"
	im = PaintWidget(v, image.Rect(0, 0, 20, 20), 0)
	assert.Equal(t, nil, checkImage([]string{
		".....",
		".....",
		".rrr.",
		".rrr.",
		".rrr.",
		".rrr.",
		".....",
		".....",
		".....",
		".....",
		"ggggg",
		"ggggg",
		".....",
		".....",
		".....",
		"..b..",
		"..b..",
		"..b..",
		".....",
		".....",
	}, im))

}

func TestVectorMainAlignCenter(t *testing.T) {
	v := Vector{
		Expanded:  true,
		Vertical:  true,
		MainAlign: "center",
		Children: []Widget{
			Box{Width: 3, Height: 4, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 5, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 1, Height: 3, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}

	// MainAlignCenter places children adjacent without spacing,
	// centered along the main axis.

	im := PaintWidget(v, image.Rect(0, 0, 100, 12), 0)
	assert.Equal(t, nil, checkImage([]string{
		".....",
		"rrr..",
		"rrr..",
		"rrr..",
		"rrr..",
		"ggggg",
		"ggggg",
		"b....",
		"b....",
		"b....",
		".....",
		".....",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 100, 11), 0)
	assert.Equal(t, nil, checkImage([]string{
		".....",
		"rrr..",
		"rrr..",
		"rrr..",
		"rrr..",
		"ggggg",
		"ggggg",
		"b....",
		"b....",
		"b....",
		".....",
	}, im))

	v.Vertical = false
	v.CrossAlign = "center"
	im = PaintWidget(v, image.Rect(0, 0, 13, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"..rrr.....b..",
		"..rrrgggggb..",
		"..rrrgggggb..",
		"..rrr........",
	}, im))
}

func TestVectorMainAlignSpaceBetween(t *testing.T) {
	v := Vector{
		Expanded:  true,
		Vertical:  false,
		MainAlign: "space_between",
		Children: []Widget{
			Box{Width: 3, Height: 4, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 5, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 1, Height: 3, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}

	// MainAlignSpaceBetween distributes space evenly between
	// children, but not before the first child or after the last
	// child.

	im := PaintWidget(v, image.Rect(0, 0, 9, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrgggggb",
		"rrrgggggb",
		"rrr.....b",
		"rrr......",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 10, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrr.gggggb",
		"rrr.gggggb",
		"rrr......b",
		"rrr.......",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 11, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrr.ggggg.b",
		"rrr.ggggg.b",
		"rrr.......b",
		"rrr........",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 12, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrr..ggggg.b",
		"rrr..ggggg.b",
		"rrr........b",
		"rrr.........",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 13, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrr..ggggg..b",
		"rrr..ggggg..b",
		"rrr.........b",
		"rrr..........",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 14, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrr...ggggg..b",
		"rrr...ggggg..b",
		"rrr..........b",
		"rrr...........",
	}, im))

	im = PaintWidget(v, image.Rect(0, 0, 15, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrr...ggggg...b",
		"rrr...ggggg...b",
		"rrr...........b",
		"rrr............",
	}, im))

	// Some extra attention to the edge case of 1 and 2 children

	v = Vector{
		Expanded:   true,
		Vertical:   false,
		MainAlign:  "space_between",
		CrossAlign: "center",
		Children: []Widget{
			Box{Width: 3, Height: 4, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 5, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
		},
	}

	im = PaintWidget(v, image.Rect(0, 0, 14, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrr...........",
		"rrr......ggggg",
		"rrr......ggggg",
		"rrr...........",
	}, im))

	v = Vector{
		Expanded:   true,
		Vertical:   false,
		MainAlign:  "space_between",
		CrossAlign: "center",
		Children: []Widget{
			Box{Width: 3, Height: 4, Color: color.RGBA{0xff, 0, 0, 0xff}},
		},
	}

	im = PaintWidget(v, image.Rect(0, 0, 14, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrr...........",
		"rrr...........",
		"rrr...........",
		"rrr...........",
	}, im))

	// And for no particular reason: many children
	v = Vector{
		Expanded:   true,
		Vertical:   false,
		MainAlign:  "space_between",
		CrossAlign: "end",
		Children: []Widget{
			Box{Width: 3, Height: 4, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 5, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 1, Height: 4, Color: color.RGBA{0, 0, 0xff, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0xff, 0xff, 0xff, 0xff}},
			Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 4, Height: 4, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}

	im = PaintWidget(v, image.Rect(0, 0, 50, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrr...................b.......................bbbb",
		"rrr...................b..............rrr......bbbb",
		"rrr.......ggggg.......b..............rrr......bbbb",
		"rrr.......ggggg.......b......ww......rrr......bbbb",
	}, im))

}

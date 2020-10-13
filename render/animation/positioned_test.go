package animation

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidbyt/pixlet/render"
)

func TestPositionedLinearCurve(t *testing.T) {
	o := AnimatedPositioned{
		Child: render.Box{
			Width:  4,
			Height: 3,
			Color:  color.RGBA{0xff, 0x00, 0x00, 0xff},
		},
		Duration: 6,
		XStart:   0,
		YStart:   0,
		XEnd:     5,
		YEnd:     5,
		Curve:    LinearCurve{},
	}

	// These frames should show the box moving from (0, 0) to (5,
	// 5), one pixel per frame (since Duration is 6, which equals
	// the number of positions).

	assert.Equal(t, 6, o.FrameCount())

	im := o.Paint(image.Rect(0, 0, 10, 6), 0)
	assert.Equal(t, nil, render.CheckImage([]string{
		"rrrr......",
		"rrrr......",
		"rrrr......",
		"..........",
		"..........",
		"..........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 6), 1)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		".rrrr.....",
		".rrrr.....",
		".rrrr.....",
		"..........",
		"..........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 6), 2)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"..........",
		"..rrrr....",
		"..rrrr....",
		"..rrrr....",
		"..........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 6), 3)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"..........",
		"..........",
		"...rrrr...",
		"...rrrr...",
		"...rrrr...",
	}, im))

	// These last two frames, the rectangle moves out of bounds.
	im = o.Paint(image.Rect(0, 0, 10, 6), 4)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"..........",
		"..........",
		"..........",
		"....rrrr..",
		"....rrrr..",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 6), 5)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"..........",
		"..........",
		"..........",
		"..........",
		".....rrrr.",
	}, im))

	// For now, the behavior for later frames is to freeze the
	// final frame, i.e. keep child at the end position.
	im = o.Paint(image.Rect(0, 0, 10, 6), 6)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"..........",
		"..........",
		"..........",
		"..........",
		".....rrrr.",
	}, im))
	im = o.Paint(image.Rect(0, 0, 10, 6), 7)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"..........",
		"..........",
		"..........",
		"..........",
		".....rrrr.",
	}, im))

}

func TestPositionedEaseIn(t *testing.T) {
	o := AnimatedPositioned{
		Child: render.Box{
			Width:  2,
			Height: 2,
			Color:  color.RGBA{0x00, 0xff, 0x00, 0xff},
		},
		Duration: 10,
		XStart:   -3,
		YStart:   1,
		XEnd:     4,
		YEnd:     1,
		Curve:    EaseOut,
	}

	// A green box appearing from off screen, moving rapidly at
	// first, then decelerating into its end position.

	im := o.Paint(image.Rect(0, 0, 10, 4), 0)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"..........",
		"..........",
		"..........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 4), 1)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		".gg.......",
		".gg.......",
		"..........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 4), 2)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"..gg......",
		"..gg......",
		"..........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 4), 3)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"...gg.....",
		"...gg.....",
		"..........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 4), 4)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"...gg.....",
		"...gg.....",
		"..........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 4), 5)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"....gg....",
		"....gg....",
		"..........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 4), 9)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"....gg....",
		"....gg....",
		"..........",
	}, im))
}

// The Delay and Hold parameters delays the start of transition and
//  holds it at its final position for a certain number of frames.
func TestPositionedDelayAndHold(t *testing.T) {
	o := AnimatedPositioned{
		Child: render.Box{
			Width:  1,
			Height: 2,
			Color:  color.RGBA{0x00, 0xff, 0x00, 0xff},
		},
		Duration: 5,
		XStart:   0,
		XEnd:     4,
		Delay:    3,
		Hold:     2,
		Curve:    LinearCurve{},
	}

	// Duration is 5 frames. On top of that, there's a 3 frame
	// delay before it starts, and it's held in its final position
	// for 2 frames, so we expect 13 frames in total.
	assert.Equal(t, 10, o.FrameCount())

	// No movement during delay
	im := o.Paint(image.Rect(0, 0, 5, 2), 0)
	assert.Equal(t, nil, render.CheckImage([]string{
		"g....",
		"g....",
	}, im))

	im = o.Paint(image.Rect(0, 0, 5, 2), 1)
	assert.Equal(t, nil, render.CheckImage([]string{
		"g....",
		"g....",
	}, im))
	im = o.Paint(image.Rect(0, 0, 5, 2), 2)
	assert.Equal(t, nil, render.CheckImage([]string{
		"g....",
		"g....",
	}, im))

	// After that, linear motion to XEnd=4
	im = o.Paint(image.Rect(0, 0, 5, 2), 3)
	assert.Equal(t, nil, render.CheckImage([]string{
		"g....",
		"g....",
	}, im))
	im = o.Paint(image.Rect(0, 0, 5, 2), 4)
	assert.Equal(t, nil, render.CheckImage([]string{
		".g...",
		".g...",
	}, im))
	im = o.Paint(image.Rect(0, 0, 5, 2), 5)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..g..",
		"..g..",
	}, im))
	im = o.Paint(image.Rect(0, 0, 5, 2), 6)
	assert.Equal(t, nil, render.CheckImage([]string{
		"...g.",
		"...g.",
	}, im))
	im = o.Paint(image.Rect(0, 0, 5, 2), 7)
	assert.Equal(t, nil, render.CheckImage([]string{
		"....g",
		"....g",
	}, im))

	// Final frame is held
	im = o.Paint(image.Rect(0, 0, 5, 2), 8)
	assert.Equal(t, nil, render.CheckImage([]string{
		"....g",
		"....g",
	}, im))
	im = o.Paint(image.Rect(0, 0, 5, 2), 9)
	assert.Equal(t, nil, render.CheckImage([]string{
		"....g",
		"....g",
	}, im))

	// Requesting frames beyond FrameCount() returns final frame
	// as well
	im = o.Paint(image.Rect(0, 0, 5, 2), 10)
	assert.Equal(t, nil, render.CheckImage([]string{
		"....g",
		"....g",
	}, im))
	im = o.Paint(image.Rect(0, 0, 5, 2), 101212)
	assert.Equal(t, nil, render.CheckImage([]string{
		"....g",
		"....g",
	}, im))
}

// If child is animated, then that animation plays while the position
// changes
func TestPositionedChildAnimation(t *testing.T) {
	// A box alternating color and size each frame
	child := render.Animation{
		Children: []render.Widget{
			render.Box{Width: 2, Height: 2, Color: color.RGBA{0xff, 0, 0, 0xff}},
			render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
		},
	}

	// Movement from lower left upwards to the right. Delay and
	// hold for 1 frame.
	o := AnimatedPositioned{
		Child:    child,
		Duration: 6,
		XStart:   0,
		YStart:   5,
		XEnd:     5,
		YEnd:     0,
		Curve:    LinearCurve{},
		Delay:    1,
		Hold:     1,
	}

	assert.Equal(t, 8, o.FrameCount())

	im := o.Paint(image.Rect(0, 0, 10, 6), 0)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"..........",
		"..........",
		"..........",
		"..........",
		"rr........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 6), 1)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"..........",
		"..........",
		"..........",
		"..........",
		"g.........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 6), 2)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"..........",
		"..........",
		"..........",
		".rr.......",
		".rr.......",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 6), 3)
	assert.Equal(t, nil, render.CheckImage([]string{
		"..........",
		"..........",
		"..........",
		"..g.......",
		"..........",
		"..........",
	}, im))

	// fast forward to the end
	im = o.Paint(image.Rect(0, 0, 10, 6), 6)
	assert.Equal(t, nil, render.CheckImage([]string{
		".....rr...",
		".....rr...",
		"..........",
		"..........",
		"..........",
		"..........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 6), 7)
	assert.Equal(t, nil, render.CheckImage([]string{
		".....g....",
		"..........",
		"..........",
		"..........",
		"..........",
		"..........",
	}, im))

	// past final frame, and past hold, animation still runs
	im = o.Paint(image.Rect(0, 0, 10, 6), 8)
	assert.Equal(t, nil, render.CheckImage([]string{
		".....rr...",
		".....rr...",
		"..........",
		"..........",
		"..........",
		"..........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 10, 6), 9)
	assert.Equal(t, nil, render.CheckImage([]string{
		".....g....",
		"..........",
		"..........",
		"..........",
		"..........",
		"..........",
	}, im))
}

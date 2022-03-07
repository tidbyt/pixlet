package animation

import (
	"image"
	"math"

	"github.com/fogleman/gg"

	"tidbyt.dev/pixlet/render"
)

// Bounce moves its child horizontally or vertically.
//
// The `bounce_direction` will be 'horizontal' and will move from right
// to left if left empty, if specified as 'vertical' the Bounce will
// will move from bottom to top.
//
// In horizontal mode the height of the Bounce will be that of its child,
// but its `width` must be specified explicitly. In vertical mode the width
// will be that of its child but the `height` must be specified explicitly.
//
// If the child's width fits fully, it will not move, unless `bounce_always`
// is set.
//
// The `pause` parameter controls the amount of frames to pause at the
// beginning and midpoint of the bounce animation.
//
// The `curve` will be 'linear' if left empty, if specified it will allow
// another easing function to be used. The predefined curves are 'linear',
// 'ease_in', 'ease_out' and 'ease_in_out'. It is also possible to specify
// a custom cubic bézier curve using the notation 'cubic-bezier(a, b, c, d)'
// with `a`, `b`, `c` and `d` being floating-point numbers.
//
// DOC(Child): Widget to potentially bounce
// DOC(Width): Width of the Bounce, required for horizontal
// DOC(Height): Height of the Bounce, required for vertical
// DOC(BounceDirection): Direction to bounce, 'vertical' or 'horizontal', default is horizontal
// DOC(BounceAlways): Bounce child, even if it fits entirely
// DOC(Pause): Pause duration at beginning and midpoint of animation, default and minimum value is 1
// DOC(Curve): Easing curve to use, default is 'linear'
//
// EXAMPLE BEGIN
// render.Bounce(
//      width=64,
//      child=render.Text("this won't fit in 64 pixels"),
// )
// EXAMPLE END
type Bounce struct {
	render.Widget
	Child           render.Widget `starlark:"child,required"`
	Width           int           `starlark:"width"`
	Height          int           `starlark:"height"`
	BounceDirection string        `starlark:"bounce_direction"`
	BounceAlways    bool          `starlark:"bounce_always"`
	Pause           int           `starlark:"pause"`
	Curve           Curve         `starlark:"curve"`
}

func (b Bounce) FrameCount() int {
	var img image.Image
	var diff int

	if b.isVertical() {
		img = b.Child.Paint(image.Rect(0, 0, render.DefaultFrameWidth, b.Height*10), 0)
		diff = b.Height - img.Bounds().Dy()
	} else {
		img = b.Child.Paint(image.Rect(0, 0, b.Width*10, render.DefaultFrameHeight), 0)
		diff = b.Width - img.Bounds().Dx()
	}

	if diff > 0 && !b.BounceAlways {
		return 1
	}

	pause := render.MaxInt(1, b.Pause)
	dist := render.AbsInt(diff)

	return (pause - 1) + dist + (pause - 1) + dist
}

func (b Bounce) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	var img image.Image
	var diff int

	if b.isVertical() {
		img = b.Child.Paint(image.Rect(0, 0, render.DefaultFrameWidth, b.Height*10), frameIdx)
		diff = b.Height - img.Bounds().Dy()
	} else {
		img = b.Child.Paint(image.Rect(0, 0, b.Width*10, render.DefaultFrameHeight), frameIdx)
		diff = b.Width - img.Bounds().Dx()
	}

	pause := render.MaxInt(1, b.Pause)
	dist := render.AbsInt(diff)

	bgnIdx := (pause - 1)
	midIdx := bgnIdx + dist
	cntIdx := midIdx + (pause - 1)
	endIdx := cntIdx + dist

	var value float64
	if diff > 0 && !b.BounceAlways {
		// child fits entirely
		value = 0.0
	} else if frameIdx < bgnIdx {
		// pause before animating
		value = b.Curve.Transform(0.0)
	} else if frameIdx < midIdx {
		// animate forwards
		percent := float64(frameIdx-bgnIdx) / float64(dist)
		value = b.Curve.Transform(percent)
	} else if frameIdx < cntIdx {
		// pause before changing direction
		value = b.Curve.Transform(1.0)
	} else if frameIdx < endIdx {
		// animate backwards (using reversed curve)
		percent := float64(frameIdx-cntIdx) / float64(dist)
		value = 1.0 - b.Curve.Reverse().Transform(percent)
	} else {
		// if more frames are requested, freeze at final frame
		value = b.Curve.Transform(0.0)
	}

	offset := render.SignInt(diff) * int(math.Round(float64(dist)*value))

	var dc *gg.Context
	if b.isVertical() {
		dc = gg.NewContext(img.Bounds().Dx(), b.Height)
		dc.DrawImage(img, 0, offset)
	} else {
		dc = gg.NewContext(b.Width, img.Bounds().Dy())
		dc.DrawImage(img, offset, 0)
	}

	return dc.Image()
}

func (b Bounce) isVertical() bool {
	return b.BounceDirection == "vertical"
}

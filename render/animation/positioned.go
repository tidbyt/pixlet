package animation

import (
	"image"
	"math"

	"github.com/fogleman/gg"

	"tidbyt.dev/pixlet/render"
)

// Animate a widget from start to end coordinates.
//
// **DEPRECATED**: Please use `animation.Transformation` instead.
//
// DOC(Child): Widget to animate
// DOC(XStart): Horizontal start coordinate
// DOC(XEnd): Horizontal end coordinate
// DOC(YStart): Vertical start coordinate
// DOC(YEnd): Vertical end coordinate
// DOC(Duration): Duration of animation in frames
// DOC(Curve): Easing curve to use, default is 'linear'
// DOC(Delay): Delay before animation in frames
// DOC(Hold): Delay after animation in frames
//
type AnimatedPositioned struct {
	render.Widget
	Child    render.Widget `starlark:"child,required"`
	XStart   int           `starlark:"x_start"`
	XEnd     int           `starlark:"x_end"`
	YStart   int           `starlark:"y_start"`
	YEnd     int           `starlark:"y_end"`
	Duration int           `starlark:"duration,required"`
	Curve    Curve         `starlark:"curve,required"`
	Delay    int           `starlark:"delay"`
	Hold     int           `starlark:"hold"`
}

func (o AnimatedPositioned) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	return bounds
}

func (o AnimatedPositioned) Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int) {
	var position float64

	if frameIdx < o.Delay {
		position = 0.0
	} else if frameIdx >= o.Delay+o.Duration {
		position = 0.9999999999
	} else {
		position = o.Curve.Transform(float64(frameIdx-o.Delay) / float64(o.Duration))
	}

	dx := 1
	if o.XStart > o.XEnd {
		dx = -1
	}
	dy := 1
	if o.YStart > o.YEnd {
		dy = -1
	}

	sx := int(math.Ceil(math.Abs(float64(o.XEnd-o.XStart)) * position))
	sy := int(math.Ceil(math.Abs(float64(o.YEnd-o.YStart)) * position))

	x := o.XStart + dx*sx
	y := o.YStart + dy*sy

	dc.Push()
	dc.Translate(float64(x), float64(y))
	o.Child.Paint(dc, bounds, frameIdx)
	dc.Pop()
}

func (o AnimatedPositioned) FrameCount() int {
	return o.Duration + o.Delay + o.Hold
}

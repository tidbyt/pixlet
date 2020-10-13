package animation

import (
	"image"
	"math"

	"github.com/fogleman/gg"

	"github.com/tidbyt/pixlet/render"
)

type AnimatedPositioned struct {
	render.Widget
	Child    render.Widget
	XStart   int
	XEnd     int
	YStart   int
	YEnd     int
	Duration int
	Curve    Curve
	Delay    int
	Hold     int
}

func (o AnimatedPositioned) Paint(bounds image.Rectangle, frameIdx int) image.Image {
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

	im := o.Child.Paint(bounds, frameIdx)

	dc := gg.NewContext(bounds.Dx(), bounds.Dy())
	dc.DrawImage(im, x, y)
	return dc.Image()
}

func (o AnimatedPositioned) FrameCount() int {
	return o.Duration + o.Delay + o.Hold
}

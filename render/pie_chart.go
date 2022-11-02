package render

import (
	"image"
	"image/color"
	"math"

	"github.com/tidbyt/gg"
)

// PieChart draws the PieChart of a circle with the given `start` and `end`
// angles, in radians, with size `diameter` and a defined `color`.
// If `fill` is passed, the PieChart will be filled in with `color`,
// to the center of the circle, like a piece of a pie.
//
// DOC(Color): Fill and stroke color
// DOC(Diameter): Diameter of the circle
// DOC(Start): Angle, in radians, where the circle's PieChart begins. 0 rad is at 3pm on a clock.
// DOC(End): Angle, in radians, where the circle's PieChart ends.
// Doc(Fill): Whether or not the PieChart should be filled in.
//
// EXAMPLE BEGIN
// render.PieChart(
//      color="#fff",
//      diameter=30,
//      start=0,
//      end=math.pi / 2,
//      fill=True,
// )
// EXAMPLE END
type PieChart struct {
	Widget

	Colors   []color.Color `starlark:"colors, required"`
	Weights  []float64     `starlark:"weights, required"`
	Diameter int           `starlark:"diameter,required"`
}

func (c PieChart) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	return image.Rect(0, 0, c.Diameter, c.Diameter)
}

func (c PieChart) Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int) {
	total := 0.0
	for _, v := range c.Weights {
		total += v
	}

	r := float64(c.Diameter) / 2

	start := 0.0
	for i, v := range c.Weights {
		end := start + v/total
		dc.SetColor(c.Colors[i%len(c.Colors)])
		dc.DrawArc(r, r, r, start*2*math.Pi, end*2*math.Pi)
		dc.LineTo(r, r)
		dc.LineTo(r+r*math.Cos(start*2*math.Pi), r+r*math.Sin(start*2*math.Pi))
		dc.Fill()
		start = end
	}
}

func (c PieChart) FrameCount() int {
	return 1
}

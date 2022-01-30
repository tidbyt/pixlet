package render

import (
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

// Circle draws a circle with the given `diameter` and `color`. If a
// `child` widget is provided, it is drawn in the center of the
// circle.
//
// If the `color` attribute is a single string, that string will be
// interpreted as an HTML-like hexadecimal color code.  If it is a
// pair of `(string, float)`, the string will be interpreted as
// an HTML-like hexadecimal color code and the float must be a value
// between 0.0 (fully transparent) and 1.0 (fully opaque) for the
// transparency of the color.
//
// DOC(Child): Widget to place in the center of the circle
// DOC(Color): Fill color
// DOC(Diameter): Diameter of the circle
//
// EXAMPLE BEGIN
// render.Circle(
//      color="#666",
//      diameter=30,
//      child=render.Circle(color="#0ff", diameter=10),
// )
// EXAMPLE END
type Circle struct {
	Widget

	Child    Widget
	Color    color.Color `starlark:"color, required"`
	Diameter int         `starlark:"diameter,required"`
}

func (c Circle) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	dc := gg.NewContext(c.Diameter, c.Diameter)
	dc.SetColor(c.Color)

	r := float64(c.Diameter) / 2
	dc.DrawCircle(r, r, r)
	dc.Fill()

	if c.Child != nil {
		center := math.Ceil(
			float64(c.Diameter) / 2,
		)

		im := c.Child.Paint(image.Rect(0, 0, c.Diameter, c.Diameter), frameIdx)
		dc.DrawImageAnchored(
			im,
			int(center),
			int(center),
			0.5,
			0.5,
		)
	}

	return dc.Image()
}

func (c Circle) FrameCount() int {
	if c.Child != nil {
		return c.Child.FrameCount()
	}
	return 1
}

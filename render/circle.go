package render

import (
	"image"
	"image/color"
	"math"

	"tidbyt.dev/pixlet/render/canvas"
)

// Circle draws a circle with the given `diameter` and `color`. If a
// `child` widget is provided, it is drawn in the center of the
// circle.
//
// DOC(Child): Widget to place in the center of the circle
// DOC(Color): Fill color
// DOC(Diameter): Diameter of the circle
//
// EXAMPLE BEGIN
// render.Circle(
//
//	color="#666",
//	diameter=30,
//	child=render.Circle(color="#0ff", diameter=10),
//
// )
// EXAMPLE END
type Circle struct {
	Widget

	Child    Widget
	Color    color.Color `starlark:"color, required"`
	Diameter int         `starlark:"diameter,required"`
}

func (c Circle) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	return image.Rect(0, 0, c.Diameter, c.Diameter)
}

func (c Circle) Paint(dc canvas.Canvas, bounds image.Rectangle, frameIdx int) {
	dc.SetColor(c.Color)

	r := float64(c.Diameter) / 2
	dc.AddCircle(r, r, r)
	dc.FillPath()

	if c.Child != nil {
		dc.Push()
		childBounds := c.Child.PaintBounds(image.Rect(0, 0, c.Diameter, c.Diameter), frameIdx)

		// This is a bit convoluted to obtain the same rounding behavior as with the old
		// local context rendering
		center := math.Ceil(
			float64(c.Diameter) / 2,
		)
		x := int(center)
		y := int(center)
		x -= int(0.5 * float64(childBounds.Size().X))
		y -= int(0.5 * float64(childBounds.Size().Y))

		dc.Translate(float64(x), float64(y))

		c.Child.Paint(dc, image.Rect(0, 0, c.Diameter, c.Diameter), frameIdx)
		dc.Pop()
	}
}

func (c Circle) FrameCount() int {
	if c.Child != nil {
		return c.Child.FrameCount()
	}
	return 1
}

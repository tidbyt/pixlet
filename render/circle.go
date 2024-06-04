package render

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/tidbyt/gg"
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
//      color="#666",
//      diameter=30,
//      child=render.Circle(color="#0ff", diameter=10),
// )
// EXAMPLE END
type Circle struct {
	Type string `starlark:"-"`

	Child    Widget
	Color    color.RGBA `starlark:"color, required"`
	Diameter int        `starlark:"diameter,required"`
}

func (c Circle) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	return image.Rect(0, 0, c.Diameter, c.Diameter)
}

func (c Circle) Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int) {
	dc.SetColor(c.Color)

	r := float64(c.Diameter) / 2
	dc.DrawCircle(r, r, r)
	dc.Fill()

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

func (c *Circle) UnmarshalJSON(data []byte) error {
	type Alias Circle
	aux := &struct {
		Child json.RawMessage
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Child != nil {
		child, err := UnmarshalWidgetJSON(aux.Child)
		if err != nil {
			return err
		}
		c.Child = child
	}
	return nil
}

func (c *Circle) MarshalJSON() ([]byte, error) {
	col := ""
	if c.Color != (color.RGBA{}) {
		r, g, b, a := c.Color.RGBA()
		col = fmt.Sprintf("#%02x%02x%02x%02x", r>>8, g>>8, b>>8, a>>8)
	}
	type Alias Circle
	aux := &struct {
		*Alias
		Color string
	}{
		Alias: (*Alias)(c),
		Color: col,
	}

	return json.Marshal(aux)
}

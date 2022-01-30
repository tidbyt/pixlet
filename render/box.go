package render

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

// A Box is a rectangular widget that can hold a child widget.
//
// Boxes are transparent unless `color` is provided. They expand to
// fill all available space, unless `width` and/or `height` is
// provided. Boxes can have a `child`, which will be centered in the
// box, and the child can be padded (via `padding`).
//
// If the `color` attribute is a single string, that string will be
// interpreted as an HTML-like hexadecimal color code.  If it is a
// pair of `(string, float)`, the string will be interpreted as
// an HTML-like hexadecimal color code and the float must be a value
// between 0.0 (fully transparent) and 1.0 (fully opaque) for the
// transparency of the color.
//
// DOC(Child): Child to center inside box
// DOC(Width): Limits Box width
// DOC(Height): Limits Box height
// DOC(Padding): Padding around the child widget
// DOC(Color): Background color
//
// EXAMPLE BEGIN
// render.Box(
//      color="#00f",
//      child=render.Box(
//           width=20,
//           height=10,
//           color="#f00",
//      )
// )
// EXAMPLE END
type Box struct {
	Widget
	Child         Widget
	Width, Height int
	Padding       int
	Color         color.Color
}

func (b Box) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	w, h := b.Width, b.Height
	if w == 0 {
		w = bounds.Dx()
	}
	if h == 0 {
		h = bounds.Dy()
	}

	dc := gg.NewContext(w, h)
	if b.Color != nil {
		dc.SetColor(b.Color)
		dc.Clear()
	}

	if b.Child != nil {
		chW := w - b.Padding*2
		chH := h - b.Padding*2

		if chW < 0 || chH < 0 {
			// padding makes the child invisible, no point painting it
		} else {
			dc.DrawRectangle(
				float64(b.Padding),
				float64(b.Padding),
				float64(chW),
				float64(chH),
			)
			dc.Clip()

			im := b.Child.Paint(image.Rect(0, 0, chW, chH), frameIdx)
			dc.DrawImageAnchored(
				im,
				w/2,
				h/2,
				0.5,
				0.5,
			)
		}
	}

	return dc.Image()
}

func (b Box) FrameCount() int {
	if b.Child != nil {
		return b.Child.FrameCount()
	}
	return 1
}

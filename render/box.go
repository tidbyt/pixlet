package render

import (
	"image"
	"image/color"

	"github.com/tidbyt/gg"
)

// A Box is a rectangular widget that can hold a child widget.
//
// Boxes are transparent unless `color` is provided. They expand to
// fill all available space, unless `width` and/or `height` is
// provided. Boxes can have a `child`, which will be centered in the
// box, and the child can be padded (via `padding`).
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

func (b Box) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	w, h := b.Width, b.Height
	if w == 0 {
		w = bounds.Dx()
	}
	if h == 0 {
		h = bounds.Dy()
	}
	return image.Rect(0, 0, w, h)
}

func (b Box) Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int) {
	w, h := b.Width, b.Height
	if w == 0 {
		w = bounds.Dx()
	}
	if h == 0 {
		h = bounds.Dy()
	}

	if b.Color != nil {
		dc.SetColor(b.Color)
		dc.DrawRectangle(0, 0, float64(w), float64(h))
		dc.Fill()
	}

	if b.Child != nil {
		chW := w - b.Padding*2
		chH := h - b.Padding*2

		if chW < 0 || chH < 0 {
			// padding makes the child invisible, no point painting it
		} else {
			dc.Push()

			dc.DrawRectangle(
				float64(b.Padding),
				float64(b.Padding),
				float64(chW),
				float64(chH),
			)
			dc.Clip()

			childBounds := b.Child.PaintBounds(image.Rect(0, 0, chW, chH), frameIdx)

			// This is a bit convoluted to obtain the same rounding behavior as with the old
			// local-context rendering
			x := w / 2
			y := h / 2
			x -= int(0.5 * float64(childBounds.Size().X))
			y -= int(0.5 * float64(childBounds.Size().Y))

			dc.Translate(float64(x), float64(y))
			b.Child.Paint(dc, image.Rect(0, 0, chW, chH), frameIdx)
			dc.Pop()
		}
	}
}

func (b Box) FrameCount() int {
	if b.Child != nil {
		return b.Child.FrameCount()
	}
	return 1
}

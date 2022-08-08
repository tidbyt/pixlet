package render

import (
	"image"

	"github.com/fogleman/gg"
)

// Sequence renders a list of child widgets in sequence.
//
// Each child widget is rendered for the duration of its
// frame count, then the next child wiget in the list will
// be rendered and so on.
//
// It comes in quite useful when chaining animations.
// If you want to know more about that, go check
// out the [animation](animation.md) documentation.
//
// DOC(Children): List of child widgets
//
// EXAMPLE BEGIN
// render.Sequence(
//   children = [
//     animation.Transformation(...),
//     animation.Transformation(...),
//     ...
//   ],
// ),
// EXAMPLE END
type Sequence struct {
	Widget

	Children []Widget `starlark:"children,required"`
}

func (s Sequence) FrameCount() int {
	fc := 0

	for _, c := range s.Children {
		fc += c.FrameCount()
	}

	return fc
}

func (s Sequence) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	fc := 0

	for _, c := range s.Children {
		if frameIdx < fc+c.FrameCount() {
			return c.PaintBounds(bounds, frameIdx-fc)
		}

		fc += c.FrameCount()
	}

	return image.Rect(0, 0, 0, 0)
}

func (s Sequence) Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int) {
	fc := 0

	for _, c := range s.Children {
		if frameIdx < fc+c.FrameCount() {
			dc.Push()
			c.Paint(dc, bounds, frameIdx-fc)
			dc.Pop()
		}

		fc += c.FrameCount()
	}
}

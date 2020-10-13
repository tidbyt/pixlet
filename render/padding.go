package render

import (
	"image"

	"github.com/fogleman/gg"
)

type Insets struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

// Padding places padding around its child.
//
// If the `pad` attribute is a single integer, that amount of padding
// will be placed on all sides of the child. If it's a 4-tuple `(left,
// top, right, bottom)`, then padding will be placed on the sides
// accordingly.
//
// DOC(Child): The Widget to place padding around
// DOC(Expanded): This is a confusing parameter
// DOC(Pad): Padding around the child
type Padding struct {
	Widget

	Child    Widget `starlark:"child,required"`
	Pad      Insets
	Expanded bool
}

func (p Padding) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	im := p.Child.Paint(
		image.Rect(0, 0, bounds.Dx()-p.Pad.Left-p.Pad.Right, bounds.Dy()-p.Pad.Top-p.Pad.Bottom),
		frameIdx,
	)

	var width, height int
	if p.Expanded {
		width = bounds.Dx()
		height = bounds.Dy()
	} else {
		width = im.Bounds().Dx() + p.Pad.Left + p.Pad.Right
		height = im.Bounds().Dy() + p.Pad.Top + p.Pad.Bottom
	}

	dc := gg.NewContext(width, height)
	dc.DrawRectangle(
		float64(p.Pad.Left),
		float64(p.Pad.Top),
		float64(width-p.Pad.Left-p.Pad.Right),
		float64(height-p.Pad.Top-p.Pad.Bottom),
	)
	dc.Clip()
	dc.DrawImage(im, p.Pad.Left, p.Pad.Top)

	return dc.Image()
}

func (p Padding) FrameCount() int {
	if p.Child != nil {
		return p.Child.FrameCount()
	}
	return 1
}

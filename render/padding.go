package render

import (
	"image"
	"image/color"

	"github.com/tidbyt/gg"
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
// DOC(Color): Background color
type Padding struct {
	Type string `starlark:"-"`

	Child    Widget `starlark:"child,required"`
	Pad      Insets
	Expanded bool
	Color    color.RGBA
}

func (p Padding) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	cb := p.Child.PaintBounds(
		image.Rect(0, 0, bounds.Dx()-p.Pad.Left-p.Pad.Right, bounds.Dy()-p.Pad.Top-p.Pad.Bottom),
		frameIdx,
	)

	var width, height int
	if p.Expanded {
		width = bounds.Dx()
		height = bounds.Dy()
	} else {
		width = cb.Dx() + p.Pad.Left + p.Pad.Right
		height = cb.Dy() + p.Pad.Top + p.Pad.Bottom
	}

	return image.Rect(0, 0, width, height)
}

func (p Padding) Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int) {
	cb := p.Child.PaintBounds(
		image.Rect(0, 0, bounds.Dx()-p.Pad.Left-p.Pad.Right, bounds.Dy()-p.Pad.Top-p.Pad.Bottom),
		frameIdx,
	)

	var width, height int
	if p.Expanded {
		width = bounds.Dx()
		height = bounds.Dy()
	} else {
		width = cb.Dx() + p.Pad.Left + p.Pad.Right
		height = cb.Dy() + p.Pad.Top + p.Pad.Bottom
	}

	if p.Color != (color.RGBA{}) {
		dc.SetColor(p.Color)
		dc.DrawRectangle(0, 0, float64(width), float64(height))
		dc.Fill()
	}

	dc.Push()

	// Some apps use negative padding as a positioning hack.
	clipLeft := p.Pad.Left
	clipTop := p.Pad.Top
	if clipLeft < 0 {
		clipLeft = 0
	}
	if clipTop < 0 {
		clipTop = 0
	}

	dc.DrawRectangle(
		float64(clipLeft),
		float64(clipTop),
		float64(width-p.Pad.Left-p.Pad.Right),
		float64(height-p.Pad.Top-p.Pad.Bottom),
	)
	dc.Clip()

	dc.Translate(float64(p.Pad.Left), float64(p.Pad.Top))

	p.Child.Paint(dc, image.Rect(0, 0, bounds.Dx()-p.Pad.Left-p.Pad.Right, bounds.Dy()-p.Pad.Top-p.Pad.Bottom),
		frameIdx)
	dc.Pop()
}

func (p Padding) FrameCount() int {
	if p.Child != nil {
		return p.Child.FrameCount()
	}
	return 1
}

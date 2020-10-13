package render

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

// WrappedText draws multi-line text.
//
// The optional `width` and `height` parameters limit the drawing
// area. If not set, WrappedText will use as much vertical and
// horizontal space as possible to fit the text.
//
// DOC(Content): The text string to draw
// DOC(Font): Desired font face
// DOC(Height): Limits height of the area on which text may be drawn
// DOC(Width): Limits width of the area on which text may be drawn
// DOC(LineSpacing): Controls spacing between lines
// DOC(Color): Desired font color
//
// EXAMPLE BEGIN
// render.WrappedText(
//       content="this is a multi-line text string",
//       width=50,
//       color="#fa0",
// )
// EXAMPLE END
type WrappedText struct {
	Widget

	Content     string `starlark:"content,required"`
	Font        string
	Height      int
	Width       int
	LineSpacing int
	Color       color.Color
}

func (tw WrappedText) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	face := Font[DefaultFontFace]
	if tw.Font != "" {
		face = Font[tw.Font]
	}

	// The bounds provided by user or parent widget
	width := tw.Width
	if width == 0 {
		width = bounds.Dx()
	}
	height := tw.Height
	if height == 0 {
		height = bounds.Dy()
	}

	// Compute size of multi line string
	//
	// NOTE: Can't use dc.MeasureMultilineString() here. It only
	// deals with texts that have actual \n in them.
	dc := gg.NewContext(width, 0)
	dc.SetFontFace(face)
	w := 0.0
	h := 0.0
	for _, line := range dc.WordWrap(tw.Content, float64(width)) {
		lw, lh := dc.MeasureString(line)
		if lw > w {
			w = lw
		}
		h += lh
	}

	// Size of drawing context
	if tw.Width != 0 {
		width = tw.Width
	} else if int(w) < bounds.Dx() {
		width = int(w)
	} else {
		width = bounds.Dx()
	}

	if tw.Height != 0 {
		height = tw.Height
	} else if int(h) < bounds.Dy() {
		height = int(h)
	} else {
		height = bounds.Dy()
	}

	metrics := face.Metrics()
	descent := metrics.Descent.Floor()

	// And draw
	dc = gg.NewContext(width, height)
	dc.SetFontFace(face)
	if tw.Color != nil {
		dc.SetColor(tw.Color)
	} else {
		dc.SetColor(DefaultFontColor)
	}

	dc.DrawStringWrapped(
		tw.Content,
		0,
		float64(-descent),
		0,
		0,
		float64(width),
		(float64(tw.LineSpacing)+dc.FontHeight())/dc.FontHeight(),
		gg.AlignLeft,
	)

	return dc.Image()
}

func (tw WrappedText) FrameCount() int {
	return 1
}

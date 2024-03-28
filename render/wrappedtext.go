package render

import (
	"image"
	"image/color"

	"github.com/tidbyt/gg"
	"tidbyt.dev/pixlet/fonts"
	"tidbyt.dev/pixlet/render/canvas"
)

// WrappedText draws multi-line text.
//
// The optional `width` and `height` parameters limit the drawing
// area. If not set, WrappedText will use as much vertical and
// horizontal space as possible to fit the text.
//
// Alignment of the text is controlled by passing one of the following `align` values:
// - `"left"`: align text to the left
// - `"center"`: align text in the center
// - `"right"`: align text to the right
//
// DOC(Content): The text string to draw
// DOC(Font): Desired font face
// DOC(Height): Limits height of the area on which text may be drawn
// DOC(Width): Limits width of the area on which text may be drawn
// DOC(LineSpacing): Controls spacing between lines
// DOC(Color): Desired font color
// DOC(Align): Text Alignment
// EXAMPLE BEGIN
// render.WrappedText(
//
//	content="this is a multi-line text string",
//	width=50,
//	color="#fa0",
//
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
	Align       string
}

func (tw WrappedText) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	if tw.Font == "" {
		tw.Font = DefaultFontFace
	}
	face := fonts.GetFont(tw.Font).Font.NewFace()
	// The bounds provided by user or parent widget
	width := tw.Width
	if width == 0 {
		width = bounds.Dx()
	}
	height := tw.Height
	if height == 0 {
		height = bounds.Dy()
	}
	linespace := float64(tw.LineSpacing)
	if linespace <= 0 {
		linespace = 0
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
		h += lh + linespace
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

	return image.Rect(0, 0, width, height)
}

func (tw WrappedText) Paint(dc canvas.Canvas, bounds image.Rectangle, frameIdx int) {
	if tw.Font == "" {
		tw.Font = DefaultFontFace
	}
	font := fonts.GetFont(tw.Font)
	// Text alignment
	align := canvas.AlignLeft
	if tw.Align == "center" {
		align = canvas.AlignCenter
	} else if tw.Align == "right" {
		align = canvas.AlignRight
	}

	width := tw.PaintBounds(bounds, frameIdx).Dx()

	dc.SetFont(font)
	if tw.Color != nil {
		dc.SetColor(tw.Color)
	} else {
		dc.SetColor(DefaultFontColor)
	}

	dc.DrawStringWrapped(
		0,
		0,
		float64(width),
		float64(tw.LineSpacing),
		tw.Content,
		align,
	)
}

func (tw WrappedText) FrameCount() int {
	return 1
}

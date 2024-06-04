package render

import (
	"image"
	"image/color"

	"github.com/tidbyt/gg"

	"golang.org/x/image/font"
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
	Type string `starlark:"-"`

	Content     string `starlark:"content,required"`
	Font        string
	Height      int
	Width       int
	LineSpacing int
	Color       color.RGBA
	Align       string

	face font.Face
}

func (tw *WrappedText) Init() error {
	if tw.Font == "" {
		tw.Font = DefaultFontFace
	}

	face, err := GetFont(tw.Font)
	if err != nil {
		return err
	}

	tw.face = face

	return nil
}

func (tw *WrappedText) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
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
	dc.SetFontFace(tw.face)
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

func (tw *WrappedText) Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int) {
	// Text alignment
	align := gg.AlignLeft
	if tw.Align == "center" {
		align = gg.AlignCenter
	} else if tw.Align == "right" {
		align = gg.AlignRight
	}

	width := tw.PaintBounds(bounds, frameIdx).Dx()

	metrics := tw.face.Metrics()
	descent := metrics.Descent.Floor()

	dc.SetFontFace(tw.face)
	if tw.Color != (color.RGBA{}) {
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
		align,
	)
}

func (tw *WrappedText) FrameCount() int {
	return 1
}

package render

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

var (
	DefaultFontFace  = "tb-8"
	DefaultFontColor = color.White
)

// Text draws a string of text on a single line.
//
// By default, the text will use the "tb-8" font, but other fonts can
// be chosen via the `font` attribute. The `height` and `offset`
// parameters allow fine tuning of the vertical layout of the
// string. Take a look at the [font documentation](fonts.md) for more
// information.
//
// DOC(Content): The text string to draw
// DOC(Font): Desired font face
// DOC(Height): Limits height of the area on which text is drawn
// DOC(Offset): Shifts position of text vertically.
// DOC(Color): Desired font color
//
// EXAMPLE BEGIN
// render.Text(content="Tidbyt!", color="#099")
// EXAMPLE END
type Text struct {
	Widget
	Content string `starlark:"content,required"`
	Font    string
	Height  int
	Offset  int
	Color   color.Color

	img image.Image
}

func (t *Text) Size() (int, int) {
	if t.img == nil {
		t.paint()
	}

	return t.img.Bounds().Dx(), t.img.Bounds().Dy()
}

func (t *Text) Paint(
	bounds image.Rectangle, frameIdx int,
) image.Image {
	if t.img == nil {
		t.paint()
	}
	return t.img
}

func (t *Text) paint() {
	face := Font[DefaultFontFace]
	if t.Font != "" {
		face = Font[t.Font]
	}

	dc := gg.NewContext(0, 0)
	dc.SetFontFace(face)

	w, _ := dc.MeasureString(t.Content)
	width := int(w)

	metrics := face.Metrics()
	ascent := metrics.Ascent.Floor()
	descent := metrics.Descent.Floor()

	height := ascent + descent
	if t.Height != 0 {
		height = t.Height
	}

	dc = gg.NewContext(width, height)
	dc.SetFontFace(face)
	if t.Color != nil {
		dc.SetColor(t.Color)
	} else {
		dc.SetColor(DefaultFontColor)
	}

	dc.DrawString(t.Content, 0, float64(height-descent-t.Offset))

	t.img = dc.Image()
}

func (t Text) FrameCount() int {
	return 1
}

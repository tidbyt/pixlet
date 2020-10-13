package render

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

type Tracer struct {
	Widget
	Path        Path
	TraceLength int
}

func (t Tracer) FrameCount() int {
	return t.Path.Length()
}

func (t Tracer) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	width, height := t.Path.Size()
	dc := gg.NewContext(width, height)

	x, y := t.Path.Point(frameIdx)

	dc.SetColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	dc.SetPixel(x, y)

	for i := 0; i < t.TraceLength; i++ {
		col := uint8(0xdd - i*(0xff/t.TraceLength))
		dc.SetColor(color.RGBA{col, col, col, 0xff})
		x, y := t.Path.Point(frameIdx - (i + 1))
		dc.SetPixel(x, y)
	}

	return dc.Image()
}

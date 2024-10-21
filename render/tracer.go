package render

import (
	"image"
	"image/color"

	"github.com/tidbyt/gg"
)

type Tracer struct {
	Type string `starlark:"-"`

	Path        Path
	TraceLength int
}

func (t Tracer) FrameCount() int {
	return t.Path.Length()
}

func (t Tracer) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	width, height := t.Path.Size()
	return image.Rect(0, 0, width, height)
}

func (t Tracer) Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int) {
	x, y := t.Path.Point(frameIdx)

	dc.SetColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	tx, ty := dc.TransformPoint(float64(x), float64(y))
	dc.SetPixel(int(tx), int(ty))

	for i := 0; i < t.TraceLength; i++ {
		col := uint8(0xdd - i*(0xff/t.TraceLength))
		dc.SetColor(color.RGBA{col, col, col, 0xff})
		x, y := t.Path.Point(frameIdx - (i + 1))
		tx, ty := dc.TransformPoint(float64(x), float64(y))
		dc.SetPixel(int(tx), int(ty))
	}
}

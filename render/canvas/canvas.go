package canvas

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"sync"

	"golang.org/x/image/font"
)

type TextAlign string

const (
	AlignLeft   TextAlign = "left"
	AlignCenter TextAlign = "center"
	AlignRight  TextAlign = "right"
)

type Canvas interface {
	// AddArc adds an arc to the current path. It has the given radius, is centered at (x, y),
	// starts at startAngle, and ends at endAngle.
	AddArc(x, y, r, startAngle, endAngle float64)

	// AddCircle adds a circle to the current path. It has the given radius and is centered at (x, y).
	AddCircle(x, y, r float64)

	// AddLineTo adds a line to the current path, from the current point to (x, y).
	AddLineTo(x, y float64)

	// AddRectangle adds a rectangle to the current path.
	AddRectangle(x, y, w, h float64)

	// Clear fills the clip with the current color.
	Clear()

	// ClipRectangle sets the clip to the intersection of the current clip and the given rectangle.
	ClipRectangle(x, y, w, h float64)

	// DrawGoImage draws the given image at (x, y).
	DrawGoImage(x, y float64, img image.Image)

	// DrawImageFromBuffer draws the given image at (x, y).
	DrawImageFromBuffer(x, y, w, h float64, img []byte)

	// DrawPixel draws a pixel at (x, y).
	DrawPixel(x, y int)

	// DrawString draws the given text at (x, y).
	DrawString(x, y float64, text string)

	// DrawStringWrapped draws the given text at (x, y), wrapping it to fit within the given width
	// and using the given spacing between lines.
	DrawStringWrapped(x, y, w, spacing float64, text string, align TextAlign)

	// FillPath fills the current path.
	FillPath()

	// FontHeight returns the height of the current font.
	FontHeight() float64

	// Image returns an image of the canvas.
	Image() image.Image

	// MeasureString returns the width and height of the given text with
	// the current font.
	MeasureString(text string) (w, h float64)

	// Pop restores the previous state of the clip and transformations from the stack.
	Pop()

	// Push saves the current state of the clip and transformations to the stack.
	Push()

	// Rotate rotates the current transformation matrix by angle radians, clockwise.
	Rotate(angle float64)

	// Scale updates the current transformation matrix to scale by (x, y).
	Scale(x, y float64)

	// SetColor sets the current color for drawing.
	SetColor(color color.Color)

	// SetFontFace sets the current font face.
	SetFontFace(fontFace font.Face)

	// TransformPoint does a thing.
	TransformPoint(x, y float64) (ax, ay float64)

	// Translate translates the canvas by (dx, dy) for future drawing operations.
	Translate(dx, dy float64)
}

var (
	encoder imageEncoder
)

type imageEncoder interface {
	Encode(w io.Writer, m image.Image) error
}

type pngBufferPool struct {
	pool sync.Pool
}

func init() {
	encoder = &png.Encoder{
		CompressionLevel: png.DefaultCompression,
		BufferPool: &pngBufferPool{
			pool: sync.Pool{},
		},
	}
}

func (p *pngBufferPool) Get() *png.EncoderBuffer {
	if buf := p.pool.Get(); buf != nil {
		return buf.(*png.EncoderBuffer)
	}

	return nil
}

func (p *pngBufferPool) Put(buf *png.EncoderBuffer) {
	p.pool.Put(buf)
}

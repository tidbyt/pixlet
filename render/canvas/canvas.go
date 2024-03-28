package canvas

import (
	"image"
	"image/color"

	"tidbyt.dev/pixlet/fonts"
)

type TextAlign string

const (
	AlignLeft   TextAlign = "left"
	AlignCenter TextAlign = "center"
	AlignRight  TextAlign = "right"
)

type Provider func(width, height int) Canvas

var (
	DefaultCanvasProvider = NewGGCanvas
)

type Canvas interface {
	// AddArc adds an arc to the current path. It has the given radius, is centered at (x, y),
	// starts at startAngle, and ends at endAngle.
	AddArc(x, y, r, startAngle, endAngle float64)

	// AddCircle adds a circle to the current path. It has the given radius and is centered at (x, y).
	AddCircle(x, y, r float64)

	// AddLineTo adds a line to the current path, from the current point to (x, y).
	AddLineTo(x, y float64)

	// AddPixel adds a pixel to the current path at (x, y).
	AddPixel(x, y int)

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

	// DrawString draws the given text at (x, y).
	DrawString(x, y float64, text string)

	// DrawStringWrapped draws the given text at (x, y), wrapping it to fit within the given width
	// and with optional extra spacing between lines.
	DrawStringWrapped(x, y, w, extraSpacing float64, text string, align TextAlign)

	// FillPath fills the current path.
	FillPath()

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

	// SetFont sets the current font.
	SetFont(font *fonts.Font)

	// Translate translates the canvas by (dx, dy) for future drawing operations.
	Translate(dx, dy float64)
}

func NewCanvas(width, height int) Canvas {
	return DefaultCanvasProvider(width, height)
}

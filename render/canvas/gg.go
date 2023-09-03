package canvas

import (
	"bytes"
	"image"
	"image/color"

	"github.com/tidbyt/gg"
	"golang.org/x/image/font"
)

type GGCanvas struct {
	dc *gg.Context
}

func NewGGCanvas(width, height int) Canvas {
	dc := gg.NewContext(width, height)
	return &GGCanvas{dc: dc}
}

func (c *GGCanvas) AddArc(x, y, r, startAngle, endAngle float64) {
	c.dc.DrawArc(x, y, r, startAngle, endAngle)
}

func (c *GGCanvas) AddCircle(x, y, r float64) {
	c.dc.DrawCircle(x, y, r)
}

func (c *GGCanvas) AddLineTo(x, y float64) {
	c.dc.LineTo(x, y)
}

func (c *GGCanvas) AddRectangle(x, y, w, h float64) {
	c.dc.DrawRectangle(x, y, w, h)
}

func (c *GGCanvas) Clear() {
	c.dc.Clear()
}

func (c *GGCanvas) ClipRectangle(x, y, w, h float64) {
	c.dc.DrawRectangle(x, y, w, h)
	c.dc.Clip()
}

func (c *GGCanvas) DrawGoImage(x, y float64, img image.Image) {
	c.dc.DrawImage(img, int(x), int(y))
}

func (c *GGCanvas) DrawImageFromBuffer(x, y, w, h float64, img []byte) {
	im, _, _ := image.Decode(bytes.NewReader(img))
	c.dc.DrawImage(im, int(x), int(y))
}

func (c *GGCanvas) DrawPixel(x, y int) {
	c.dc.SetPixel(x, y)
}

func (c *GGCanvas) DrawString(x, y float64, text string) {
	c.dc.DrawString(text, x, y)
}

func (c *GGCanvas) DrawStringWrapped(x, y, w, spacing float64, text string, align TextAlign) {
	var alignFlag gg.Align
	switch align {
	case AlignLeft:
		alignFlag = gg.AlignLeft
	case AlignCenter:
		alignFlag = gg.AlignCenter
	case AlignRight:
		alignFlag = gg.AlignRight
	}

	c.dc.DrawStringWrapped(text, x, y, 0, 0, w, spacing, alignFlag)
}

func (c *GGCanvas) FillPath() {
	c.dc.Fill()
}

func (c *GGCanvas) FontHeight() float64 {
	return c.dc.FontHeight()
}

func (c *GGCanvas) Image() image.Image {
	return c.dc.Image()
}

func (c *GGCanvas) MeasureString(text string) (w, h float64) {
	return c.dc.MeasureString(text)
}

func (c *GGCanvas) Pop() {
	c.dc.Pop()
}

func (c *GGCanvas) Push() {
	c.dc.Push()
}

func (c *GGCanvas) Rotate(angle float64) {
	c.dc.Rotate(angle)
}

func (c *GGCanvas) Scale(x, y float64) {
	c.dc.Scale(x, y)
}

func (c *GGCanvas) SetColor(color color.Color) {
	c.dc.SetColor(color)
}

func (c *GGCanvas) SetFontFace(fontFace font.Face) {
	c.dc.SetFontFace(fontFace)
}

func (c *GGCanvas) TransformPoint(x, y float64) (ax, ay float64) {
	return c.dc.TransformPoint(x, y)
}

func (c *GGCanvas) Translate(dx, dy float64) {
	c.dc.Translate(dx, dy)
}

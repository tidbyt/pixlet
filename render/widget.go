package render

import (
	"image"
)

// A Widget is a self-contained object that can render itself as an image.
type Widget interface {
	Paint(bounds image.Rectangle, frameIdx int) image.Image
	FrameCount() int
}

type WidgetStaticSize interface {
	Size() (int, int)
}

// Computes a (mod m). Useful for handling frameIdx > num available
// frames in Widget.Paint()
func ModInt(a, m int) int {
	a = a % m
	if a < 0 {
		a += m
	}
	return a
}

package render

import (
	"image"
)

// A Widget is a self-contained object that can render itself as an image.
type Widget interface {
	Paint(bounds image.Rectangle, frameIdx int) image.Image
	FrameCount() int
}

// Widgets can require initialization
type WidgetWithInit interface {
	Init() error
}

// WidgetStaticSize has inherent size and width known before painting.
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

// Computes the maximum frame count of a slice of widgets.
func MaxFrameCount(widgets []Widget) int {
	m := 1

	for _, w := range widgets {
		if c := w.FrameCount(); c > m {
			m = c
		}
	}

	return m
}

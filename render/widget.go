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

// Returns the smaller of x or y.
//
// Exists to avoid having to cast to 'float64' to use 'math.Min'.
func MinInt(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

// Returns the larger of x or y.
//
// Exists to avoid having to cast to 'float64' to use 'math.Max'.
func MaxInt(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

// Returns the absolute value of x.
//
// Exists to avoid having to cast to 'float64' to use 'math.Abs'.
func AbsInt(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

// Returns -1 for a negative x, +1 for a positive x and 0 if x is 0.
func SignInt(x int) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	} else {
		return 0
	}
}

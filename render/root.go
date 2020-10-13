package render

import (
	"github.com/fogleman/gg"
	"image"
	"image/color"
)

const (
	// DefaultFrameWidth is the normal width for a frame.
	DefaultFrameWidth = 64

	// DefaultFrameHeight is the normal height for a frame.
	DefaultFrameHeight = 32
)

// Every Widget tree has a Root.
//
// The child widget, and all its descendants, will be drawn on a 64x32
// canvas. Root places its child in the upper left corner of the
// canvas.
//
// If the tree contains animated widgets, the resulting animation will
// run with _delay_ milliseconds per frame.
type Root struct {
	Child Widget
	Delay int32
}

// Paint renders the child widget onto the frame. It doesn't do
// any resizing or alignment.
func (r Root) Paint(solidBackground bool) []image.Image {
	numFrames := r.Child.FrameCount()
	frames := make([]image.Image, numFrames)

	for i := 0; i < numFrames; i++ {
		dc := gg.NewContext(DefaultFrameWidth, DefaultFrameHeight)
		if solidBackground {
			dc.SetColor(color.Black)
			dc.Clear()
		}
		im := r.Child.Paint(image.Rect(0, 0, DefaultFrameWidth, DefaultFrameHeight), i)
		dc.DrawImage(im, 0, 0)
		frames[i] = dc.Image()
	}
	return frames
}

// PaintRoots draws >=1 Roots which must all have the same dimensions.
func PaintRoots(solidBackground bool, roots ...Root) []image.Image {
	var images []image.Image
	for _, r := range roots {
		images = append(images, r.Paint(solidBackground)...)
	}

	return images
}

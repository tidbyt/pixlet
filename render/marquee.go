package render

import (
	"github.com/fogleman/gg"
	"image"
)

// Marquee scrolls its child horizontally.
//
// The height of the Marquee will be that of its child, but its
// `width` must be specified explicitly. If the child's width fits
// fully, it will not scroll. Otherwise, it will be scrolled right to
// left.
//
// DOC(Child): Widget to potentially scroll
// DOC(Width): Width of the Marquee
//
// EXAMPLE BEGIN
// render.Marquee(
//      width=64,
//      child=render.Text("this won't fit in 64 pixels"),
// )
// EXAMPLE END
type Marquee struct {
	Widget
	Child Widget `starlark:"child,required"`
	Width int    `starlark:"width,required"`
}

func (m Marquee) FrameCount() int {
	im := m.Child.Paint(image.Rect(0, 0, m.Width*2, DefaultFrameHeight), 0)
	cw := im.Bounds().Dx()

	if cw <= m.Width {
		return 1
	}

	return m.Width*2 + cw
}

func (m Marquee) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	// We'll only scroll frame 0 of the child. Scrolling an
	// animation would be madness.
	im := m.Child.Paint(image.Rect(0, 0, m.Width*2, bounds.Dy()), 0)

	cw := im.Bounds().Dx()

	var offset int
	if cw <= m.Width {
		// child fits entirely and we don't want to scroll it anyway
		offset = 0
	} else if frameIdx <= cw+m.Width {
		// first, scroll the entire child in to and then out
		// of view
		offset = m.Width - frameIdx
	} else if frameIdx <= cw+m.Width*2 {
		// then, scroll back into view
		offset = m.Width - (frameIdx - cw - m.Width)
	} else {
		// if more than FrameCount frames are requested,
		// freeze marquee at final frame
		offset = 0
	}

	dc := gg.NewContext(m.Width, im.Bounds().Dy())
	dc.DrawImage(im, offset, 0)

	return dc.Image()
}

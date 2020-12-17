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
// The `offset_start` and `offset_end` parameters control the position
// of the child in the beginning and the end of the animation.
//
// DOC(Child): Widget to potentially scroll
// DOC(Width): Width of the Marquee
// DOC(OffsetStart): Position of child at beginning of animation
// DOC(OffsetEnd): Position of child at end of animation
//
// EXAMPLE BEGIN
// render.Marquee(
//      width=64,
//      child=render.Text("this won't fit in 64 pixels"),
//      offset_start=5,
//      offset_end=32,
// )
// EXAMPLE END
type Marquee struct {
	Widget
	Child       Widget `starlark:"child,required"`
	Width       int    `starlark:"width,required"`
	OffsetStart int    `starlark:"offset_start"`
	OffsetEnd   int    `starlark:"offset_end"`
}

func (m Marquee) FrameCount() int {
	im := m.Child.Paint(image.Rect(0, 0, m.Width*2, DefaultFrameHeight), 0)
	cw := im.Bounds().Dx()

	if cw <= m.Width {
		return 1
	}

	offstart := m.OffsetStart
	if offstart < -cw {
		offstart = -cw
	}

	offend := m.OffsetEnd
	if offend < -cw {
		offend = -cw
	}

	return cw + offstart + m.Width - offend + 1
}

func (m Marquee) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	// We'll only scroll frame 0 of the child. Scrolling an
	// animation would be madness.
	im := m.Child.Paint(image.Rect(0, 0, m.Width*2, bounds.Dy()), 0)

	cw := im.Bounds().Dx()

	offstart := m.OffsetStart
	if offstart < -cw {
		offstart = -cw
	}

	offend := m.OffsetEnd
	if offend < -cw {
		offend = -cw
	}

	loopIdx := cw + offstart
	endIdx := cw + offstart + m.Width - offend

	var offset int
	if cw <= m.Width {
		// child fits entirely and we don't want to scroll it anyway
		offset = 0
	} else if frameIdx <= loopIdx {
		// first scroll child out of view
		offset = offstart - frameIdx
	} else if frameIdx <= endIdx {
		// then, scroll back into view
		offset = offend + (endIdx - frameIdx)
	} else {
		// if more than FrameCount frames are requested,
		// freeze marquee at final frame
		offset = offend
	}

	dc := gg.NewContext(m.Width, im.Bounds().Dy())
	dc.DrawImage(im, offset, 0)

	return dc.Image()
}

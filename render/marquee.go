package render

import (
	"github.com/fogleman/gg"
	"image"
)

// Marquee scrolls its child horizontally or vertically.
//
// The `scroll_direction` will be 'horizontal' and will scroll from right
// to left if left empty, if specified as 'vertical' the Marquee will
// scroll from bottom to top.
//
// The `behavior` will be 'scroll' and will scroll the child out of view
// and back into the view if left empty, if specified as 'alternate' the
// Marquee will scroll the child out of the view until it reaches its end,
// then switch directions.
//
// In horizontal mode the height of the Marquee will be that of its child,
// but its `width` must be specified explicitly. In vertical mode the width
// will be that of its child but the `height` must be specified explicitly.
//
// If the child's width fits fully, it will only scroll if `scroll_always`
// is set.
//
// The `offset_start` and `offset_end` parameters control the position
// of the child in the beginning and the end of the animation.
// This is only supported if `behavior` is set to `scroll`.
//
// The `pause_start` and `pause_midway` parameters control the number of
// frames to pause at the beginning and at the midway point of the animation.
// This is only supported if `behavior` is set to `alternate`.
//
// DOC(Child): Widget to potentially scroll
// DOC(Width): Width of the Marquee, required for horizontal
// DOC(Height): Height of the Marquee, required for vertical
// DOC(Behavior): Behavior of scrolling, 'scroll' or 'alternate', default is scroll
// DOC(OffsetStart): Position of child at beginning of animation
// DOC(OffsetEnd): Position of child at end of animation
// DOC(PauseStart): Pause at beginning of animation, default and minumum is 1
// DOC(PauseMidway): Pause at midway point of animation, default and minumum is 1
// DOC(ScrollDirection): Direction to scroll, 'vertical' or 'horizontal', default is horizontal
// DOC(ScrollAlways): Scroll child, even if it fits entirely
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
	Child           Widget `starlark:"child,required"`
	Width           int    `starlark:"width"`
	Height          int    `starlark:"height"`
	OffsetStart     int    `starlark:"offset_start"`
	OffsetEnd       int    `starlark:"offset_end"`
	Behavior        string `starlark:"behavior"`
	ScrollDirection string `starlark:"scroll_direction"`
	ScrollAlways    bool   `starlark:"scroll_always"`
	PauseStart      int    `starlark:"pause_start"`
	PauseMidway     int    `starlark:"pause_midway"`
}

// avoid having to cast to float64 to use 'math.Max'
func MaxInt(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func (m Marquee) FrameCount() int {
	var im image.Image
	var cw int
	var size int
	if m.isVertical() {
		im = m.Child.Paint(image.Rect(0, 0, DefaultFrameWidth, m.Height*10), 0)
		cw = im.Bounds().Dy()
		size = m.Height
	} else {
		im = m.Child.Paint(image.Rect(0, 0, m.Width*10, DefaultFrameHeight), 0)
		cw = im.Bounds().Dx()
		size = m.Width
	}

	if cw <= size && !m.ScrollAlways {
		return 1
	}

	if m.isAlternating() {
		return m.AlternatingFrameCount(cw, size)
	} else {
		return m.ScrollingFrameCount(cw, size)
	}
}

func (m Marquee) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	var im image.Image
	var cw int
	var size int
	if m.isVertical() {
		// We'll only scroll frame 0 of the child. Scrolling an
		// animation would be madness.
		im = m.Child.Paint(image.Rect(0, 0, bounds.Dx(), m.Height*10), 0)
		cw = im.Bounds().Dy()
		size = m.Height
	} else {
		im = m.Child.Paint(image.Rect(0, 0, m.Width*10, bounds.Dy()), 0)
		cw = im.Bounds().Dx()
		size = m.Width
	}

	offstart := m.OffsetStart
	if offstart < -cw {
		offstart = -cw
	}

	offend := m.OffsetEnd
	if offend < -cw {
		offend = -cw
	}

	var offset int
	if m.isAlternating() {
		offset = m.AlternatingOffset(cw, size, frameIdx)
	} else {
		offset = m.ScrollingOffset(cw, size, frameIdx)
	}

	var dc *gg.Context
	if m.isVertical() {
		dc = gg.NewContext(im.Bounds().Dx(), m.Height)
		dc.DrawImage(im, 0, offset)
	} else {
		dc = gg.NewContext(m.Width, im.Bounds().Dy())
		dc.DrawImage(im, offset, 0)
	}

	return dc.Image()
}

func (m Marquee) ScrollingFrameCount(cw int, size int) int {
	offstart := m.OffsetStart
	if offstart < -cw {
		offstart = -cw
	}

	offend := m.OffsetEnd
	if offend < -cw {
		offend = -cw
	}

	return cw + offstart + size - offend + 1
}

func (m Marquee) ScrollingOffset(cw int, size int, frameIdx int) int {
	offstart := m.OffsetStart
	if offstart < -cw {
		offstart = -cw
	}

	offend := m.OffsetEnd
	if offend < -cw {
		offend = -cw
	}

	loopIdx := cw + offstart
	endIdx := cw + offstart + size - offend

	if cw <= size && !m.ScrollAlways {
		// child fits entirely and we don't want to scroll it anyway
		return 0
	} else if frameIdx <= loopIdx {
		// first scroll child out of view
		return offstart - frameIdx
	} else if frameIdx <= endIdx {
		// then, scroll back into view
		return offend + (endIdx - frameIdx)
	} else {
		// if more than FrameCount frames are requested,
		// freeze marquee at final frame
		return offend
	}
}

func (m Marquee) AlternatingFrameCount(cw int, size int) int {
	if cw == size {
		return 1
	}

	// calculate excess size (too much or too little space)
	var diff int
	if cw < size {
		diff = size - cw
	} else {
		diff = cw - size
	}

	// pause values need to be at least 1
	pauseStart := MaxInt(1, m.PauseStart)
	pauseMidway := MaxInt(1, m.PauseMidway)

	// subtract one from diff, as the start and midway points
	// are accounted for by the pause values
	return pauseStart + (diff - 1) + pauseMidway + (diff - 1)
}

func (m Marquee) AlternatingOffset(cw int, size int, frameIdx int) int {
	if cw == size {
		return 0
	}

	var diff int
	var dir int
	if cw < size {
		diff = size - cw
		dir = 1
	} else {
		diff = cw - size
		dir = -1
	}

	// pause values need to be at least 1
	pauseStart := MaxInt(1, m.PauseStart)
	pauseMidway := MaxInt(1, m.PauseMidway)

	// subtract 1 to adjust for index starting at 0
	startIdx := pauseStart - 1
	midwayIdx := startIdx + diff
	continueIdx := midwayIdx + pauseMidway
	endIdx := continueIdx + diff

	if cw <= size && !m.ScrollAlways {
		// child fits entirely and we don't want to scroll it anyway
		return 0
	} else if frameIdx < startIdx {
		// paused before starting animation
		return 0
	} else if frameIdx < midwayIdx {
		// scroll child into one direction
		return dir * (frameIdx - startIdx)
	} else if frameIdx < continueIdx {
		// paused before changing direction
		return dir * (midwayIdx - startIdx)
	} else if frameIdx < endIdx {
		// scroll child into the other direction
		return dir * ((diff - 1) - (frameIdx - continueIdx))
	} else {
		// if more than FrameCount frames are requested,
		// freeze marquee at final frame
		return 0
	}
}

func (m Marquee) isAlternating() bool {
	return m.Behavior == "alternate"
}

func (m Marquee) isVertical() bool {
	return m.ScrollDirection == "vertical"
}

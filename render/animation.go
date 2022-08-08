package render

import (
	"image"

	"github.com/tidbyt/gg"
)

// Animations turns a list of children into an animation, where each
// child is a frame.
//
// FIXME: Behaviour when children themselves are animated is a bit
// weird. Think and fix.
//
// DOC(Children): Children to use as frames in the animation
//
// EXAMPLE BEGIN
// render.Animation(
//      children=[
//           render.Box(width=10, height=10, color="#300"),
//           render.Box(width=12, height=12, color="#500"),
//           render.Box(width=14, height=14, color="#700"),
//           render.Box(width=16, height=16, color="#900"),
//           render.Box(width=18, height=18, color="#b00"),
//      ],
// )
// EXAMPLE END
type Animation struct {
	Widget
	Children []Widget
}

func (a Animation) FrameCount() int {
	return len(a.Children)
}

func (a Animation) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	if frameIdx > len(a.Children) {
		frameIdx %= len(a.Children)
		if frameIdx < 0 {
			frameIdx += len(a.Children)
		}
	}

	return a.Children[ModInt(frameIdx, len(a.Children))].PaintBounds(bounds, frameIdx)
}

func (a Animation) Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int) {
	if frameIdx > len(a.Children) {
		frameIdx %= len(a.Children)
		if frameIdx < 0 {
			frameIdx += len(a.Children)
		}
	}

	a.Children[ModInt(frameIdx, len(a.Children))].Paint(dc, bounds, frameIdx)
}

package render

import (
	"image"

	"tidbyt.dev/pixlet/render/canvas"
)

// Stack draws its children on top of each other.
//
// Just like a stack of pancakes, except with Widgets instead of
// pancakes. The Stack will be given a width and height sufficient to
// fit all its children.
//
// DOC(Children): Widgets to stack
//
// EXAMPLE BEGIN
// render.Stack(
//
//	children=[
//	     render.Box(width=50, height=25, color="#911"),
//	     render.Text("hello there"),
//	     render.Box(width=4, height=32, color="#119"),
//	],
//
// )
// EXAMPLE END
type Stack struct {
	Widget
	Children []Widget `starlark:"children,required"`
}

func (s Stack) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	width, height := 0, 0
	for _, child := range s.Children {
		cb := child.PaintBounds(bounds, frameIdx)
		imW, imH := cb.Dx(), cb.Dy()
		if imW > width {
			width = imW
		}
		if imH > height {
			height = imH
		}
	}

	if width > bounds.Dx() {
		width = bounds.Dx()
	}
	if height > bounds.Dy() {
		height = bounds.Dy()
	}

	return image.Rect(0, 0, width, height)
}

func (s Stack) Paint(dc canvas.Canvas, bounds image.Rectangle, frameIdx int) {
	for _, child := range s.Children {
		dc.Push()
		child.Paint(dc, bounds, frameIdx)
		dc.Pop()
	}
}

func (s Stack) FrameCount() int {
	return MaxFrameCount(s.Children)
}

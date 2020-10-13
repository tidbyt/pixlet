package render

import (
	"image"

	"github.com/fogleman/gg"
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
//      children=[
//           render.Box(width=50, height=25, color="#911"),
//           render.Text("hello there"),
//           render.Box(width=4, height=32, color="#119"),
//      ],
// )
// EXAMPLE END
type Stack struct {
	Widget
	Children []Widget `starlark:"children,required"`
}

func (s Stack) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	width, height := 0, 0
	images := make([]image.Image, 0, len(s.Children))
	for _, child := range s.Children {
		im := child.Paint(bounds, frameIdx)
		imW, imH := im.Bounds().Dx(), im.Bounds().Dy()
		if imW > width {
			width = imW
		}
		if imH > height {
			height = imH
		}
		images = append(images, im)
	}

	if width > bounds.Dx() {
		width = bounds.Dx()
	}
	if height > bounds.Dy() {
		height = bounds.Dy()
	}

	dc := gg.NewContext(width, height)
	for _, im := range images {
		dc.DrawImage(im, 0, 0)
	}
	return dc.Image()
}

func (s Stack) FrameCount() int {
	n := 1
	for _, child := range s.Children {
		fc := child.FrameCount()
		if fc > n {
			n = fc
		}
	}
	return n
}

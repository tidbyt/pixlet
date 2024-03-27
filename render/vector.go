package render

import (
	"image"

	"tidbyt.dev/pixlet/render/canvas"
)

// A vector draws its children either vertically or horizontally (like
// a row or a column).
//
// A vector has a main axis along which children are draw. The main
// axis is either horizontal or vertical (i.e. a row or a
// column). MainAlign controls how children are placed along this
// axis. CrossAlign controls placement orthogonally to the main axis.
type Vector struct {
	Widget

	Children   []Widget
	MainAlign  string `starlark: "main_align"`
	CrossAlign string `starlark: "cross_align"`
	Expanded   bool
	Vertical   bool
}

func (v Vector) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	// (dx, dy) determines the orientation of this Vector
	dx, dy := 1, 0
	if v.Vertical {
		dx, dy = 0, 1
	}

	boundsW := bounds.Dx()
	boundsH := bounds.Dy()

	// Paint as many children as we can fit. Compute their max and
	// total width and height as we go along.
	maxW, maxH := 0, 0
	sumW, sumH := 0, 0
	for _, child := range v.Children {
		cb := child.PaintBounds(image.Rect(0, 0, boundsW-dx*sumW, boundsH-dy*sumH), frameIdx)

		imW := cb.Dx()
		imH := cb.Dy()

		sumW += imW
		if imW > maxW {
			maxW = imW
		}
		sumH += imH
		if imH > maxH {
			maxH = imH
		}

		// This checks if we've overflowed the main axis
		if sumW*dx >= boundsW || sumH*dy >= boundsH {
			break
		}
	}

	// Compute the final dimensions of the vector. If the vector
	// is expanded, then it will span the full bounds along the
	// main axis. Otherwise, it will be the size of its children.
	// Along the cross axis, size will be the max of the
	// children. However, in both cases, total size can never
	// exceed the available bounds.
	width := dx*sumW + dy*maxW
	height := dx*maxH + dy*sumH
	if v.Expanded {
		width = dx*boundsW + dy*maxW
		height = dx*maxH + dy*boundsH
	}
	if height > boundsH {
		height = boundsH
	}
	if width > boundsW {
		width = boundsW
	}

	return image.Rect(0, 0, width, height)
}

func (v Vector) Paint(dc canvas.Canvas, bounds image.Rectangle, frameIdx int) {
	// (dx, dy) determines the orientation of this Vector
	dx, dy := 1, 0
	if v.Vertical {
		dx, dy = 0, 1
	}

	boundsW := bounds.Dx()
	boundsH := bounds.Dy()

	// Paint as many children as we can fit. Compute their max and
	// total width and height as we go along.
	maxW, maxH := 0, 0
	sumW, sumH := 0, 0
	childrenBounds := make([]image.Rectangle, 0, len(v.Children))
	for _, child := range v.Children {
		cb := child.PaintBounds(image.Rect(0, 0, boundsW-dx*sumW, boundsH-dy*sumH), frameIdx)

		imW := cb.Dx()
		imH := cb.Dy()

		sumW += imW
		if imW > maxW {
			maxW = imW
		}
		sumH += imH
		if imH > maxH {
			maxH = imH
		}

		childrenBounds = append(childrenBounds, cb)

		// This checks if we've overflowed the main axis
		if sumW*dx >= boundsW || sumH*dy >= boundsH {
			break
		}
	}

	// Compute the final dimensions of the vector. If the vector
	// is expanded, then it will span the full bounds along the
	// main axis. Otherwise, it will be the size of its children.
	// Along the cross axis, size will be the max of the
	// children. However, in both cases, total size can never
	// exceed the available bounds.
	width := dx*sumW + dy*maxW
	height := dx*maxH + dy*sumH
	if v.Expanded {
		width = dx*boundsW + dy*maxW
		height = dx*maxH + dy*boundsH
	}
	if height > boundsH {
		height = boundsH
	}
	if width > boundsW {
		width = boundsW
	}

	// These control position and spacing across main axis
	offset := 0
	spacing := 0
	spacingResidual := 0

	// The amount of space we have to play with
	remaining := (dx*(width-sumW) + dy*(height-sumH))
	if remaining < 0 {
		remaining = 0
	}

	switch v.MainAlign {
	case "start":
		// all = 0
	case "end":
		offset = dx*(width-sumW) + dy*(height-sumH)
		if offset < 0 {
			offset = 0
		}
	case "space_evenly":
		spacing = remaining / (len(childrenBounds) + 1)
		spacingResidual = remaining % (len(childrenBounds) + 1)
		offset = spacing
	case "space_around":
		spacing = remaining / len(childrenBounds)
		spacingResidual = remaining % len(childrenBounds)
		offset = spacing / 2
	case "center":
		offset = remaining / 2
	case "space_between":
		n := len(childrenBounds)
		if n > 1 {
			spacing = remaining / (n - 1)
			spacingResidual = remaining % (n - 1)
			if spacingResidual > 0 {
				offset = -1
				spacingResidual += 1
			}
		}
	}

	maxW, maxH = 0, 0
	sumW, sumH = 0, 0

	// Draw the children
	for i, cb := range childrenBounds {
		imW := cb.Dx()
		imH := cb.Dy()
		child := v.Children[i]

		// Residual space gets distributed 1 pixel at a time
		if spacingResidual > 0 {
			offset += 1
			spacingResidual -= 1
		}

		// Cross axis position depends on cross axis alignment
		crossOffset := 0
		switch v.CrossAlign {
		case "start":
			// crossOffset = 0
		case "center":
			crossOffset = (dx*(height-imH) + dy*(width-imW)) / 2
		case "end":
			crossOffset = dx*(height-imH) + dy*(width-imW)
		}

		dc.Push()
		dc.Translate(float64(dx*offset+dy*crossOffset), float64(dx*crossOffset+dy*offset))

		dc.ClipRectangle(
			float64(0),
			float64(0),
			float64(cb.Dx()),
			float64(cb.Dy()),
		)

		child.Paint(dc, image.Rect(0, 0, boundsW-dx*sumW, boundsH-dy*sumH), frameIdx)
		dc.Pop()

		sumW += imW
		if imW > maxW {
			maxW = imW
		}
		sumH += imH
		if imH > maxH {
			maxH = imH
		}

		offset += dx*imW + dy*imH + spacing

		if offset >= dx*boundsW+dy*boundsH {
			break
		}
	}
}

func (v Vector) FrameCount() int {
	return MaxFrameCount(v.Children)
}

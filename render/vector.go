package render

import (
	"github.com/fogleman/gg"
	"image"
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

func (v Vector) Paint(bounds image.Rectangle, frameIdx int) image.Image {
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
	images := make([]image.Image, 0, len(v.Children))
	for _, child := range v.Children {
		im := child.Paint(image.Rect(0, 0, boundsW-dx*sumW, boundsH-dy*sumH), frameIdx)

		imW := im.Bounds().Dx()
		imH := im.Bounds().Dy()

		sumW += imW
		if imW > maxW {
			maxW = imW
		}
		sumH += imH
		if imH > maxH {
			maxH = imH
		}

		images = append(images, im)

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
		spacing = remaining / (len(images) + 1)
		spacingResidual = remaining % (len(images) + 1)
		offset = spacing
	case "space_around":
		spacing = remaining / len(images)
		spacingResidual = remaining % len(images)
		offset = spacing / 2
	case "center":
		offset = remaining / 2
	case "space_between":
		n := len(images)
		if n > 1 {
			spacing = remaining / (n - 1)
			spacingResidual = remaining % (n - 1)
			if spacingResidual > 0 {
				offset = -1
				spacingResidual += 1
			}
		}
	}

	// Draw the children
	dc := gg.NewContext(width, height)
	for _, im := range images {
		imW := im.Bounds().Dx()
		imH := im.Bounds().Dy()

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

		dc.DrawImage(im, dx*offset+dy*crossOffset, dx*crossOffset+dy*offset)
		offset += dx*imW + dy*imH + spacing

		if offset >= dx*boundsW+dy*boundsH {
			break
		}
	}

	return dc.Image()
}

func (v Vector) FrameCount() int {
	n := 1
	for _, child := range v.Children {
		cn := child.FrameCount()
		if cn > n {
			n = cn
		}
	}
	return n
}

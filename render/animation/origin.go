package animation

import (
	"image"
)

// An anchor point to use for scaling and rotation transforms.
//
// DOC(X): Horizontal anchor point
// DOC(Y): Vertical anchor point
//
type Origin struct {
	X NumberOrPercentage `starlark:"x,required"`
	Y NumberOrPercentage `starlark:"y,required"`
}

func (self Origin) Transform(bounds image.Rectangle) Vec2f {
	return Vec2f{
		self.X.Transform(bounds.Dx()),
		self.Y.Transform(bounds.Dy()),
	}
}

var DefaultOrigin = Origin{
	X: Percentage{0.5},
	Y: Percentage{0.5},
}

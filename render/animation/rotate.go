package animation

import (
	"github.com/fogleman/gg"
)

// Transform by rotating by a given angle in degrees.
//
// DOC(Angle): Angle to rotate by in degrees
//
type Rotate struct {
	Angle float64 `starlark:"angle,required"`
}

func (self Rotate) Apply(ctx *gg.Context, origin Vec2f, rounding Rounding) {
	ctx.RotateAbout(gg.Radians(self.Angle), origin.X, origin.Y)
}

func (self Rotate) Interpolate(other Transform, progress float64) (result Transform, ok bool) {
	if other, ok := other.(Rotate); ok {
		return Rotate{Lerp(self.Angle, other.Angle, progress)}, true
	}

	return RotateDefault, false
}

var RotateDefault = Rotate{0.0}

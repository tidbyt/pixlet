package animation

import (
	"math"

	"tidbyt.dev/pixlet/render/canvas"
)

// Transform by rotating by a given angle in degrees.
//
// DOC(Angle): Angle to rotate by in degrees
type Rotate struct {
	Angle float64 `starlark:"angle,required"`
}

func (self Rotate) Apply(ctx canvas.Canvas, origin Vec2f, rounding Rounding) {
	ctx.Translate(origin.X, origin.Y)
	ctx.Rotate(self.Angle * math.Pi / 180)
	ctx.Translate(-origin.X, -origin.Y)
}

func (self Rotate) Interpolate(other Transform, progress float64) (result Transform, ok bool) {
	if other, ok := other.(Rotate); ok {
		return Rotate{Lerp(self.Angle, other.Angle, progress)}, true
	}

	return RotateDefault, false
}

var RotateDefault = Rotate{0.0}

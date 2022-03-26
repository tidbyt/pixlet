package animation

import (
	"github.com/fogleman/gg"
)

// Transform by scaling by a given factor.
//
// DOC(X): Horizontal scale factor
// DOC(Y): Vertical scale factor
//
type Scale struct {
	Vec2f
}

func (self Scale) Apply(ctx *gg.Context, origin Vec2f, rounding Rounding) {
	ctx.ScaleAbout(self.X, self.Y, origin.X, origin.Y)
}

func (self Scale) Interpolate(other Transform, progress float64) (result Transform, ok bool) {
	if other, ok := other.(Scale); ok {
		return Scale{self.Lerp(other.Vec2f, progress)}, true
	}

	return ScaleDefault, false
}

var ScaleDefault = Scale{Vec2f{1.0, 1.0}}

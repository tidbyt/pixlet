package animation

import (
	"github.com/tidbyt/gg"
)

// Transform by translating by a given offset.
//
// DOC(X): Horizontal offset
// DOC(Y): Vertical offset
//
type Translate struct {
	Vec2f
}

func (self Translate) Apply(ctx *gg.Context, origin Vec2f, rounding Rounding) {
	ctx.Translate(rounding.Apply(self.X), rounding.Apply(self.Y))
}

func (self Translate) Interpolate(other Transform, progress float64) (result Transform, ok bool) {
	if other, ok := other.(Translate); ok {
		return Translate{self.Lerp(other.Vec2f, progress)}, true
	}

	return TranslateDefault, false
}

var TranslateDefault = Translate{Vec2f{0.0, 0.0}}

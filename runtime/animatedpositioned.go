package runtime

import (
	"fmt"

	"github.com/mitchellh/hashstructure"
	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/render"
	"tidbyt.dev/pixlet/render/animation"
)

type AnimatedPositioned struct {
	Widget
	animation.AnimatedPositioned
}

func newAnimatedPositioned(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		child    starlark.Value
		xStart   starlark.Int
		xEnd     starlark.Int
		yStart   starlark.Int
		yEnd     starlark.Int
		duration starlark.Int
		curve    starlark.String
		delay    starlark.Int
		hold     starlark.Int
	)

	if err := starlark.UnpackArgs(
		"AnimatedPositioned",
		args, kwargs,
		"child", &child,
		"x_start?", &xStart,
		"x_end?", &xEnd,
		"y_start?", &yStart,
		"y_end?", &yEnd,
		"duration", &duration,
		"curve", &curve,
		"delay?", &delay,
		"hold?", &hold,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Row: %s", err)
	}

	childWidget, ok := child.(Widget)
	if !ok {
		return nil, fmt.Errorf("invalid type for child: %s (expected a Widget)", child.Type())
	}

	ap := AnimatedPositioned{}
	ap.Child = childWidget.AsRenderWidget()
	ap.XStart = int(xStart.BigInt().Int64())
	ap.XEnd = int(xEnd.BigInt().Int64())
	ap.YStart = int(yStart.BigInt().Int64())
	ap.YEnd = int(yEnd.BigInt().Int64())
	ap.Duration = int(duration.BigInt().Int64())
	ap.Delay = int(delay.BigInt().Int64())
	ap.Hold = int(hold.BigInt().Int64())

	switch curveName := curve.GoString(); curveName {
	case "linear":
		ap.Curve = animation.LinearCurve{}
	case "ease_in":
		ap.Curve = animation.EaseIn
	case "ease_out":
		ap.Curve = animation.EaseOut
	case "ease_in_out":
		ap.Curve = animation.EaseInOut
	default:
		return nil, fmt.Errorf("unknown curve %s", curveName)
	}

	return ap, nil
}

func (ap AnimatedPositioned) AsRenderWidget() render.Widget {
	return ap.AnimatedPositioned
}

func (ap AnimatedPositioned) AttrNames() []string {
	return []string{}
}

func (ap AnimatedPositioned) Attr(name string) (starlark.Value, error) {
	return nil, nil
}

func (ap AnimatedPositioned) String() string       { return "AnimatedPositioned(...)" }
func (ap AnimatedPositioned) Type() string         { return "AnimatedPositioned" }
func (ap AnimatedPositioned) Freeze()              {}
func (ap AnimatedPositioned) Truth() starlark.Bool { return true }

func (ap AnimatedPositioned) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(ap, nil)
	return uint32(sum), err
}

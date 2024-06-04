package animation_runtime

// Code generated by runtime/gen. DO NOT EDIT.

import (
	"fmt"
	"sync"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"tidbyt.dev/pixlet/render"
	"tidbyt.dev/pixlet/render/animation"
	"tidbyt.dev/pixlet/runtime/modules/render_runtime"
)

type AnimationModule struct {
	once   sync.Once
	module starlark.StringDict
}

var animationModule = AnimationModule{}

func LoadAnimationModule() (starlark.StringDict, error) {
	animationModule.once.Do(func() {
		animationModule.module = starlark.StringDict{
			"animation": &starlarkstruct.Module{
				Name: "render",
				Members: starlark.StringDict{

					"AnimatedPositioned": starlark.NewBuiltin("AnimatedPositioned", newAnimatedPositioned),

					"Keyframe": starlark.NewBuiltin("Keyframe", newKeyframe),

					"Origin": starlark.NewBuiltin("Origin", newOrigin),

					"Rotate": starlark.NewBuiltin("Rotate", newRotate),

					"Scale": starlark.NewBuiltin("Scale", newScale),

					"Transformation": starlark.NewBuiltin("Transformation", newTransformation),

					"Translate": starlark.NewBuiltin("Translate", newTranslate),
				},
			},
		}
	})

	return animationModule.module, nil
}

type AnimatedPositioned struct {
	render_runtime.Widget

	animation.AnimatedPositioned

	starlarkChild starlark.Value

	starlarkCurve starlark.Value

	frame_count *starlark.Builtin
}

func newAnimatedPositioned(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		child    starlark.Value
		duration starlark.Int
		curve    starlark.Value
		x_start  starlark.Int
		x_end    starlark.Int
		y_start  starlark.Int
		y_end    starlark.Int
		delay    starlark.Int
		hold     starlark.Int
	)

	if err := starlark.UnpackArgs(
		"AnimatedPositioned",
		args, kwargs,
		"child", &child,
		"duration", &duration,
		"curve", &curve,
		"x_start?", &x_start,
		"x_end?", &x_end,
		"y_start?", &y_start,
		"y_end?", &y_end,
		"delay?", &delay,
		"hold?", &hold,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for AnimatedPositioned: %s", err)
	}

	w := &AnimatedPositioned{}

	if child != nil {
		childWidget, ok := child.(render_runtime.Widget)
		if !ok {
			return nil, fmt.Errorf(
				"invalid type for child: %s (expected Widget)",
				child.Type(),
			)
		}
		w.Child = childWidget.AsRenderWidget()
		w.starlarkChild = child
	}

	w.Duration = int(duration.BigInt().Int64())

	w.starlarkCurve = curve
	if curve == nil {
		w.Curve = animation.DefaultCurve
	} else if val, err := CurveFromStarlark(curve); err == nil {
		w.Curve = val
	} else {
		return nil, err
	}

	w.XStart = int(x_start.BigInt().Int64())

	w.XEnd = int(x_end.BigInt().Int64())

	w.YStart = int(y_start.BigInt().Int64())

	w.YEnd = int(y_end.BigInt().Int64())

	w.Delay = int(delay.BigInt().Int64())

	w.Hold = int(hold.BigInt().Int64())

	w.frame_count = starlark.NewBuiltin("frame_count", animatedpositionedFrameCount)

	return w, nil
}

func (w *AnimatedPositioned) AsRenderWidget() render.Widget {
	w.AnimatedPositioned.Type = "AnimatedPositioned"
	return &w.AnimatedPositioned
}

func (w *AnimatedPositioned) AttrNames() []string {
	return []string{
		"child", "duration", "curve", "x_start", "x_end", "y_start", "y_end", "delay", "hold",
	}
}

func (w *AnimatedPositioned) Attr(name string) (starlark.Value, error) {
	switch name {

	case "child":

		return w.starlarkChild, nil

	case "duration":

		return starlark.MakeInt(int(w.Duration)), nil

	case "curve":

		return w.starlarkCurve, nil

	case "x_start":

		return starlark.MakeInt(int(w.XStart)), nil

	case "x_end":

		return starlark.MakeInt(int(w.XEnd)), nil

	case "y_start":

		return starlark.MakeInt(int(w.YStart)), nil

	case "y_end":

		return starlark.MakeInt(int(w.YEnd)), nil

	case "delay":

		return starlark.MakeInt(int(w.Delay)), nil

	case "hold":

		return starlark.MakeInt(int(w.Hold)), nil

	case "frame_count":
		return w.frame_count.BindReceiver(w), nil

	default:
		return nil, nil
	}
}

func (w *AnimatedPositioned) String() string       { return "AnimatedPositioned(...)" }
func (w *AnimatedPositioned) Type() string         { return "AnimatedPositioned" }
func (w *AnimatedPositioned) Freeze()              {}
func (w *AnimatedPositioned) Truth() starlark.Bool { return true }

func (w *AnimatedPositioned) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

func animatedpositionedFrameCount(
	thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {

	w := b.Receiver().(*AnimatedPositioned)
	count := w.FrameCount()

	return starlark.MakeInt(count), nil
}

type Keyframe struct {
	animation.Keyframe

	starlarkPercentage starlark.Value

	starlarkTransforms *starlark.List

	starlarkCurve starlark.Value
}

func newKeyframe(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		percentage starlark.Value
		transforms *starlark.List
		curve      starlark.Value
	)

	if err := starlark.UnpackArgs(
		"Keyframe",
		args, kwargs,
		"percentage", &percentage,
		"transforms", &transforms,
		"curve?", &curve,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Keyframe: %s", err)
	}

	w := &Keyframe{}

	w.starlarkPercentage = percentage
	if val, err := PercentageFromStarlark(percentage); err == nil {
		w.Percentage = val
	} else {
		return nil, err
	}

	w.starlarkTransforms = transforms
	for i := 0; i < transforms.Len(); i++ {
		switch transformsVal := transforms.Index(i).(type) {
		case *Translate:
			w.Transforms = append(w.Transforms, transformsVal.Translate)
		case *Scale:
			w.Transforms = append(w.Transforms, transformsVal.Scale)
		case *Rotate:
			w.Transforms = append(w.Transforms, transformsVal.Rotate)
		default:
			return nil, fmt.Errorf("expected transform, but got '%s'", transformsVal.Type())
		}
	}

	w.starlarkCurve = curve
	if curve == nil {
		w.Curve = animation.DefaultCurve
	} else if val, err := CurveFromStarlark(curve); err == nil {
		w.Curve = val
	} else {
		return nil, err
	}

	return w, nil
}

func (w *Keyframe) AttrNames() []string {
	return []string{
		"percentage", "transforms", "curve",
	}
}

func (w *Keyframe) Attr(name string) (starlark.Value, error) {
	switch name {

	case "percentage":

		return w.starlarkPercentage, nil

	case "transforms":

		return w.starlarkTransforms, nil

	case "curve":

		return w.starlarkCurve, nil

	default:
		return nil, nil
	}
}

func (w *Keyframe) String() string       { return "Keyframe(...)" }
func (w *Keyframe) Type() string         { return "Keyframe" }
func (w *Keyframe) Freeze()              {}
func (w *Keyframe) Truth() starlark.Bool { return true }

func (w *Keyframe) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Origin struct {
	animation.Origin

	starlarkX starlark.Value

	starlarkY starlark.Value
}

func newOrigin(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		x starlark.Value
		y starlark.Value
	)

	if err := starlark.UnpackArgs(
		"Origin",
		args, kwargs,
		"x", &x,
		"y", &y,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Origin: %s", err)
	}

	w := &Origin{}

	w.starlarkX = x
	if val, err := PercentageFromStarlark(x); err == nil {
		w.X = val
	} else {
		return nil, err
	}

	w.starlarkY = y
	if val, err := PercentageFromStarlark(y); err == nil {
		w.Y = val
	} else {
		return nil, err
	}

	return w, nil
}

func (w *Origin) AttrNames() []string {
	return []string{
		"x", "y",
	}
}

func (w *Origin) Attr(name string) (starlark.Value, error) {
	switch name {

	case "x":

		return w.starlarkX, nil

	case "y":

		return w.starlarkY, nil

	default:
		return nil, nil
	}
}

func (w *Origin) String() string       { return "Origin(...)" }
func (w *Origin) Type() string         { return "Origin" }
func (w *Origin) Freeze()              {}
func (w *Origin) Truth() starlark.Bool { return true }

func (w *Origin) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Rotate struct {
	animation.Rotate

	starlarkAngle starlark.Value
}

func newRotate(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		angle starlark.Value
	)

	if err := starlark.UnpackArgs(
		"Rotate",
		args, kwargs,
		"angle", &angle,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Rotate: %s", err)
	}

	w := &Rotate{}

	w.starlarkAngle = angle
	if val, ok := starlark.AsFloat(w.starlarkAngle); ok {
		w.Angle = val
	} else {
		return nil, fmt.Errorf("expected number, but got: %s", w.starlarkAngle.String())
	}

	return w, nil
}

func (w *Rotate) AttrNames() []string {
	return []string{
		"angle",
	}
}

func (w *Rotate) Attr(name string) (starlark.Value, error) {
	switch name {

	case "angle":

		return w.starlarkAngle, nil

	default:
		return nil, nil
	}
}

func (w *Rotate) String() string       { return "Rotate(...)" }
func (w *Rotate) Type() string         { return "Rotate" }
func (w *Rotate) Freeze()              {}
func (w *Rotate) Truth() starlark.Bool { return true }

func (w *Rotate) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Scale struct {
	animation.Scale

	starlarkX starlark.Value

	starlarkY starlark.Value
}

func newScale(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		x starlark.Value
		y starlark.Value
	)

	if err := starlark.UnpackArgs(
		"Scale",
		args, kwargs,
		"x", &x,
		"y", &y,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Scale: %s", err)
	}

	w := &Scale{}

	w.starlarkX = x
	if val, ok := starlark.AsFloat(w.starlarkX); ok {
		w.X = val
	} else {
		return nil, fmt.Errorf("expected number, but got: %s", w.starlarkX.String())
	}

	w.starlarkY = y
	if val, ok := starlark.AsFloat(w.starlarkY); ok {
		w.Y = val
	} else {
		return nil, fmt.Errorf("expected number, but got: %s", w.starlarkY.String())
	}

	return w, nil
}

func (w *Scale) AttrNames() []string {
	return []string{
		"x", "y",
	}
}

func (w *Scale) Attr(name string) (starlark.Value, error) {
	switch name {

	case "x":

		return w.starlarkX, nil

	case "y":

		return w.starlarkY, nil

	default:
		return nil, nil
	}
}

func (w *Scale) String() string       { return "Scale(...)" }
func (w *Scale) Type() string         { return "Scale" }
func (w *Scale) Freeze()              {}
func (w *Scale) Truth() starlark.Bool { return true }

func (w *Scale) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Transformation struct {
	render_runtime.Widget

	animation.Transformation

	starlarkChild starlark.Value

	starlarkKeyframes *starlark.List

	starlarkOrigin starlark.Value

	starlarkDirection starlark.String

	starlarkFillMode starlark.String

	starlarkRounding starlark.String

	frame_count *starlark.Builtin
}

func newTransformation(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		child          starlark.Value
		keyframes      *starlark.List
		duration       starlark.Int
		delay          starlark.Int
		width          starlark.Int
		height         starlark.Int
		origin         starlark.Value
		direction      starlark.String
		fill_mode      starlark.String
		rounding       starlark.String
		wait_for_child starlark.Bool
	)

	if err := starlark.UnpackArgs(
		"Transformation",
		args, kwargs,
		"child", &child,
		"keyframes", &keyframes,
		"duration", &duration,
		"delay?", &delay,
		"width?", &width,
		"height?", &height,
		"origin?", &origin,
		"direction?", &direction,
		"fill_mode?", &fill_mode,
		"rounding?", &rounding,
		"wait_for_child?", &wait_for_child,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Transformation: %s", err)
	}

	w := &Transformation{}

	if child != nil {
		childWidget, ok := child.(render_runtime.Widget)
		if !ok {
			return nil, fmt.Errorf(
				"invalid type for child: %s (expected Widget)",
				child.Type(),
			)
		}
		w.Child = childWidget.AsRenderWidget()
		w.starlarkChild = child
	}

	w.starlarkKeyframes = keyframes
	for i := 0; i < keyframes.Len(); i++ {
		if val, ok := keyframes.Index(i).(*Keyframe); ok {
			w.Keyframes = append(w.Keyframes, val.Keyframe)
		} else {
			return nil, fmt.Errorf("invalid type for keyframes: %s (expected Keyframe)", keyframes.Type())
		}
	}

	w.Duration = int(duration.BigInt().Int64())

	w.Delay = int(delay.BigInt().Int64())

	w.Width = int(width.BigInt().Int64())

	w.Height = int(height.BigInt().Int64())

	w.starlarkOrigin = origin
	if origin == nil {
		w.Origin = animation.DefaultOrigin
	} else if val, ok := origin.(*Origin); ok {
		w.Origin = val.Origin
	} else {
		return nil, fmt.Errorf("invalid type for origin: %s (expected Origin)", origin.Type())
	}

	w.starlarkDirection = direction
	switch direction {
	case "normal":
		w.Direction = animation.DirectionNormal
	case "reverse":
		w.Direction = animation.DirectionReverse
	case "alternate":
		w.Direction = animation.DirectionAlternate
	case "alternate-reverse":
		w.Direction = animation.DirectionAlternateReverse
	case "":
		w.Direction = animation.DefaultDirection
	default:
		return nil, fmt.Errorf("invalid type for direction: %s (expected 'normal', 'reverse', 'alternate' or 'alternate-reverse')", direction.Type())
	}

	w.starlarkFillMode = fill_mode
	switch fill_mode {
	case "forwards":
		w.FillMode = animation.FillModeForwards{}
	case "backwards":
		w.FillMode = animation.FillModeBackwards{}
	case "":
		w.FillMode = animation.DefaultFillMode
	default:
		return nil, fmt.Errorf("invalid type for fill_mode: %s (expected 'forwards' or 'backwards')", fill_mode.Type())
	}

	w.starlarkRounding = rounding
	switch rounding {
	case "round":
		w.Rounding = animation.Round{}
	case "floor":
		w.Rounding = animation.RoundFloor{}
	case "ceil":
		w.Rounding = animation.RoundCeil{}
	case "none":
		w.Rounding = animation.RoundNone{}
	case "":
		w.Rounding = animation.DefaultRounding
	default:
		return nil, fmt.Errorf("invalid type for rounding: %s (expected 'round', 'floor', 'ceil' or 'none')", rounding.Type())
	}

	w.WaitForChild = bool(wait_for_child)

	w.frame_count = starlark.NewBuiltin("frame_count", transformationFrameCount)

	if err := w.Init(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Transformation) AsRenderWidget() render.Widget {
	w.Transformation.Type = "Transformation"
	return &w.Transformation
}

func (w *Transformation) AttrNames() []string {
	return []string{
		"child", "keyframes", "duration", "delay", "width", "height", "origin", "direction", "fill_mode", "rounding", "wait_for_child",
	}
}

func (w *Transformation) Attr(name string) (starlark.Value, error) {
	switch name {

	case "child":

		return w.starlarkChild, nil

	case "keyframes":

		return w.starlarkKeyframes, nil

	case "duration":

		return starlark.MakeInt(int(w.Duration)), nil

	case "delay":

		return starlark.MakeInt(int(w.Delay)), nil

	case "width":

		return starlark.MakeInt(int(w.Width)), nil

	case "height":

		return starlark.MakeInt(int(w.Height)), nil

	case "origin":

		return w.starlarkOrigin, nil

	case "direction":

		return w.starlarkDirection, nil

	case "fill_mode":

		return w.starlarkFillMode, nil

	case "rounding":

		return w.starlarkRounding, nil

	case "wait_for_child":

		return starlark.Bool(w.WaitForChild), nil

	case "frame_count":
		return w.frame_count.BindReceiver(w), nil

	default:
		return nil, nil
	}
}

func (w *Transformation) String() string       { return "Transformation(...)" }
func (w *Transformation) Type() string         { return "Transformation" }
func (w *Transformation) Freeze()              {}
func (w *Transformation) Truth() starlark.Bool { return true }

func (w *Transformation) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

func transformationFrameCount(
	thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {

	w := b.Receiver().(*Transformation)
	count := w.FrameCount()

	return starlark.MakeInt(count), nil
}

type Translate struct {
	animation.Translate

	starlarkX starlark.Value

	starlarkY starlark.Value
}

func newTranslate(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		x starlark.Value
		y starlark.Value
	)

	if err := starlark.UnpackArgs(
		"Translate",
		args, kwargs,
		"x", &x,
		"y", &y,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Translate: %s", err)
	}

	w := &Translate{}

	w.starlarkX = x
	if val, ok := starlark.AsFloat(w.starlarkX); ok {
		w.X = val
	} else {
		return nil, fmt.Errorf("expected number, but got: %s", w.starlarkX.String())
	}

	w.starlarkY = y
	if val, ok := starlark.AsFloat(w.starlarkY); ok {
		w.Y = val
	} else {
		return nil, fmt.Errorf("expected number, but got: %s", w.starlarkY.String())
	}

	return w, nil
}

func (w *Translate) AttrNames() []string {
	return []string{
		"x", "y",
	}
}

func (w *Translate) Attr(name string) (starlark.Value, error) {
	switch name {

	case "x":

		return w.starlarkX, nil

	case "y":

		return w.starlarkY, nil

	default:
		return nil, nil
	}
}

func (w *Translate) String() string       { return "Translate(...)" }
func (w *Translate) Type() string         { return "Translate" }
func (w *Translate) Freeze()              {}
func (w *Translate) Truth() starlark.Bool { return true }

func (w *Translate) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

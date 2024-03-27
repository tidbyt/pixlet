package animation

import (
	"fmt"
	"image"
	"sort"

	"tidbyt.dev/pixlet/render"
	"tidbyt.dev/pixlet/render/canvas"
)

func makeKeyframe(p float64) Keyframe {
	return Keyframe{
		Percentage: Percentage{p},
		Transforms: make([]Transform, 0),
		Curve:      LinearCurve{},
	}
}

// Create default '0%' and '100%' keyframes.
var defaultKeyframeFrom = makeKeyframe(0.0)
var defaultKeyframeTo = makeKeyframe(1.0)
var defaultKeyframes = []Keyframe{defaultKeyframeFrom, defaultKeyframeTo}

// Sort keyframes and add 0% and 100% keyframes in case they are missing.
//
// See: https://www.w3.org/TR/css-animations-1/#keyframes
func processKeyframes(arr []Keyframe) []Keyframe {
	// If no keyframes are specified, return default keyframes.
	if len(arr) == 0 {
		return defaultKeyframes
	}

	// Sort keyframes by increasing order in time.
	sort.SliceStable(arr, func(i, j int) bool {
		return arr[i].Percentage.Value < arr[j].Percentage.Value
	})

	// Prepend default '0%' keyframe, if not specified.
	if arr[0].Percentage.Value != 0.0 {
		arr = append([]Keyframe{defaultKeyframeFrom}, arr...)
	}

	// Append default '100%' keyframe, if not specified.
	if arr[len(arr)-1].Percentage.Value != 1.0 {
		arr = append(arr, defaultKeyframeTo)
	}

	return arr
}

// Find the adjacent keyframes given a percentage value.
func findKeyframes(arr []Keyframe, p float64) (Keyframe, Keyframe, error) {
	if len(arr) < 2 {
		return defaultKeyframeFrom, defaultKeyframeTo,
			fmt.Errorf("invalid number of keyframes: %d (expected at least 2)", len(arr))
	}

	if p < 0.0 || p > 1.0 {
		return defaultKeyframeFrom, defaultKeyframeTo,
			fmt.Errorf("invalid range for percentage: %f (expected number in range [0.0, 1.0])", p)
	}

	for i := 0; i < len(arr)-1; i++ {
		p0 := arr[i].Percentage.Value
		p1 := arr[i+1].Percentage.Value

		if p0 <= p && p <= p1 {
			return arr[i], arr[i+1], nil
		}
	}

	return defaultKeyframeFrom, defaultKeyframeTo,
		fmt.Errorf("failed to find adjacent keyframes for percentage: %f (this should be unreachable!?)", p)
}

// Transformation makes it possible to animate a child widget by
// transitioning between transforms which are applied to the child wiget.
//
// It supports animating translation, scale and rotation of its child.
//
// If you have used CSS transforms and animations before, some of the
// following concepts will be familiar to you.
//
// Keyframes define a list of transforms to apply at a specific point in
// time, which is given as a percentage of the total animation duration.
//
// A keyframe is created via `animation.Keyframe(percentage, transforms, curve)`.
//
// The `percentage` specifies its point in time and can be expressed as
// a floating point number in the range `0.0` to `1.0`.
//
// In case a keyframe at percentage 0% or 100% is missing, a default
// keyframe without transforms and with a "linear" easing curve is inserted.
//
// As the animation progresses, transforms defined by the previous and
// next keyframe will be interpolated to determine the transform to apply
// at the current frame.
//
// The `duration` and `delay` of the animation are expressed as a number
// of frames.
//
// By default a transform `origin` of `animation.Origin(0.5, 0.5)` is used,
// which defines the anchor point for scaling and rotation to be exactly the
// center of the child widget. A different `origin` can be specified by
// providing a custom `animation.Origin`.
//
// The animation `direction` defaults to `normal`, playing the animation
// forwards. Other possible values are `reverse` to play it backwards,
// `alternate` to play it forwards, then backwards or `alternate-reverse`
// to play it backwards, then forwards.
//
// The animation `fill_mode` defaults to `forwards`, and controls which
// transforms will be applied to the child widget after the animation
// finishes. A value of `forwards` will retain the transforms of the last
// keyframe, while a value of `backwards` will rever to the transforms
// of the first keyframe.
//
// When translating the child widget on the X- or Y-axis, it often is
// desireable to round to even integers, which can be controlled via
// `rounding`, which defaults to `round`. Possible values are `round` to
// round to the nearest integer, `floor` to round down, `ceil` to round
// up or `none` to not perform any rounding. Rounding only is applied for
// translation transforms, but not to scaling or rotation transforms.
//
// If `wait_for_child` is set to `True`, the animation will finish and
// then wait for all child frames to play before restarting. If it is set
// to `False`, it will not wait.
//
// DOC(Child): Widget to animate
// DOC(Keyframes): List of animation keyframes
// DOC(Duration): Duration of animation (in frames)
// DOC(Delay): Duration to wait before animation (in frames)
// DOC(Width): Width of the animation canvas
// DOC(Height): Height of the animation canvas
// DOC(Origin): Origin for transforms, default is '50%, 50%'
// DOC(Direction): Direction of the animation, default is 'normal'
// DOC(FillMode): Fill mode of the animation, default is 'forwards'
// DOC(Rounding): Rounding to use for interpolated translation coordinates (not used for scale and rotate), default is 'round'
// DOC(WaitForChild): Wait for all child frames to play after finishing
//
// EXAMPLE BEGIN
// animation.Transformation(
//
//	child = render.Box(render.Circle(diameter = 6, color = "#0f0")),
//	duration = 100,
//	delay = 0,
//	origin = animation.Origin(0.5, 0.5),
//	direction = "alternate",
//	fill_mode = "forwards",
//	keyframes = [
//	  animation.Keyframe(
//	    percentage = 0.0,
//	    transforms = [animation.Rotate(0), animation.Translate(-10, 0), animation.Rotate(0)],
//	    curve = "ease_in_out",
//	  ),
//	  animation.Keyframe(
//	    percentage = 1.0,
//	    transforms = [animation.Rotate(360), animation.Translate(-10, 0), animation.Rotate(-360)],
//	  ),
//	],
//
// ),
// EXAMPLE END
type Transformation struct {
	render.Widget

	Child        render.Widget `starlark:"child,required"`
	Keyframes    []Keyframe    `starlark:"keyframes,required"`
	Duration     int           `starlark:"duration,required"`
	Delay        int           `starlark:"delay"`
	Width        int           `starlark:"width"`
	Height       int           `starlark:"height"`
	Origin       Origin        `starlark:"origin"`
	Direction    Direction     `starlark:"direction"`
	FillMode     FillMode      `starlark:"fill_mode"`
	Rounding     Rounding      `starlark:"rounding"`
	WaitForChild bool          `starlark:"wait_for_child"`
}

func (self *Transformation) Init() error {
	self.Keyframes = processKeyframes(self.Keyframes)

	return nil
}

func (self *Transformation) FrameCount() int {
	fc := self.Direction.FrameCount(self.Delay, self.Duration)
	cfc := self.Child.FrameCount()

	if self.WaitForChild && cfc > fc {
		return cfc
	}

	return fc
}

func (self *Transformation) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	w, h := self.Width, self.Height

	if w == 0 {
		w = bounds.Dx()
	}

	if h == 0 {
		h = bounds.Dy()
	}
	return image.Rect(0, 0, w, h)
}

func (self *Transformation) Paint(dc canvas.Canvas, bounds image.Rectangle, frameIdx int) {
	bounds = self.PaintBounds(bounds, frameIdx)
	cb := self.Child.PaintBounds(bounds, frameIdx)

	// As the origin might have been specified in relative units,
	// transform it given the child widget bounds.
	origin := self.Origin.Transform(cb)

	// Calculate the overall animation progress.
	progress := self.Direction.Progress(
		self.Delay,
		self.Duration,
		self.FillMode.Value(),
		frameIdx,
	)

	dc.Push()

	// Find the adjacent keyframes to interpolate between.
	if from, to, err := findKeyframes(self.Keyframes, progress); err == nil {
		// Rescale animation progress to progress between keyframes and apply easing curve.
		progress = Rescale(from.Percentage.Value, to.Percentage.Value, 0.0, 1.0, progress)
		progress = from.Curve.Transform(progress)

		// Interpolate between transforms and apply them in order.
		if transforms, ok := InterpolateTransforms(from.Transforms, to.Transforms, progress); ok {
			for _, transform := range transforms {
				transform.Apply(dc, origin, self.Rounding)
			}
		}
	}

	self.Child.Paint(dc, bounds, frameIdx)

	dc.Pop()
}

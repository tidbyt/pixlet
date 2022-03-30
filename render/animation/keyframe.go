package animation

// A keyframe defining specific point in time in the animation.
//
// The keyframe _percentage_ can is expressed as a floating point value between `0.0` and `1.0`.
//
// DOC(Percentage): Percentage of the time at which this keyframe occurs through the animation.
// DOC(Transforms): List of transforms at this keyframe to interpolate to or from.
// DOC(Curve): Easing curve to use, default is 'linear'
//
type Keyframe struct {
	Percentage Percentage  `starlark:"percentage,required"`
	Transforms []Transform `starlark:"transforms,required"`
	Curve      Curve       `starlark:"curve"`
}

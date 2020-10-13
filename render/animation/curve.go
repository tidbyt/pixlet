package animation

import "math"

var EaseIn = CubicBezierCurve{0.3, 0, 1, 1}
var EaseOut = CubicBezierCurve{0, 0, 0, 1}

// TODO: figure out if what curve to use here. unless we're going back
// to Ivo's curve (0.3, 0, 0, 1), make sure to update the unit tests
//
// var EaseInOut = CubicBezierCurve{0.3, 0, 0, 1}
var EaseInOut = CubicBezierCurve{0.65, 0, 0.35, 1}

type Curve interface {
	Transform(t float64) float64
}

// Linear curve moving from 0 to 1 (wait for it...) linearly
type LinearCurve struct{}

func (lc LinearCurve) Transform(t float64) float64 {
	return t
}

// Bezier curve defined by a, b, c and d.
type CubicBezierCurve struct {
	a, b, c, d float64
}

func (cb CubicBezierCurve) Transform(t float64) float64 {
	start, end := 0.0, 1.0
	epsilon := 0.0001

	for {
		mid := start + (end-start)/2
		x := cb.computeBezier(mid, cb.a, cb.c)
		if math.Abs(x-t) < epsilon {
			return cb.computeBezier(mid, cb.b, cb.d)
		}
		if x < t {
			start = mid
		} else {
			end = mid
		}
	}

	return math.NaN()
}

func (cb CubicBezierCurve) computeBezier(t, e, f float64) float64 {
	return 3*e*(1-t)*(1-t)*t + 3*f*(1-t)*t*t + t*t*t
}

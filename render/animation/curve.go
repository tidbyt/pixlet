package animation

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"go.starlark.net/starlark"
)

var EaseIn = CubicBezierCurve{0.3, 0, 1, 1}
var EaseOut = CubicBezierCurve{0, 0, 0, 1}

var DefaultCurve = LinearCurve{}

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

// Custom curve implemented as a starlark function
type CustomCurve struct {
	Function *starlark.Function
}

func (cc CustomCurve) Transform(t float64) float64 {
	r, err := starlark.Call(&starlark.Thread{}, cc.Function, starlark.Tuple{starlark.Float(t)}, nil)
	if err != nil {
		fmt.Printf("Error calling curve function %s: %s\n", cc.Function.String(), err.Error())
		return math.NaN()
	}

	f, ok := starlark.AsFloat(r)
	if !ok {
		fmt.Printf("Curve function did not return a floating point value!\n")
		return math.NaN()
	}

	return f
}

var cubicBezierRe = regexp.MustCompile(
	`^cubic-bezier\(` +
		`(?P<a>[+-]?([0-9]*\.)?[0-9]+), ` +
		`(?P<b>[+-]?([0-9]*\.)?[0-9]+), ` +
		`(?P<c>[+-]?([0-9]*\.)?[0-9]+), ` +
		`(?P<d>[+-]?([0-9]*\.)?[0-9]+)` +
		`\)$`)

func ParseCurve(str string) (Curve, error) {
	match := cubicBezierRe.FindStringSubmatch(str)
	if match != nil {
		result := make(map[string]string)

		for i, name := range cubicBezierRe.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		a, _ := strconv.ParseFloat(result["a"], 64)
		b, _ := strconv.ParseFloat(result["b"], 64)
		c, _ := strconv.ParseFloat(result["c"], 64)
		d, _ := strconv.ParseFloat(result["d"], 64)

		return CubicBezierCurve{a, b, c, d}, nil
	}

	switch str {
	case "linear":
		return LinearCurve{}, nil
	case "ease_in":
		return EaseIn, nil
	case "ease_out":
		return EaseOut, nil
	case "ease_in_out":
		return EaseInOut, nil
	default:
		return LinearCurve{}, fmt.Errorf("%s is not a valid curve string", str)
	}
}

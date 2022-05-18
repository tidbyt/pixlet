package animation

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/render"
)

func TestTransformationTranslate(t *testing.T) {
	o := Transformation{
		Child: render.Column{
			Children: []render.Widget{
				render.Box{Width: 3, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
				render.Box{Width: 3, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
				render.Box{Width: 3, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		Keyframes: []Keyframe{
			{
				Percentage: Percentage{0.0},
				Curve:      LinearCurve{},
				Transforms: []Transform{
					Translate{Vec2f{X: 0.0, Y: 0.0}},
				},
			},
			{
				Percentage: Percentage{1.0},
				Curve:      LinearCurve{},
				Transforms: []Transform{
					Translate{Vec2f{X: 5.0, Y: 5.0}},
				},
			},
		},
		Duration:     6,
		Delay:        0,
		Width:        5,
		Height:       5,
		Origin:       Origin{X: Percentage{0.5}, Y: Percentage{0.5}},
		Direction:    DefaultDirection,
		FillMode:     DefaultFillMode,
		Rounding:     DefaultRounding,
		WaitForChild: false,
	}

	// These frames should show the box moving diagonally out of frame.
	assert.Equal(t, 6, o.FrameCount())

	im := o.Paint(image.Rect(0, 0, 5, 5), 0)
	assert.Equal(t, nil, render.CheckImage([]string{
		"rrr..",
		"ggg..",
		"bbb..",
		".....",
		".....",
	}, im))

	im = o.Paint(image.Rect(0, 0, 5, 5), 1)
	assert.Equal(t, nil, render.CheckImage([]string{
		".....",
		".rrr.",
		".ggg.",
		".bbb.",
		".....",
	}, im))

	im = o.Paint(image.Rect(0, 0, 5, 5), 2)
	assert.Equal(t, nil, render.CheckImage([]string{
		".....",
		".....",
		"..rrr",
		"..ggg",
		"..bbb",
	}, im))

	im = o.Paint(image.Rect(0, 0, 5, 5), 3)
	assert.Equal(t, nil, render.CheckImage([]string{
		".....",
		".....",
		".....",
		"...rr",
		"...gg",
	}, im))

	im = o.Paint(image.Rect(0, 0, 5, 5), 4)
	assert.Equal(t, nil, render.CheckImage([]string{
		".....",
		".....",
		".....",
		".....",
		"....r",
	}, im))

	im = o.Paint(image.Rect(0, 0, 5, 5), 5)
	assert.Equal(t, nil, render.CheckImage([]string{
		".....",
		".....",
		".....",
		".....",
		".....",
	}, im))
}

func TestTransformationScale(t *testing.T) {
	o := Transformation{
		// Choosing only red, as scaling will interpolate between colors...
		Child: render.Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
		Keyframes: []Keyframe{
			{
				Percentage: Percentage{0.0},
				Curve:      LinearCurve{},
				Transforms: []Transform{
					Scale{Vec2f{X: 1.0, Y: 1.0}},
				},
			},
			{
				Percentage: Percentage{0.5},
				Curve:      LinearCurve{},
				Transforms: []Transform{
					Scale{Vec2f{X: 2.0, Y: 2.0}},
				},
			},
			{
				Percentage: Percentage{1.0},
				Curve:      LinearCurve{},
				Transforms: []Transform{
					Scale{Vec2f{X: 3.0, Y: 3.0}},
				},
			},
		},
		Duration:     3,
		Delay:        0,
		Width:        9,
		Height:       9,
		Origin:       Origin{X: Percentage{0.0}, Y: Percentage{0.0}},
		Direction:    DefaultDirection,
		FillMode:     DefaultFillMode,
		Rounding:     DefaultRounding,
		WaitForChild: false,
	}

	// These frames should show the box scaling from 1x to 3x.
	assert.Equal(t, 3, o.FrameCount())

	im := o.Paint(image.Rect(0, 0, 9, 9), 0)
	assert.Equal(t, nil, render.CheckImage([]string{
		"rrr......",
		"rrr......",
		"rrr......",
		".........",
		".........",
		".........",
		".........",
		".........",
		".........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 9, 9), 1)
	assert.Equal(t, nil, render.CheckImage([]string{
		"rrrrrr...",
		"rrrrrr...",
		"rrrrrr...",
		"rrrrrr...",
		"rrrrrr...",
		"rrrrrr...",
		".........",
		".........",
		".........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 9, 9), 2)
	assert.Equal(t, nil, render.CheckImage([]string{
		"rrrrrrrrr",
		"rrrrrrrrr",
		"rrrrrrrrr",
		"rrrrrrrrr",
		"rrrrrrrrr",
		"rrrrrrrrr",
		"rrrrrrrrr",
		"rrrrrrrrr",
		"rrrrrrrrr",
	}, im))

}

func TestTransformationRotate(t *testing.T) {
	o := Transformation{
		Child: render.Column{
			Children: []render.Widget{
				render.Box{Width: 3, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
				render.Box{Width: 3, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
				render.Box{Width: 3, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		Keyframes: []Keyframe{
			{
				Percentage: Percentage{0.0},
				Curve:      LinearCurve{},
				Transforms: []Transform{
					Rotate{Angle: 0.0},
				},
			},
			{
				Percentage: Percentage{1.0},
				Curve:      LinearCurve{},
				Transforms: []Transform{
					Rotate{Angle: 360.0},
				},
			},
		},
		Duration:     5,
		Delay:        0,
		Width:        3,
		Height:       3,
		Origin:       DefaultOrigin,
		Direction:    DefaultDirection,
		FillMode:     DefaultFillMode,
		Rounding:     DefaultRounding,
		WaitForChild: false,
	}

	// These frames should show the box rotating 90 degrees each frame.
	assert.Equal(t, 5, o.FrameCount())

	im := o.Paint(image.Rect(0, 0, 3, 3), 0)
	assert.Equal(t, nil, render.CheckImage([]string{
		"rrr",
		"ggg",
		"bbb",
	}, im))

	im = o.Paint(image.Rect(0, 0, 3, 3), 1)
	assert.Equal(t, nil, render.CheckImage([]string{
		"bgr",
		"bgr",
		"bgr",
	}, im))

	im = o.Paint(image.Rect(0, 0, 3, 3), 2)
	assert.Equal(t, nil, render.CheckImage([]string{
		"bbb",
		"ggg",
		"rrr",
	}, im))

	im = o.Paint(image.Rect(0, 0, 3, 3), 3)
	assert.Equal(t, nil, render.CheckImage([]string{
		"rgb",
		"rgb",
		"rgb",
	}, im))

	im = o.Paint(image.Rect(0, 0, 3, 3), 4)
	assert.Equal(t, nil, render.CheckImage([]string{
		"rrr",
		"ggg",
		"bbb",
	}, im))
}

func TestTransformationAll(t *testing.T) {
	// Checking colors with anti-aliasing (when scaling) is more complex.
	ic := render.ImageChecker{Palette: map[string]color.RGBA{
		"█": {0xff, 0xff, 0xff, 0xff},
		"▓": {0x80, 0x80, 0x80, 0x80},
		"▒": {0x7f, 0x7f, 0x7f, 0x7f},
		"░": {0x40, 0x40, 0x40, 0x40},
		".": {0, 0, 0, 0},
		"●": {0, 0, 0, 0xff},
		"◉": {0, 0, 0, 0x80},
		"◎": {0, 0, 0, 0x7f},
		"○": {0, 0, 0, 0x40},
	}}

	o := Transformation{
		// █.●
		// ...
		// ●.█
		Child: render.Box{
			Child: render.Column{
				Children: []render.Widget{
					render.Row{
						Children: []render.Widget{
							render.Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0xff, 0xff, 0xff}},
							render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0, 0}},
							render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0, 0xff}},
						},
					},
					render.Box{Width: 3, Height: 1, Color: color.RGBA{0, 0, 0, 0}},
					render.Row{
						Children: []render.Widget{
							render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0, 0xff}},
							render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0, 0}},
							render.Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0xff, 0xff, 0xff}},
						},
					},
				},
			}},
		Keyframes: []Keyframe{
			{
				Percentage: Percentage{0.0},
				Curve:      LinearCurve{},
				Transforms: []Transform{
					Translate{Vec2f{X: -3.0, Y: -3.0}},
					Scale{Vec2f{X: 1.0, Y: 1.0}},
					Rotate{Angle: 0.0},
				},
			},
			{
				Percentage: Percentage{0.75},
				Curve:      LinearCurve{},
				Transforms: []Transform{
					Translate{Vec2f{X: 0.0, Y: 0.0}},
					Scale{Vec2f{X: 1.0, Y: 1.0}},
					Rotate{Angle: 270.0},
				},
			},
			{
				Percentage: Percentage{1.0},
				Curve:      LinearCurve{},
				Transforms: []Transform{
					Translate{Vec2f{X: 1.0, Y: 1.0}},
					Scale{Vec2f{X: 2.0, Y: 2.0}},
					Rotate{Angle: 360.0},
				},
			},
		},
		Duration:     5,
		Delay:        0,
		Width:        9,
		Height:       9,
		Origin:       DefaultOrigin,
		Direction:    DefaultDirection,
		FillMode:     DefaultFillMode,
		Rounding:     DefaultRounding,
		WaitForChild: false,
	}

	// These frames should show the four "corners" being,
	// translated, rotated and in the end scaled to 2x.
	assert.Equal(t, 5, o.FrameCount())

	im := o.Paint(image.Rect(0, 0, 9, 9), 0)
	assert.Equal(t, nil, ic.Check([]string{
		"█.●......",
		".........",
		"●.█......",
		".........",
		".........",
		".........",
		".........",
		".........",
		".........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 9, 9), 1)
	assert.Equal(t, nil, ic.Check([]string{
		".........",
		".●.█.....",
		".........",
		".█.●.....",
		".........",
		".........",
		".........",
		".........",
		".........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 3, 3), 2)
	assert.Equal(t, nil, ic.Check([]string{
		".........",
		".........",
		"..█.●....",
		".........",
		"..●.█....",
		".........",
		".........",
		".........",
		".........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 3, 3), 3)
	assert.Equal(t, nil, ic.Check([]string{
		".........",
		".........",
		".........",
		"...●.█...",
		".........",
		"...█.●...",
		".........",
		".........",
		".........",
	}, im))

	im = o.Paint(image.Rect(0, 0, 3, 3), 4)

	assert.Equal(t, nil, ic.Check([]string{
		".........",
		".........",
		"..░▓░.○◉○",
		"..▓█▓.◉●◎",
		"..░▓░.○◉○",
		".........",
		"..○◎○.░▓░",
		"..◎●◉.▓█▓",
		"..○◎○.░▓░",
	}, im))
}

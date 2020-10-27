package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlotComputeLimits(t *testing.T) {
	p := Plot{
		X: []float64{3.14, 3.56, 3.9},
		Y: []float64{1.62, 2.7, 2.9},
	}

	check := func(xMin, xMax, yMin, yMax float64) {
		xA, xB, yA, yB := p.computeLimits()
		assert.Equal(t, xMin, xA)
		assert.Equal(t, xMax, xB)
		assert.Equal(t, yMin, yA)
		assert.Equal(t, yMax, yB)
	}

	// Without any limits set, data's min and max are used
	check(3.14, 3.9, 1.62, 2.9)

	// XLimMin below, within and above data
	p.XLimMin = new(float64)
	*p.XLimMin = 3.0
	check(3.0, 3.9, 1.62, 2.9)
	*p.XLimMin = 3.2
	check(3.2, 3.9, 1.62, 2.9)
	*p.XLimMin = 4.0
	check(4.0, 4.5, 1.62, 2.9)

	// XLimMax above, within and below data
	p.XLimMin = nil
	p.XLimMax = new(float64)
	*p.XLimMax = 4.1
	check(3.14, 4.1, 1.62, 2.9)
	*p.XLimMax = 3.2
	check(3.14, 3.2, 1.62, 2.9)
	*p.XLimMax = -17
	check(-17.5, -17, 1.62, 2.9)

	// YLimMin below, within and above data
	p.XLimMax = nil
	p.YLimMin = new(float64)
	*p.YLimMin = 1.0
	check(3.14, 3.9, 1, 2.9)
	*p.YLimMin = 2.0
	check(3.14, 3.9, 2.0, 2.9)
	p.YLimMin = new(float64)
	*p.YLimMin = 3.0
	check(3.14, 3.9, 3.0, 3.5)

	// YLimMax above, within and below data
	p.YLimMin = nil
	p.YLimMax = new(float64)
	*p.YLimMax = 3.14
	check(3.14, 3.9, 1.62, 3.14)
	*p.YLimMax = 2.0
	check(3.14, 3.9, 1.62, 2.0)
	*p.YLimMax = 1.0
	check(3.14, 3.9, 0.5, 1.0)

	// All limits in conjunction
	p.XLimMin = new(float64)
	p.XLimMax = new(float64)
	p.YLimMin = new(float64)
	p.YLimMax = new(float64)
	*p.XLimMin = 3.0
	*p.XLimMax = 4.0
	*p.YLimMin = 1.0
	*p.YLimMax = 3.0
	check(3.0, 4.0, 1.0, 3.0)
	*p.XLimMin = 3.3
	*p.XLimMax = 3.4
	*p.YLimMin = 2.1
	*p.YLimMax = 2.2
	check(3.3, 3.4, 2.1, 2.2)

	// No limits with single Y value centers vertically
	p.XLimMin = nil
	p.XLimMax = nil
	p.YLimMin = nil
	p.YLimMax = nil
	p.Y = []float64{3.14, 3.14, 3.14}
	check(3.14, 3.9, 3.14-0.5, 3.14+0.5)

	// No limits with single X value places points on left hand
	// side
	p.X = []float64{2, 2, 2}
	check(2, 2.5, 3.14-0.5, 3.14+0.5)

}

// Tests of the internal translatePoints() method
func TestPlotTranslatePoints(t *testing.T) {
	p := Plot{
		Width:  10,
		Height: 10,
	}

	p.X = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	p.Y = []float64{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}
	assert.Equal(t, []PathPoint{
		PathPoint{0, 9},
		PathPoint{1, 8},
		PathPoint{2, 7},
		PathPoint{3, 6},
		PathPoint{4, 5},
		PathPoint{5, 4},
		PathPoint{6, 3},
		PathPoint{7, 2},
		PathPoint{8, 1},
		PathPoint{9, 0},
	}, p.translatePoints())
	assert.Equal(t, 9, p.invThreshold)

	// Zoom in with XLim/YLim so that half the points fall outside
	// of view.
	p.XLimMin = new(float64)
	*p.XLimMin = 2
	p.XLimMax = new(float64)
	*p.XLimMax = 6
	p.YLimMin = new(float64)
	*p.YLimMin = 4
	p.YLimMax = new(float64)
	*p.YLimMax = 12

	// The points with X=2,3,4,5,6 will be mapped onto the 10x10
	// canvas. The lowest falls on 0 and the highest on 9. Since
	// they're equidistant, the stride between them must be 9/4 = 2.25.
	assert.Equal(t, []PathPoint{
		PathPoint{-5, 14}, // -4.5
		PathPoint{-2, 11}, // -2.25
		PathPoint{0, 9},   // 0
		PathPoint{2, 7},   // 2.25
		PathPoint{5, 4},   // 4.5
		PathPoint{7, 2},   // 6.75
		PathPoint{9, 0},   // 9
		PathPoint{11, -2}, // 11.25
		PathPoint{14, -5}, // 13.5
		PathPoint{16, -7}, // 15.75
	}, p.translatePoints())
	assert.Equal(t, 14, p.invThreshold)
}

func TestPlotFlatLine(t *testing.T) {
	ic := ImageChecker{
		palette: map[string]color.RGBA{
			"1": color.RGBA{0xff, 0xff, 0xff, 0xff},
			".": color.RGBA{0, 0, 0, 0},
		},
	}

	// Flatline
	p := Plot{
		Width:  10,
		Height: 5,
		X:      []float64{0, 9},
		Y:      []float64{47, 47},
	}
	assert.Equal(t, nil, ic.Check([]string{
		"..........",
		"..........",
		"1111111111",
		"..........",
		"..........",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))

	// Extend view to the eft
	p.XLimMin = new(float64)
	*p.XLimMin = -10
	assert.Equal(t, nil, ic.Check([]string{
		"..........",
		"..........",
		".....11111",
		"..........",
		"..........",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))

}

func TestPlotVerticalLine(t *testing.T) {
	ic := ImageChecker{
		palette: map[string]color.RGBA{
			"1": color.RGBA{0xff, 0xff, 0xff, 0xff},
			".": color.RGBA{0, 0, 0, 0},
		},
	}

	// Flatline
	p := Plot{
		Width:  10,
		Height: 5,
		X:      []float64{37, 37},
		Y:      []float64{1, -3},
	}

	assert.Equal(t, nil, ic.Check([]string{
		"1.........",
		"1.........",
		"1.........",
		"1.........",
		"1.........",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))
}

func TestPlotJaggedLine(t *testing.T) {
	ic := ImageChecker{
		palette: map[string]color.RGBA{
			"1": color.RGBA{0xff, 0xff, 0xff, 0xff},
			".": color.RGBA{0, 0, 0, 0},
		},
	}

	p := Plot{
		Width:  10,
		Height: 5,
		X:      []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Y:      []float64{1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
	}

	assert.Equal(t, nil, ic.Check([]string{
		"....1....1",
		"...11...1.",
		"..1..1.1..",
		".1...11...",
		"1....1....",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))

	// Make it bigger
	p.Width = 20
	p.Height = 10
	assert.Equal(t, nil, ic.Check([]string{
		"........1..........1",
		".......11.........1.",
		"......1..1.......1..",
		".....1...1......1...",
		"....1....1.....1....",
		"...1......1...1.....",
		"...1......1...1.....",
		"..1.......1..1......",
		".1.........11.......",
		"1..........1........",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))

	// Zoom in on the second valley
	p.XLimMin = new(float64)
	p.XLimMax = new(float64)
	p.YLimMin = new(float64)
	p.YLimMax = new(float64)
	*p.XLimMin = 5
	*p.XLimMax = 7
	*p.YLimMin = 1
	*p.YLimMax = 5

	assert.Equal(t, nil, ic.Check([]string{
		"1...................",
		".1..................",
		"..1.................",
		"...1................",
		"....1...............",
		".....11.............",
		".......1............",
		"........1........111",
		".........1...1111...",
		"..........111.......",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))

}

func TestPlotFewPoints(t *testing.T) {
	ic := ImageChecker{
		palette: map[string]color.RGBA{
			"1": color.RGBA{0xff, 0xff, 0xff, 0xff},
			".": color.RGBA{0, 0, 0, 0},
		},
	}

	p := Plot{
		Width:  10,
		Height: 10,
		X:      []float64{100, 200, 200, 100, 100, 200, 200, 100},
		Y:      []float64{-10, -10, -20, -20, -10, -20, -10, -20},
	}

	assert.Equal(t, nil, ic.Check([]string{
		"1111111111",
		"11......11",
		"1.1....1.1",
		"1..1..1..1",
		"1...11...1",
		"1...11...1",
		"1..1..1..1",
		"1.1....1.1",
		"11......11",
		"1111111111",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))

}

func TestPlotInvertedColor(t *testing.T) {
	ic := ImageChecker{
		palette: map[string]color.RGBA{
			"1": color.RGBA{0x0, 0xff, 0x0, 0xff},
			"2": color.RGBA{0xff, 0, 0, 0xff},
			".": color.RGBA{0, 0, 0, 0},
		},
	}

	p := Plot{
		Width:         10,
		Height:        5,
		X:             []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		Y:             []float64{1, 1, 0, 0, -1, -1, 1, 1, 0, 0},
		Color:         &color.RGBA{0, 0xff, 0, 0xff},
		ColorInverted: &color.RGBA{0xff, 0, 0, 0xff},
	}
	assert.Equal(t, nil, ic.Check([]string{
		"11....11..",
		"..1...1.1.",
		"..11..1.11",
		"....22....",
		"....22....",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))

	p.Width = 20
	p.Height = 10
	assert.Equal(t, nil, ic.Check([]string{
		"111..........111....",
		"...1.........1..1...",
		"...1.........1..1...",
		"....1.......1....1..",
		"....111.....1....111",
		"......2.....2.......",
		".......2....2.......",
		".......2...2........",
		"........2..2........",
		"........2222........",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))

}

func TestPlotSurfaceFill(t *testing.T) {
	ic := ImageChecker{
		palette: map[string]color.RGBA{
			"1": color.RGBA{0x0, 0xff, 0x0, 0xff},
			",": color.RGBA{0x0, 0xff, 0x0, 0x55},
			"2": color.RGBA{0xff, 0, 0, 0xff},
			":": color.RGBA{0xff, 0, 0, 0x55},
			".": color.RGBA{0, 0, 0, 0},
		},
	}

	// Fill with single color
	p := Plot{
		Width:  20,
		Height: 10,
		X:      []float64{0, 1, 2, 3, 4},
		Y:      []float64{5, 5, -1, -1, 2},
		Color:  &color.RGBA{0, 0xff, 0, 0xff},
		Fill:   true,
	}
	assert.Equal(t, nil, ic.Check([]string{
		"111111..............",
		",,,,,,1.............",
		",,,,,,1.............",
		",,,,,,,1............",
		",,,,,,,1...........1",
		",,,,,,,,1.........1,",
		",,,,,,,,1........1,,",
		",,,,,,,,,1......1,,,",
		".........1,,,,,1....",
		"..........11111.....",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))

	// Fil with ColorInverted
	p.ColorInverted = &color.RGBA{0xff, 0, 0, 0xff}
	assert.Equal(t, nil, ic.Check([]string{
		"111111..............",
		",,,,,,1.............",
		",,,,,,1.............",
		",,,,,,,1............",
		",,,,,,,1...........1",
		",,,,,,,,1.........1,",
		",,,,,,,,1........1,,",
		",,,,,,,,,1......1,,,",
		".........2:::::2....",
		"..........22222.....",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))
}

func TestPlotXLim(t *testing.T) {
	ic := ImageChecker{
		palette: map[string]color.RGBA{
			"1": color.RGBA{0xff, 0xff, 0xff, 0xff},
			".": color.RGBA{0, 0, 0, 0},
		},
	}

	p := Plot{
		Width:   20,
		Height:  10,
		X:       []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Y:       []float64{1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
		XLimMin: new(float64),
		XLimMax: new(float64),
	}

	// No change when min/max matches that of data
	*p.XLimMin = 1
	*p.XLimMax = 10

	// More space on the right
	*p.XLimMin = 1
	*p.XLimMax = 15
	assert.Equal(t, nil, ic.Check([]string{
		".....1......1.......",
		".....1......1.......",
		"....11.....1........",
		"....1.1....1........",
		"...1..1...1.........",
		"..1...1..1..........",
		"..1...1..1..........",
		".1.....11...........",
		".1.....11...........",
		"1......1............",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))

	// And on the left
	*p.XLimMin = -4
	*p.XLimMax = 15
	assert.Equal(t, nil, ic.Check([]string{
		".........1....1.....",
		".........1....1.....",
		"........11...1......",
		"........11...1......",
		".......1.1..1.......",
		".......1..1.1.......",
		"......1...11........",
		"......1...11........",
		"......1...11........",
		".....1....1.........",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))

	// And then do the opposite to "zoom in"
	*p.XLimMin = 3
	*p.XLimMax = 8
	assert.Equal(t, nil, ic.Check([]string{
		".......11...........",
		".....11.1...........",
		"...11....1..........",
		".11......1..........",
		"1........1.........1",
		"..........1......11.",
		"..........1.....1...",
		"..........1...11....",
		"...........111......",
		"...........1........",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))
}

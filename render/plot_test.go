package render

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var Empty = [2]float64{math.NaN(), math.NaN()}

func TestPlotComputeLimits(t *testing.T) {
	p := Plot{
		Data: [][2]float64{
			{3.14, 1.62},
			{3.56, 2.7},
			{3.9, 2.9},
		},
		XLim: Empty,
		YLim: Empty,
	}

	check := func(xMin, xMax, yMin, yMax float64) {
		xA, xB, yA, yB := p.computeLimits()
		assert.Equal(t, xMin, xA)
		assert.Equal(t, xMax, xB)
		assert.Equal(t, yMin, yA)
		assert.Equal(t, yMax, yB)
	}

	reset := func() {
		p.XLim = Empty
		p.YLim = Empty
	}

	// Without any limits set, data's min and max are used
	check(3.14, 3.9, 1.62, 2.9)

	// XLim min below, within and above data
	reset()
	p.XLim[0] = 3.0
	check(3.0, 3.9, 1.62, 2.9)
	p.XLim[0] = 3.2
	check(3.2, 3.9, 1.62, 2.9)
	p.XLim[0] = 4.0
	check(4.0, 4.5, 1.62, 2.9)

	// XLim max above, within and below data
	reset()
	p.XLim[1] = 4.1
	check(3.14, 4.1, 1.62, 2.9)
	p.XLim[1] = 3.2
	check(3.14, 3.2, 1.62, 2.9)
	p.XLim[1] = -17
	check(-17.5, -17, 1.62, 2.9)

	// YLim min below, within and above data
	reset()
	p.YLim[0] = 1.0
	check(3.14, 3.9, 1, 2.9)
	p.YLim[0] = 2.0
	check(3.14, 3.9, 2.0, 2.9)
	p.YLim[0] = 3.0
	check(3.14, 3.9, 3.0, 3.5)

	// YLim max above, within and below data
	reset()
	p.YLim[1] = 3.14
	check(3.14, 3.9, 1.62, 3.14)
	p.YLim[1] = 2.0
	check(3.14, 3.9, 1.62, 2.0)
	p.YLim[1] = 1.0
	check(3.14, 3.9, 0.5, 1.0)

	// All limits in conjunction
	p.XLim[0] = 3.0
	p.XLim[1] = 4.0
	p.YLim[0] = 1.0
	p.YLim[1] = 3.0
	check(3.0, 4.0, 1.0, 3.0)
	p.XLim[0] = 3.3
	p.XLim[1] = 3.4
	p.YLim[0] = 2.1
	p.YLim[1] = 2.2
	check(3.3, 3.4, 2.1, 2.2)

	// No limits with single Y value centers vertically
	reset()
	p.Data = [][2]float64{
		{3.14, 3.14},
		{3.56, 3.14},
		{3.9, 3.14},
	}
	check(3.14, 3.9, 3.14-0.5, 3.14+0.5)

	// No limits with single X value places points on left hand
	// side
	reset()
	p.Data = [][2]float64{
		{2, 3.14},
		{2, 3.14},
		{2, 3.14},
	}
	check(2, 2.5, 3.14-0.5, 3.14+0.5)

}

// Tests of the internal translatePoints() method
func TestPlotTranslatePoints(t *testing.T) {
	p := Plot{
		Width:  10,
		Height: 10,
		XLim:   Empty,
		YLim:   Empty,
	}

	p.Data = [][2]float64{
		{0, 0},
		{1, 2},
		{2, 4},
		{3, 6},
		{4, 8},
		{5, 10},
		{6, 12},
		{7, 14},
		{8, 16},
		{9, 18},
	}
	assert.Equal(t, []PathPoint{
		{0, 9},
		{1, 8},
		{2, 7},
		{3, 6},
		{4, 5},
		{5, 4},
		{6, 3},
		{7, 2},
		{8, 1},
		{9, 0},
	}, p.translatePoints())
	assert.Equal(t, 9, p.invThreshold)

	// Zoom in with XLim/YLim so that half the points fall outside
	// of view.
	p.XLim[0] = 2
	p.XLim[1] = 6
	p.YLim[0] = 4
	p.YLim[1] = 12

	// The points with X=2,3,4,5,6 will be mapped onto the 10x10
	// canvas. The lowest falls on 0 and the highest on 9. Since
	// they're equidistant, the stride between them must be 9/4 = 2.25.
	assert.Equal(t, []PathPoint{
		{-5, 14}, // -4.5
		{-2, 11}, // -2.25
		{0, 9},   // 0
		{2, 7},   // 2.25
		{5, 4},   // 4.5
		{7, 2},   // 6.75
		{9, 0},   // 9
		{11, -2}, // 11.25
		{14, -5}, // 13.5
		{16, -7}, // 15.75
	}, p.translatePoints())
	assert.Equal(t, 14, p.invThreshold)
}

func TestPlotFlatLine(t *testing.T) {
	ic := ImageChecker{
		Palette: map[string]color.RGBA{
			"1": {0xff, 0xff, 0xff, 0xff},
			".": {0, 0, 0, 0},
		},
	}

	// Flatline
	p := Plot{
		Width:  10,
		Height: 5,
		Data:   [][2]float64{{0, 47}, {9, 47}},
		XLim:   Empty,
		YLim:   Empty,
	}
	assert.Equal(t, nil, ic.Check([]string{
		"..........",
		"..........",
		"1111111111",
		"..........",
		"..........",
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))

	// Extend view to the left
	p.XLim[0] = -10
	assert.Equal(t, nil, ic.Check([]string{
		"..........",
		"..........",
		".....11111",
		"..........",
		"..........",
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))

}

func TestPlotVerticalLine(t *testing.T) {
	ic := ImageChecker{
		Palette: map[string]color.RGBA{
			"1": color.RGBA{0xff, 0xff, 0xff, 0xff},
			".": color.RGBA{0, 0, 0, 0},
		},
	}

	// Flatline
	p := Plot{
		Width:  10,
		Height: 5,
		Data:   [][2]float64{{37, 1}, {37, -3}},
		XLim:   Empty,
		YLim:   Empty,
	}

	assert.Equal(t, nil, ic.Check([]string{
		"1.........",
		"1.........",
		"1.........",
		"1.........",
		"1.........",
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))
}

func TestPlotJaggedLine(t *testing.T) {
	ic := ImageChecker{
		Palette: map[string]color.RGBA{
			"1": {0xff, 0xff, 0xff, 0xff},
			".": {0, 0, 0, 0},
		},
	}

	p := Plot{
		Width:  10,
		Height: 5,
		Data: [][2]float64{
			{1, 1},
			{2, 2},
			{3, 3},
			{4, 4},
			{5, 5},
			{6, 1},
			{7, 2},
			{8, 3},
			{9, 4},
			{10, 5},
		},
		XLim: Empty,
		YLim: Empty,
	}

	assert.Equal(t, nil, ic.Check([]string{
		"....1....1",
		"...11...1.",
		"..1..1.1..",
		".1...11...",
		"1....1....",
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))

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
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))

	// Zoom in on the second valley
	p.XLim[0] = 5
	p.XLim[1] = 7
	p.YLim[0] = 1
	p.YLim[1] = 5

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
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))

}

func TestPlotFewPoints(t *testing.T) {
	ic := ImageChecker{
		Palette: map[string]color.RGBA{
			"1": {0xff, 0xff, 0xff, 0xff},
			".": {0, 0, 0, 0},
		},
	}

	p := Plot{
		Width:  10,
		Height: 10,
		Data: [][2]float64{
			{100, -10},
			{200, -10},
			{200, -20},
			{100, -20},
			{100, -10},
			{200, -20},
			{200, -10},
			{100, -20},
		},
		XLim: Empty,
		YLim: Empty,
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
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))

}

func TestPlotInvertedColor(t *testing.T) {
	ic := ImageChecker{
		Palette: map[string]color.RGBA{
			"1": {0x0, 0xff, 0x0, 0xff},
			"2": {0xff, 0, 0, 0xff},
			".": {0, 0, 0, 0},
		},
	}

	p := Plot{
		Width:  10,
		Height: 5,
		Data: [][2]float64{
			{0, 1},
			{1, 1},
			{2, 0},
			{3, 0},
			{4, -1},
			{5, -1},
			{6, 1},
			{7, 1},
			{8, 0},
			{9, 0},
		},
		XLim:          Empty,
		YLim:          Empty,
		Color:         &color.RGBA{0, 0xff, 0, 0xff},
		ColorInverted: &color.RGBA{0xff, 0, 0, 0xff},
	}
	assert.Equal(t, nil, ic.Check([]string{
		"11....11..",
		"..1...1.1.",
		"..11..1.11",
		"....22....",
		"....22....",
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))

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
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))

}

func TestPlotSurfaceFill(t *testing.T) {
	ic := ImageChecker{
		Palette: map[string]color.RGBA{
			"1": {0x0, 0xff, 0x0, 0xff},
			",": {0x0, 0x55, 0x0, 0xff},
			"2": {0xff, 0, 0, 0xff},
			":": {0x55, 0, 0, 0xff},
			".": {0, 0, 0, 0},
		},
	}

	// Fill with single color
	p := Plot{
		Width:  20,
		Height: 10,
		Data: [][2]float64{
			{0, 5},
			{1, 5},
			{2, -1},
			{3, -1},
			{4, 2},
		},
		XLim:  Empty,
		YLim:  Empty,
		Color: &color.RGBA{0, 0xff, 0, 0xff},
		Fill:  true,
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
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))

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
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))
}

func TestPlotXLim(t *testing.T) {
	ic := ImageChecker{
		Palette: map[string]color.RGBA{
			"1": {0xff, 0xff, 0xff, 0xff},
			".": {0, 0, 0, 0},
		},
	}

	p := Plot{
		Width:  20,
		Height: 10,
		Data: [][2]float64{
			{1, 1},
			{2, 2},
			{3, 3},
			{4, 4},
			{5, 5},
			{6, 1},
			{7, 2},
			{8, 3},
			{9, 4},
			{10, 5},
		},
		XLim: Empty,
		YLim: Empty,
	}

	// No change when min/max matches that of data
	p.XLim[0] = 1
	p.XLim[1] = 10

	// More space on the right
	p.XLim[0] = 1
	p.XLim[1] = 15
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
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))

	// And on the left
	p.XLim[0] = -4
	p.XLim[1] = 15
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
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))

	// And then do the opposite to "zoom in"
	p.XLim[0] = 3
	p.XLim[1] = 8
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
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))
}

func TestPlotScatter(t *testing.T) {
	ic := ImageChecker{
		Palette: map[string]color.RGBA{
			"1": {0xff, 0xff, 0xff, 0xff},
			".": {0, 0, 0, 0},
		},
	}

	// Flatline
	p := Plot{
		Width:     10,
		Height:    5,
		Data:      [][2]float64{{0, 47}, {9, 47}},
		XLim:      Empty,
		YLim:      Empty,
		ChartType: "scatter",
	}
	assert.Equal(t, nil, ic.Check([]string{
		"..........",
		"..........",
		"1........1",
		"..........",
		"..........",
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))

	// Extend view to the left
	p.XLim[0] = -10
	assert.Equal(t, nil, ic.Check([]string{
		"..........",
		"..........",
		".....1...1",
		"..........",
		"..........",
	}, PaintWidget(p, image.Rect(0, 0, 100, 100), 0)))

}

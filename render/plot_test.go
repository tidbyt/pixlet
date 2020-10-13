package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Warning: The expected results below are really just copy-paste of
// the actual output of Plot when it was in a state that "looked
// ok". I tried to get this thing to be super predictable and pixel
// perfect, but it's really tricky with a) floating point math and b)
// Bresenham's line drawing algorithm, so this will have to do.

func TestPlot(t *testing.T) {
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
		"........1...........",
		".......11.........11",
		"......1..1.......1..",
		".....1...1......1...",
		".....1...1......1...",
		"....1.....1....1....",
		"...1......1...1.....",
		"..1.......1..1......",
		".1.........11.......",
		"1..........1........",
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
	assert.Equal(t, nil, ic.Check([]string{
		"........1...........",
		".......11.........11",
		"......1..1.......1..",
		".....1...1......1...",
		".....1...1......1...",
		"....1.....1....1....",
		"...1......1...1.....",
		"..1.......1..1......",
		".1.........11.......",
		"1..........1........",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))

	// More space on the right
	*p.XLimMin = 1
	*p.XLimMax = 15
	assert.Equal(t, nil, ic.Check([]string{
		".....1......1.......",
		".....1......1.......",
		"....11.....1........",
		"...1..1....1........",
		"...1..1...1.........",
		"..1...1...1.........",
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
		".......1..1.1.......",
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
		"....1....1..........",
		"..11.....1..........",
		".1........1.........",
		"1.........1........1",
		"...........1.....11.",
		"...........1...11...",
		"............111.....",
		"............1.......",
	}, p.Paint(image.Rect(0, 0, 100, 100), 0)))
}

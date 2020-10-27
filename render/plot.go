package render

import (
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

var DefaultPlotColor = color.RGBA{0xff, 0xff, 0xff, 0xff}

// surface fill gets line color with this alpha
var FillAlpha uint8 = 0x55

type Plot struct {
	Widget

	// Coordinates of points to plot
	X []float64
	Y []float64

	// Overall size of the plot
	Height int
	Width  int

	// Primary line color
	Color *color.RGBA

	// Optional line color for Y-values below 0
	ColorInverted *color.RGBA

	// Optional limit on X and Y axis
	XLimMin *float64
	XLimMax *float64
	YLimMin *float64
	YLimMax *float64

	// If true, also paint surface between line and X-axis
	Fill bool

	invThreshold int
}

// Computes X and Y limits
func (p *Plot) computeLimits() (float64, float64, float64, float64) {

	// If all limits are set by user, no computation is required
	if p.XLimMin != nil && p.XLimMax != nil && p.YLimMin != nil && p.YLimMax != nil {
		return *p.XLimMin, *p.XLimMax, *p.YLimMin, *p.YLimMax
	}

	// Otherwise we'll need min/max of X and Y
	minX, maxX, minY, maxY := p.X[0], p.X[0], p.Y[0], p.Y[0]
	for i := 1; i < len(p.X); i++ {
		if p.X[i] < minX {
			minX = p.X[i]
		}
		if p.X[i] > maxX {
			maxX = p.X[i]
		}
		if p.Y[i] < minY {
			minY = p.Y[i]
		}
		if p.Y[i] > maxY {
			maxY = p.Y[i]
		}
	}

	// Limits not set by user will default to the min/max of the
	// data, so that it all fits on canvas.
	xLimMin := minX
	xLimMax := maxX
	yLimMin := minY
	yLimMax := maxY
	if p.XLimMin != nil {
		xLimMin = *p.XLimMin
	}
	if p.XLimMax != nil {
		xLimMax = *p.XLimMax
	}
	if p.YLimMin != nil {
		yLimMin = *p.YLimMin
	}
	if p.YLimMax != nil {
		yLimMax = *p.YLimMax
	}

	// The inferred limits can be non-sensical if user provides
	// only the min or max of a limit. In these cases, we take the
	// provided limit and add an arbitraty +-0.5 to create limits
	// that result in all points displayed "off-screen".
	if xLimMax < xLimMin {
		if p.XLimMin == nil {
			xLimMin = xLimMax - 0.5
		} else {
			xLimMax = xLimMin + 0.5
		}
	}
	if yLimMax < yLimMin {
		if p.YLimMin == nil {
			yLimMin = yLimMax - 0.5
		} else {
			yLimMax = yLimMin + 0.5
		}
	}

	// If all X or all Y are equal, then the default limits would
	// have min==max, which is non-sensical.
	if xLimMin == xLimMax {
		// Place points furthest left on canvas
		xLimMin = minX
		xLimMax = minX + 0.5
	}
	if yLimMin == yLimMax {
		// Place points in vertical center
		yLimMin = minY - 0.5
		yLimMax = minY + 0.5
	}

	return xLimMin, xLimMax, yLimMin, yLimMax
}

// Maps the points in X and Y to positions on the canvas
func (p *Plot) translatePoints() []PathPoint {
	xLimMin, xLimMax, yLimMin, yLimMax := p.computeLimits()

	// Translate
	points := make([]PathPoint, len(p.X))
	for i := 0; i < len(p.X); i++ {
		nX := (p.X[i] - xLimMin) / (xLimMax - xLimMin)
		nY := (p.Y[i] - yLimMin) / (yLimMax - yLimMin)
		points[i] = PathPoint{
			X: int(math.Round(nX * float64(p.Width-1))),
			Y: p.Height - 1 - int(math.Round(nY*float64(p.Height-1))),
		}
	}
	p.invThreshold = p.Height - 1 - int(math.Round(((0-yLimMin)/(yLimMax-yLimMin))*float64(p.Height-1)))

	return points
}

func (p Plot) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	dc := gg.NewContext(p.Width, p.Height)

	// Set line and fill colors
	col := color.RGBA{0xff, 0xff, 0xff, 0xff}
	if p.Color != nil {
		col = *p.Color
	}
	colInv := col
	if p.ColorInverted != nil {
		colInv = *p.ColorInverted
	}
	fillCol := col
	fillCol.A = FillAlpha
	fillColInv := colInv
	fillColInv.A = FillAlpha

	pl := &PolyLine{Vertices: p.translatePoints()}

	// the optional surface fill
	for i := 0; p.Fill && i < pl.Length(); i++ {
		x, y := pl.Point(i)
		if x < 0 || x >= p.Width || y < 0 || y >= p.Height {
			continue
		}
		if y > p.invThreshold {
			dc.SetColor(fillColInv)
			for ; y != p.invThreshold; y-- {
				dc.SetPixel(x, y)
			}
		} else {
			dc.SetColor(fillCol)
			for ; y <= p.invThreshold; y++ {
				dc.SetPixel(x, y)
			}
		}
	}

	// the line itself
	for i := 0; i < pl.Length(); i++ {
		x, y := pl.Point(i)
		if y > p.invThreshold {
			dc.SetColor(colInv)
		} else {
			dc.SetColor(col)
		}
		dc.SetPixel(x, y)
	}

	return dc.Image()
}

func (p Plot) FrameCount() int {
	return 1
}

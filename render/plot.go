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

	// If true, also pain surface between line and X-axis
	Fill bool

	invThreshold int
}

// translatePoints() maps the points in X and Y to fit in Width x
// Height pixels. it sets invThreshold to the height corresponding to
// Y=0.
func (p *Plot) translatePoints() []PathPoint {
	// Find min/max X and Y
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

	// Xlim and YLim simply overrides actual min/max
	if p.XLimMin != nil {
		minX = *p.XLimMin
	}
	if p.XLimMax != nil {
		maxX = *p.XLimMax
	}
	if p.YLimMin != nil {
		minY = *p.YLimMin
	}
	if p.YLimMax != nil {
		maxY = *p.YLimMax
	}

	// Add a small epsilon to maxY so normalization interval
	// becomes half open. This in turn keeps points at the maximum
	// Y-value from being mapped to p.Height, which would be out
	// of bounds.
	epsilon := (maxY - minY) / (float64(p.Height) * 100)
	maxY += epsilon

	// Normalize to [0,1) and compute pixel position
	points := make([]PathPoint, len(p.X))
	for i := 0; i < len(p.X); i++ {
		nX := (p.X[i] - minX) / (maxX - minX)
		nY := (p.Y[i] - minY) / (maxY - minY)
		points[i].X = int(math.Floor(nX * float64(p.Width)))
		points[i].Y = p.Height - 1 - int(math.Floor(nY*float64(p.Height)))
	}

	// In the same way, translate Y=0
	p.invThreshold = p.Height - 1 - int(math.Floor((0-minY/(maxY-minY))*float64(p.Height)))

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

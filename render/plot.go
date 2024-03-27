package render

import (
	"image"
	"image/color"
	"math"

	"tidbyt.dev/pixlet/render/canvas"
)

var DefaultPlotColor = color.RGBA{0xff, 0xff, 0xff, 0xff}

// surface fill gets line color dampened by this factor
var FillDampFactor uint8 = 0x55

// Plot is a widget that draws a data series.
//
// DOC(Data): A list of 2-tuples of numbers
// DOC(Width): Limits Plot width
// DOC(Height): Limits Plot height
// DOC(Color): Line color, default is '#fff'
// DOC(ColorInverted): Line color for Y-values below 0
// DOC(XLim): Limit X-axis to a range
// DOC(YLim): Limit Y-axis to a range
// DOC(Fill): Paint surface between line and X-axis
// DOC(FillColor): Fill color for Y-values above 0
// DOC(FillColorInverted): Fill color for Y-values below 0
// DOC(ChartType): Specifies the type of chart to render, "scatter" or "line", default is "line"
//
// EXAMPLE BEGIN
// render.Plot(
//
//	data = [
//	  (0, 3.35),
//	  (1, 2.15),
//	  (2, 2.37),
//	  (3, -0.31),
//	  (4, -3.53),
//	  (5, 1.31),
//	  (6, -1.3),
//	  (7, 4.60),
//	  (8, 3.33),
//	  (9, 5.92),
//	],
//	width = 64,
//	height = 32,
//	color = "#0f0",
//	color_inverted = "#f00",
//	x_lim = (0, 9),
//	y_lim = (-5, 7),
//	fill = True,
//
// ),
// EXAMPLE END
type Plot struct {
	Widget

	// Coordinates of points to plot
	Data [][2]float64 `starlark:"data,required"`

	// Overall size of the plot
	Width  int `starlark:"width,required"`
	Height int `starlark:"height,required"`

	// Primary line color
	Color color.Color `starlark:"color"`

	// Optional line color for Y-values below 0
	ColorInverted color.Color `starlark:"color_inverted"`

	// Optional limit on X and Y axis
	XLim [2]float64 `starlark:"x_lim"`
	YLim [2]float64 `starlark:"y_lim"`

	// If true, also paint surface between line and X-axis
	Fill bool `starlark:"fill"`

	// Optional, default "line". If set to "scatter", the line connecting dots will not be drawn
	ChartType string `starlark:"chart_type"`

	// Optional fill color for Y-values above 0
	FillColor color.Color `starlark:"fill_color"`

	// Optional fill color for Y-values below 0
	FillColorInverted color.Color `starlark:"fill_color_inverted"`

	invThreshold int
}

// Computes X and Y limits
func (p *Plot) computeLimits() (float64, float64, float64, float64) {

	// If all limits are set by user, no computation is required
	if !math.IsNaN(p.XLim[0]) && !math.IsNaN(p.XLim[1]) &&
		!math.IsNaN(p.YLim[0]) && !math.IsNaN(p.YLim[1]) {
		return p.XLim[0], p.XLim[1], p.YLim[0], p.YLim[1]
	}

	// Otherwise we'll need min/max of X and Y
	pt := p.Data[0]
	minX, maxX, minY, maxY := pt[0], pt[0], pt[1], pt[1]
	for i := 1; i < len(p.Data); i++ {
		pt = p.Data[i]
		if pt[0] < minX {
			minX = pt[0]
		}
		if pt[0] > maxX {
			maxX = pt[0]
		}
		if pt[1] < minY {
			minY = pt[1]
		}
		if pt[1] > maxY {
			maxY = pt[1]
		}
	}

	// Limits not set by user will default to the min/max of the
	// data, so that it all fits on canvas.
	xLimMin := minX
	xLimMax := maxX
	yLimMin := minY
	yLimMax := maxY
	if !math.IsNaN(p.XLim[0]) {
		xLimMin = p.XLim[0]
	}
	if !math.IsNaN(p.XLim[1]) {
		xLimMax = p.XLim[1]
	}
	if !math.IsNaN(p.YLim[0]) {
		yLimMin = p.YLim[0]
	}
	if !math.IsNaN(p.YLim[1]) {
		yLimMax = p.YLim[1]
	}

	// The inferred limits can be non-sensical if user provides
	// only the min or max of a limit. In these cases, we take the
	// provided limit and add an arbitraty +-0.5 to create limits
	// that result in all points displayed "off-screen".
	if xLimMax < xLimMin {
		if math.IsNaN(p.XLim[0]) {
			xLimMin = xLimMax - 0.5
		} else {
			xLimMax = xLimMin + 0.5
		}
	}
	if yLimMax < yLimMin {
		if math.IsNaN(p.YLim[0]) {
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
	points := make([]PathPoint, len(p.Data))
	for i := 0; i < len(p.Data); i++ {
		pt := p.Data[i]
		nX := (pt[0] - xLimMin) / (xLimMax - xLimMin)
		nY := (pt[1] - yLimMin) / (yLimMax - yLimMin)
		points[i] = PathPoint{
			X: int(math.Round(nX * float64(p.Width-1))),
			Y: p.Height - 1 - int(math.Round(nY*float64(p.Height-1))),
		}
	}
	p.invThreshold = p.Height - 1 - int(math.Round(((0-yLimMin)/(yLimMax-yLimMin))*float64(p.Height-1)))

	return points
}

func dampenColor(c color.Color, a uint8) color.Color {
	r, g, b, _ := c.RGBA()
	return color.RGBA{uint8(r * uint32(a) / 255), uint8(g * uint32(a) / 255), uint8(b * uint32(a) / 255), 0xFF}
}

func (p Plot) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	return image.Rect(0, 0, p.Width, p.Height)
}

func (p Plot) Paint(dc canvas.Canvas, bounds image.Rectangle, frameIdx int) {
	// Set line and fill colors
	var col color.Color
	col = color.RGBA{0xff, 0xff, 0xff, 0xff}
	if p.Color != nil {
		col = p.Color
	}
	colInv := col
	if p.ColorInverted != nil {
		colInv = p.ColorInverted
	}

	fillCol := dampenColor(col, FillDampFactor)
	if p.FillColor != nil {
		fillCol = p.FillColor
	}

	fillColInv := dampenColor(colInv, FillDampFactor)
	if p.FillColorInverted != nil {
		fillColInv = p.FillColorInverted
	}

	pl := &PolyLine{Vertices: p.translatePoints()}

	// the optional surface fill
	for i := 0; p.Fill && i < pl.Length(); i++ {
		x, y := pl.Point(i)
		if x < 0 || x >= p.Width || y < 0 || y >= p.Height {
			continue
		}
		if y > p.invThreshold {
			dc.SetColor(fillColInv)
			for ; y != p.invThreshold && y >= 0; y-- {
				tx, ty := dc.TransformPoint(float64(x), float64(y))
				dc.DrawPixel(int(tx), int(ty))
			}
		} else {
			dc.SetColor(fillCol)
			for ; y <= p.invThreshold && y <= p.Height; y++ {
				tx, ty := dc.TransformPoint(float64(x), float64(y))
				dc.DrawPixel(int(tx), int(ty))
			}
		}
	}

	if p.ChartType == "scatter" {
		points := p.translatePoints()
		for i := 0; i < len(points); i++ {
			point := points[i]
			if point.Y > p.invThreshold {
				dc.SetColor(colInv)
			} else {
				dc.SetColor(col)
			}
			dc.DrawPixel(int(point.X), int(point.Y))
		}
	} else {
		// the line itself
		for i := 0; i < pl.Length(); i++ {
			x, y := pl.Point(i)
			if y > p.invThreshold {
				dc.SetColor(colInv)
			} else {
				dc.SetColor(col)
			}
			tx, ty := dc.TransformPoint(float64(x), float64(y))
			dc.DrawPixel(int(tx), int(ty))
		}
	}
}

func (p Plot) FrameCount() int {
	return 1
}

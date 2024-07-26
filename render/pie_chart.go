package render

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/tidbyt/gg"
)

// PieChart draws a circular pie chart of size `diameter`. It takes two
// arguments for the data: parallel lists `colors` and `weights` representing
// the shading and relative sizes of each data entry.
//
// DOC(Colors): List of color hex codes
// DOC(Weights): List of numbers corresponding to the relative size of each color
// DOC(Diameter): Diameter of the circle
//
// EXAMPLE BEGIN
// render.PieChart(
//      colors = [ "#fff", "#0f0", "#00f" ],
//      weights  = [ 180, 135, 45 ],
//      diameter = 30,
// )
// EXAMPLE END
type PieChart struct {
	Type string `starlark:"-"`

	Colors   []color.RGBA `starlark:"colors, required"`
	Weights  []float64    `starlark:"weights, required"`
	Diameter int          `starlark:"diameter,required"`
}

func (c PieChart) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	return image.Rect(0, 0, c.Diameter, c.Diameter)
}

func (c PieChart) Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int) {
	total := 0.0
	for _, v := range c.Weights {
		total += v
	}

	r := float64(c.Diameter) / 2

	start := 0.0
	for i, v := range c.Weights {
		end := start + v/total
		dc.SetColor(c.Colors[i%len(c.Colors)])
		dc.DrawArc(r, r, r, start*2*math.Pi, end*2*math.Pi)
		dc.LineTo(r, r)
		dc.LineTo(r+r*math.Cos(start*2*math.Pi), r+r*math.Sin(start*2*math.Pi))
		dc.Fill()
		start = end
	}
}

func (c PieChart) FrameCount() int {
	return 1
}

func (c *PieChart) MarshalJSON() ([]byte, error) {
	type Alias PieChart
	aux := &struct {
		*Alias
		Colors []string
	}{
		Alias:  (*Alias)(c),
		Colors: make([]string, len(c.Colors)),
	}

	for i, col := range c.Colors {
		r, g, b, a := col.RGBA()
		aux.Colors[i] = fmt.Sprintf("#%02x%02x%02x%02x", r>>8, g>>8, b>>8, a>>8)
	}

	return json.Marshal(aux)
}

func (c *PieChart) UnmarshalJSON(data []byte) error {
	type Alias PieChart
	aux := &struct {
		*Alias
		Colors []string
	}{
		Alias: (*Alias)(c),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	c.Colors = make([]color.RGBA, len(aux.Colors))
	for i, col := range aux.Colors {
		c.Colors[i], err = ParseColor(col)
		if err != nil {
			return err
		}
	}

	return nil
}

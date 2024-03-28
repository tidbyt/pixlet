package render

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"tidbyt.dev/pixlet/render/canvas"
)

type Starfield struct {
	Widget

	Child  Widget
	Color  color.Color
	Width  int
	Height int
	stars  []*Star
}

type Star struct {
	X     float64
	Y     float64
	D     float64
	V     float64
	PrevX float64
	PrevY float64
}

func (s *Star) Tick() {
	s.PrevX = s.X
	s.PrevY = s.Y
	s.X = s.X / s.D
	s.Y = s.Y / s.D
	s.D *= 1 - s.V
}

// TODO: use PolyLine from path.go instead
func DrawLine(dc canvas.Canvas, x0, y0, x1, y1 int) {
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}

	sx := -1
	if x0 < x1 {
		sx = 1
	}

	dy := y1 - y0
	if dy > 0 {
		dy = -dy
	}

	sy := -1
	if y0 < y1 {
		sy = 1
	}

	if dx == 0 {
		for ; y0 != y1; y0 += sy {
			dc.DrawPixel(x0, y0)
		}
		dc.DrawPixel(x0, y0)
		return
	}

	if dy == 0 {
		for ; x0 != x1; x0 += sx {
			dc.DrawPixel(x0, y0)
		}
		dc.DrawPixel(x0, y0)
		return
	}

	err := dx + dy

	dc.DrawPixel(x0, y0)
	for x0 != x1 && y0 != y1 {
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
		dc.DrawPixel(x0, y0)
	}
}

func (s *Starfield) Init() error {
	s.initStars(60)
	return nil
}

func (s *Starfield) initStars(numStars int) {
	s.stars = make([]*Star, numStars)
	for i := 0; i < numStars; i++ {
		s.stars[i] = &Star{
			X: 2*rand.Float64() - 1,
			Y: 2*rand.Float64() - 1,
			D: 1.0,  //rand.Float64(),
			V: 0.01, //rand.Float64()
		}
		if i < 10 {
			fmt.Println(s.stars[i].X, s.stars[i].Y)
		}
	}
}

func (s *Starfield) PaintBounds(bounds image.Rectangle, frameIdx int) image.Image {
	return image.Rect(0, 0, 64, 32)
}

func (s *Starfield) Paint(dc canvas.Canvas, bounds image.Rectangle, frameIdx int) {
	Black := color.RGBA{R: 0, G: 0, B: 0, A: 0xff}

	dc.SetColor(Black)
	dc.Clear()

	for j := 0; j < len(s.stars); j++ {
		s.stars[j].Tick()
		if math.Abs(s.stars[j].X) > 1.0 || math.Abs(s.stars[j].Y) > 1.0 {
			s.stars[j].X = rand.NormFloat64() * 0.3
			s.stars[j].Y = rand.NormFloat64() * 0.3
			s.stars[j].PrevX = 0
			s.stars[j].PrevY = 0
			s.stars[j].D = 0.9 //rand.Float64()
			continue
		}

		pX := int(math.Round(float64(bounds.Dx()) * (s.stars[j].PrevX + 0.5)))
		pY := int(math.Round(float64(bounds.Dy()) * (s.stars[j].PrevY + 0.5)))
		X := int(math.Round(float64(bounds.Dx()) * (s.stars[j].X + 0.5)))
		Y := int(math.Round(float64(bounds.Dy()) * (s.stars[j].Y + 0.5)))

		if pX != 0 && pY != 0 && (pX != X || pY != Y) {
			dc.SetColor(color.RGBA{0x22, 0x22, 0x22, 0xff})
			DrawLine(dc, pX, pY, X, Y)
		}
		dc.SetColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
		dc.DrawPixel(X, Y)
	}

	dc.SetColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
}

func (s *Starfield) FrameCount() int {
	return 300
}

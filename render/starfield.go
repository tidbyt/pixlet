package render

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type Starfield struct {
	Widget

	Child  Widget
	Color  color.Color
	Width  int
	Height int
	frames []image.Image
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
func DrawLine(dc *gg.Context, x0, y0, x1, y1 int) {
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
			dc.SetPixel(x0, y0)
		}
		dc.SetPixel(x0, y0)
		return
	}

	if dy == 0 {
		for ; x0 != x1; x0 += sx {
			dc.SetPixel(x0, y0)
		}
		dc.SetPixel(x0, y0)
		return
	}

	err := dx + dy

	dc.SetPixel(x0, y0)
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
		dc.SetPixel(x0, y0)
	}
}

func AnimateStarfield(width, height, numStars, numFrames int) []image.Image {
	//White := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	Black := color.RGBA{R: 0, G: 0, B: 0, A: 0xff}

	star := make([]*Star, numStars)
	for i := 0; i < numStars; i++ {
		star[i] = &Star{
			X: 2*rand.Float64() - 1,
			Y: 2*rand.Float64() - 1,
			D: 1.0,  //rand.Float64(),
			V: 0.01, //rand.Float64()
		}
		if i < 10 {
			fmt.Println(star[i].X, star[i].Y)
		}
	}

	frames := make([]image.Image, 0, numFrames)
	for i := 0; i < numFrames; i++ {
		dc := gg.NewContext(width, height)
		dc.SetColor(Black)
		dc.Clear()

		for j := 0; j < len(star); j++ {
			star[j].Tick()
			if math.Abs(star[j].X) > 1.0 || math.Abs(star[j].Y) > 1.0 {
				star[j].X = rand.NormFloat64() * 0.3
				star[j].Y = rand.NormFloat64() * 0.3
				star[j].PrevX = 0
				star[j].PrevY = 0
				star[j].D = 0.9 //rand.Float64()
				continue
			}

			pX := int(math.Round(float64(width) * (star[j].PrevX + 0.5)))
			pY := int(math.Round(float64(height) * (star[j].PrevY + 0.5)))
			X := int(math.Round(float64(width) * (star[j].X + 0.5)))
			Y := int(math.Round(float64(height) * (star[j].Y + 0.5)))

			if pX != 0 && pY != 0 && (pX != X || pY != Y) {
				dc.SetColor(color.RGBA{0x22, 0x22, 0x22, 0xff})
				DrawLine(dc, pX, pY, X, Y)
			}
			dc.SetColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
			dc.SetPixel(X, Y)
		}

		dc.SetColor(color.RGBA{0xff, 0xff, 0xff, 0xff})

		frames = append(frames, dc.Image())
	}

	return frames

}

func (s *Starfield) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	if len(s.frames) == 0 {
		s.frames = AnimateStarfield(bounds.Dx(), bounds.Dy(), 40, 300)
	}

	return s.frames[frameIdx]
}

func (s *Starfield) FrameCount() int {
	if len(s.frames) == 0 {
		s.frames = AnimateStarfield(64, 32, 60, 300)
	}

	return len(s.frames)
}

package animation

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/render"
)

func TestBounceLinearHorizontalChildFits(t *testing.T) {
	b := Bounce{
		Child: render.Row{
			Children: []render.Widget{
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		Width:           6,
		BounceDirection: "horizontal",
		Curve:           LinearCurve{},
	}
	im := image.Rect(0, 0, 100, 100)

	assert.Equal(t, 1, b.FrameCount())
	assert.Equal(t, nil, render.CheckImage([]string{"rgb..."}, b.Paint(im, 0)))
}

func TestBounceAlwaysLinearHorizontalChildFits(t *testing.T) {
	b := Bounce{
		Child: render.Row{
			Children: []render.Widget{
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		Width:           6,
		BounceDirection: "horizontal",
		BounceAlways:    true,
		Curve:           LinearCurve{},
	}
	im := image.Rect(0, 0, 100, 100)

	assert.Equal(t, 6, b.FrameCount())
	assert.Equal(t, nil, render.CheckImage([]string{"rgb..."}, b.Paint(im, 0)))
	assert.Equal(t, nil, render.CheckImage([]string{".rgb.."}, b.Paint(im, 1)))
	assert.Equal(t, nil, render.CheckImage([]string{"..rgb."}, b.Paint(im, 2)))
	assert.Equal(t, nil, render.CheckImage([]string{"...rgb"}, b.Paint(im, 3)))
	assert.Equal(t, nil, render.CheckImage([]string{"..rgb."}, b.Paint(im, 4)))
	assert.Equal(t, nil, render.CheckImage([]string{".rgb.."}, b.Paint(im, 5)))
}

func TestBounceAlwaysLinearHorizontalChildFitsPause(t *testing.T) {
	b := Bounce{
		Child: render.Row{
			Children: []render.Widget{
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		Width:           6,
		BounceDirection: "horizontal",
		BounceAlways:    true,
		Pause:           3,
		Curve:           LinearCurve{},
	}
	im := image.Rect(0, 0, 100, 100)

	assert.Equal(t, 10, b.FrameCount())
	assert.Equal(t, nil, render.CheckImage([]string{"rgb..."}, b.Paint(im, 0)))
	assert.Equal(t, nil, render.CheckImage([]string{"rgb..."}, b.Paint(im, 1)))
	assert.Equal(t, nil, render.CheckImage([]string{"rgb..."}, b.Paint(im, 2)))
	assert.Equal(t, nil, render.CheckImage([]string{".rgb.."}, b.Paint(im, 3)))
	assert.Equal(t, nil, render.CheckImage([]string{"..rgb."}, b.Paint(im, 4)))
	assert.Equal(t, nil, render.CheckImage([]string{"...rgb"}, b.Paint(im, 5)))
	assert.Equal(t, nil, render.CheckImage([]string{"...rgb"}, b.Paint(im, 6)))
	assert.Equal(t, nil, render.CheckImage([]string{"...rgb"}, b.Paint(im, 7)))
	assert.Equal(t, nil, render.CheckImage([]string{"..rgb."}, b.Paint(im, 8)))
	assert.Equal(t, nil, render.CheckImage([]string{".rgb.."}, b.Paint(im, 9)))
}

func TestBounceAlwaysEaseInHorizontalChildFits(t *testing.T) {
	b := Bounce{
		Child: render.Row{
			Children: []render.Widget{
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		Width:           15,
		BounceDirection: "horizontal",
		BounceAlways:    true,
		Curve:           EaseIn,
	}
	im := image.Rect(0, 0, 100, 100)

	assert.Equal(t, 24, b.FrameCount())
	assert.Equal(t, nil, render.CheckImage([]string{"rgb............"}, b.Paint(im, 0)))
	assert.Equal(t, nil, render.CheckImage([]string{"rgb............"}, b.Paint(im, 1)))
	assert.Equal(t, nil, render.CheckImage([]string{".rgb..........."}, b.Paint(im, 2)))
	assert.Equal(t, nil, render.CheckImage([]string{"..rgb.........."}, b.Paint(im, 3)))
	assert.Equal(t, nil, render.CheckImage([]string{"..rgb.........."}, b.Paint(im, 4)))
	assert.Equal(t, nil, render.CheckImage([]string{"...rgb........."}, b.Paint(im, 5)))
	assert.Equal(t, nil, render.CheckImage([]string{"....rgb........"}, b.Paint(im, 6)))
	assert.Equal(t, nil, render.CheckImage([]string{"......rgb......"}, b.Paint(im, 7)))
	assert.Equal(t, nil, render.CheckImage([]string{".......rgb....."}, b.Paint(im, 8)))
	assert.Equal(t, nil, render.CheckImage([]string{"........rgb...."}, b.Paint(im, 9)))
	assert.Equal(t, nil, render.CheckImage([]string{".........rgb..."}, b.Paint(im, 10)))
	assert.Equal(t, nil, render.CheckImage([]string{"...........rgb."}, b.Paint(im, 11)))
	assert.Equal(t, nil, render.CheckImage([]string{"............rgb"}, b.Paint(im, 12)))
	assert.Equal(t, nil, render.CheckImage([]string{"...........rgb."}, b.Paint(im, 13)))
	assert.Equal(t, nil, render.CheckImage([]string{".........rgb..."}, b.Paint(im, 14)))
	assert.Equal(t, nil, render.CheckImage([]string{"........rgb...."}, b.Paint(im, 15)))
	assert.Equal(t, nil, render.CheckImage([]string{".......rgb....."}, b.Paint(im, 16)))
	assert.Equal(t, nil, render.CheckImage([]string{"......rgb......"}, b.Paint(im, 17)))
	assert.Equal(t, nil, render.CheckImage([]string{"....rgb........"}, b.Paint(im, 18)))
	assert.Equal(t, nil, render.CheckImage([]string{"...rgb........."}, b.Paint(im, 19)))
	assert.Equal(t, nil, render.CheckImage([]string{"..rgb.........."}, b.Paint(im, 20)))
	assert.Equal(t, nil, render.CheckImage([]string{"..rgb.........."}, b.Paint(im, 21)))
	assert.Equal(t, nil, render.CheckImage([]string{".rgb..........."}, b.Paint(im, 22)))
	assert.Equal(t, nil, render.CheckImage([]string{"rgb............"}, b.Paint(im, 23)))
}

func TestBounceAlwaysEaseOutHorizontalChildFits(t *testing.T) {
	b := Bounce{
		Child: render.Row{
			Children: []render.Widget{
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		Width:           15,
		BounceDirection: "horizontal",
		BounceAlways:    true,
		Curve:           EaseOut,
	}
	im := image.Rect(0, 0, 100, 100)

	assert.Equal(t, 24, b.FrameCount())
	assert.Equal(t, nil, render.CheckImage([]string{"rgb............"}, b.Paint(im, 0)))
	assert.Equal(t, nil, render.CheckImage([]string{".....rgb......."}, b.Paint(im, 1)))
	assert.Equal(t, nil, render.CheckImage([]string{".......rgb....."}, b.Paint(im, 2)))
	assert.Equal(t, nil, render.CheckImage([]string{"........rgb...."}, b.Paint(im, 3)))
	assert.Equal(t, nil, render.CheckImage([]string{".........rgb..."}, b.Paint(im, 4)))
	assert.Equal(t, nil, render.CheckImage([]string{"..........rgb.."}, b.Paint(im, 5)))
	assert.Equal(t, nil, render.CheckImage([]string{"...........rgb."}, b.Paint(im, 6)))
	assert.Equal(t, nil, render.CheckImage([]string{"...........rgb."}, b.Paint(im, 7)))
	assert.Equal(t, nil, render.CheckImage([]string{"...........rgb."}, b.Paint(im, 8)))
	assert.Equal(t, nil, render.CheckImage([]string{"............rgb"}, b.Paint(im, 9)))
	assert.Equal(t, nil, render.CheckImage([]string{"............rgb"}, b.Paint(im, 10)))
	assert.Equal(t, nil, render.CheckImage([]string{"............rgb"}, b.Paint(im, 11)))
	assert.Equal(t, nil, render.CheckImage([]string{"............rgb"}, b.Paint(im, 12)))
	assert.Equal(t, nil, render.CheckImage([]string{"............rgb"}, b.Paint(im, 13)))
	assert.Equal(t, nil, render.CheckImage([]string{"............rgb"}, b.Paint(im, 14)))
	assert.Equal(t, nil, render.CheckImage([]string{"............rgb"}, b.Paint(im, 15)))
	assert.Equal(t, nil, render.CheckImage([]string{"...........rgb."}, b.Paint(im, 16)))
	assert.Equal(t, nil, render.CheckImage([]string{"...........rgb."}, b.Paint(im, 17)))
	assert.Equal(t, nil, render.CheckImage([]string{"...........rgb."}, b.Paint(im, 18)))
	assert.Equal(t, nil, render.CheckImage([]string{"..........rgb.."}, b.Paint(im, 19)))
	assert.Equal(t, nil, render.CheckImage([]string{".........rgb..."}, b.Paint(im, 20)))
	assert.Equal(t, nil, render.CheckImage([]string{"........rgb...."}, b.Paint(im, 21)))
	assert.Equal(t, nil, render.CheckImage([]string{".......rgb....."}, b.Paint(im, 22)))
	assert.Equal(t, nil, render.CheckImage([]string{".....rgb......."}, b.Paint(im, 23)))
}

func TestBounceLinearHorizontal(t *testing.T) {
	b := Bounce{
		Child: render.Row{
			Children: []render.Widget{
				render.Box{Width: 4, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
				render.Box{Width: 4, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
				render.Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		Width:           6,
		BounceDirection: "horizontal",
		Curve:           LinearCurve{},
	}
	im := image.Rect(0, 0, 100, 100)

	assert.Equal(t, 12, b.FrameCount())
	assert.Equal(t, nil, render.CheckImage([]string{"rrrrgg"}, b.Paint(im, 0)))
	assert.Equal(t, nil, render.CheckImage([]string{"rrrggg"}, b.Paint(im, 1)))
	assert.Equal(t, nil, render.CheckImage([]string{"rrgggg"}, b.Paint(im, 2)))
	assert.Equal(t, nil, render.CheckImage([]string{"rggggb"}, b.Paint(im, 3)))
	assert.Equal(t, nil, render.CheckImage([]string{"ggggbb"}, b.Paint(im, 4)))
	assert.Equal(t, nil, render.CheckImage([]string{"gggbbb"}, b.Paint(im, 5)))
	assert.Equal(t, nil, render.CheckImage([]string{"ggbbbb"}, b.Paint(im, 6)))
	assert.Equal(t, nil, render.CheckImage([]string{"gggbbb"}, b.Paint(im, 7)))
	assert.Equal(t, nil, render.CheckImage([]string{"ggggbb"}, b.Paint(im, 8)))
	assert.Equal(t, nil, render.CheckImage([]string{"rggggb"}, b.Paint(im, 9)))
	assert.Equal(t, nil, render.CheckImage([]string{"rrgggg"}, b.Paint(im, 10)))
	assert.Equal(t, nil, render.CheckImage([]string{"rrrggg"}, b.Paint(im, 11)))

}

func TestBounceLinearVertical(t *testing.T) {
	b := Bounce{
		Child: render.Column{
			Children: []render.Widget{
				render.Box{Width: 1, Height: 4, Color: color.RGBA{0xff, 0, 0, 0xff}},
				render.Box{Width: 1, Height: 4, Color: color.RGBA{0, 0xff, 0, 0xff}},
				render.Box{Width: 1, Height: 4, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		Height:          6,
		BounceDirection: "vertical",
		Curve:           LinearCurve{},
	}
	im := image.Rect(0, 0, 100, 100)

	assert.Equal(t, 12, b.FrameCount())
	assert.Equal(t, nil, render.CheckImage([]string{"r", "r", "r", "r", "g", "g"}, b.Paint(im, 0)))
	assert.Equal(t, nil, render.CheckImage([]string{"r", "r", "r", "g", "g", "g"}, b.Paint(im, 1)))
	assert.Equal(t, nil, render.CheckImage([]string{"r", "r", "g", "g", "g", "g"}, b.Paint(im, 2)))
	assert.Equal(t, nil, render.CheckImage([]string{"r", "g", "g", "g", "g", "b"}, b.Paint(im, 3)))
	assert.Equal(t, nil, render.CheckImage([]string{"g", "g", "g", "g", "b", "b"}, b.Paint(im, 4)))
	assert.Equal(t, nil, render.CheckImage([]string{"g", "g", "g", "b", "b", "b"}, b.Paint(im, 5)))
	assert.Equal(t, nil, render.CheckImage([]string{"g", "g", "b", "b", "b", "b"}, b.Paint(im, 6)))
	assert.Equal(t, nil, render.CheckImage([]string{"g", "g", "g", "b", "b", "b"}, b.Paint(im, 7)))
	assert.Equal(t, nil, render.CheckImage([]string{"g", "g", "g", "g", "b", "b"}, b.Paint(im, 8)))
	assert.Equal(t, nil, render.CheckImage([]string{"r", "g", "g", "g", "g", "b"}, b.Paint(im, 9)))
	assert.Equal(t, nil, render.CheckImage([]string{"r", "r", "g", "g", "g", "g"}, b.Paint(im, 10)))
	assert.Equal(t, nil, render.CheckImage([]string{"r", "r", "r", "g", "g", "g"}, b.Paint(im, 11)))
}

func TestBounceAlwaysLinearHorizontalChildFitsAnimation(t *testing.T) {
	b := Bounce{
		Child: render.Animation{
			Children: []render.Widget{
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
				render.Box{Width: 1, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			},
		},
		Width:           6,
		BounceDirection: "horizontal",
		BounceAlways:    true,
		Curve:           LinearCurve{},
	}
	im := image.Rect(0, 0, 100, 100)

	assert.Equal(t, 10, b.FrameCount())
	assert.Equal(t, nil, render.CheckImage([]string{"r....."}, b.Paint(im, 0)))
	assert.Equal(t, nil, render.CheckImage([]string{".g...."}, b.Paint(im, 1)))
	assert.Equal(t, nil, render.CheckImage([]string{"..r..."}, b.Paint(im, 2)))
	assert.Equal(t, nil, render.CheckImage([]string{"...g.."}, b.Paint(im, 3)))
	assert.Equal(t, nil, render.CheckImage([]string{"....r."}, b.Paint(im, 4)))
	assert.Equal(t, nil, render.CheckImage([]string{".....g"}, b.Paint(im, 5)))
	assert.Equal(t, nil, render.CheckImage([]string{"....r."}, b.Paint(im, 6)))
	assert.Equal(t, nil, render.CheckImage([]string{"...g.."}, b.Paint(im, 7)))
	assert.Equal(t, nil, render.CheckImage([]string{"..r..."}, b.Paint(im, 8)))
	assert.Equal(t, nil, render.CheckImage([]string{".g...."}, b.Paint(im, 9)))
}

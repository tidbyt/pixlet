package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarqueeNoScrollHorizontal(t *testing.T) {
	m := Marquee{
		Width: 6,
		Child: Row{
			Children: []Widget{
				Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
				Box{Width: 2, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
				Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
	}

	mv := Marquee{
		Height: 3,
		Child: Row{
			Children: []Widget{
				Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
				Box{Width: 2, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
				Box{Width: 1, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
		ScrollDirection: "vertical",
	}

	// Child fits so there's just 1 single frame
	assert.Equal(t, 1, m.FrameCount())
	assert.Equal(t, 1, mv.FrameCount())
	im := m.Paint(image.Rect(0, 0, 100, 100), 0)
	imv := mv.Paint(image.Rect(0, 0, 100, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrggb",
		"rrrgg.",
		"rrr...",
	}, im))
	assert.Equal(t, nil, checkImage([]string{
		"rrrggb",
		"rrrgg.",
		"rrr...",
	}, imv))
}

// The addition of OffsetStart and OffsetEnd changes the default
// behaviour of Marquee. Passing start==width ann end==0 mimics the
// old default.
func TestMarqueeOldBehavior(t *testing.T) {
	m := Marquee{
		Width:       6,
		OffsetStart: 6,
		OffsetEnd:   0,
		Child: Row{
			Children: []Widget{
				Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
				Box{Width: 3, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
				Box{Width: 3, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
			},
		},
	}

	// The child's 9 pixels will be scrolled into view (7 frames),
	// scrolled out of view (9 frames) and then finally scrolled
	// back into view again (6 frames). 22 frames in total.
	assert.Equal(t, 22, m.FrameCount())

	// Scrolling into view
	assert.Equal(t, nil, checkImage([]string{
		"......",
		"......",
		"......",
	}, m.Paint(image.Rect(0, 0, 100, 100), 0)))

	assert.Equal(t, nil, checkImage([]string{
		"....rr",
		"....rr",
		"....rr",
	}, m.Paint(image.Rect(0, 0, 100, 100), 2)))

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, m.Paint(image.Rect(0, 0, 100, 100), 6)))

	// Scrolling out of view
	assert.Equal(t, nil, checkImage([]string{
		"rgggbb",
		"rggg..",
		"r.....",
	}, m.Paint(image.Rect(0, 0, 100, 100), 8)))

	assert.Equal(t, nil, checkImage([]string{
		"b.....",
		"......",
		"......",
	}, m.Paint(image.Rect(0, 0, 100, 100), 14)))

	assert.Equal(t, nil, checkImage([]string{
		"......",
		"......",
		"......",
	}, m.Paint(image.Rect(0, 0, 100, 100), 15)))

	// Scrolling back into view
	assert.Equal(t, nil, checkImage([]string{
		"...rrr",
		"...rrr",
		"...rrr",
	}, m.Paint(image.Rect(0, 0, 100, 100), 18)))

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, m.Paint(image.Rect(0, 0, 100, 100), 21)))

	// Later frames keep it fixed in the last frame. This makes
	// multiple simultaneous marquees look nice when they've
	// different length.

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, m.Paint(image.Rect(0, 0, 100, 100), 22)))

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, m.Paint(image.Rect(0, 0, 100, 100), 26)))

	assert.Equal(t, nil, checkImage([]string{
		"rrrggg",
		"rrrggg",
		"rrr...",
	}, m.Paint(image.Rect(0, 0, 100, 100), 100000)))
}

func TestMarqueeOffsetStart(t *testing.T) {
	child := Row{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	m := Marquee{
		Width: 6,
		Child: child,
	}
	im := image.Rect(0, 0, 100, 100)

	// OffsetStart affects the initial position of the child
	m.OffsetStart = 2
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, m.Paint(im, 1)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, m.Paint(im, 2)))
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, m.Paint(im, 3)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, m.Paint(im, 4)))
	assert.Equal(t, nil, checkImage([]string{"bbbb.."}, m.Paint(im, 5)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, m.Paint(im, 6)))
	assert.Equal(t, nil, checkImage([]string{"bb...."}, m.Paint(im, 7)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, m.Paint(im, 8)))
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 9)))

	assert.Equal(t, nil, checkImage([]string{".....r"}, m.Paint(im, 10)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, m.Paint(im, 11)))
	assert.Equal(t, nil, checkImage([]string{"...rgg"}, m.Paint(im, 12)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, m.Paint(im, 13)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, m.Paint(im, 14)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, m.Paint(im, 15)))
	assert.Equal(t, 16, m.FrameCount())

	// Negative OffsetStart
	m.OffsetStart = -2
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{"bbbb.."}, m.Paint(im, 1)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, m.Paint(im, 2)))
	assert.Equal(t, nil, checkImage([]string{"bb...."}, m.Paint(im, 3)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, m.Paint(im, 4)))
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 5)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, m.Paint(im, 6)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, m.Paint(im, 7)))
	assert.Equal(t, nil, checkImage([]string{"...rgg"}, m.Paint(im, 8)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, m.Paint(im, 9)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, m.Paint(im, 10)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, m.Paint(im, 11)))
	assert.Equal(t, 12, m.FrameCount())

	// Overly negative OffsetStart is truncated to child width
	m.OffsetStart = -1000
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, m.Paint(im, 1)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, m.Paint(im, 2)))
	assert.Equal(t, 7, m.FrameCount())
	m.OffsetStart = -7
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, m.Paint(im, 1)))
	assert.Equal(t, 7, m.FrameCount())
	m.OffsetStart = -8
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, m.Paint(im, 1)))
	assert.Equal(t, 7, m.FrameCount())
	m.OffsetStart = -6
	assert.Equal(t, nil, checkImage([]string{"b....."}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 1)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, m.Paint(im, 2)))
	assert.Equal(t, 8, m.FrameCount())
}

func TestMarqueeOffsetEnd(t *testing.T) {
	child := Row{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 2, Height: 1, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 4, Height: 1, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	m := Marquee{
		Width: 6,
		Child: child,
	}
	im := image.Rect(0, 0, 100, 100)

	// OffsetEnd affects the final position of the child
	m.OffsetEnd = 2
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, m.Paint(im, 1)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, m.Paint(im, 2)))
	assert.Equal(t, nil, checkImage([]string{"bbbb.."}, m.Paint(im, 3)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, m.Paint(im, 4)))
	assert.Equal(t, nil, checkImage([]string{"bb...."}, m.Paint(im, 5)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, m.Paint(im, 6)))
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 7)))
	assert.Equal(t, nil, checkImage([]string{".....r"}, m.Paint(im, 8)))
	assert.Equal(t, nil, checkImage([]string{"....rg"}, m.Paint(im, 9)))
	assert.Equal(t, nil, checkImage([]string{"...rgg"}, m.Paint(im, 10)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, m.Paint(im, 11)))
	assert.Equal(t, 12, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, m.Paint(im, 12)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, m.Paint(im, 13)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, m.Paint(im, 1024)))

	// Negative offset places child outside of marquee
	m.OffsetEnd = -4
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, m.Paint(im, 1)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, m.Paint(im, 2)))
	assert.Equal(t, nil, checkImage([]string{"...rgg"}, m.Paint(im, 10)))
	assert.Equal(t, nil, checkImage([]string{"..rggb"}, m.Paint(im, 11)))
	assert.Equal(t, nil, checkImage([]string{".rggbb"}, m.Paint(im, 12)))
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, m.Paint(im, 13)))
	assert.Equal(t, nil, checkImage([]string{"ggbbbb"}, m.Paint(im, 14)))
	assert.Equal(t, nil, checkImage([]string{"gbbbb."}, m.Paint(im, 15)))
	assert.Equal(t, nil, checkImage([]string{"bbbb.."}, m.Paint(im, 16)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, m.Paint(im, 17)))
	assert.Equal(t, 18, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, m.Paint(im, 18)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, m.Paint(im, 19)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, m.Paint(im, 1024)))

	// Very negative offset is truncated to width of child
	m.OffsetEnd = -133
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{"bbb..."}, m.Paint(im, 17)))
	assert.Equal(t, nil, checkImage([]string{"bb...."}, m.Paint(im, 18)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, m.Paint(im, 19)))
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 20)))
	assert.Equal(t, 21, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 21)))
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 22)))
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 23)))

	// OffsetEnd >= width means it doesn't scroll back
	m.OffsetEnd = 6
	assert.Equal(t, nil, checkImage([]string{"rggbbb"}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{"b....."}, m.Paint(im, 6)))
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 7)))
	assert.Equal(t, 8, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 8)))
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 9)))
	assert.Equal(t, nil, checkImage([]string{"......"}, m.Paint(im, 1024)))

}

func TestMarqueeVerticalScroll(t *testing.T) {
	child := Column{
		Children: []Widget{
			Box{Width: 1, Height: 1, Color: color.RGBA{0xff, 0, 0, 0xff}},
			Box{Width: 1, Height: 2, Color: color.RGBA{0, 0xff, 0, 0xff}},
			Box{Width: 1, Height: 4, Color: color.RGBA{0, 0, 0xff, 0xff}},
		},
	}
	m := Marquee{
		Height:          6,
		Child:           child,
		ScrollDirection: "vertical",
	}
	im := image.Rect(0, 0, 100, 100)

	// OffsetEnd affects the final position of the child
	m.OffsetStart = 2
	assert.Equal(t, nil, checkImage([]string{
		".",
		".",
		"r",
		"g",
		"g",
		"b",
	}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{
		".",
		"r",
		"g",
		"g",
		"b",
		"b",
	}, m.Paint(im, 1)))
	assert.Equal(t, nil, checkImage([]string{
		"r",
		"g",
		"g",
		"b",
		"b",
		"b",
	}, m.Paint(im, 2)))
	assert.Equal(t, nil, checkImage([]string{
		"g",
		"g",
		"b",
		"b",
		"b",
		"b",
	}, m.Paint(im, 3)))
	assert.Equal(t, nil, checkImage([]string{
		"g",
		"b",
		"b",
		"b",
		"b",
		".",
	}, m.Paint(im, 4)))
	assert.Equal(t, nil, checkImage([]string{
		"b",
		"b",
		"b",
		"b",
		".",
		".",
	}, m.Paint(im, 5)))
	assert.Equal(t, nil, checkImage([]string{
		"b",
		"b",
		"b",
		".",
		".",
		".",
	}, m.Paint(im, 6)))
	assert.Equal(t, nil, checkImage([]string{
		"b",
		"b",
		".",
		".",
		".",
		".",
	}, m.Paint(im, 7)))
	assert.Equal(t, nil, checkImage([]string{
		"b",
		".",
		".",
		".",
		".",
		".",
	}, m.Paint(im, 8)))
	assert.Equal(t, nil, checkImage([]string{
		".",
		".",
		".",
		".",
		".",
		".",
	}, m.Paint(im, 9)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, m.Paint(im, 10)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", "r", "g"}, m.Paint(im, 11)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", "r", "g", "g"}, m.Paint(im, 12)))
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, m.Paint(im, 13)))
	assert.Equal(t, nil, checkImage([]string{".", "r", "g", "g", "b", "b"}, m.Paint(im, 14)))
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, m.Paint(im, 15)))
	assert.Equal(t, 16, m.FrameCount())

	// Negative OffsetStart
	m.OffsetStart = -2
	assert.Equal(t, nil, checkImage([]string{"g", "b", "b", "b", "b", "."}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", "b", ".", "."}, m.Paint(im, 1)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", ".", ".", "."}, m.Paint(im, 2)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", ".", ".", ".", "."}, m.Paint(im, 3)))
	assert.Equal(t, nil, checkImage([]string{"b", ".", ".", ".", ".", "."}, m.Paint(im, 4)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, m.Paint(im, 5)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, m.Paint(im, 6)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", "r", "g"}, m.Paint(im, 7)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", "r", "g", "g"}, m.Paint(im, 8)))
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, m.Paint(im, 9)))
	assert.Equal(t, nil, checkImage([]string{".", "r", "g", "g", "b", "b"}, m.Paint(im, 10)))
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, m.Paint(im, 11)))
	assert.Equal(t, 12, m.FrameCount())

	// Overly negative OffsetStart is truncated to child width
	m.OffsetStart = -1000
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, m.Paint(im, 1)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", "r", "g"}, m.Paint(im, 2)))
	assert.Equal(t, 7, m.FrameCount())
	m.OffsetStart = -7
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, m.Paint(im, 1)))
	assert.Equal(t, 7, m.FrameCount())
	m.OffsetStart = -8
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, m.Paint(im, 1)))
	assert.Equal(t, 7, m.FrameCount())
	m.OffsetStart = -6
	assert.Equal(t, nil, checkImage([]string{"b", ".", ".", ".", ".", "."}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, m.Paint(im, 1)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, m.Paint(im, 2)))
	assert.Equal(t, 8, m.FrameCount())

	// OffsetEnd affects the final position of the child
	m.OffsetStart = 0
	m.OffsetEnd = 2
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{"g", "g", "b", "b", "b", "b"}, m.Paint(im, 1)))
	assert.Equal(t, nil, checkImage([]string{"g", "b", "b", "b", "b", "."}, m.Paint(im, 2)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", "b", ".", "."}, m.Paint(im, 3)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", ".", ".", "."}, m.Paint(im, 4)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", ".", ".", ".", "."}, m.Paint(im, 5)))
	assert.Equal(t, nil, checkImage([]string{"b", ".", ".", ".", ".", "."}, m.Paint(im, 6)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, m.Paint(im, 7)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "r"}, m.Paint(im, 8)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", "r", "g"}, m.Paint(im, 9)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", "r", "g", "g"}, m.Paint(im, 10)))
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, m.Paint(im, 11)))
	assert.Equal(t, 12, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, m.Paint(im, 12)))
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, m.Paint(im, 13)))
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, m.Paint(im, 1024)))

	// Negative offset places child outside of marquee
	m.OffsetEnd = -4
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{"g", "g", "b", "b", "b", "b"}, m.Paint(im, 1)))
	assert.Equal(t, nil, checkImage([]string{"g", "b", "b", "b", "b", "."}, m.Paint(im, 2)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", "r", "g", "g"}, m.Paint(im, 10)))
	assert.Equal(t, nil, checkImage([]string{".", ".", "r", "g", "g", "b"}, m.Paint(im, 11)))
	assert.Equal(t, nil, checkImage([]string{".", "r", "g", "g", "b", "b"}, m.Paint(im, 12)))
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, m.Paint(im, 13)))
	assert.Equal(t, nil, checkImage([]string{"g", "g", "b", "b", "b", "b"}, m.Paint(im, 14)))
	assert.Equal(t, nil, checkImage([]string{"g", "b", "b", "b", "b", "."}, m.Paint(im, 15)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", "b", ".", "."}, m.Paint(im, 16)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", ".", ".", "."}, m.Paint(im, 17)))
	assert.Equal(t, 18, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", ".", ".", "."}, m.Paint(im, 18)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", ".", ".", "."}, m.Paint(im, 19)))
	assert.Equal(t, nil, checkImage([]string{"b", "b", "b", ".", ".", "."}, m.Paint(im, 1024)))

	// OffsetEnd >= width means it doesn't scroll back
	m.OffsetEnd = 6
	assert.Equal(t, nil, checkImage([]string{"r", "g", "g", "b", "b", "b"}, m.Paint(im, 0)))
	assert.Equal(t, nil, checkImage([]string{"b", ".", ".", ".", ".", "."}, m.Paint(im, 6)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, m.Paint(im, 7)))
	assert.Equal(t, 8, m.FrameCount())
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, m.Paint(im, 8)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, m.Paint(im, 9)))
	assert.Equal(t, nil, checkImage([]string{".", ".", ".", ".", ".", "."}, m.Paint(im, 1024)))
}

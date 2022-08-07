package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextDefault(t *testing.T) {
	text := &Text{Content: "A"}
	text.Init()
	im := PaintWidget(text, image.Rect(0, 0, 0, 0), 0)
	assert.Equal(t, nil, checkImage([]string{
		".....",
		".ww..",
		"w..w.",
		"w..w.",
		"wwww.",
		"w..w.",
		"w..w.",
		".....",
	}, im))
	w, h := text.Size()
	assert.Equal(t, 5, w)
	assert.Equal(t, 8, h)

	text = &Text{Content: "j!ÑÖ"}
	text.Init()
	im = PaintWidget(text, image.Rect(0, 0, 0, 0), 0)
	assert.Equal(t, nil, checkImage([]string{
		"...." + ".." + ".w.w." + "w..w.",
		"..w." + "w." + "w.w.." + ".....",
		"...." + "w." + "w..w." + ".ww..",
		"..w." + "w." + "ww.w." + "w..w.",
		"..w." + "w." + "w.ww." + "w..w.",
		"..w." + ".." + "w..w." + "w..w.",
		"w.w." + "w." + "w..w." + ".ww..",
		".w.." + ".." + "....." + ".....",
	}, im))
	w, h = text.Size()
	assert.Equal(t, 16, w)
	assert.Equal(t, 8, h)
}

func TestTextParameters(t *testing.T) {
	text := &Text{
		Content: "ᚠӠ",
		Font:    "6x13",
	}
	text.Init()

	im := PaintWidget(text, image.Rect(0, 0, 0, 0), 0)
	assert.Equal(t, nil, checkImage([]string{
		"......" + "......",
		"......" + "......",
		".w..w." + "wwwww.",
		".w.w.." + "....w.",
		".ww..." + "...w..",
		".w..w." + "..w...",
		".w.w.." + ".www..",
		".ww..." + "....w.",
		".w...." + "....w.",
		".w...." + "w...w.",
		".w...." + ".www..",
		"......" + "......",
		"......" + "......",
	}, im))
	w, h := text.Size()
	assert.Equal(t, 12, w)
	assert.Equal(t, 13, h)

	// Green text, pushed down 2 pixels and capped at height 10
	text = &Text{
		Content: "ᚠӠ",
		Font:    "6x13",
		Color:   color.RGBA{0, 0xff, 0, 0xff},
		Offset:  -2,
		Height:  10,
	}
	text.Init()
	im = PaintWidget(text, image.Rect(0, 0, 0, 0), 0)
	assert.Equal(t, nil, checkImage([]string{
		"......" + "......",
		".g..g." + "ggggg.",
		".g.g.." + "....g.",
		".gg..." + "...g..",
		".g..g." + "..g...",
		".g.g.." + ".ggg..",
		".gg..." + "....g.",
		".g...." + "....g.",
		".g...." + "g...g.",
		".g...." + ".ggg..",
	}, im))
	w, h = text.Size()
	assert.Equal(t, w, 12)
	assert.Equal(t, h, 10)
}

// Make sure the fonts render as expected
func TestTextFonts(t *testing.T) {

	// Content is chosen to extend below baseline, above cap
	// height and to include variable width characters

	text := &Text{
		Content: "QqÖ!",
		Font:    "6x13",
	}
	text.Init()

	im := PaintWidget(text, image.Rect(0, 0, 0, 0), 0)
	assert.Equal(t, nil, checkImage([]string{
		"......" + "......" + "......" + "......",
		"......" + "......" + ".w.w.." + "......",
		".www.." + "......" + ".w.w.." + "..w...",
		"w...w." + "......" + "......" + "..w...",
		"w...w." + "......" + ".www.." + "..w...",
		"w...w." + ".wwww." + "w...w." + "..w...",
		"w...w." + "w...w." + "w...w." + "..w...",
		"w...w." + "w...w." + "w...w." + "..w...",
		"w...w." + "w...w." + "w...w." + "..w...",
		"w.w.w." + ".wwww." + "w...w." + "......",
		".www.." + "....w." + ".www.." + "..w...",
		"....w." + "....w." + "......" + "......",
		"......" + "....w." + "......" + "......",
	}, im))
	w, h := text.Size()
	assert.Equal(t, 6*4, w)
	assert.Equal(t, 13, h)

	text = &Text{
		Content: "QqÖ!",
		Font:    "Dina_r400-6",
	}
	text.Init()
	im = PaintWidget(text, image.Rect(0, 0, 0, 0), 0)
	assert.Equal(t, nil, checkImage([]string{
		"......" + "......" + ".w.w.." + "......",
		"......" + "......" + "......" + "......",
		".www.." + "......" + ".www.." + "..w...",
		"w...w." + ".wwww." + "w...w." + "..w...",
		"w...w." + "w...w." + "w...w." + "..w...",
		"w...w." + "w...w." + "w...w." + "..w...",
		"w..w.." + "w...w." + "w...w." + "......",
		".ww.w." + ".wwww." + ".www.." + "..w...",
		"....w." + "....w." + "......" + "......",
		"......" + "....w." + "......" + "......",
	}, im))
	w, h = text.Size()
	assert.Equal(t, 6*4, w)
	assert.Equal(t, 10, h)

	text = &Text{
		Content: "QqÖ!",
		Font:    "5x8",
	}
	text.Init()
	im = PaintWidget(text, image.Rect(0, 0, 0, 0), 0)
	assert.Equal(t, nil, checkImage([]string{
		"....." + "....." + "w..w." + ".....",
		".ww.." + "....." + "....." + "..w..",
		"w..w." + "....." + ".ww.." + "..w..",
		"w..w." + ".www." + "w..w." + "..w..",
		"ww.w." + "w..w." + "w..w." + "..w..",
		"w.ww." + ".www." + "w..w." + ".....",
		".ww.." + "...w." + ".ww.." + "..w..",
		"...w." + "...w." + "....." + ".....",
	}, im))
	w, h = text.Size()
	assert.Equal(t, 5*4, w)
	assert.Equal(t, 8, h)

	text = &Text{
		Content: "QqÖ!",
		Font:    "tb-8",
	}
	text.Init()
	im = PaintWidget(text, image.Rect(0, 0, 0, 0), 0)
	assert.Equal(t, nil, checkImage([]string{
		"....." + "....." + "w..w." + "..",
		".ww.." + "....." + "....." + "w.",
		"w..w." + "....." + ".ww.." + "w.",
		"w..w." + ".www." + "w..w." + "w.",
		"ww.w." + "w..w." + "w..w." + "w.",
		"w.ww." + ".www." + "w..w." + "..",
		".ww.." + "...w." + ".ww.." + "w.",
		"...w." + "...w." + "....." + "..",
	}, im))
	w, h = text.Size()
	assert.Equal(t, 17, w)
	assert.Equal(t, 8, h)
}

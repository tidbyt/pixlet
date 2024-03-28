package render

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrappedTextWithBounds(t *testing.T) {
	text := &WrappedText{Content: "AB CD."}
	assert.NoError(t, text.Init())

	// Sufficient space to fit on single line
	im := PaintWidget(text, image.Rect(0, 0, 25, 8), 0)
	assert.Equal(t, nil, checkImage([]string{
		"....." + "........" + "....." + ".......",
		".ww.." + "www....." + ".ww.." + "www....",
		"w..w." + "w..w...." + "w..w." + "w..w...",
		"w..w." + "www....." + "w...." + "w..w...",
		"wwww." + "w..w...." + "w...." + "w..w...",
		"w..w." + "w..w...." + "w..w." + "w..w...",
		"w..w." + "www....." + ".ww.." + "www..w.",
		"....." + "........" + "....." + ".......",
	}, im))

	// Reduce avaialable width and it wraps
	im = PaintWidget(text, image.Rect(0, 0, 21, 16), 0)
	assert.Equal(t, nil, checkImage([]string{
		"....." + ".......",
		".ww.." + "www....",
		"w..w." + "w..w...",
		"w..w." + "www....",
		"wwww." + "w..w...",
		"w..w." + "w..w...",
		"w..w." + "www....",
		"....." + ".......",
		"....." + ".......",
		".ww.." + "www....",
		"w..w." + "w..w...",
		"w...." + "w..w...",
		"w...." + "w..w...",
		"w..w." + "w..w...",
		".ww.." + "www..w.",
		"....." + ".......",
	}, im))

	// Overflow is cut off
	im = PaintWidget(text, image.Rect(0, 0, 7, 12), 0)
	assert.Equal(t, nil, checkImage([]string{
		"....." + "..",
		".ww.." + "ww",
		"w..w." + "w.",
		"w..w." + "ww",
		"wwww." + "w.",
		"w..w." + "w.",
		"w..w." + "ww",
		"....." + "..",
		"....." + "..",
		".ww.." + "ww",
		"w..w." + "w.",
		"w...." + "w.",
	}, im))
}

func TestWrappedTextWithsize(t *testing.T) {
	// Weight and Height parameters override the bounds
	text := &WrappedText{Content: "AB CD.", Width: 7, Height: 12}
	assert.NoError(t, text.Init())
	im := PaintWidget(text, image.Rect(0, 0, 40, 40), 0)
	assert.Equal(t, nil, checkImage([]string{
		"....." + "..",
		".ww.." + "ww",
		"w..w." + "w.",
		"w..w." + "ww",
		"wwww." + "w.",
		"w..w." + "w.",
		"w..w." + "ww",
		"....." + "..",
		"....." + "..",
		".ww.." + "ww",
		"w..w." + "w.",
		"w...." + "w.",
	}, im))

	// Height can be overridden separately
	text = &WrappedText{Content: "AB CD.", Height: 12}
	assert.NoError(t, text.Init())
	im = PaintWidget(text, image.Rect(0, 0, 9, 40), 0)
	assert.Equal(t, nil, checkImage([]string{
		"....." + "....",
		".ww.." + "www.",
		"w..w." + "w..w",
		"w..w." + "www.",
		"wwww." + "w..w",
		"w..w." + "w..w",
		"w..w." + "www.",
		"....." + "....",
		"....." + "....",
		".ww.." + "www.",
		"w..w." + "w..w",
		"w...." + "w..w",
	}, im))

	// Ditto for Width
	text = &WrappedText{Content: "AB CD.", Width: 3}
	assert.NoError(t, text.Init())
	im = PaintWidget(text, image.Rect(0, 0, 9, 5), 0)
	assert.Equal(t, nil, checkImage([]string{
		"...",
		".ww",
		"w..",
		"w..",
		"www",
	}, im))
}

func TestWrappedTextLineSpacing(t *testing.T) {

	// Single pixel line space
	text := &WrappedText{Content: "AB CD.", LineSpacing: 1}
	assert.NoError(t, text.Init())
	im := PaintWidget(text, image.Rect(0, 0, 21, 16), 0)
	assert.Equal(t, nil, checkImage([]string{
		"....." + ".......",
		".ww.." + "www....",
		"w..w." + "w..w...",
		"w..w." + "www....",
		"wwww." + "w..w...",
		"w..w." + "w..w...",
		"w..w." + "www....",
		"....." + ".......",
		"....." + ".......", // extra line
		"....." + ".......",
		".ww.." + "www....",
		"w..w." + "w..w...",
		"w...." + "w..w...",
		"w...." + "w..w...",
		"w..w." + "w..w...",
		".ww.." + "www..w.",
	}, im))

	// Add another one
	text = &WrappedText{Content: "AB CD.", LineSpacing: 2}
	assert.NoError(t, text.Init())
	im = PaintWidget(text, image.Rect(0, 0, 21, 16), 0)
	assert.Equal(t, nil, checkImage([]string{
		"....." + ".......",
		".ww.." + "www....",
		"w..w." + "w..w...",
		"w..w." + "www....",
		"wwww." + "w..w...",
		"w..w." + "w..w...",
		"w..w." + "www....",
		"....." + ".......",
		"....." + ".......", // extra line
		"....." + ".......", // and here
		"....." + ".......",
		".ww.." + "www....",
		"w..w." + "w..w...",
		"w...." + "w..w...",
		"w...." + "w..w...",
		"w..w." + "w..w...",
	}, im))
}

func TestWrappedTextAlignment(t *testing.T) {
	// Default to left align.
	text := &WrappedText{Content: "AB CD."}
	assert.NoError(t, text.Init())
	im := PaintWidget(text, image.Rect(0, 0, 21, 16), 0)
	assert.Equal(t, nil, checkImage([]string{
		"......." + ".....",
		".ww..ww" + "w....",
		"w..w.w." + ".w...",
		"w..w.ww" + "w....",
		"wwww.w." + ".w...",
		"w..w.w." + ".w...",
		"w..w.ww" + "w....",
		"......." + ".....",
		"......." + ".....",
		".ww..ww" + "w....",
		"w..w.w." + ".w...",
		"w....w." + ".w...",
		"w....w." + ".w...",
		"w..w.w." + ".w...",
		".ww..ww" + "w..w.",
		"......." + ".....",
	}, im))

	// Right alignment.
	text = &WrappedText{Content: "AB CD.", Align: "right"}
	assert.NoError(t, text.Init())
	im = PaintWidget(text, image.Rect(0, 0, 21, 16), 0)
	assert.Equal(t, nil, checkImage([]string{
		"......." + ".....",
		"...ww.." + "www..",
		"..w..w." + "w..w.",
		"..w..w." + "www..",
		"..wwww." + "w..w.",
		"..w..w." + "w..w.",
		"..w..w." + "www..",
		"......." + ".....",
		"......." + ".....",
		".ww..ww" + "w....",
		"w..w.w." + ".w...",
		"w....w." + ".w...",
		"w....w." + ".w...",
		"w..w.w." + ".w...",
		".ww..ww" + "w..w.",
		"......." + ".....",
	}, im))

	// Center alignment.
	text = &WrappedText{Content: "AB CD.", Align: "center"}
	assert.NoError(t, text.Init())
	im = PaintWidget(text, image.Rect(0, 0, 21, 16), 0)
	assert.Equal(t, nil, checkImage([]string{
		"......." + ".....",
		"..ww..w" + "ww...",
		".w..w.w" + "..w..",
		".w..w.w" + "ww...",
		".wwww.w" + "..w..",
		".w..w.w" + "..w..",
		".w..w.w" + "ww...",
		"......." + ".....",
		"......." + ".....",
		".ww..ww" + "w....",
		"w..w.w." + ".w...",
		"w....w." + ".w...",
		"w....w." + ".w...",
		"w..w.w." + ".w...",
		".ww..ww" + "w..w.",
		"......." + ".....",
	}, im))
}

func TestWrappedTextMissingFont(t *testing.T) {
	text := &WrappedText{Content: "AB CD.", Font: "missing"}
	assert.Error(t, text.Init())
}

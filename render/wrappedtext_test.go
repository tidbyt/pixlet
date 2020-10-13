package render

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrappedTextWithBounds(t *testing.T) {
	text := WrappedText{Content: "AB CD."}

	// Sufficient space to fit on single line
	im := text.Paint(image.Rect(0, 0, 25, 8), 0)
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
	im = text.Paint(image.Rect(0, 0, 21, 16), 0)
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
	im = text.Paint(image.Rect(0, 0, 7, 12), 0)
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

func TestWrappedTextWithSize(t *testing.T) {
	// Weight and Height parameters override the bounds
	text := WrappedText{Content: "AB CD.", Width: 7, Height: 12}
	im := text.Paint(image.Rect(0, 0, 40, 40), 0)
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
	text = WrappedText{Content: "AB CD.", Height: 12}
	im = text.Paint(image.Rect(0, 0, 9, 40), 0)
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
	text = WrappedText{Content: "AB CD.", Width: 3}
	im = text.Paint(image.Rect(0, 0, 9, 5), 0)
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
	text := WrappedText{Content: "AB CD.", LineSpacing: 1}
	im := text.Paint(image.Rect(0, 0, 21, 17), 0)
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
	text = WrappedText{Content: "AB CD.", LineSpacing: 2}
	im = text.Paint(image.Rect(0, 0, 21, 17), 0)
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
		"w..w." + "w..w...", // truncation here
	}, im))
}

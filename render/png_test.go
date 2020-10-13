package render

import (
	"encoding/base64"
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

// A base4 encoded PNG depicting a red rectangle with a centered red
// plus sign on a transparent background.
const testPNG = "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAMCAYAAABbayygAAAAOUlEQVQoU2P8z8Dwn4EIwAhSyMjAwIhPLVgNukKYDciaaawQl6dATkCxmmiFMF8PgGeICnAiYpABACrQO/WD80OVAAAAAElFTkSuQmCC"

func TestPNG(t *testing.T) {
	raw, _ := base64.StdEncoding.DecodeString(testPNG)
	png := &PNG{Src: string(raw)}

	// Size of PNG is independent of bounds
	im := png.Paint(image.Rect(0, 0, 0, 0), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrrrrrrr",
		"r........r",
		"r...rr...r",
		"r...rr...r",
		"r...rr...r",
		"r.rrrrrr.r",
		"r.rrrrrr.r",
		"r...rr...r",
		"r...rr...r",
		"r...rr...r",
		"r........r",
		"rrrrrrrrrr",
	}, im))
	w, h := png.Size()
	assert.Equal(t, 10, w)
	assert.Equal(t, 12, h)

	im = png.Paint(image.Rect(0, 0, 100, 100), 0)
	assert.Equal(t, nil, checkImage([]string{
		"rrrrrrrrrr",
		"r........r",
		"r...rr...r",
		"r...rr...r",
		"r...rr...r",
		"r.rrrrrr.r",
		"r.rrrrrr.r",
		"r...rr...r",
		"r...rr...r",
		"r...rr...r",
		"r........r",
		"rrrrrrrrrr",
	}, im))
	w, h = png.Size()
	assert.Equal(t, 10, w)
	assert.Equal(t, 12, h)
}

// Check that scaled image is scaled, but don't bother checking
// individual pixels.
func TestPNGScale(t *testing.T) {
	raw, _ := base64.StdEncoding.DecodeString(testPNG)
	png := &PNG{Src: string(raw), Width: 5, Height: 6}

	w, h := png.Size()
	assert.Equal(t, 5, w)
	assert.Equal(t, 6, h)
	png.Paint(image.Rect(0, 0, 0, 0), 0)
}

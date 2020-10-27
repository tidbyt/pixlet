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

func TestImage(t *testing.T) {
	raw, _ := base64.StdEncoding.DecodeString(testPNG)
	img := &Image{Src: string(raw)}

	// Size of Image is independent of bounds
	im := img.Paint(image.Rect(0, 0, 0, 0), 0)
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
	w, h := img.Size()
	assert.Equal(t, 10, w)
	assert.Equal(t, 12, h)

	im = img.Paint(image.Rect(0, 0, 100, 100), 0)
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
	w, h = img.Size()
	assert.Equal(t, 10, w)
	assert.Equal(t, 12, h)
}

// Check that scaled image is scaled, but don't bother checking
// individual pixels.
func TestImageScale(t *testing.T) {
	raw, _ := base64.StdEncoding.DecodeString(testPNG)
	img := &Image{Src: string(raw), Width: 5, Height: 6}

	w, h := img.Size()
	assert.Equal(t, 5, w)
	assert.Equal(t, 6, h)
	im := img.Paint(image.Rect(0, 0, 0, 0), 0)
	assert.Equal(t, 5, im.Bounds().Dx())
	assert.Equal(t, 6, im.Bounds().Dy())
}

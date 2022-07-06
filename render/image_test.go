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
	img.Init()

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
	img.Init()

	w, h := img.Size()
	assert.Equal(t, 5, w)
	assert.Equal(t, 6, h)
	im := img.Paint(image.Rect(0, 0, 0, 0), 0)
	assert.Equal(t, 5, im.Bounds().Dx())
	assert.Equal(t, 6, im.Bounds().Dy())
}

// Check that scaled image is scaled
// maintaining aspect ratio when only width is provided
// but don't bother checking individual pixels.
func TestImageScaleAspectRatioWidth(t *testing.T) {
	raw, _ := base64.StdEncoding.DecodeString(testPNG)
	img := &Image{Src: string(raw), Width: 5}
	img.Init()

	w, h := img.Size()
	assert.Equal(t, 5, w)
	assert.Equal(t, 6, h)
	im := img.Paint(image.Rect(0, 0, 0, 0), 0)
	assert.Equal(t, 5, im.Bounds().Dx())
	assert.Equal(t, 6, im.Bounds().Dy())
}

// Check that scaled image is scaled
// maintaining aspect ratio when only height is provided
// but don't bother checking individual pixels.
func TestImageScaleAspectRatioHeight(t *testing.T) {
	raw, _ := base64.StdEncoding.DecodeString(testPNG)
	img := &Image{Src: string(raw), Height: 6}
	img.Init()

	w, h := img.Size()
	assert.Equal(t, 5, w)
	assert.Equal(t, 6, h)
	im := img.Paint(image.Rect(0, 0, 0, 0), 0)
	assert.Equal(t, 5, im.Bounds().Dx())
	assert.Equal(t, 6, im.Bounds().Dy())
}

func TestImageAnimatedGif(t *testing.T) {
	// Animated 5x4 GIF with 4 frames:
	//
	// frame 0: ..x..
	//          x....
	//          .x...
	//          ...x.
	//
	// Subsequent frames shift pixels right, overflowing into the
	// next row.
	//
	// GIF has no disposal method set, and a delay of 1230 ms

	const testGIF = "R0lGODlhBQAEAPAAAAAAAAAAACH5BAF7AAAAIf8LTkVUU0NBUEUyLjADAQAAACwAAAAABQAEAAACBgRiaLmLBQAh+QQBewAAACwAAAAABQAEAAACBYRzpqhXACH5BAF7AAAALAAAAAAFAAQAAAIGDG6Qp8wFACH5BAF7AAAALAAAAAAFAAQAAAIGRIBnyMoFADs="

	raw, _ := base64.StdEncoding.DecodeString(testGIF)
	img := &Image{Src: string(raw)}
	img.Init()

	w, h := img.Size()
	assert.Equal(t, 5, w)
	assert.Equal(t, 4, h)
	assert.Equal(t, 1230, img.Delay)

	// 4 frames in this animation
	assert.Equal(t, 4, img.FrameCount())

	// black pixels moving right
	assert.Equal(t, nil, checkImage([]string{
		"..x..",
		"x....",
		".x...",
		"...x.",
	}, img.Paint(image.Rect(0, 0, 100, 100), 0)))

	// since no disposal method is set, subsequent frames should
	// draw on top of first frame
	assert.Equal(t, nil, checkImage([]string{
		"..xx.",
		"xx...",
		".xx..",
		"...xx",
	}, img.Paint(image.Rect(0, 0, 100, 100), 1)))
	assert.Equal(t, nil, checkImage([]string{
		"x.xxx",
		"xxx..",
		".xxx.",
		"...xx",
	}, img.Paint(image.Rect(0, 0, 100, 100), 2)))
	assert.Equal(t, nil, checkImage([]string{
		"xxxxx",
		"xxxx.",
		".xxxx",
		"...xx",
	}, img.Paint(image.Rect(0, 0, 100, 100), 3)))

	// loops after the last frame
	assert.Equal(t, nil, checkImage([]string{
		"..x..",
		"x....",
		".x...",
		"...x.",
	}, img.Paint(image.Rect(0, 0, 100, 100), 4)))
	assert.Equal(t, nil, checkImage([]string{
		"..xx.",
		"xx...",
		".xx..",
		"...xx",
	}, img.Paint(image.Rect(0, 0, 100, 100), 5)))
}

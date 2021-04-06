package render

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"

	// register image formats
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"github.com/nfnt/resize"
)

// Image renders the binary image data passed via `src`. Supported
// formats include PNG, JPEG and GIF.
//
// If `width` or `height` are set, the image will be scaled
// accordingly, with nearest neighbor interpolation. Otherwise the
// image's original dimensions are used.
//
// If the image data encodes an animated GIF, the Image instance will
// also be animated. Frame delay (in milliseconds) can be read from
// the `delay` attribute.
//
// DOC(Src): Binary image data
// DOC(Width): Scale image to this width
// DOC(Height): Scale image to this height
// DOC(Delay): (Read-only) Frame delay in ms, for animated GIFs
type Image struct {
	Widget
	Src           string `starlark:"src,required"`
	Width, Height int
	Delay         int `starlark:"delay,readonly"`

	imgs []image.Image
}

func (p *Image) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	return p.imgs[ModInt(frameIdx, len(p.imgs))]
}

func (p *Image) Size() (int, int) {
	return p.imgs[0].Bounds().Dx(), p.imgs[0].Bounds().Dy()
}

func (p *Image) FrameCount() int {
	return len(p.imgs)
}

func (p *Image) Init() error {
	var w, h int

	gifImg, err := gif.DecodeAll(bytes.NewReader([]byte(p.Src)))
	if err == nil {
		p.Delay = gifImg.Delay[0] * 10
		for _, im := range gifImg.Image {
			imRGBA := image.NewRGBA(image.Rect(0, 0, im.Bounds().Dx(), im.Bounds().Dy()))
			draw.Draw(imRGBA, imRGBA.Bounds(), im, image.Point{0, 0}, draw.Src)
			p.imgs = append(p.imgs, imRGBA)
		}
		w = p.imgs[0].Bounds().Dx()
		h = p.imgs[0].Bounds().Dy()
	} else {
		im, _, err := image.Decode(bytes.NewReader([]byte(p.Src)))
		if err != nil {
			return fmt.Errorf("decoding image data: %v", err)
		}

		p.imgs = []image.Image{im}
		w = im.Bounds().Dx()
		h = im.Bounds().Dy()
	}

	if p.Width != 0 || p.Height != 0 {
		nw, nh := p.Width, p.Height
		if nw == 0 {
			nw = w
		}
		if nh == 0 {
			nh = h
		}

		for i := 0; i < len(p.imgs); i++ {
			p.imgs[i] = resize.Resize(uint(nw), uint(nh), p.imgs[i], resize.NearestNeighbor)
		}
	}

	return nil
}

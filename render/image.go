package render

import (
	"bytes"
	"image"
	"log"

	// register image formats
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

// Image renders the binary image data passed via `src`. Supported
// formats include PNG, JPEG and GIF.
//
// If `width` or `height` are set, the image will be scaled
// accordingly, with nearest neighbor interpolation. Otherwise the
// image's original dimensions are used.
//
// DOC(Src): Binary image data
// DOC(Width): Scale image to this width
// DOC(Height): Scale image to this height
type Image struct {
	Widget
	Src           string `starlark:"src,required"`
	Width, Height int

	img image.Image
}

func (p *Image) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	if p.img == nil {
		p.paint()
	}
	return p.img
}

func (p *Image) Size() (int, int) {
	if p.img == nil {
		p.paint()
	}
	return p.img.Bounds().Dx(), p.img.Bounds().Dy()
}

func (p *Image) paint() {
	im, _, err := image.Decode(bytes.NewReader([]byte(p.Src)))
	if err != nil {
		log.Println("error decoding Image:", err)
		p.img = image.NewRGBA(image.Rect(0, 0, 0, 0))
		return
	}

	w, h := im.Bounds().Dx(), im.Bounds().Dy()

	if p.Width == 0 && p.Height == 0 || p.Width == w && p.Height == h {
		p.img = im
		return
	}

	nw, nh := p.Width, p.Height

	if nw == 0 {
		nw = w
	}
	if nh == 0 {
		nh = h
	}

	im = resize.Resize(uint(nw), uint(nh), im, resize.NearestNeighbor)

	dc := gg.NewContext(nw, nh)
	dc.DrawImage(im, 0, 0)

	p.img = dc.Image()
}

func (p *Image) FrameCount() int {
	return 1
}

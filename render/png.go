package render

import (
	"bytes"
	"image"
	"image/png"
	"log"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

// PNG is a widget for rendering a pre-existing PNG.
//
// WARNING: This widget will likely be removed in the near future. Use
// Image instead.
type PNG struct {
	Widget
	Src           string `starlark:"src,required"`
	Width, Height int

	img image.Image
}

func (p *PNG) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	if p.img == nil {
		p.paint()
	}
	return p.img
}

func (p *PNG) Size() (int, int) {
	if p.img == nil {
		p.paint()
	}
	return p.img.Bounds().Dx(), p.img.Bounds().Dy()
}

func (p *PNG) paint() {
	im, err := png.Decode(bytes.NewReader([]byte(p.Src)))
	if err != nil {
		log.Println("error decoding PNG:", err)
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

func (p *PNG) FrameCount() int {
	return 1
}

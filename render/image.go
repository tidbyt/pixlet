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

	"github.com/nfnt/resize"
	"github.com/tidbyt/go-libwebp/webp"
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

func (p *Image) InitFromWebP(data []byte) error {
	decoder, err := webp.NewAnimationDecoder(data)
	if err != nil {
		return fmt.Errorf("creating animation decoder: %v", err)
	}

	img, err := decoder.Decode()
	if err != nil {
		return fmt.Errorf("decoding image data: %v", err)
	}

	p.Delay = img.Timestamp[0]
	for _, im := range img.Image {
		p.imgs = append(p.imgs, im)
	}

	return nil
}

func (p *Image) InitFromGIF(data []byte) error {
	// GIF support is quite limited and can't handle different
	// positioned and sized frames, as well as disposal types.
	//
	// This means that many more optimized GIFs will not render
	// correctly, with frame contents jumping around and previous
	// frame contents always getting disposed, instead of kept.
	//
	// Unfortunatley the 'image/gif' package does not even expose
	// frame positions, making it hard to implement these features.
	//
	// Consider using WebP instead.
	img, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decoding image data: %v", err)
	}

	p.Delay = img.Delay[0] * 10
	for _, im := range img.Image {
		imRGBA := image.NewRGBA(image.Rect(0, 0, im.Bounds().Dx(), im.Bounds().Dy()))
		draw.Draw(imRGBA, imRGBA.Bounds(), im, image.Point{0, 0}, draw.Src)
		p.imgs = append(p.imgs, imRGBA)
	}

	return nil
}

func (p *Image) InitFromImage(data []byte) error {
	im, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decoding image data: %v", err)
	}

	p.imgs = []image.Image{im}

	return nil
}

func (p *Image) Init() error {
	err := p.InitFromWebP([]byte(p.Src))
	if err != nil {
		err = p.InitFromGIF([]byte(p.Src))
		if err != nil {
			err = p.InitFromImage([]byte(p.Src))
		}
	}

	if err != nil {
		return err
	}

	w := p.imgs[0].Bounds().Dx()
	h := p.imgs[0].Bounds().Dy()

	if p.Width != 0 || p.Height != 0 {
		nw, nh := p.Width, p.Height
		if nw == 0 {
			// scale width, maintaining original aspect ratio
			nw = int(float64(nh)*(float64(w)/float64(h)))
		}
		if nh == 0 {
			// scale height, maintaining original aspect ratio
			nh = int(float64(nw)*(float64(h)/float64(w)))
		}

		for i := 0; i < len(p.imgs); i++ {
			p.imgs[i] = resize.Resize(uint(nw), uint(nh), p.imgs[i], resize.NearestNeighbor)
		}
	}

	return nil
}

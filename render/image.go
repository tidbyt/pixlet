package render

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"

	// register image formats
	_ "image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/tidbyt/gg"
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

func (p *Image) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	return p.imgs[ModInt(frameIdx, len(p.imgs))].Bounds()
}

func (p *Image) Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int) {
	dc.DrawImage(p.imgs[ModInt(frameIdx, len(p.imgs))], 0, 0)
}

func (p *Image) Size() (int, int) {
	return p.imgs[0].Bounds().Dx(), p.imgs[0].Bounds().Dy()
}

func (p *Image) FrameCount() int {
	return len(p.imgs)
}

func (p *Image) InitFromGIF(data []byte) error {
	// Consider using WebP instead.
	img, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decoding image data: %v", err)
	}

	p.Delay = img.Delay[0] * 10

	var prev_src *image.Paletted
	disposal_length := len(img.Disposal)
	compositing_op := draw.Src

	last := image.NewRGBA(image.Rect(0, 0, img.Config.Width, img.Config.Height))

	for index, src := range img.Image {
		bounds := img.Image[index].Bounds()
		disposal_method := img.Disposal[index]
		is_disposal_previous := disposal_method == gif.DisposalPrevious

		// if the frame is DisposalPrevious
		// reset to the last non-DisposalPrevious frame
		if is_disposal_previous && prev_src != nil {
			draw.Draw(last, last.Bounds(), prev_src, image.ZP, draw.Over)
		}

		// if this is a non-DisposalPrevious frame
		// and the next frame is DisposalPrevious
		// store the src to reset before the next frame draws
		if !is_disposal_previous && index+1 < disposal_length && img.Disposal[index+1] == gif.DisposalPrevious {
			prev_src = src
		}

		draw.Draw(last, bounds, img.Image[index], image.Point{bounds.Min.X, bounds.Min.Y}, compositing_op)
		frame := *last
		frame.Pix = make([]uint8, len(last.Pix))
		copy(frame.Pix, last.Pix)

		// if this is a non-DisposalPrevious frame
		// set the compositing operation to Over
		if !is_disposal_previous {
			compositing_op = draw.Over
		}

		// if the frame is DisposalBackground
		// remove the frame pixels
		if disposal_method == gif.DisposalBackground {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					last.Set(x, y, color.Transparent)
				}
			}
		}

		p.imgs = append(p.imgs, &frame)
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

func (p *Image) InitFromSVG(data []byte) error {
	svgData, _ := oksvg.ReadIconStream(bytes.NewReader(data), oksvg.StrictErrorMode)
	w := int(svgData.ViewBox.W)
	h := int(svgData.ViewBox.H)

	if w == 0 && h == 0 {
		return errors.New("decoding svg data failed")
	}

	svgData.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	svgData.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, rgba, nil)

	if err != nil {
		return err
	}

	err = p.InitFromImage([]byte(buf.Bytes()))
	if err != nil {
		return err
	}

	return nil
}

func (p *Image) Init() error {
	err := p.InitFromWebP([]byte(p.Src))
	if err != nil {
		err = p.InitFromGIF([]byte(p.Src))
		if err != nil {
			err = p.InitFromSVG([]byte(p.Src))
			if err != nil {
				err = p.InitFromImage([]byte(p.Src))
			}
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
			nw = int(float64(nh) * (float64(w) / float64(h)))
		}
		if nh == 0 {
			// scale height, maintaining original aspect ratio
			nh = int(float64(nw) * (float64(h) / float64(w)))
		}

		for i := 0; i < len(p.imgs); i++ {
			p.imgs[i] = resize.Resize(uint(nw), uint(nh), p.imgs[i], resize.NearestNeighbor)
		}
	}

	return nil
}

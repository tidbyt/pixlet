//go:build !js && !wasm

package render

import (
	"fmt"

	"github.com/tidbyt/go-libwebp/webp"
)

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
		p.imgs = append(p.imgs, imageContainer{Image: im})
	}

	return nil
}

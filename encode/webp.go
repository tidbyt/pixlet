//go:build !js && !wasm

package encode

import (
	"fmt"
	"time"

	"github.com/tidbyt/go-libwebp/webp"
)

// Renders a screen to WebP. Optionally pass filters for
// postprocessing each individual frame.
func (s *Screens) EncodeWebP(maxDuration int, filters ...ImageFilter) ([]byte, error) {
	images, err := s.render(filters...)
	if err != nil {
		return nil, err
	}

	if len(images) == 0 {
		return []byte{}, nil
	}

	bounds := images[0].Bounds()
	anim, err := webp.NewAnimationEncoder(
		bounds.Dx(),
		bounds.Dy(),
		WebPKMin,
		WebPKMax,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "initializing encoder", err)
	}
	defer anim.Close()

	remainingDuration := time.Duration(maxDuration) * time.Millisecond
	for _, im := range images {
		frameDuration := time.Duration(s.delay) * time.Millisecond

		if maxDuration > 0 {
			if frameDuration > remainingDuration {
				frameDuration = remainingDuration
			}
			remainingDuration -= frameDuration
		}

		if err := anim.AddFrame(im, frameDuration); err != nil {
			return nil, fmt.Errorf("%s: %w", "adding frame", err)
		}

		if maxDuration > 0 && remainingDuration <= 0 {
			break
		}
	}

	buf, err := anim.Assemble()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "encoding animation", err)
	}

	return buf, nil
}

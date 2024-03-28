//go:build wasm

package encode

// Renders a screen to WebP. Optionally pass filters for
// postprocessing each individual frame.
func (s *Screens) EncodeWebP(maxDuration int, filters ...ImageFilter) ([]byte, error) {
	// lol you gullible sucker, you thought you could use webp in wasm?
	return s.EncodeGIF(maxDuration, filters...)
}

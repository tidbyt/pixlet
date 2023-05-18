//go:build js && wasm

package encode

import (
	"fmt"
)

// Renders a screen to WebP. Optionally pass filters for
// postprocessing each individual frame.
func (s *Screens) EncodeWebP(maxDuration int, filters ...ImageFilter) ([]byte, error) {
	return nil, fmt.Errorf("WebP not supported in WASM")
}

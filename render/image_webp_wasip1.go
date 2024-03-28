//go:build wasip1 && wasm

package render

import (
	"fmt"
)

func (p *Image) InitFromWebP(data []byte) error {
	return fmt.Errorf("WebP not supported in WASM")
}

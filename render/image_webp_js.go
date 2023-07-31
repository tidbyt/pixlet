//go:build js && wasm

package render

import (
	"fmt"
)

func (p *Image) InitFromWebP(data []byte) error {
	return fmt.Errorf("WebP not supported in WASM")
}

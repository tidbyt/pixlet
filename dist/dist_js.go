//go:build js && wasm

package dist

import "embed"

// dummy values not used in wasm build
var (
	Static embed.FS
	Index  []byte
)

//go:build js && wasm

package browser

import (
	wasmhttp "github.com/nlepage/go-wasm-http-server"
	"log"
)

func (b *Browser) serveHTTP() error {
	log.Printf("listening via wasmhttp")
	wasmhttp.Serve(b.r)
	return nil
}

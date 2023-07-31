//go:build js && wasm

package watcher

import (
	"fmt"
)

func (w *Watcher) Run() error {
	return fmt.Errorf("file watching not supported in WASM")
}

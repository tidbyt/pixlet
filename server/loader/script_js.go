//go:build js && wasm

package loader

import (
	"fmt"
	"io"
	"net/http"

	"tidbyt.dev/pixlet/runtime"
)

func loadScript(appID string, filename string) (*runtime.Applet, error) {
	res, err := http.Get(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file %s: %w", filename, err)
	}
	defer res.Body.Close()

	src, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return runtime.NewApplet(appID, src)
}

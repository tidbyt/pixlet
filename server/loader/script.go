//go:build !js && !wasm

package loader

import (
	"fmt"
	"os"

	"tidbyt.dev/pixlet/runtime"
)

func loadScript(appID string, filename string) (*runtime.Applet, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return runtime.NewApplet(appID, src)
}

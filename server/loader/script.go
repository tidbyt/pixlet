//go:build !js && !wasm

package loader

import (
	"fmt"
	"io/ioutil"

	"tidbyt.dev/pixlet/runtime"
)

func loadScript(applet *runtime.Applet, appID string, filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	err = applet.Load(appID, filename, src, nil)
	if err != nil {
		return fmt.Errorf("failed to load applet: %w", err)
	}

	return nil
}

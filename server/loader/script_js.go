//go:build js && wasm

package loader

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"tidbyt.dev/pixlet/runtime"
)

func loadScript(applet *runtime.Applet, filename string) error {
	res, err := http.Get(filename)
	if err != nil {
		return fmt.Errorf("failed to fetch file %s: %w", filename, err)
	}
	defer res.Body.Close()

	src, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	err = applet.Load(filename, src, nil)
	if err != nil {
		return fmt.Errorf("failed to load applet: %w", err)
	}

	return nil
}

package loader

import (
	"io/fs"

	"tidbyt.dev/pixlet/runtime"
)

func loadScript(appID string, fs fs.FS) (*runtime.Applet, error) {
	return runtime.NewAppletFromFS(appID, fs)
}

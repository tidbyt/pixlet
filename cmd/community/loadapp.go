package community

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/runtime"
	"tidbyt.dev/pixlet/tools"
)

var LoadAppCmd = &cobra.Command{
	Use:     "load-app <path>",
	Short:   "Validates an app can be successfully loaded in our runtime.",
	Example: `pixlet community load-app examples/clock`,
	Long:    `This command ensures an app can be loaded into our runtime successfully.`,
	Args:    cobra.ExactArgs(1),
	RunE:    LoadApp,
}

func LoadApp(cmd *cobra.Command, args []string) error {
	path := args[0]

	// check if path exists, and whether it is a directory or a file
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", path, err)
	}

	var fs fs.FS
	if info.IsDir() {
		fs = os.DirFS(path)
	} else {
		if !strings.HasSuffix(path, ".star") {
			return fmt.Errorf("script file must have suffix .star: %s", path)
		}

		fs = tools.NewSingleFileFS(path)
	}

	cache := runtime.NewInMemoryCache()
	runtime.InitHTTP(cache)
	runtime.InitCache(cache)

	if _, err := runtime.NewAppletFromFS(filepath.Base(path), fs, runtime.WithPrintDisabled()); err != nil {
		return fmt.Errorf("failed to load applet: %w", err)
	}

	return nil
}

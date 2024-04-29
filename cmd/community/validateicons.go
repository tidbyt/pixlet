package community

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"tidbyt.dev/pixlet/icons"
	"tidbyt.dev/pixlet/runtime"
	"tidbyt.dev/pixlet/schema"
	"tidbyt.dev/pixlet/tools"
)

var ValidateIconsCmd = &cobra.Command{
	Use:     "validate-icons <path>",
	Short:   "Validates the schema icons used are available in our mobile app.",
	Example: `pixlet community validate-icons examples/schema_hello_world`,
	Long: `This command determines if the icons selected in your app schema are supported
by our mobile app.`,
	Args: cobra.ExactArgs(1),
	RunE: ValidateIcons,
}

func ValidateIcons(cmd *cobra.Command, args []string) error {
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

	applet, err := runtime.NewAppletFromFS(filepath.Base(path), fs, runtime.WithPrintDisabled())
	if err != nil {
		return fmt.Errorf("failed to load applet: %w", err)
	}

	s := schema.Schema{}
	js := applet.SchemaJSON
	if len(js) == 0 {
		return nil
	}

	err = json.Unmarshal(js, &s)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	for _, field := range s.Fields {
		if field.Icon == "" {
			continue
		}

		if _, ok := icons.IconsMap[field.Icon]; !ok {
			return fmt.Errorf("app '%s' contains unknown icon: '%s'", applet.ID, field.Icon)
		}
	}

	return nil
}

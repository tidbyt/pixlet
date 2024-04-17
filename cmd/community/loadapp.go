package community

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/runtime"
)

var LoadAppCmd = &cobra.Command{
	Use:     "load-app <filespec>",
	Short:   "Validates an app can be successfully loaded in our runtime.",
	Example: `  pixlet community load-app app.star`,
	Long:    `This command ensures an app can be loaded into our runtime successfully.`,
	Args:    cobra.ExactArgs(1),
	RunE:    LoadApp,
}

func LoadApp(cmd *cobra.Command, args []string) error {
	script := args[0]

	if !strings.HasSuffix(script, ".star") {
		return fmt.Errorf("script file must have suffix .star: %s", script)
	}

	src, err := os.ReadFile(script)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", script, err)
	}

	cache := runtime.NewInMemoryCache()
	runtime.InitHTTP(cache)
	runtime.InitCache(cache)

	if _, err := runtime.NewApplet(script, src, runtime.WithPrintDisabled()); err != nil {
		return fmt.Errorf("failed to load applet: %w", err)
	}

	return nil
}

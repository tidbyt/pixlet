package community

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
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

	src, err := ioutil.ReadFile(script)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", script, err)
	}

	cache := runtime.NewInMemoryCache()
	runtime.InitHTTP(cache)
	runtime.InitCache(cache)

	// Remove the print function from the starlark thread.
	initializers := []runtime.ThreadInitializer{}
	initializers = append(initializers, func(thread *starlark.Thread) *starlark.Thread {
		thread.Print = func(thread *starlark.Thread, msg string) {}
		return thread
	})

	applet := runtime.Applet{}
	err = applet.LoadWithInitializers("", script, src, nil, initializers...)
	if err != nil {
		return fmt.Errorf("failed to load applet: %w", err)
	}

	return nil
}

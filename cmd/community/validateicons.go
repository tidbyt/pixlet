package community

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/icons"
	"tidbyt.dev/pixlet/runtime"
	"tidbyt.dev/pixlet/schema"
)

var ValidateIconsCmd = &cobra.Command{
	Use:     "validate-icons <pathspec>",
	Short:   "Validates the schema icons used are available in our mobile app.",
	Example: `  pixlet community validate-icons app.star`,
	Long: `This command determines if the icons selected in your app schema are supported
by our mobile app.`,
	Args: cobra.ExactArgs(1),
	RunE: ValidateIcons,
}

func ValidateIcons(cmd *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("couldn't open app: %w", err)
	}
	defer f.Close()

	src, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read app %s: %w", args[0], err)
	}

	cache := runtime.NewInMemoryCache()
	runtime.InitHTTP(cache)
	runtime.InitCache(cache)

	applet := runtime.Applet{}
	err = applet.Load(args[0], src, nil)
	if err != nil {
		return fmt.Errorf("failed to load applet: %w", err)
	}

	s := schema.Schema{}
	schemaStr := applet.GetSchema()
	if schemaStr == "" {
		return nil
	}

	err = json.Unmarshal([]byte(schemaStr), &s)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	for _, field := range s.Fields {
		if field.Icon == "" {
			continue
		}

		if _, ok := icons.IconsMap[field.Icon]; !ok {
			return fmt.Errorf("app '%s' contains unknown icon: '%s'", applet.Filename, field.Icon)
		}
	}

	return nil
}

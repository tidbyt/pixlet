package community

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/manifest"
)

var ValidateManifestCmd = &cobra.Command{
	Use:     "validate-manifest <pathspec>",
	Short:   "Validates an app manifest is ready for publishing",
	Example: `  pixlet community validate-manifest manifest.yaml`,
	Long: `This command determines if your app manifest is configured properly by
validating the contents of each field.`,
	Args: cobra.ExactArgs(1),
	RunE: ValidateManifest,
}

func ValidateManifest(cmd *cobra.Command, args []string) error {
	fileName := filepath.Base(args[0])
	if fileName != "manifest.yaml" && fileName != "manifest.yml" {
		return fmt.Errorf("supplied manifest must be named manifest.yaml or manifest.yml")
	}

	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("couldn't open manifest: %w", err)
	}
	defer f.Close()

	m, err := manifest.LoadManifest(f)
	if err != nil {
		return fmt.Errorf("couldn't load manifest: %w", err)
	}

	err = m.Validate()
	if err != nil {
		return fmt.Errorf("couldn't validate manifest: %w", err)
	}

	return nil
}

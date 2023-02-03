package community

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/manifest"
)

var ValidateManifestAppFileName string

func init() {
	ValidateManifestCmd.Flags().StringVarP(&ValidateManifestAppFileName, "app-file-name", "a", "", "ensures the app file name is the same as the manifest")
}

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
	if fileName != manifest.ManifestFileName {
		return fmt.Errorf("supplied manifest must be named %s", manifest.ManifestFileName)
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

	if ValidateManifestAppFileName != "" && m.FileName != ValidateManifestAppFileName {
		return fmt.Errorf("app name doesn't match: %s != %s", ValidateManifestAppFileName, m.FileName)
	}

	return nil
}

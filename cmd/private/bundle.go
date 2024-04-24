package private

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/bundle"
)

var bundleOutput string

func init() {
	BundleCmd.Flags().StringVarP(&bundleOutput, "output", "o", "./", "output directory for the bundle")
}

var BundleCmd = &cobra.Command{
	Use:     "bundle",
	Short:   "Creates a new app bundle",
	Example: `  pixlet bundle ./my-app`,
	Long: `This command will create a new app bundle from an app directory. The directory
should contain an app manifest and source file. The output of this command will
be a gzip compressed tar file that can be uploaded to Tidbyt for deployment.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bundleInput := args[0]
		info, err := os.Stat(bundleInput)
		if err != nil {
			return fmt.Errorf("input directory invalid: %w", err)
		}

		if !info.IsDir() {
			return fmt.Errorf("input must be a directory")
		}

		info, err = os.Stat(bundleOutput)
		if err != nil {
			return fmt.Errorf("output directory invalid: %w", err)
		}

		if !info.IsDir() {
			return fmt.Errorf("output must be a directory")
		}

		ab, err := bundle.FromDir(bundleInput)
		if err != nil {
			return fmt.Errorf("could not init bundle: %w", err)
		}

		return ab.WriteBundleToPath(bundleOutput)
	},
}

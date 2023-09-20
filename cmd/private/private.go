package private

import (
	"github.com/spf13/cobra"
)

func init() {
	PrivateCmd.AddCommand(BundleCmd)
	PrivateCmd.AddCommand(UploadCmd)
	PrivateCmd.AddCommand(DeployCmd)
	PrivateCmd.AddCommand(ListCmd)
}

var PrivateCmd = &cobra.Command{
	Use:   "private",
	Short: "Utilities to manage private apps",
	Long: `The private subcommand provides a set of utilities for managing
private apps, including team apps and individual users private apps.`,
}

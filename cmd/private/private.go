package private

import (
	"github.com/spf13/cobra"
)

func init() {
	PrivateCmd.AddCommand(CreateCmd)
	PrivateCmd.AddCommand(BundleCmd)
	PrivateCmd.AddCommand(UploadCmd)
	PrivateCmd.AddCommand(DeployCmd)
	PrivateCmd.AddCommand(DeleteCmd)
	PrivateCmd.AddCommand(ListCmd)
	PrivateCmd.AddCommand(LogsCmd)
}

var PrivateCmd = &cobra.Command{
	Use:   "private",
	Short: "Utilities to manage private apps",
	Long: `The private subcommand provides a set of utilities for managing
private apps. Requires Tidbyt Plus or Tidbyt for Teams.`,
}

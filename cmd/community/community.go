package community

import (
	"github.com/spf13/cobra"
)

func init() {
	CommunityCmd.AddCommand(CreateManifestCmd)
	CommunityCmd.AddCommand(TargetDeterminatorCmd)
	CommunityCmd.AddCommand(ValidateManifestCmd)
}

var CommunityCmd = &cobra.Command{
	Use:   "community",
	Short: "Utilities to manage the community repo",
	Long: `The community subcommand provides a set of utilities for managing the
community repo. This subcommand should be considered slightly unstable in that
we may determine a utility here should move to a more generalizable tool.`,
}

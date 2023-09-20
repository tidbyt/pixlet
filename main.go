package main

import (
	"os"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/cmd"
	"tidbyt.dev/pixlet/cmd/community"
	"tidbyt.dev/pixlet/cmd/private"
)

var (
	rootCmd = &cobra.Command{
		Use:          "pixlet",
		Short:        "pixel graphics rendering",
		Long:         "Pixlet renders graphics for pixel devices, like Tidbyt",
		SilenceUsage: true,
	}
)

func init() {
	rootCmd.AddCommand(cmd.ServeCmd)
	rootCmd.AddCommand(cmd.RenderCmd)
	rootCmd.AddCommand(cmd.PushCmd)
	rootCmd.AddCommand(cmd.EncryptCmd)
	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.AddCommand(cmd.ProfileCmd)
	rootCmd.AddCommand(cmd.LoginCmd)
	rootCmd.AddCommand(cmd.DevicesCmd)
	rootCmd.AddCommand(cmd.ListCmd)
	rootCmd.AddCommand(cmd.DeleteCmd)
	rootCmd.AddCommand(cmd.FormatCmd)
	rootCmd.AddCommand(cmd.LintCmd)
	rootCmd.AddCommand(cmd.CheckCmd)
	rootCmd.AddCommand(cmd.SetAuthCmd)
	rootCmd.AddCommand(community.CommunityCmd)
	rootCmd.AddCommand(private.PrivateCmd)

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

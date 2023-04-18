package main

import (
	"os"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/cmd"
	"tidbyt.dev/pixlet/cmd/community"
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
	rootCmd.AddCommand(cmd.FormatCmd)
	rootCmd.AddCommand(cmd.LintCmd)
	rootCmd.AddCommand(cmd.CreateCmd)
	rootCmd.AddCommand(cmd.CheckCmd)
	rootCmd.AddCommand(cmd.BundleCmd)
	rootCmd.AddCommand(cmd.UploadCmd)
	rootCmd.AddCommand(cmd.DeployCmd)
	rootCmd.AddCommand(cmd.SetAuthCmd)
	rootCmd.AddCommand(cmd.TestCmd)
	rootCmd.AddCommand(community.CommunityCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

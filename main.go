package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "pixlet",
	Short: "pixel graphics rendering",
	Long:  "Pixlet renders graphics for pixel devices, like Tidbyt",
}

func init() {
	rootCmd.AddCommand(cmd.ServeCmd)
	rootCmd.AddCommand(cmd.RenderCmd)
	rootCmd.AddCommand(cmd.PushCmd)
	rootCmd.AddCommand(cmd.EncryptCmd)
	rootCmd.AddCommand(cmd.VersionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

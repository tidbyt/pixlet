package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of Pixlet",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Pixlet version: %s\n", Version)
	},
}

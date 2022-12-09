package cmd

import (
	"fmt"
	"os"

	"github.com/bazelbuild/buildtools/buildifier/utils"
	"github.com/bazelbuild/buildtools/differ"
	"github.com/spf13/cobra"
)

func init() {
	LintCmd.Flags().BoolVarP(&vflag, "verbose", "v", false, "print verbose information to standard error")
	LintCmd.Flags().BoolVarP(&rflag, "recursive", "r", false, "find starlark files recursively")
	LintCmd.Flags().BoolVarP(&dryRunFlag, "dry-run", "d", false, "no code modifications")
	LintCmd.Flags().StringVarP(&format, "format", "f", "", "diagnostics format: text or json (default text)")
}

var LintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lints Pixlet apps.",
	Run:   lintCmd,
}

func lintCmd(cmd *cobra.Command, args []string) {
	mode := "check"
	lint := "warn"

	if err := utils.ValidateFormat(&format, &mode); err != nil {
		fmt.Fprintf(os.Stderr, "buildifier: %s\n", err)
		os.Exit(2)
	}

	dflag := false
	if err := utils.ValidateModes(&mode, &lint, &dflag); err != nil {
		fmt.Fprintf(os.Stderr, "buildifier: %s\n", err)
		os.Exit(2)
	}

	differ, deprecationWarning := differ.Find()
	if deprecationWarning && mode == "diff" {
		fmt.Fprintf(os.Stderr, "buildifier: selecting diff program with the BUILDIFIER_DIFF, BUILDIFIER_MULTIDIFF, and DISPLAY environment variables is deprecated, use flags -diff_command and -multi_diff instead\n")
	}
	diff = differ

	exitCode := runBuildifier(args, lint, mode, format, rflag, vflag)
	os.Exit(exitCode)
}

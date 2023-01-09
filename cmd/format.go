package cmd

import (
	"os"

	"github.com/bazelbuild/buildtools/differ"
	"github.com/spf13/cobra"
)

func init() {
	FormatCmd.Flags().BoolVarP(&vflag, "verbose", "v", false, "print verbose information to standard error")
	FormatCmd.Flags().BoolVarP(&rflag, "recursive", "r", false, "find starlark files recursively")
	FormatCmd.Flags().BoolVarP(&dryRunFlag, "dry-run", "d", false, "display a diff of formatting changes without modification")
	FormatCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "output format, text or json")
}

var FormatCmd = &cobra.Command{
	Use:   "format <pathspec>...",
	Short: "Formats Tidbyt apps",
	Example: `  pixlet format app.star
  pixlet format app.star --dry-run
  pixlet format --recursive ./`,
	Long: `The format command provides a code formatter for Tidbyt apps. By default, it
will format your starlark source code in line. If you wish you see the output
before applying, add the --dry-run flag.`,
	Args: cobra.MinimumNArgs(1),
	Run:  formatCmd,
}

func formatCmd(cmd *cobra.Command, args []string) {
	// Lint refers to the lint mode for buildifier, with the options being off,
	// warn, or fix. For pixlet format, we don't want to lint at all.
	lint := "off"

	// Mode refers to formatting mode for buildifier, with the options being
	// check, diff, or fix. For the pixlet format command, we want to fix the
	// resolvable issue by default and provide a dry run flag to be able to
	// diff the changes before fixing them.
	mode := "fix"
	if dryRunFlag {
		mode = "diff"
	}

	// Copied from the buildifier source, we need to supply a diff program for
	// the differ.
	differ, _ := differ.Find()
	diff = differ

	// Run buildifier and exit with the returned exit code.
	exitCode := runBuildifier(args, lint, mode, outputFormat, rflag, vflag)
	os.Exit(exitCode)
}

package cmd

import (
	"os"

	"github.com/bazelbuild/buildtools/differ"
	"github.com/spf13/cobra"
)

func init() {
	LintCmd.Flags().BoolVarP(&vflag, "verbose", "v", false, "print verbose information to standard error")
	LintCmd.Flags().BoolVarP(&rflag, "recursive", "r", false, "find starlark files recursively")
	LintCmd.Flags().BoolVarP(&fixFlag, "fix", "f", false, "automatically fix resolvable lint issues")
	LintCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "output format, text or json")
}

var LintCmd = &cobra.Command{
	Use: "lint <pathspec>...",
	Example: `  pixlet lint app.star
  pixlet lint --recursive --fix ./`,
	Short: "Lints Tidbyt apps",
	Long: `The lint command provides a linter for Tidbyt apps. It's capable of linting a
file, a list of files, or directory with the recursive option. Additionally, it
provides an option to automatically fix resolvable linter issues.`,
	Args: cobra.MinimumNArgs(1),
	Run:  lintCmd,
}

func lintCmd(cmd *cobra.Command, args []string) {
	// Mode refers to formatting mode for buildifier, with the options being
	// check, diff, or fix. For the pixlet lint command, we only want to check
	// formatting.
	mode := "check"

	// Lint refers to the lint mode for buildifier, with the options being off,
	// warn, or fix. For pixlet lint, we want to warn by default but offer a
	// flag to automatically fix resolvable issues.
	lint := "warn"

	// If the fix flag is enabled, the lint command should both format and lint.
	if fixFlag {
		mode = "fix"
		lint = "fix"
	}

	// Copied from the buildifier source, we need to supply a diff program for
	// the differ.
	differ, _ := differ.Find()
	diff = differ

	// TODO: We currently offer misspelling protection in the community repo
	// for app manifests. We'll want to consider adding additional spelling
	// support to pixlet lint to ensure typos in apps don't make it to
	// production.

	// Run buildifier and exit with the returned exit code.
	exitCode := runBuildifier(args, lint, mode, outputFormat, rflag, vflag)
	os.Exit(exitCode)
}

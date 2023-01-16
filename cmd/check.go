package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/bazelbuild/buildtools/buildifier/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	CheckCmd.Flags().BoolVarP(&rflag, "recursive", "r", false, "find apps recursively")
}

var CheckCmd = &cobra.Command{
	Use:     "check <pathspec>...",
	Example: `  pixlet check app.star`,
	Short:   "Checks if an app is ready to publish",
	Long: `The check command runs a series of checks to ensure your app is ready
to publish in the community repo. Every failed check will have a solution
provided. If your app fails a check, try the provided solution and reach out on
Discord if you get stuck.`,
	Args: cobra.MinimumNArgs(1),
	RunE: checkCmd,
}

func checkCmd(cmd *cobra.Command, args []string) error {
	// Use the same logic as buildifier to find relevant Tidbyt apps.
	apps := args
	if rflag {
		discovered, err := utils.ExpandDirectories(&args)
		if err != nil {
			return fmt.Errorf("could not discover apps using recursive flag: %w", err)
		}
		apps = discovered
	}

	// Check every app.
	foundIssue := false
	for _, app := range apps {
		// Check app formatting.
		dryRunFlag = true
		err := formatCmd(cmd, []string{app})
		if err != nil {
			foundIssue = true
			failure(app, fmt.Errorf("app is not formatted correctly: %w", err), fmt.Sprintf("try `pixlet format %s`", app))
			continue
		}

		// Create temporary file for app rendering.
		f, err := os.CreateTemp("", "")
		if err != nil {
			return fmt.Errorf("could not create temp file for rendering, check your system: %w", err)
		}
		defer os.Remove(f.Name())

		// Check if app will render.
		silenceOutput = true
		output = f.Name()
		err = render(cmd, []string{app})
		if err != nil {
			foundIssue = true
			failure(app, fmt.Errorf("app failed to render: %w", err), fmt.Sprintf("try `pixlet render %s` and resolve any errors", app))
			continue
		}

		// Check if app is linted.
		outputFormat = "off"
		err = lintCmd(cmd, []string{app})
		if err != nil {
			foundIssue = true
			failure(app, fmt.Errorf("app has lint warnings: %w", err), fmt.Sprintf("try `pixlet lint --fix %s`", app))
			continue
		}

		success(app)
	}

	if foundIssue {
		return fmt.Errorf("one or more apps failed checks")
	}

	return nil
}

func success(app string) {
	c := color.New(color.FgGreen)
	c.Printf("✔️ %s\n", app)
}

func failure(app string, err error, sol string) {
	c := color.New(color.FgRed)
	c.Printf("✖ %s\n", app)

	// Ensure multiline errors are properly indented.
	multilineError := strings.Split(err.Error(), "\n")
	for index, line := range multilineError {
		if index == 0 {
			continue
		}

		// The builtin starlark Backtrace function prints the last line at an
		// awkward indentation level. This check helps keep the failure indented
		// at one more level to ensure it's even more clear what is broken.
		if strings.Contains(line, "Error in") && index == len(multilineError)-1 {
			multilineError[index] = fmt.Sprintf("      %s", line)
		} else {
			multilineError[index] = fmt.Sprintf("  %s", line)
		}
	}
	problem := strings.Join(multilineError, "\n")

	fmt.Printf("  ▪️ Problem: %v\n", problem)
	fmt.Printf("  ▪️ Solution: %v\n", sol)
}

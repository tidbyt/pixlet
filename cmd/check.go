package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bazelbuild/buildtools/buildifier/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/cmd/community"
	"tidbyt.dev/pixlet/manifest"
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
	} else {
		for _, app := range apps {
			if filepath.Ext(app) != ".star" {
				return fmt.Errorf("only starlark source files or directories with the recursive flag are supported")
			}
		}
	}

	// TODO: this needs to be parallelized.

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

		// TODO: Check if app will render once we are able to enable target
		// determination.

		// Check if an app can load.
		err = community.LoadApp(cmd, []string{app})
		if err != nil {
			foundIssue = true
			failure(app, fmt.Errorf("app failed to load: %w", err), "try `pixlet community load-app` and resolve any runtime issues")
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

		// Ensure icons are valid.
		err = community.ValidateIcons(cmd, []string{app})
		if err != nil {
			foundIssue = true
			failure(app, fmt.Errorf("app has invalid icons: %w", err), "try `pixlet community list-icons` for the full list of valid icons")
			continue
		}

		// Check app manifest exists
		dir := filepath.Dir(app)
		if !doesManifestExist(dir) {
			foundIssue = true
			failure(app, fmt.Errorf("couldn't find app manifest: %w", err), fmt.Sprintf("try `pixlet community create-manifest %s`", filepath.Join(dir, manifest.ManifestFileName)))
			continue
		}

		// Validate manifest.
		manifestFile := filepath.Join(dir, manifest.ManifestFileName)
		community.ValidateManifestAppFileName = filepath.Base(app)
		err = community.ValidateManifest(cmd, []string{manifestFile})
		if err != nil {
			foundIssue = true
			failure(app, fmt.Errorf("manifest didn't validate: %w", err), "try correcting the validation issue by updating your manifest")
			continue
		}

		// Check spelling.
		community.SilentSpelling = true
		err = community.SpellCheck(cmd, []string{manifestFile})
		if err != nil {
			foundIssue = true
			failure(app, fmt.Errorf("manifest contains spelling errors: %w", err), fmt.Sprintf("try `pixlet community spell-check --fix %s`", manifestFile))
			continue
		}
		// TODO: enable spell check for apps once we can run it successfully
		// against the community repo.

		// If we're here, the app and manifest are good to go!
		success(app)
	}

	if foundIssue {
		return fmt.Errorf("one or more apps failed checks")
	}

	return nil
}

func doesManifestExist(dir string) bool {
	file := filepath.Join(dir, manifest.ManifestFileName)
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		return false
	}

	return true
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
		if (strings.Contains(line, "Error in") || strings.Contains(line, "Error:")) && index == len(multilineError)-1 {
			multilineError[index] = fmt.Sprintf("      %s", line)
		} else {
			multilineError[index] = fmt.Sprintf("  %s", line)
		}
	}
	problem := strings.Join(multilineError, "\n")

	fmt.Printf("  ▪️ Problem: %v\n", problem)
	fmt.Printf("  ▪️ Solution: %v\n", sol)
}

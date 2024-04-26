package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/cmd/community"
	"tidbyt.dev/pixlet/manifest"
	"tidbyt.dev/pixlet/tools"
)

const MaxRenderTime = 1000000000 // 1000ms

var CheckCmd = &cobra.Command{
	Use:     "check <path>...",
	Example: `pixlet check examples/clock`,
	Short:   "Check if an app is ready to publish",
	Long: `Check if an app is ready to publish.

The path argument should be the path to the Pixlet app to check. The
app can be a single file with the .star extension, or a directory
containing multiple Starlark files and resources.

The check command runs a series of checks to ensure your app is ready
to publish in the community repo. Every failed check will have a solution
provided. If your app fails a check, try the provided solution and reach out on
Discord if you get stuck.`,
	Args: cobra.MinimumNArgs(1),
	RunE: checkCmd,
}

func checkCmd(cmd *cobra.Command, args []string) error {
	// check every path.
	foundIssue := false
	for _, path := range args {
		// check if path exists, and whether it is a directory or a file
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to stat %s: %w", path, err)
		}

		var fsys fs.FS
		var baseDir string
		if info.IsDir() {
			fsys = os.DirFS(path)
			baseDir = path
		} else {
			if !strings.HasSuffix(path, ".star") {
				return fmt.Errorf("script file must have suffix .star: %s", path)
			}

			fsys = tools.NewSingleFileFS(path)
			baseDir = filepath.Dir(path)
		}

		// run format and lint on *.star files in the fs
		fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() || !strings.HasSuffix(p, ".star") {
				return nil
			}

			realPath := filepath.Join(baseDir, p)

			dryRunFlag = true
			if err := formatCmd(cmd, []string{realPath}); err != nil {
				foundIssue = true
				failure(p, fmt.Errorf("app is not formatted correctly: %w", err), fmt.Sprintf("try `pixlet format %s`", realPath))
			}

			outputFormat = "off"
			err = lintCmd(cmd, []string{realPath})
			if err != nil {
				foundIssue = true
				failure(p, fmt.Errorf("app has lint warnings: %w", err), fmt.Sprintf("try `pixlet lint --fix %s`", realPath))
			}

			return nil
		})

		// Check if an app can load.
		err = community.LoadApp(cmd, []string{path})
		if err != nil {
			foundIssue = true
			failure(path, fmt.Errorf("app failed to load: %w", err), "try `pixlet community load-app` and resolve any runtime issues")
			continue
		}

		// Ensure icons are valid.
		err = community.ValidateIcons(cmd, []string{path})
		if err != nil {
			foundIssue = true
			failure(path, fmt.Errorf("app has invalid icons: %w", err), "try `pixlet community list-icons` for the full list of valid icons")
			continue
		}

		// Check app manifest exists
		if !doesManifestExist(baseDir) {
			foundIssue = true
			failure(path, fmt.Errorf("couldn't find app manifest"), fmt.Sprintf("try `pixlet community create-manifest %s`", filepath.Join(baseDir, manifest.ManifestFileName)))
			continue
		}

		// Validate manifest.
		manifestFile := filepath.Join(baseDir, manifest.ManifestFileName)
		community.ValidateManifestAppFileName = filepath.Base(path)
		err = community.ValidateManifest(cmd, []string{manifestFile})
		if err != nil {
			foundIssue = true
			failure(path, fmt.Errorf("manifest didn't validate: %w", err), "try correcting the validation issue by updating your manifest")
			continue
		}

		// Create temporary file for app rendering.
		f, err := os.CreateTemp("", "")
		if err != nil {
			return fmt.Errorf("could not create temp file for rendering, check your system: %w", err)
		}
		defer os.Remove(f.Name())

		// Check if app renders.
		silenceOutput = true
		output = f.Name()
		err = render(cmd, []string{path})
		if err != nil {
			foundIssue = true
			failure(path, fmt.Errorf("app failed to render: %w", err), "try `pixlet render` and resolve any runtime issues")
			continue
		}

		// Check performance.
		p, err := ProfileApp(path, map[string]string{})
		if err != nil {
			return fmt.Errorf("could not profile app: %w", err)
		}
		if p.DurationNanos > MaxRenderTime {
			foundIssue = true
			failure(
				path,
				fmt.Errorf("app takes too long to render %s", time.Duration(p.DurationNanos)),
				fmt.Sprintf("try optimizing your app using `pixlet profile %s` to get it under %s", path, time.Duration(MaxRenderTime)),
			)
			continue
		}

		// If we're here, the app and manifest are good to go!
		success(path)
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

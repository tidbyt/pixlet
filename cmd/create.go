package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/cmd/community"
	"tidbyt.dev/pixlet/tools/generator"
	"tidbyt.dev/pixlet/tools/repo"
)

var (
	appDir string
	usePackageName bool
)

func init() {
	CreateCmd.Flags().StringVarP(&appDir,"appDir","","","Path for created app (when not in Tidbyt repos)")
	CreateCmd.Flags().BoolVarP(&usePackageName,"usePackageName","",false,"Create app in a folder using PackageName (when not in Tidbyt repos)")
}
// CreateCmd prompts the user for info and generates a new app.
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new app",
	Long:  `This command will prompt for all of the information we need to generate a new Tidbyt app.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the current working directory.
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("app creation failed, something went wrong with your local filesystem: %w", err)
		}

		// Determine what type of app this is an what the root should be.
		var root string
		var appType generator.AppType
		if repo.IsInRepo(cwd, "community") {
			appType = generator.Community
			root, err = repo.RepoRoot(cwd)
			if err != nil {
				return fmt.Errorf("app creation failed, something went wrong with your community repo: %w", err)
			}
		} else if repo.IsInRepo(cwd, "tidbyt") {
			appType = generator.Internal
			root, err = repo.RepoRoot(cwd)
			if err != nil {
				return fmt.Errorf("app creation failed, something went wrong with your tidbyt repo: %w", err)
			}
		} else if (appDir!="" || usePackageName) {
			appType = generator.LocalCustom
			root = filepath.Join(cwd,appDir)
		} else {
			appType = generator.Local
			root = cwd
		}

		// Prompt the user for input.
		app, err := community.ManifestPrompt()
		if err != nil {
			return fmt.Errorf("app creation, couldn't get user input: %w", err)
		}

		// Append packageName to create root folder too
		if (usePackageName && (appType==generator.Local || appType==generator.LocalCustom)){
			root = filepath.Join(root,app.PackageName)	
		} 
		// Generate app.
		g, err := generator.NewGenerator(appType, root)
		if err != nil {
			return fmt.Errorf("app creation failed %w", err)
		}
		absolutePath, err := g.GenerateApp(app)
		if err != nil {
			return fmt.Errorf("app creation failed: %w", err)
		}

		// Get the relative path from where the user started. Note, we're not
		// using the root here, given the root can be git repo specific.
		relativePath, err := filepath.Rel(cwd, absolutePath)
		if err != nil {
			return fmt.Errorf("app was created, but we don't know where: %w", err)
		}

		// Let the user know where the app is and how to use it.
		fmt.Println("")
		fmt.Println("App created at:")
		fmt.Printf("\t%s\n", absolutePath)
		fmt.Println("")
		fmt.Println("To start the app, run:")
		fmt.Printf("\tpixlet serve %s\n", relativePath)
		fmt.Println("")
		fmt.Println("For docs, head to:")
		fmt.Printf("\thttps://tidbyt.dev\n")
		return nil
	},
}

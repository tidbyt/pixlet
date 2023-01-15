package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/tools/generator"
	"tidbyt.dev/pixlet/tools/manifest"
	"tidbyt.dev/pixlet/tools/repo"
)

// CreateCmd prompts the user for info and generates a new app.
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new app.",
	Long:  `This command will prompt for all of the information we need to generate a new app in this repo.`,
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
		} else {
			appType = generator.Local
			root = cwd
		}

		// Get the name of the app.
		namePrompt := promptui.Prompt{
			Label:    "Name (what do you want to call your app?)",
			Validate: manifest.ValidateName,
		}
		name, err := namePrompt.Run()
		if err != nil {
			return fmt.Errorf("app creation failed %w", err)
		}

		// Get the summary of the app.
		summaryPrompt := promptui.Prompt{
			Label:    "Summary (what's the short and sweet of what this app does?)",
			Validate: manifest.ValidateSummary,
		}
		summary, err := summaryPrompt.Run()
		if err != nil {
			return fmt.Errorf("app creation failed %w", err)
		}

		// Get the description of the app.
		descPrompt := promptui.Prompt{
			Label:    "Description (what's the long form of what this app does?)",
			Validate: manifest.ValidateDesc,
		}
		desc, err := descPrompt.Run()
		if err != nil {
			return fmt.Errorf("app creation failed %w", err)
		}

		// Get the author of the app.
		authorPrompt := promptui.Prompt{
			Label:    "Author (your name or your Github handle)",
			Validate: manifest.ValidateAuthor,
		}
		author, err := authorPrompt.Run()
		if err != nil {
			return fmt.Errorf("app creation failed %w", err)
		}

		// Generate app.
		g, err := generator.NewGenerator(appType, root)
		if err != nil {
			return fmt.Errorf("app creation failed %w", err)
		}
		app := &manifest.Manifest{
			ID:          manifest.GenerateID(name),
			Name:        name,
			Summary:     summary,
			Desc:        desc,
			Author:      author,
			FileName:    manifest.GenerateFileName(name),
			PackageName: manifest.GeneratePackageName(name),
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

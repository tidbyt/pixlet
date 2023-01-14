package cmd

import (
	"fmt"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {
		if !repo.IsInRepo("community") {
			fmt.Println("app creation failed, try again from the community repo")
			os.Exit(1)
		}

		// Get the name of the app.
		namePrompt := promptui.Prompt{
			Label:    "Name (what do you want to call your app?)",
			Validate: manifest.ValidateName,
		}
		name, err := namePrompt.Run()
		if err != nil {
			fmt.Printf("app creation failed %v\n", err)
			os.Exit(1)
		}

		// Get the summary of the app.
		summaryPrompt := promptui.Prompt{
			Label:    "Summary (what's the short and sweet of what this app does?)",
			Validate: manifest.ValidateSummary,
		}
		summary, err := summaryPrompt.Run()
		if err != nil {
			fmt.Printf("app creation failed %v\n", err)
			os.Exit(1)
		}

		// Get the description of the app.
		descPrompt := promptui.Prompt{
			Label:    "Description (what's the long form of what this app does?)",
			Validate: manifest.ValidateDesc,
		}
		desc, err := descPrompt.Run()
		if err != nil {
			fmt.Printf("app creation failed %v\n", err)
			os.Exit(1)
		}

		// Get the author of the app.
		authorPrompt := promptui.Prompt{
			Label:    "Author (your name or your Github handle)",
			Validate: manifest.ValidateAuthor,
		}
		author, err := authorPrompt.Run()
		if err != nil {
			fmt.Printf("app creation failed %v\n", err)
			os.Exit(1)
		}

		// Generate app.
		g, err := generator.NewGenerator()
		if err != nil {
			fmt.Printf("app creation failed %v\n", err)
			os.Exit(1)
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
		err = g.GenerateApp(app)
		if err != nil {
			fmt.Printf("app creation failed %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("pixlet serve apps/%s/%s\n", app.PackageName, app.FileName)
	},
}

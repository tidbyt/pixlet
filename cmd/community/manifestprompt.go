//go:build !js && !wasm

package community

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"tidbyt.dev/pixlet/manifest"
)

func ManifestPrompt() (*manifest.Manifest, error) {
	// Get the name of the app.
	namePrompt := promptui.Prompt{
		Label:    "Name (what do you want to call your app?)",
		Validate: manifest.ValidateName,
	}
	name, err := namePrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("app creation failed %w", err)
	}

	// Get the summary of the app.
	summaryPrompt := promptui.Prompt{
		Label:    "Summary (what's the short and sweet of what this app does?)",
		Validate: manifest.ValidateSummary,
	}
	summary, err := summaryPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("app creation failed %w", err)
	}

	// Get the description of the app.
	descPrompt := promptui.Prompt{
		Label:    "Description (what's the long form of what this app does?)",
		Validate: manifest.ValidateDesc,
	}
	desc, err := descPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("app creation failed %w", err)
	}

	// Get the author of the app.
	authorPrompt := promptui.Prompt{
		Label:    "Author (your name or your Github handle)",
		Validate: manifest.ValidateAuthor,
	}
	author, err := authorPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("app creation failed %w", err)
	}

	return &manifest.Manifest{
		ID:      manifest.GenerateID(name),
		Name:    name,
		Summary: summary,
		Desc:    desc,
		Author:  author,
	}, nil
}

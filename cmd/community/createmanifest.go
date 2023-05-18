//go:build !js && !wasm

package community

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/manifest"
)

var CreateManifestCmd = &cobra.Command{
	Use:     "create-manifest <pathspec>",
	Short:   "Creates an app manifest from a prompt",
	Example: `  pixlet community create-manifest manifest.yaml`,
	Long:    `This command creates an app manifest by asking a series of prompts.`,
	Args:    cobra.ExactArgs(1),
	RunE:    CreateManifest,
}

func init() {
	CommunityCmd.AddCommand(CreateManifestCmd)
}

func CreateManifest(cmd *cobra.Command, args []string) error {
	fileName := filepath.Base(args[0])
	if fileName != manifest.ManifestFileName {
		return fmt.Errorf("supplied manifest must be named %s", manifest.ManifestFileName)
	}

	f, err := os.Create(args[0])
	if err != nil {
		return fmt.Errorf("couldn't open manifest: %w", err)
	}
	defer f.Close()

	m, err := ManifestPrompt()
	if err != nil {
		return fmt.Errorf("failed prompt: %w", err)
	}

	err = m.WriteManifest(f)
	if err != nil {
		return fmt.Errorf("couldn't write manifest: %w", err)
	}

	return nil
}

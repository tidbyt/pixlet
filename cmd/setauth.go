package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"tidbyt.dev/pixlet/cmd/config"
)

var SetAuthCmd = &cobra.Command{
	Use:     "set-auth",
	Short:   "Sets a custom access token in the private pixlet config.",
	Example: `  pixlet set-auth <token_json>`,
	Long: `This command sets a custom access token for use in subsequent runs. Normal users
should not need this - use 'pixlet login' instead.`,
	Args: cobra.ExactArgs(1),
	RunE: SetAuth,
}

func SetAuth(cmd *cobra.Command, args []string) error {
	authJSON := args[0]
	tok := &oauth2.Token{}
	err := json.Unmarshal([]byte(authJSON), tok)
	if err != nil {
		return fmt.Errorf("could not load auth JSON: %w", err)
	}

	config.PrivateConfig.Set("token", tok)
	if err := config.PrivateConfig.WriteConfig(); err != nil {
		return fmt.Errorf("could not persist auth token in config: %w", err)
	}

	return nil
}

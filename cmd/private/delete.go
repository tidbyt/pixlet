package private

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/cmd/config"
)

var deleteURL string
var deleteAppID string

func init() {
	DeleteCmd.Flags().StringVarP(&deleteURL, "url", "u", "https://api.tidbyt.com", "base URL of Tidbyt API")
	DeleteCmd.Flags().StringVarP(&deleteAppID, "app", "a", "", "ID of app to delete")
	DeleteCmd.MarkFlagRequired("app")
}

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a private app",
	Long:  `Deletes a private app, and attempt to uninstall it from owner's devices.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiToken := config.OAuthTokenFromConfig(cmd.Context())
		if apiToken == "" {
			return fmt.Errorf("login with `pixlet login` or use `pixlet set-auth` to configure auth")
		}

		requestURL := fmt.Sprintf("%s/v0/apps/%s", deleteURL, deleteAppID)
		req, err := http.NewRequest("DELETE", requestURL, nil)
		if err != nil {
			return fmt.Errorf("could not create http request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))

		client := http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("could not make HTTP request to %s: %w", requestURL, err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not read response body: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("request returned status %d with message: %s", resp.StatusCode, body)
		}

		return nil
	},
}

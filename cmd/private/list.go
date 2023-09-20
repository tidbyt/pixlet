package private

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/cmd/config"
)

type TidbytApp struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Version        string `json:"version"`
	Private        bool   `json:"private"`
	OrganizationID string `json:"organizationID"`
}

var listURL string

func init() {
	ListCmd.Flags().StringVarP(&listURL, "url", "u", "https://api.tidbyt.com", "base URL of Tidbyt API")
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists private apps",
	Long:  `Lists private apps, including team apps.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiToken := config.OAuthTokenFromConfig(cmd.Context())
		if apiToken == "" {
			return fmt.Errorf("login with `pixlet login` or use `pixlet set-auth` to configure auth")
		}

		requestURL := fmt.Sprintf("%s/v0/apps", listURL)
		req, err := http.NewRequest("GET", requestURL, nil)
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
		if resp.StatusCode != 200 {
			return fmt.Errorf("request returned status %d with message: %s", resp.StatusCode, body)
		}

		if err != nil {
			return fmt.Errorf("could not read response body: %w", err)
		}

		var apps struct {
			Apps []TidbytApp `json:"apps"`
		}
		err = json.Unmarshal(body, &apps)
		if err != nil {
			return fmt.Errorf("could not parse response body: %w", err)
		}

		for _, app := range apps.Apps {
			if !app.Private {
				continue
			}
			fmt.Printf("%s - %s", app.ID, app.Name)
			if app.OrganizationID != "" {
				fmt.Printf(" - team=%s", app.OrganizationID)
			}
			fmt.Printf("\n")
		}

		return nil
	},
}

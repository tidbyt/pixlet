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
	Name           string `json:"name,omitempty"`
	Version        string `json:"version,omitempty"`
	Private        bool   `json:"private,omitempty"`
	OrganizationID string `json:"organizationID,omitempty"`
}

type TidbytAppVersion struct {
	ID      string `json:"id"`
	Created string `json:"created,omitempty"`
}

var listURL string

func init() {
	ListCmd.Flags().StringVarP(&listURL, "url", "u", "https://api.tidbyt.com", "base URL of Tidbyt API")
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists private apps and versions",
	Long:  `Lists private apps, or available versions of a single private app.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiToken := config.OAuthTokenFromConfig(cmd.Context())
		if apiToken == "" {
			return fmt.Errorf("login with `pixlet login` or use `pixlet set-auth` to configure auth")
		}

		appID := ""
		if len(args) > 0 {
			appID = args[0]
		}

		var requestURL string
		if appID != "" {
			requestURL = fmt.Sprintf("%s/v0/apps/%s/versions", listURL, appID)
		} else {
			requestURL = fmt.Sprintf("%s/v0/apps", listURL)
		}

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

		if appID != "" {
			err = listVersions(body)
			if err != nil {
				return fmt.Errorf("could not list versions: %w", err)
			}
		} else {
			err = listApps(body)
			if err != nil {
				return fmt.Errorf("could not list apps: %w", err)
			}
		}

		return nil
	},
}

func listApps(body []byte) error {
	var apps struct {
		Apps []TidbytApp `json:"apps"`
	}
	err := json.Unmarshal(body, &apps)
	if err != nil {
		return fmt.Errorf("could not parse response body: %w", err)
	}

	for _, app := range apps.Apps {
		if !app.Private {
			continue
		}
		appJson, err := json.Marshal(app)
		if err != nil {
			return fmt.Errorf("could not marshal app: %w", err)
		}
		fmt.Println(string(appJson))
	}

	return nil
}

func listVersions(body []byte) error {
	var versions struct {
		Versions []TidbytAppVersion `json:"versions"`
	}
	err := json.Unmarshal(body, &versions)
	if err != nil {
		return fmt.Errorf("could not parse response body: %w", err)
	}

	for _, version := range versions.Versions {
		versionJson, err := json.Marshal(version)
		if err != nil {
			return fmt.Errorf("could not marshal version: %w", err)
		}
		fmt.Println(string(versionJson))
	}

	return nil
}

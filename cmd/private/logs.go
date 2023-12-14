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

type TidbytLogsResponse struct {
	Lines []struct {
		Timestamp string `json:"timestamp"`
		Message   string `json:"message"`
	} `json:"lines"`
}

var logsURL string
var logsAppID string

func init() {
	LogsCmd.Flags().StringVarP(&logsURL, "url", "u", "https://api.tidbyt.com", "base URL of Tidbyt API")
	LogsCmd.Flags().StringVarP(&logsAppID, "app", "a", "", "app ID to list versions for")
	LogsCmd.MarkFlagRequired("app")
}

var LogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Logs for private app",
	Long:  `Prints recent log lines for a private app`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiToken := config.OAuthTokenFromConfig(cmd.Context())
		if apiToken == "" {
			return fmt.Errorf("login with `pixlet login` or use `pixlet set-auth` to configure auth")
		}

		if logsAppID == "" {
			return fmt.Errorf("must specify app ID")
		}

		requestURL := fmt.Sprintf("%s/v0/apps/%s/logs", logsURL, logsAppID)
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

		var logs TidbytLogsResponse
		err = json.Unmarshal(body, &logs)
		if err != nil {
			return fmt.Errorf("could not parse response body: %w", err)
		}

		for _, line := range logs.Lines {
			fmt.Printf("%s %s\n", line.Timestamp, line.Message)
		}

		return nil
	},
}

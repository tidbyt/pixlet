package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var deployVersion string
var deployAppID string
var deployURL string

type TidbytAppDeploy struct {
	AppID   string `json:"appID"`
	Version string `json:"version"`
}

func init() {
	DeployCmd.Flags().StringVarP(&deployAppID, "app", "a", "", "app ID of the bundle to deploy")
	DeployCmd.MarkFlagRequired("app")
	DeployCmd.Flags().StringVarP(&deployVersion, "version", "v", "", "version of the bundle to deploy")
	DeployCmd.MarkFlagRequired("version")

	DeployCmd.Flags().StringVarP(&deployURL, "url", "u", "https://api.tidbyt.com", "base URL of the remote bundle store")
}

var DeployCmd = &cobra.Command{
	Use:     "deploy",
	Short:   "Deploys an app to production (internal only)",
	Example: `  pixlet deploy --app fuzzy-clock --version v0.0.1`,
	Long: `This command will deploy an app to production in the Tidbyt backend. Note, this
command is for internal use only at the moment, and normal API tokens will not
be able to deploy apps. We fully intend to make this command generally available
once our backend can support public deploys.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiToken := oauthTokenFromConfig(cmd.Context())
		if apiToken == "" {
			return fmt.Errorf("login with `pixlet login` or use `pixlet set-auth` to configure auth")
		}

		if deployAppID == "" {
			return fmt.Errorf("app must not be blank")
		}

		if deployVersion == "" {
			return fmt.Errorf("version must not be blank")
		}

		d := &TidbytAppDeploy{
			AppID:   deployAppID,
			Version: deployVersion,
		}

		b, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("could not create http request: %w", err)
		}

		requestURL := fmt.Sprintf("%s/v0/apps/%s/deploy", deployURL, deployAppID)
		req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(b))
		if err != nil {
			return fmt.Errorf("could not create http request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))

		client := http.Client{
			Timeout: 30 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("could not make HTTP request to %s: %w", requestURL, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("request returned status %d with message: %s", resp.StatusCode, body)
		}

		return nil
	},
}

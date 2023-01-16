package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

const (
	TidbytAPIPush = "https://api.tidbyt.com/v0/devices/%s/push"
	APITokenEnv   = "TIDBYT_API_TOKEN"
)

var (
	apiToken       string
	installationID string
	background     bool
)

type TidbytPushJSON struct {
	DeviceID       string `json:"deviceID"`
	Image          string `json:"image"`
	InstallationID string `json:"installationID"`
	Background     bool   `json:"background"`
}

func init() {
	PushCmd.Flags().StringVarP(&apiToken, "api-token", "t", "", "Tidbyt API token")
	PushCmd.Flags().StringVarP(&installationID, "installation-id", "i", "", "Give your installation an ID to keep it in the rotation")
	PushCmd.Flags().BoolVarP(&background, "background", "b", false, "Don't immediately show the image on the device")
}

var PushCmd = &cobra.Command{
	Use:   "push [device ID] [webp image]",
	Short: "Render a Pixlet script and push the WebP output to a Tidbyt",
	Args:  cobra.MinimumNArgs(2),
	RunE:  push,
}

func push(cmd *cobra.Command, args []string) error {
	deviceID := args[0]
	image := args[1]

	// TODO (mark): This is better served as a flag, but I don't want to break
	// folks in the short term. We should consider dropping this as an arguement
	// in a future release.
	if len(args) == 3 {
		installationID = args[2]
	}

	if apiToken == "" {
		apiToken = os.Getenv(APITokenEnv)
	}

	if apiToken == "" {
		apiToken = oauthTokenFromConfig(cmd.Context())
	}

	if apiToken == "" {
		return fmt.Errorf("blank Tidbyt API token (use `pixlet login`, set $%s or pass with --api-token)", APITokenEnv)
	}

	imageData, err := ioutil.ReadFile(image)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", image, err)
	}

	payload, err := json.Marshal(
		TidbytPushJSON{
			DeviceID:       deviceID,
			Image:          base64.StdEncoding.EncodeToString(imageData),
			InstallationID: installationID,
			Background:     background,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(TidbytAPIPush, deviceID),
		bytes.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("creating POST request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("pushing to API: %w", err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Tidbyt API returned status %s\n", resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		return fmt.Errorf("Tidbyt API returned status: %s", resp.Status)
	}

	return nil
}

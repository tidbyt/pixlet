package main

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
	apiToken string
)

type TidbytPushJSON struct {
	DeviceID string `json:"deviceID"`
	Image    string `json:"image"`
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringVarP(&apiToken, "api-token", "", "", "Tidbyt API token")
}

var pushCmd = &cobra.Command{
	Use:   "push [device ID] [webp image]",
	Short: "Pushes a webp image to a Tidbyt device",
	Args:  cobra.MinimumNArgs(2),
	Run:   push,
}

func push(cmd *cobra.Command, args []string) {
	deviceID := args[0]
	image := args[1]

	if apiToken == "" {
		apiToken = os.Getenv(APITokenEnv)
	}

	if apiToken == "" {
		fmt.Printf("blank Tidbyt API token (set $%s or pass with --api-token)\n", APITokenEnv)
		os.Exit(1)
	}

	imageData, err := ioutil.ReadFile(image)
	if err != nil {
		fmt.Printf("failed to read file %s: %v\n", image, err)
		os.Exit(1)
	}

	payload, err := json.Marshal(
		TidbytPushJSON{
			DeviceID: deviceID,
			Image:    base64.StdEncoding.EncodeToString(imageData),
		},
	)
	if err != nil {
		fmt.Printf("failed to marshal json: %v\n", err)
		os.Exit(1)
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(TidbytAPIPush, deviceID),
		bytes.NewReader(payload),
	)
	if err != nil {
		fmt.Printf("creating POST request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("pushing to API: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Tidbyt API returned status %s\n", resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		os.Exit(1)
	}

	os.Exit(0)
}

//go:build !js && !wasm

package private

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/cmd/community"
	"tidbyt.dev/pixlet/cmd/config"
	"tidbyt.dev/pixlet/tools/generator"
)

var createOrg string
var createURL string
var createDir string

type TidbytCreateAppRequest struct {
	OrganizationID string `json:"organizationID"`
	Private        bool   `json:"private"`
}

type TidbytCreateAppReply struct {
	AppID string `json:"appID"`
}

func init() {
	CreateCmd.Flags().StringVarP(&createOrg, "org", "o", "", "organization to create the app in")
	CreateCmd.Flags().StringVarP(&createURL, "url", "u", "https://api.tidbyt.com", "base URL of Tidbyt API")
	CreateCmd.Flags().StringVarP(&createDir, "app-dir", "a", ".", "directory to create the app in")
}

// CreateCmd prompts the user for info and generates a new app.
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new app",
	Long:  `This command will prompt for all of the information we need to generate a new Tidbyt app. No flags are necessary unless you are creating a private app, which is only available with our Tidbyt For Teams offering.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Make sure user is authenticated
		apiToken := config.OAuthTokenFromConfig(cmd.Context())
		if apiToken == "" {
			return fmt.Errorf("login with `pixlet login` or use `pixlet set-auth` to configure auth")
		}

		// Prompt the user for input
		app, err := community.ManifestPrompt()
		if err != nil {
			return fmt.Errorf("app creation, couldn't get user input: %w", err)
		}

		// Create a private app in backend
		app.ID, err = createPrivateApp(apiToken, createOrg)
		if err != nil {
			if strings.Contains(err.Error(), "user is not authorized to create apps") {
				return fmt.Errorf("user is not authorized to create apps for organization %s, please reach out to your Tidbyt For Teams account representative to enable this feature for your account", createOrg)
			}

			return fmt.Errorf("remote app creation failed: %w", err)
		}

		// Generate app
		g, err := generator.NewGenerator(generator.Local, createDir)
		if err != nil {
			return fmt.Errorf("app creation failed %w", err)
		}
		starlarkPath, err := g.GenerateApp(app)
		if err != nil {
			return fmt.Errorf("app creation failed: %w", err)
		}

		starlarkPathAbs, err := filepath.Abs(starlarkPath)
		if err != nil {
			return fmt.Errorf("app was created, but we don't know where: %w", err)
		}

		// Let the user know where the app is and how to use it.
		fmt.Println("")
		fmt.Println("App created at:")
		fmt.Printf("\t%s\n", starlarkPathAbs)
		fmt.Println("")
		fmt.Println("To start the app, run:")
		fmt.Printf("\tpixlet serve %s\n", starlarkPath)
		fmt.Println("")
		fmt.Println("For docs, head to:")
		fmt.Printf("\thttps://tidbyt.dev\n")
		fmt.Println("")
		fmt.Println("To upload and deploy your app:")
		if createDir == "." {
			fmt.Printf("\tpixlet private upload\n")
		} else {
			fmt.Printf("\tpixlet private upload --app-dir %s\n", createDir)
		}

		return nil
	},
}

// Creates private app for organization. If org is blank, attempts to
// create a per-user private app.
func createPrivateApp(apiToken string, org string) (string, error) {
	createAppRequest := &TidbytCreateAppRequest{
		OrganizationID: org,
		Private:        org == "",
	}

	b, err := json.Marshal(createAppRequest)
	if err != nil {
		return "", fmt.Errorf("could not create http request: %w", err)
	}

	requestURL := fmt.Sprintf("%s/v0/apps", createURL)
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(b))
	if err != nil {
		return "", fmt.Errorf("could not create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not make HTTP request to %s: %w", requestURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request returned status %d with message: %s", resp.StatusCode, body)
	}

	createAppReply := &TidbytCreateAppReply{}
	err = json.NewDecoder(resp.Body).Decode(&createAppReply)
	if err != nil {
		return "", fmt.Errorf("could not decode response: %w", err)
	}

	return createAppReply.AppID, nil
}

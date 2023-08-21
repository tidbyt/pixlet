//go:build !js && !wasm

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/cmd/community"
	"tidbyt.dev/pixlet/tools/generator"
	"tidbyt.dev/pixlet/tools/repo"
)

var createPrivate bool
var createOrg string
var createURL string

type TidbytCreateAppRequest struct {
	OrganizationID string `json:"organizationID"`
}

type TidbytCreateAppReply struct {
	AppID string `json:"appID"`
}

func init() {
	CreateCmd.Flags().StringVarP(&createOrg, "org", "o", "", "organization to create the app in")
	CreateCmd.Flags().BoolVarP(&createPrivate, "private", "p", false, "create a private app")
	CreateCmd.Flags().StringVarP(&createURL, "url", "u", "https://api.tidbyt.com", "base URL of the remote bundle store")
}

// CreateCmd prompts the user for info and generates a new app.
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new app",
	Long:  `This command will prompt for all of the information we need to generate a new Tidbyt app. No flags are necessary unless you are creating a private app, which is only available with our Tidbyt For Teams offering.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the current working directory.
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("app creation failed, something went wrong with your local filesystem: %w", err)
		}

		// Determine what type of app this is an what the root should be.
		var root string
		var appType generator.AppType
		if repo.IsInRepo(cwd, "community") {
			appType = generator.Community
			root, err = repo.RepoRoot(cwd)
			if err != nil {
				return fmt.Errorf("app creation failed, something went wrong with your community repo: %w", err)
			}
		} else if repo.IsInRepo(cwd, "tidbyt") {
			appType = generator.Internal
			root, err = repo.RepoRoot(cwd)
			if err != nil {
				return fmt.Errorf("app creation failed, something went wrong with your tidbyt repo: %w", err)
			}
		} else {
			appType = generator.Local
			root = cwd
		}

		// Prompt the user for input.
		app, err := community.ManifestPrompt()
		if err != nil {
			return fmt.Errorf("app creation, couldn't get user input: %w", err)
		}

		if createPrivate {
			// create a private app
			apiToken := oauthTokenFromConfig(cmd.Context())
			if apiToken == "" {
				return fmt.Errorf("login with `pixlet login` or use `pixlet set-auth` to configure auth")
			}

			if createOrg == "" {
				return fmt.Errorf("organization must not be blank")
			}

			app.ID, err = createPrivateApp(apiToken, createOrg)
			if err != nil {
				return fmt.Errorf("remote app creation failed: %w", err)
			}
		}

		// Generate app.
		g, err := generator.NewGenerator(appType, root)
		if err != nil {
			return fmt.Errorf("app creation failed %w", err)
		}
		absolutePath, err := g.GenerateApp(app)
		if err != nil {
			return fmt.Errorf("app creation failed: %w", err)
		}

		// Get the relative path from where the user started. Note, we're not
		// using the root here, given the root can be git repo specific.
		relativePath, err := filepath.Rel(cwd, absolutePath)
		if err != nil {
			return fmt.Errorf("app was created, but we don't know where: %w", err)
		}

		// Let the user know where the app is and how to use it.
		fmt.Println("")
		fmt.Println("App created at:")
		fmt.Printf("\t%s\n", absolutePath)
		fmt.Println("")
		fmt.Println("To start the app, run:")
		fmt.Printf("\tpixlet serve %s\n", relativePath)
		fmt.Println("")
		fmt.Println("For docs, head to:")
		fmt.Printf("\thttps://tidbyt.dev\n")

		if createPrivate {
			fmt.Println("")
			fmt.Println("To deploy your app:")
			fmt.Printf("\tpixlet bundle ./\n")
			fmt.Printf("\tpixlet upload bundle.tar.gz --app %s --version v0.0.1\n", app.ID)
			fmt.Printf("\tpixlet deploy --app %s --version v0.0.1\n", app.ID)
		}
		return nil
	},
}

func createPrivateApp(apiToken string, org string) (string, error) {
	createAppRequest := &TidbytCreateAppRequest{
		OrganizationID: org,
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

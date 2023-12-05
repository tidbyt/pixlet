package private

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/bundle"
	"tidbyt.dev/pixlet/cmd/config"
)

var uploadVersion string
var uploadDir string
var uploadURL string
var uploadSkipDeploy bool

var defaultVersion = fmt.Sprintf("%d", time.Now().Unix())

type TidbytBundleUpload struct {
	AppID   string `json:"appID"`
	Version string `json:"version"`
	Bundle  string `json:"bundle"`
}

func init() {
	UploadCmd.Flags().StringVarP(&uploadDir, "app-dir", "d", ".", "app directory to upload")
	UploadCmd.Flags().StringVarP(&uploadVersion, "version", "v", defaultVersion, "version of the app")
	UploadCmd.Flags().StringVarP(&uploadURL, "url", "u", "https://api.tidbyt.com", "base URL of Tidbyt API")
	UploadCmd.Flags().BoolVarP(&uploadSkipDeploy, "skip-deploy", "s", false, "skip deploying the bundle after uploading")
}

var UploadCmd = &cobra.Command{
	Use:     "upload",
	Short:   "Uploads an app to Tidbyt",
	Example: "  pixlet private upload  --app-dir app/startrek/ --version v0.0.1",
	Long: `This command uploads your private app. By default, the app is assumed
to be in the current directory, is given timestamp as version, and is deployed.
These defaults can be overridden with the --app-dir, --version, and --skip-deploy
flags.`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Some sanity checks
		if uploadVersion == "" {
			return fmt.Errorf("version must not be blank")
		}
		if uploadDir == "" {
			return fmt.Errorf("app-dir must not be blank")
		}
		info, err := os.Stat(uploadDir)
		if err != nil {
			return fmt.Errorf("bad app-dir: %w", err)
		}
		if !info.IsDir() {
			return fmt.Errorf("app-dir must be a directory")
		}

		// Create bundle
		buf := &bytes.Buffer{}
		ab, err := bundle.InitFromPath(uploadDir)
		if err != nil {
			return fmt.Errorf("could not init bundle: %w", err)
		}
		err = ab.WriteBundle(buf)
		if err != nil {
			return err
		}

		// Authenticate
		apiToken := config.OAuthTokenFromConfig(cmd.Context())
		if apiToken == "" {
			return fmt.Errorf("login with `pixlet login` or use `pixlet set-auth` to configure auth")
		}

		// Marshal request
		uploadBundle := &TidbytBundleUpload{
			AppID:   ab.Manifest.ID,
			Version: uploadVersion,
			Bundle:  base64.StdEncoding.EncodeToString(buf.Bytes()),
		}
		body, err := json.Marshal(uploadBundle)
		if err != nil {
			return fmt.Errorf("could not marshal request: %w", err)
		}

		// Upload
		requestURL := fmt.Sprintf("%s/v0/apps/%s/upload", uploadURL, ab.Manifest.ID)
		req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("could not create upload request: %w", err)
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

		// Possibly done
		if uploadSkipDeploy {
			fmt.Printf("Uploaded version %s of app %s\n", uploadVersion, ab.Manifest.ID)
			return nil
		}

		// Otherwise: deploy

		// Build deploy request
		d := &TidbytAppDeploy{
			AppID:   ab.Manifest.ID,
			Version: uploadVersion,
		}
		body, err = json.Marshal(d)
		if err != nil {
			return fmt.Errorf("could not create deploy request: %w", err)
		}

		// Deploy
		requestURL = fmt.Sprintf("%s/v0/apps/%s/deploy", uploadURL, ab.Manifest.ID)
		req, err = http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(body))
		if err != nil {
			return fmt.Errorf("could not create http request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))

		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("could not make HTTP request to %s: %w", requestURL, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("request returned status %d with message: %s", resp.StatusCode, body)
		}

		fmt.Printf("Uploaded and deployed version %s of app %s\n", uploadVersion, ab.Manifest.ID)

		return nil
	},
}

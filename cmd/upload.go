package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/bundle"
)

var uploadVersion string
var uploadAppID string
var uploadURL string
var uploadAPIToken string

func init() {
	UploadCmd.Flags().StringVarP(&uploadAppID, "app", "a", "", "app ID of the bundle to upload")
	UploadCmd.MarkFlagRequired("app")
	UploadCmd.Flags().StringVarP(&uploadVersion, "version", "v", "", "version of the bundle to upload")
	UploadCmd.MarkFlagRequired("version")
	UploadCmd.Flags().StringVarP(&uploadAPIToken, "token", "t", "", "API token to use when uploading the bundle to the remote store")
	UploadCmd.MarkFlagRequired("token")

	UploadCmd.Flags().StringVarP(&uploadURL, "url", "u", "https://api.tidbyt.com", "base URL of the remote bundle store")
}

var UploadCmd = &cobra.Command{
	Use:     "upload",
	Short:   "Uploads an app bundle to Tidbyt (internal only)",
	Example: `  pixlet upload bundle.tar.gz --app fuzzy-clock --version v0.0.1 --token {{ api_token }}`,
	Long: `This command will upload an app bundle (see pixlet bundle) using the specified
app ID and version. Note, this is for internal use only at the moment, and
normal API tokens will not work with this command. We fully intend to make this
command public once our backend is well positioned to support it.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bundleFile := args[0]
		info, err := os.Stat(bundleFile)
		if err != nil {
			return fmt.Errorf("input bundle file invalid: %w", err)
		}

		if info.IsDir() {
			return fmt.Errorf("input bundle must be a file")
		}

		if !strings.HasSuffix(bundleFile, "tar.gz") {
			return fmt.Errorf("input bundle format is not correct, did you create it with `pixlet bundle`?")
		}

		f, err := os.Open(bundleFile)
		if err != nil {
			return fmt.Errorf("could not open bundle: %w", err)
		}
		defer f.Close()

		// Load the bundle to ensure it's valid.
		ab, err := bundle.LoadBundle(f)
		if err != nil {
			return fmt.Errorf("could not load bundle: %w", err)
		}

		// Re-write the bundle to ensure it's standard.
		buf := &bytes.Buffer{}
		err = ab.WriteBundle(buf)
		if err != nil {
			return fmt.Errorf("could not re-create bundle: %w", err)
		}

		// TODO: check URL with mats.
		requestURL := fmt.Sprintf("%s/v0/apps/%s/%s", uploadURL, uploadAppID, uploadVersion)
		req, err := http.NewRequest(http.MethodPost, requestURL, buf)
		if err != nil {
			return fmt.Errorf("could not create http request: %w", err)
		}

		req.Header.Set("Content-Type", "application/gzip")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))

		client := http.Client{
			Timeout: 30 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("could not make HTTP request to %s: %w", requestURL, err)
		}
		defer resp.Body.Close()

		statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
		if !statusOK {
			return fmt.Errorf("remote returned an error with code: %d", resp.StatusCode)
		}

		return nil
	},
}

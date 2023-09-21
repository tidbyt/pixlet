package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"tidbyt.dev/pixlet/cmd/config"
)

const (
	TidbytAPIList = "https://api.tidbyt.com/v0/devices/%s/installations"
)

type TidbytInstallationJSON struct {
	Id    string `json:"id"`
	AppId string `json:"appID"`
}

type TidbytInstallationListJSON struct {
	Installations []TidbytInstallationJSON `json:"installations"`
}

func init() {
	ListCmd.Flags().StringVarP(&apiToken, "api-token", "t", "", "Tidbyt API token")
}

var ListCmd = &cobra.Command{
	Use:   "list [device ID]",
	Short: "Lists all apps installed on a Tidbyt",
	Args:  cobra.MinimumNArgs(1),
	RunE:  listInstallations,
}

func listInstallations(cmd *cobra.Command, args []string) error {
	deviceID := args[0]

	if apiToken == "" {
		apiToken = os.Getenv(APITokenEnv)
	}

	if apiToken == "" {
		apiToken = config.OAuthTokenFromConfig(cmd.Context())
	}

	if apiToken == "" {
		return fmt.Errorf("blank Tidbyt API token (use `pixlet login`, set $%s or pass with --api-token)", APITokenEnv)
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(TidbytAPIList, deviceID), nil)
	if err != nil {
		return fmt.Errorf("creating GET request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("listing installations from API: %w", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Printf("Tidbyt API returned status %s\n", resp.Status)
		fmt.Println(string(body))
		return fmt.Errorf("Tidbyt API returned status: %s", resp.Status)
	}

	var installations TidbytInstallationListJSON
	err = json.Unmarshal(body, &installations)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 22, 8, 0, '\t', 0)
	defer w.Flush()

	for _, inst := range installations.Installations {
		fmt.Fprintf(w, "%s\t%s\n", inst.AppId, inst.Id)
	}

	return nil
}

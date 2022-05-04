package browser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	TidbytAPIPush = "https://api.tidbyt.com/v0/devices/%s/push"
)

type TidbytPushJSON struct {
	DeviceID       string `json:"deviceID"`
	Image          string `json:"image"`
	InstallationID string `json:"installationID"`
	Background     bool   `json:"background"`
}

func (b *Browser) pushHandler(w http.ResponseWriter, r *http.Request) {
	var (
		deviceID       string
		apiToken       string
		installationID string
		background     bool
	)

	var result map[string]interface{}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}

	json.Unmarshal(bodyBytes, &result)

	config := make(map[string]string)
	for k, val := range result {
		switch k {
		case "deviceID":
			deviceID = val.(string)
		case "apiToken":
			apiToken = val.(string)
		case "installationID":
			installationID = val.(string)
		case "background":
			background = val.(string) == "true"
		default:
			config[k] = val.(string)
		}
	}

	webp, err := b.loader.LoadApplet(config)

	payload, err := json.Marshal(
		TidbytPushJSON{
			DeviceID:       deviceID,
			Image:          webp,
			InstallationID: installationID,
			Background:     background,
		},
	)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(TidbytAPIPush, deviceID),
		bytes.NewReader(payload),
	)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Tidbyt API returned status %s\n", resp.Status)
		w.WriteHeader(resp.StatusCode)
		fmt.Fprintln(w, err)

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

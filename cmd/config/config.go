package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	OAuthCallbackAddr = "localhost:8085"
)

var (
	PrivateConfig = viper.New()

	OAuthConf = &oauth2.Config{
		ClientID: "d8ae7ea0-4a1a-46b0-b556-6d742687223a",
		Scopes:   []string{"device", "offline_access", "app-admin"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.tidbyt.com/oauth2/auth",
			TokenURL: "https://login.tidbyt.com/oauth2/token",
		},
		RedirectURL: fmt.Sprintf("http://%s", OAuthCallbackAddr),
	}
)

func init() {
	if ucd, err := os.UserConfigDir(); err == nil {
		configPath := filepath.Join(ucd, "tidbyt")

		if err := os.MkdirAll(configPath, os.ModePerm); err == nil {
			PrivateConfig.AddConfigPath(configPath)
		}
	}

	PrivateConfig.SetConfigName("private")
	PrivateConfig.SetConfigType("yaml")
	PrivateConfig.SetConfigPermissions(0600)

	PrivateConfig.SafeWriteConfig()
	PrivateConfig.ReadInConfig()
}

func OAuthTokenFromConfig(ctx context.Context) string {
	if !PrivateConfig.IsSet("token") {
		return ""
	}

	var tok oauth2.Token
	if err := PrivateConfig.UnmarshalKey("token", &tok); err != nil {
		fmt.Println("unmarshaling API token from config:", err)
		os.Exit(1)
	}

	if !tok.Valid() {
		// probably expired, try to refresh
		ts := OAuthConf.TokenSource(ctx, &tok)
		refreshed, err := ts.Token()
		if err != nil {
			fmt.Println("refreshing API token:", err)
			os.Exit(1)
		}

		tok = *refreshed
		PrivateConfig.Set("token", tok)
		PrivateConfig.WriteConfig()
	}

	return tok.AccessToken
}

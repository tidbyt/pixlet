package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var (
	privateConfig = viper.New()
)

func init() {
	if ucd, err := os.UserConfigDir(); err == nil {
		configPath := filepath.Join(ucd, "tidbyt")

		if err := os.MkdirAll(configPath, os.ModePerm); err == nil {
			privateConfig.AddConfigPath(configPath)
		}
	}

	privateConfig.SetConfigName("private")
	privateConfig.SetConfigType("yaml")
	privateConfig.SetConfigPermissions(0600)

	privateConfig.SafeWriteConfig()
	privateConfig.ReadInConfig()
}

func oauthTokenFromConfig(ctx context.Context) string {
	if !privateConfig.IsSet("token") {
		return ""
	}

	var tok oauth2.Token
	if err := privateConfig.UnmarshalKey("token", &tok); err != nil {
		fmt.Println("unmarshaling API token from config:", err)
		os.Exit(1)
	}

	if !tok.Valid() {
		// probably expired, try to refresh
		ts := oauthConf.TokenSource(ctx, &tok)
		refreshed, err := ts.Token()
		if err != nil {
			fmt.Println("refreshing API token:", err)
			os.Exit(1)
		}

		tok = *refreshed
		privateConfig.Set("token", tok)
		privateConfig.WriteConfig()
	}

	return tok.AccessToken
}

package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
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

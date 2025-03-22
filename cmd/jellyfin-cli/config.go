package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"codeberg.org/jfenske/jellyfin-cli/internal/libraries/api"
)

func loadConfiguration(atom zap.AtomicLevel, logger *zap.SugaredLogger) (api.JellyfinApiConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/jellyfin-cli/")
	viper.AddConfigPath("$HOME/.config/jellyfin-cli/")
	viper.AddConfigPath(".")

	viper.SetDefault("api.base_url", "http://127.0.0.1:8096")

	err := viper.ReadInConfig()
	if err != nil {
		return api.JellyfinApiConfig{}, fmt.Errorf("failed to read config: %w", err)
	}

	if rawLevel := viper.GetString("logging.level"); rawLevel != "" {
		if level, err := zap.ParseAtomicLevel(strings.ToUpper(rawLevel)); err != nil {
			logger.Warnw("invalid logging level", "level", rawLevel)
		} else {
			atom.SetLevel(level.Level())
		}
	}

	return api.JellyfinApiConfig{
		BaseUrl:       viper.GetString("api.base_url"),
		Token:         viper.GetString("api.token"),
		SkipSslVerify: viper.GetBool("api.insecure"),
	}, nil
}

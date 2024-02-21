package main

import (
	"fmt"

	"github.com/spf13/viper"

	"codeberg.org/jfenske/jellyfin-cli/api"
)

func loadConfiguration() (*api.JellyfinApiConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/jellyfin-cli/")
	viper.AddConfigPath("$HOME/.config/jellyfin-cli/")
	viper.AddConfigPath(".")

	viper.SetDefault("api.base_url", "http://127.0.0.1:8096")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &api.JellyfinApiConfig{
		BaseUrl:       viper.GetString("api.base_url"),
		Token:         viper.GetString("api.token"),
		SkipSslVerify: viper.GetBool("api.insecure"),
	}, nil
}

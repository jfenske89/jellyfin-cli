package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/jfenske89/jellyfin-cli/pkg/client"
	"github.com/jfenske89/jellyfin-cli/pkg/models"
)

// Version is set during build
var Version = "dev"

var (
	cfgFile string
	config  models.Config
	logger  *zap.SugaredLogger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "jellyfin-cli",
	Short:   "Interact Jellyfin from the command line",
	Long:    `Jellyfin CLI is a command-line tool for interacting with a Jellyfin server.`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/jellyfin-cli/config.yaml)")
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")

	// Bind flags to viper
	err := viper.BindPFlag("output.json", rootCmd.PersistentFlags().Lookup("json"))
	if err != nil {
		logger.Errorw("failed to bind flags", "error", err.Error())
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Initialize logger
	var atomicLevel zap.AtomicLevel
	logger, atomicLevel = initLogger()

	// Set the global atomic level
	globalAtomicLevel = atomicLevel
	defer func() {
		// Sync logger on exit
		_ = logger.Sync()
	}()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in default locations
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.config/jellyfin-cli/")
		viper.AddConfigPath("/etc/jellyfin-cli/")
	}

	// Set defaults
	viper.SetDefault("api.base_url", "http://127.0.0.1:8096")
	viper.SetDefault("logging.level", "INFO")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Errorw("Failed to read config file", "error", err)
		} else {
			logger.Warnw("No config file found, using defaults")
		}
	} else {
		// Set Jellyfin config from viper
		config.Jellyfin.BaseURL = viper.GetString("api.base_url")
		config.Jellyfin.Token = viper.GetString("api.token")
		config.Jellyfin.SkipSSLVerify = viper.GetBool("api.insecure")

		// Update logger level based on config
		config.Logging.Level = viper.GetString("logging.level")
		updateLogLevel(config.Logging.Level)

		logger.Debugw("Using config file", "file", viper.ConfigFileUsed())
	}
}

// getClient returns a new Jellyfin API client using the current configuration
func getClient() client.Client {
	return client.NewClient(config.Jellyfin, logger)
}

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jfenske89/jellyfin-cli/pkg/models"
)

// librariesCmd represents the libraries command
var librariesCmd = &cobra.Command{
	Use:   "libraries",
	Short: "List library folders on the Jellyfin server",
	Long:  `List all library folders (virtual folders) on the Jellyfin server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get client
		client := getClient()

		// Get command flags
		outputJSON, _ := cmd.Flags().GetBool("json")

		// Get library folders
		libraries, err := client.ListLibraryFolders(cmd.Context(), nil)
		if err != nil {
			return fmt.Errorf("failed to list library folders: %w", err)
		}

		// Output
		if outputJSON {
			outputLibrariesJSON(libraries)
		} else {
			outputLibrariesText(libraries)
		}

		return nil
	},
}

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the library",
	Long:  `Refresh the Jellyfin library to scan for new content.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get client
		client := getClient()

		// Refresh library
		err := client.RefreshLibrary(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to refresh library: %w", err)
		}

		fmt.Println("Library refresh initiated successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(librariesCmd)

	// Add refresh as a subcommand of libraries
	librariesCmd.AddCommand(refreshCmd)
}

// outputLibrariesText outputs libraries in human-readable format
func outputLibrariesText(libraries []models.LibraryFolder) {
	if len(libraries) == 0 {
		fmt.Println("No library folders found")
		return
	}

	fmt.Println("Library Folders:")
	for _, library := range libraries {
		fmt.Printf(" - %s (Type: %s)\n", library.Name, library.CollectionType)
		if library.RefreshStatus != "" {
			fmt.Printf("   Status: %s\n", library.RefreshStatus)
		}
	}
}

// outputLibrariesJSON outputs libraries in JSON format
func outputLibrariesJSON(libraries []models.LibraryFolder) {
	jsonBytes, err := json.MarshalIndent(libraries, "", "  ")
	if err != nil {
		logger.Errorw("Failed to marshal libraries to JSON", "error", err)
		return
	}

	fmt.Println(string(jsonBytes))
}

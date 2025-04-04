package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jfenske89/jellyfin-cli/pkg/models"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for content on the Jellyfin server",
	Long: `Search for content on the Jellyfin server using a text query.
	
You can filter results by type using the --type flag.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get client
		client := getClient()

		// Get command flags
		itemType, _ := cmd.Flags().GetString("type")
		limit, _ := cmd.Flags().GetInt("limit")
		outputJSON, _ := cmd.Flags().GetBool("json")

		// Combine all args into a single search query
		query := strings.Join(args, " ")

		// Set up parameters
		params := make(map[string]string)
		if itemType != "" {
			params["includeItemTypes"] = itemType
		}
		if limit > 0 {
			params["limit"] = fmt.Sprintf("%d", limit)
		}

		// Search
		results, err := client.Search(cmd.Context(), query, params)
		if err != nil {
			return fmt.Errorf("failed to search: %w", err)
		}

		// Output
		if outputJSON {
			outputSearchJSON(results)
		} else {
			outputSearchText(results, query)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Add local flags
	searchCmd.Flags().StringP("type", "t", "", "Filter by item type (Movie, Series, Episode, etc.)")
	searchCmd.Flags().IntP("limit", "l", 10, "Limit the number of results")
}

// outputSearchText outputs search results in human-readable format
func outputSearchText(results *models.SearchResponse, query string) {
	if results == nil || len(results.SearchHints) == 0 {
		fmt.Printf("No results found for '%s'\n", query)
		return
	}

	fmt.Printf("Search Results for '%s' (Found: %d):\n", query, results.TotalHints)

	for i, hint := range results.SearchHints {
		// Format result based on type
		switch hint.Type {
		case "Movie":
			year := ""
			if hint.ProductYear > 0 {
				year = fmt.Sprintf(" (%d)", hint.ProductYear)
			}
			fmt.Printf(" %d. [Movie] %s%s\n", i+1, hint.Name, year)

		case "Series":
			year := ""
			if hint.ProductYear > 0 {
				year = fmt.Sprintf(" (%d)", hint.ProductYear)
			}
			fmt.Printf(" %d. [Series] %s%s\n", i+1, hint.Name, year)

		case "Episode":
			episodeInfo := ""
			if hint.SeasonNum > 0 && hint.EpisodeNum > 0 {
				episodeInfo = fmt.Sprintf(" (S%02dE%02d)", hint.SeasonNum, hint.EpisodeNum)
			}
			fmt.Printf(" %d. [Episode] %s - %s%s\n", i+1, hint.SeriesName, hint.Name, episodeInfo)

		default:
			fmt.Printf(" %d. [%s] %s\n", i+1, hint.Type, hint.Name)
		}
	}
}

// outputSearchJSON outputs search results in JSON format
func outputSearchJSON(results *models.SearchResponse) {
	jsonBytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		logger.Errorw("Failed to marshal search results to JSON", "error", err)
		return
	}

	fmt.Println(string(jsonBytes))
}

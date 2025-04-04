package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/jfenske89/jellyfin-cli/pkg/models"
)

// activityCmd represents the activity command
var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "List activity logs from the Jellyfin server",
	Long: `List recent activity logs from the Jellyfin server.
	
You can limit the number of results using the --limit flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get client
		client := getClient()

		// Get command flags
		limit, _ := cmd.Flags().GetInt("limit")
		outputJSON, _ := cmd.Flags().GetBool("json")

		// Set up parameters
		params := make(map[string]string)
		if limit > 0 {
			params["limit"] = strconv.Itoa(limit)
		}

		// Get activity logs
		logs, err := client.ListActivityLogs(cmd.Context(), params)
		if err != nil {
			return fmt.Errorf("failed to list activity logs: %w", err)
		}

		// Output
		if outputJSON {
			outputActivityJSON(logs)
		} else {
			outputActivityText(logs)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(activityCmd)

	// Add local flags
	activityCmd.Flags().IntP("limit", "l", 10, "Limit the number of results")
}

// outputActivityText outputs activity logs in human-readable format
func outputActivityText(logs *models.ActivityLog) {
	if logs == nil || len(logs.Items) == 0 {
		fmt.Println("No activity logs found")
		return
	}

	fmt.Printf("Activity Logs (Total: %d):\n", logs.TotalCount)

	for i, item := range logs.Items {
		timeAgo := humanize.RelTime(time.Now(), item.DateCreatedUTC, "", "ago")
		fmt.Printf(" %d. [%s] %s - %s (%s)\n",
			i+1,
			item.Severity,
			item.Name,
			item.ShortOverview,
			timeAgo)

		// Print the full overview if it's different from the short overview
		if item.Overview != "" && item.Overview != item.ShortOverview {
			fmt.Printf("    %s\n", item.Overview)
		}
	}
}

// outputActivityJSON outputs activity logs in JSON format
func outputActivityJSON(logs *models.ActivityLog) {
	jsonBytes, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		logger.Errorw("Failed to marshal activity logs to JSON", "error", err)
		return
	}

	fmt.Println(string(jsonBytes))
}

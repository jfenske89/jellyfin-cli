package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/jfenske89/jellyfin-cli/pkg/models"
)

// sessionsCmd represents the sessions command
var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "List active sessions on the Jellyfin server",
	Long: `List active or all sessions on the Jellyfin server.
	
By default, it shows all sessions. Use the --active flag to show only active sessions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get client
		client := getClient()

		// Get command flags
		active, _ := cmd.Flags().GetBool("active")
		outputJSON, _ := cmd.Flags().GetBool("json")

		// Set up parameters
		params := make(map[string]string)
		if active {
			params["activeWithinSeconds"] = "600"
		}

		// Get sessions
		sessions, err := client.ListSessions(cmd.Context(), params)
		if err != nil {
			return fmt.Errorf("failed to list sessions: %w", err)
		}

		// Filter active sessions if needed
		if active {
			activeSessions := make([]models.Session, 0, len(sessions))
			for _, s := range sessions {
				if time.Since(s.LastActivityUTC) <= 10*time.Minute {
					activeSessions = append(activeSessions, s)
				}
			}
			sessions = activeSessions
		}

		// Output
		if outputJSON {
			outputSessionsJSON(sessions)
		} else {
			outputSessionsText(sessions)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sessionsCmd)

	// Add local flags
	sessionsCmd.Flags().BoolP("active", "a", false, "Only show active sessions")
}

// outputSessionsText outputs sessions in human-readable format
func outputSessionsText(sessions []models.Session) {
	if len(sessions) == 0 {
		fmt.Println("No sessions found")
		return
	}

	fmt.Println("Sessions:")
	for _, session := range sessions {
		duration := humanize.RelTime(time.Now(), session.LastActivityUTC, "", "ago")
		fmt.Printf(" - %s on %s (%s)\n", session.UserName, session.DeviceName, duration)
	}
}

// outputSessionsJSON outputs sessions in JSON format
func outputSessionsJSON(sessions []models.Session) {
	jsonBytes, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		logger.Errorw("Failed to marshal sessions to JSON", "error", err)
		return
	}

	fmt.Println(string(jsonBytes))
}

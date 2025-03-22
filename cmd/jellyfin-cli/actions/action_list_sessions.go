package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"go.uber.org/zap"

	"codeberg.org/jfenske/jellyfin-cli/internal/libraries/api"
)

type listSessionsExecutorImpl struct {
	client api.JellyfinApiClient
	logger *zap.SugaredLogger
}

func NewListSessionsExecutor(client api.JellyfinApiClient, logger *zap.SugaredLogger) Executor {
	return &listSessionsExecutorImpl{
		client: client,
		logger: logger,
	}
}

func (e *listSessionsExecutorImpl) Run(ctx context.Context, options map[string]string) error {
	output := "text"
	if val, ok := options["output"]; ok && val != "" {
		output = val
	}

	getParameters := make(map[string]string)
	if val, ok := options["active-only"]; ok && val != "" {
		if activeOnly, _ := strconv.ParseBool(val); activeOnly {
			getParameters["activeWithinSeconds"] = "600"
		}
	}

	sessions, err := e.client.ListSessions(ctx, getParameters)
	if err != nil {
		return err
	}

	// The active only GET parameter doesn't always work
	if _, ok := getParameters["activeWithinSeconds"]; ok {
		sessions = slices.DeleteFunc(sessions, func(session api.Session) bool {
			return time.Since(session.LastActivityDate()) > 10*time.Minute
		})
	}

	switch output {
	case "text":
		if len(sessions) == 0 {
			fmt.Println("No sessions")
		} else {
			fmt.Println("Sessions:")

			for _, session := range sessions {
				duration := humanize.RelTime(time.Now(), session.LastActivityDate(), "", "ago")
				fmt.Printf(" - %s (%s) %s\n", session.UserName(), session.DeviceName(), duration)
			}
		}

	case "json":
		if jsonBytes, err := json.Marshal(sessions); err != nil {
			return fmt.Errorf("failed to encode sessions: %w", err)
		} else {
			fmt.Println(string(jsonBytes))
		}
	}

	return nil
}

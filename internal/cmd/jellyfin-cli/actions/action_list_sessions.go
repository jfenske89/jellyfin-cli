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

	"codeberg.org/jfenske/jellyfin-cli/api"
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
	var output = "text"
	if val, ok := options["output"]; ok && val != "" {
		output = val
	}

	var activeOnly bool
	if val, ok := options["active-only"]; ok && val != "" {
		activeOnly, _ = strconv.ParseBool(val)
	}

	if sessions, err := e.client.ListSessions(ctx); err != nil {
		return err
	} else {
		if activeOnly {
			sessions = slices.DeleteFunc[[]api.Session, api.Session](sessions, func(session api.Session) bool {
				return time.Since(session.GetLastActivityDate()) > 10*time.Minute
			})
		}

		switch output {
		case "text":
			if len(sessions) == 0 {
				fmt.Println("No sessions")
			} else {
				fmt.Println("Sessions:")

				for _, session := range sessions {
					duration := humanize.RelTime(time.Now(), session.GetLastActivityDate(), "", "ago")
					fmt.Printf(" - %s (%s) %s\n", session.GetUserName(), session.GetDeviceName(), duration)
				}
			}

		case "json":
			if jsonBytes, err := json.Marshal(sessions); err != nil {
				return fmt.Errorf("failed to encode sessions: %w", err)
			} else {
				fmt.Println(string(jsonBytes))
			}
		}
	}

	return nil
}

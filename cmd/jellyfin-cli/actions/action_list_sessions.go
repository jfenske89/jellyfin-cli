package actions

import (
	"context"
	"flag"
	"fmt"
	"slices"
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

func (e *listSessionsExecutorImpl) Run(ctx context.Context) error {
	options := parseFlags(map[string]any{
		"active": flag.Bool("active", false, "only list active sessions"),
	})

	output := "text"
	if val, ok := options["output"].(string); ok && val != "" {
		output = val
	}

	getParameters := make(map[string]string)
	if val, ok := options["active"].(bool); ok && val {
		getParameters["activeWithinSeconds"] = "600"
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

	return writeResponse(sessions, output, e.formatText)
}

func (e *listSessionsExecutorImpl) formatText(sessions []api.Session) {
	if len(sessions) == 0 {
		fmt.Println("No sessions")
		return
	}

	fmt.Println("Sessions:")

	for _, session := range sessions {
		duration := humanize.RelTime(time.Now(), session.LastActivityDate(), "", "ago")
		fmt.Printf(" - %s on %s (%s)\n", session.UserName(), session.DeviceName(), duration)
	}
}

package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"go.uber.org/zap"

	"codeberg.org/jfenske/jellyfin-cli/internal/libraries/api"
)

type listActivityExecutorImpl struct {
	client api.JellyfinApiClient
	logger *zap.SugaredLogger
}

func NewListActivityExecutor(client api.JellyfinApiClient, logger *zap.SugaredLogger) Executor {
	return &listActivityExecutorImpl{
		client: client,
		logger: logger,
	}
}

func (e *listActivityExecutorImpl) Run(ctx context.Context, options map[string]string) error {
	output := "text"

	if val, ok := options["output"]; ok && val != "" {
		output = val
	}

	getParameters := make(map[string]string)

	activityLog, err := e.client.ListActivityLogs(ctx, getParameters)
	if err != nil {
		return err
	}

	logs := activityLog.Items()

	switch output {
	case "text":
		if len(logs) == 0 {
			fmt.Println("No activity logs")
		} else {
			fmt.Println("Activity log:")

			for _, log := range logs {
				duration := humanize.RelTime(time.Now(), log.Date(), "", "ago")
				fmt.Printf(" - %s (%s)\n", log.Name(), duration)
			}
		}

	case "json":
		if jsonBytes, err := json.Marshal(activityLog); err != nil {
			return fmt.Errorf("failed to encode activity logs: %w", err)
		} else {
			fmt.Println(string(jsonBytes))
		}
	}

	return nil
}

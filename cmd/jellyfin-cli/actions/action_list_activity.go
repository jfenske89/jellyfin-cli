package actions

import (
	"context"
	"flag"
	"fmt"
	"strconv"
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

func (e *listActivityExecutorImpl) Run(ctx context.Context) error {
	options := parseFlags(
		map[string]any{
			"limit": flag.Int64("limit", 0, "Limit the number of activity items"),
		},
	)

	output := "text"
	if val, ok := options["output"].(string); ok && val != "" {
		output = val
	}

	getParameters := make(map[string]string)
	if val, ok := options["limit"].(int64); ok && val > 0 {
		getParameters["limit"] = strconv.FormatInt(val, 10)
	}

	activityLog, err := e.client.ListActivityLogs(ctx, getParameters)
	if err != nil {
		return err
	}

	return writeResponse(activityLog.Items(), output, e.formatText)
}

func (e *listActivityExecutorImpl) formatText(logs []api.ActivityLogItem) {
	if len(logs) == 0 {
		fmt.Println("No activity logs")
		return
	}

	fmt.Println("Activity log:")

	for _, log := range logs {
		duration := humanize.RelTime(time.Now(), log.Date(), "", "ago")
		fmt.Printf(" - %s (%s)\n", log.Name(), duration)
	}
}

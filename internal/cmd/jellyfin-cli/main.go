package main

import (
	"context"
	"os"

	"codeberg.org/jfenske/jellyfin-cli/api"
	"codeberg.org/jfenske/jellyfin-cli/internal/cmd/jellyfin-cli/actions"
)

func main() {
	logger, atom := buildLogger()
	defer func() {
		// ignore errors for stderr fsync
		_ = logger.Sync()
	}()

	config, err := loadConfiguration(atom, logger)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	client := api.NewJellyfinApiClient(config, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: parse CLI arguments to determine action
	var executor actions.Executor
	var action = actions.ListSessions

	switch action {
	case actions.ListSessions:
		executor = actions.NewListSessionsExecutor(client, logger)

	default:
		usage()
		os.Exit(0)
	}

	if executor != nil {
		// TODO: parse action options (for example: --active-only or --output=json)
		options := make(map[string]string)
		if err := executor.Run(ctx, options); err != nil {
			logger.Fatalw(
				"failed to execute action",
				"action", action,
				"error", err.Error(),
			)
		}
	}
}

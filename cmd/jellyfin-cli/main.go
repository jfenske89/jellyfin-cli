package main

import (
	"context"
	"os"
	"slices"

	"codeberg.org/jfenske/jellyfin-cli/cmd/jellyfin-cli/actions"
	"codeberg.org/jfenske/jellyfin-cli/internal/libraries/api"
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

	// Get the command from arguments first
	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}

	// Extract the command from the first argument
	command := os.Args[1]
	os.Args = slices.Delete(os.Args, 1, 2)

	// Get an executor for the command
	var executor actions.Executor
	var action string

	switch command {
	case actions.ListSessions:
		executor = actions.NewListSessionsExecutor(client, logger)

	case actions.ListActivity:
		executor = actions.NewListActivityExecutor(client, logger)

	case actions.ListLibraryFolders:
		executor = actions.NewListLibraryFoldersExecutor(client, logger)

	case actions.Search:
		executor = actions.NewSearch(client, logger)

	case actions.RefreshLibrary:
		executor = actions.NewRefreshLibraryExecutor(client, logger)

	default:
		usage()
		os.Exit(0)
	}

	err = executor.Run(ctx)
	if err != nil {
		logger.Fatalw(
			"failed to execute action",
			"action", action,
			"error", err.Error(),
		)
	}
}

package main

import (
	"context"
	"os"
	"slices"

	"codeberg.org/jfenske/jellyfin-cli/api"
	"codeberg.org/jfenske/jellyfin-cli/cmd/jellyfin-cli/actions"
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

	var executor actions.Executor
	var action string

	// TODO: add improved parsing logic for actions
	if slices.Contains(os.Args, "list-sessions") {
		action = actions.ListSessions
	} else if slices.Contains(os.Args, "list-library-folders") {
		action = actions.ListLibraryFolders
	}

	switch action {
	case actions.ListSessions:
		executor = actions.NewListSessionsExecutor(client, logger)

	case actions.ListLibraryFolders:
		executor = actions.NewListLibraryFoldersExecutor(client, logger)

	default:
		usage()
		os.Exit(0)
	}

	if executor != nil {
		// TODO: add improved parsing logic for action options, which would be specific per action
		options := make(map[string]string)
		if slices.Contains(os.Args, "--active") {
			options["active-only"] = "1"
		}

		if slices.Contains(os.Args, "--output=json") || slices.Contains(os.Args, "--json") {
			options["output"] = "json"
		}

		if err := executor.Run(ctx, options); err != nil {
			logger.Fatalw(
				"failed to execute action",
				"action", action,
				"error", err.Error(),
			)
		}
	}
}

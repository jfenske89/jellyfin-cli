package main

import (
	"context"
	"log"

	"codeberg.org/jfenske/jellyfin-cli/api"
	"codeberg.org/jfenske/jellyfin-cli/internal/cmd/jellyfin-cli/actions"
)

func main() {
	log.SetFlags(0)

	var config api.JellyfinApiConfig
	if result, err := loadConfiguration(); err != nil {
		log.Fatalf(err.Error())
	} else {
		config = *result
	}

	client := api.NewJellyfinApiClient(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: parse CLI arguments to determine action
	var executor actions.Executor
	var action = actions.ListSessions

	switch action {
	case actions.ListSessions:
		executor = actions.NewListSessionsExecutor(client)

	default:
		// TODO: show help and exit
	}

	if executor != nil {
		// TODO: parse action arguments
		arguments := make(map[string]interface{})
		if err := executor.Run(ctx, arguments); err != nil {
			log.Fatalf("failed to %s: %s", action, err.Error())
		}
	}
}

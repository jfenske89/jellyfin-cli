package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"codeberg.org/jfenske/jellyfin-cli/api"
	"codeberg.org/jfenske/jellyfin-cli/internal/cmd/jellyfin-cli/actions"
)

func main() {
	log.SetFlags(0)

	if err := godotenv.Load(".env"); err != nil {
		// maybe the variables were loaded another way
		log.Printf("Warning: failed to load .env: %s", err.Error())
	}

	if os.Getenv("HOST") == "" {
		// use a default host
		os.Setenv("HOST", "127.0.0.1:8096")
	}

	if os.Getenv("TOKEN") == "" {
		// the http request will probably fail
		log.Printf("Warning: no TOKEN provided!")
	}

	// TODO: support reading from a configuration file too
	config := api.JellyfinApiConfig{
		BaseUrl: os.Getenv("BASE_URL"),
		Token:   os.Getenv("TOKEN"),
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

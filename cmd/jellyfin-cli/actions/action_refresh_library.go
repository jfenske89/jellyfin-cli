package actions

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"codeberg.org/jfenske/jellyfin-cli/internal/libraries/api"
)

type refreshLibraryExecutorImpl struct {
	client api.JellyfinApiClient
	logger *zap.SugaredLogger
}

func NewRefreshLibraryExecutor(client api.JellyfinApiClient, logger *zap.SugaredLogger) Executor {
	return &refreshLibraryExecutorImpl{
		client: client,
		logger: logger,
	}
}

func (e *refreshLibraryExecutorImpl) Run(ctx context.Context) error {
	err := e.client.RefreshLibrary(ctx)
	if err != nil {
		return err
	}

	fmt.Println("OK")

	return nil
}

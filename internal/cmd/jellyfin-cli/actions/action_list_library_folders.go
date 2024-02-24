package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"go.uber.org/zap"

	"codeberg.org/jfenske/jellyfin-cli/api"
)

type listLibraryFoldersExecutorImpl struct {
	client api.JellyfinApiClient
	logger *zap.SugaredLogger
}

func NewListLibraryFoldersExecutor(client api.JellyfinApiClient, logger *zap.SugaredLogger) Executor {
	return &listLibraryFoldersExecutorImpl{
		client: client,
		logger: logger,
	}
}

func (e *listLibraryFoldersExecutorImpl) Run(ctx context.Context, options map[string]string) error {
	output := "text"

	if val, ok := options["output"]; ok && val != "" {
		output = val
	}

	getParameters := make(map[string]string)

	if folders, err := e.client.ListLibraryFolders(ctx, getParameters); err != nil {
		return err
	} else {
		switch output {
		case "text":
			if len(folders) == 0 {
				fmt.Println("No library folders")
			} else {
				// sort by collection type and name, but prefer libraries for movies and tvshows
				weighted := map[string]int{
					"movies":  1,
					"tvshows": 2,
				}

				slices.SortFunc(folders, func(a, b api.LibraryFolder) int {
					aCollectionType := a.CollectionType()
					bCollectionType := b.CollectionType()

					aWeight := weighted[aCollectionType]

					if aWeight == 0 {
						aWeight = 9
					}

					bWeight := weighted[bCollectionType]

					if bWeight == 0 {
						bWeight = 9
					}

					aName := a.Name()
					bName := b.Name()

					aSort := fmt.Sprintf("%d#%s#%s", aWeight, aCollectionType, aName)
					bSort := fmt.Sprintf("%d#%s#%s", bWeight, bCollectionType, bName)

					return strings.Compare(aSort, bSort)
				})

				fmt.Println("Library folders:")

				var lastCollectionType string

				for _, folder := range folders {
					if lastCollectionType != folder.CollectionType() {
						if lastCollectionType != "" {
							fmt.Printf("\n")
						}

						lastCollectionType = folder.CollectionType()
						fmt.Printf("- %s:\n", lastCollectionType)
					}

					fmt.Printf("   - %s (%s)\n", folder.Name(), folder.ItemId())
				}
			}

		case "json":
			if jsonBytes, err := json.Marshal(folders); err != nil {
				return fmt.Errorf("failed to encode library folders: %w", err)
			} else {
				fmt.Println(string(jsonBytes))
			}
		}
	}

	return nil
}

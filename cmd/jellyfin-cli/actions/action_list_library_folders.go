package actions

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"go.uber.org/zap"

	"codeberg.org/jfenske/jellyfin-cli/internal/libraries/api"
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

func (e *listLibraryFoldersExecutorImpl) Run(ctx context.Context) error {
	options := parseFlags(make(map[string]any))

	output := "text"
	if val, ok := options["output"].(string); ok && val != "" {
		output = val
	}

	getParameters := make(map[string]string)

	folders, err := e.client.ListLibraryFolders(ctx, getParameters)
	if err != nil {
		return err
	}

	return writeResponse(folders, output, e.formatText)
}

func (e *listLibraryFoldersExecutorImpl) formatText(folders []api.LibraryFolder) {
	if len(folders) == 0 {
		fmt.Println("No library folders")
		return
	}

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

		fmt.Printf("   - %s <%s>\n", folder.Name(), folder.ItemId())
	}
}

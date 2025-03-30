package actions

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"codeberg.org/jfenske/jellyfin-cli/internal/libraries/api"
)

type searchExecutorImpl struct {
	client api.JellyfinApiClient
	logger *zap.SugaredLogger
}

func NewSearch(client api.JellyfinApiClient, logger *zap.SugaredLogger) Executor {
	return &searchExecutorImpl{
		client: client,
		logger: logger,
	}
}

func (e *searchExecutorImpl) Run(ctx context.Context) error {
	options := parseFlags(
		map[string]any{
			"term":  flag.String("term", "", "The required search term"),
			"limit": flag.Int64("limit", 0, "Limit the number of activity items"),
		},
	)

	output := "text"
	if val, ok := options["output"].(string); ok && val != "" {
		output = val
	}

	var term string
	getParameters := make(map[string]string)
	if val, ok := options["term"].(string); ok && val != "" {
		term = val
	} else {
		return fmt.Errorf("term is required for search")
	}

	if val, ok := options["limit"].(int64); ok && val > 0 {
		getParameters["limit"] = strconv.FormatInt(val, 10)
	}

	response, err := e.client.Search(ctx, term, getParameters)
	if err != nil {
		return err
	}

	return writeResponse(response, output, e.formatText)
}

func (e *searchExecutorImpl) formatText(response api.SearchResponse) {
	results := response.SearchHints()
	if len(response.SearchHints()) == 0 {
		fmt.Println("No search results")
		return
	}

	fmt.Println("Search results:")
	for _, result := range results {
		fmt.Printf(" - %s (%s) <%s>\n", result.Name(), result.Type(), result.GetId())
	}
}

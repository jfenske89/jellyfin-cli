package api_test

import (
	"testing"

	"codeberg.org/jfenske/jellyfin-cli/internal/libraries/api"
)

func TestSearchResponse(t *testing.T) {
	response := api.NewSearchResponse(
		[]byte(`{
			"TotalRecordCount": 2,
			"SearchHints": [
				{"Id": "hint1", "Name": "first-hint", "Type": "movie"},
				{"Id": "hint2", "Name": "second-hint", "Type": "series"}
			]
		}`),
	)

	if got := response.TotalRecordCount(); got != 2 {
		t.Errorf("TotalRecordCount() = %d; want 2", got)
	}

	hints := response.SearchHints()
	if got := len(hints); got != 2 {
		t.Errorf("len(SearchHints()) = %d; want 2", got)
	}

	firstHint := hints[0]
	if got := firstHint.GetId(); got != "hint1" {
		t.Errorf("SearchHints()[0].GetId() = %s; want hint1", got)
	}
	if got := firstHint.Name(); got != "first-hint" {
		t.Errorf("SearchHints()[0].Name() = %s; want first-hint", got)
	}
	if got := firstHint.Type(); got != "movie" {
		t.Errorf("SearchHints()[0].Type() = %s; want movie", got)
	}

	secondHint := hints[1]
	if got := secondHint.GetId(); got != "hint2" {
		t.Errorf("SearchHints()[1].GetId() = %s; want hint2", got)
	}
	if got := secondHint.Name(); got != "second-hint" {
		t.Errorf("SearchHints()[1].Name() = %s; want second-hint", got)
	}
	if got := secondHint.Type(); got != "series" {
		t.Errorf("SearchHints()[1].Type() = %s; want series", got)
	}
}

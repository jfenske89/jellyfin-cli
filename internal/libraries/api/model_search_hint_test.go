package api_test

import (
	"testing"

	"codeberg.org/jfenske/jellyfin-cli/internal/libraries/api"
)

func TestSearchHint(t *testing.T) {
	hint := api.NewSearchHint(
		[]byte(`{"Id": "test-id", "Name": "test-name", "Type": "test-type"}`),
	)

	if got := hint.GetId(); got != "test-id" {
		t.Errorf("GetId() = %s; want test-id", got)
	}

	if got := hint.Name(); got != "test-name" {
		t.Errorf("Name() = %s; want test-name", got)
	}

	if got := hint.Type(); got != "test-type" {
		t.Errorf("Type() = %s; want test-type", got)
	}
}

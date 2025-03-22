package api_test

import (
	"testing"
	"time"

	"codeberg.org/jfenske/jellyfin-cli/internal/libraries/api"
)

func TestActivityLogItem(t *testing.T) {
	date := time.Now().Add(30 * time.Second)

	session := api.NewActivityLogItem(
		[]byte(`{"Id": 1, "Name": "test-name", "Overview": "test-overview", "ShortOverview": "test-short-overview", "Type": "test-type", "ItemId": "test-item-id", "Date": "` + date.Format("2006-01-02T15:04:05.9999999Z") + `", "UserId": "test-user-id", "Severity": "test-severity"}`),
	)

	if got := session.Id(); got != 1 {
		t.Errorf("Id() = %d; want 1", got)
	}

	if got := session.Name(); got != "test-name" {
		t.Errorf("Name() = %s; want test-name", got)
	}

	if got := session.Overview(); got != "test-overview" {
		t.Errorf("Overview() = %s; want test-overview", got)
	}

	if got := session.ShortOverview(); got != "test-short-overview" {
		t.Errorf("ShortOverview() = %s; want test-short-overview", got)
	}

	if got := session.Type(); got != "test-type" {
		t.Errorf("Type() = %s; want test-type", got)
	}

	if got := session.ItemId(); got != "test-item-id" {
		t.Errorf("ItemId() = %s; want test-item-id", got)
	}

	if got := session.Date(); got.Format(time.DateTime) != date.Format(time.DateTime) {
		t.Errorf("Date() = %v; want %v", got, date)
	}

	if got := session.UserId(); got != "test-user-id" {
		t.Errorf("UserId() = %s; want test-user-id", got)
	}

	if got := session.Severity(); got != "test-severity" {
		t.Errorf("Severity() = %s; want test-severity", got)
	}

	session = api.NewActivityLogItem([]byte(`{"Date": "2024-02-18T16:31:11.9906841Z"}`))

	if got := session.Date().Format(time.DateTime); got != "2024-02-18 16:31:11" {
		t.Errorf("Date() = %v; want 2024-02-18 16:31:11", got)
	}
}

package api_test

import (
	"testing"
	"time"

	"codeberg.org/jfenske/jellyfin-cli/internal/libraries/api"
)

func TestActivityLog(t *testing.T) {
	date := time.Now().Add(30 * time.Second)

	log := api.NewActivityLog(
		[]byte(`{"TotalRecordCount": 1, "StartIndex": 2, "Items": [{"Id": 1, "Name": "test-name", "Overview": "test-overview", "ShortOverview": "test-short-overview", "Type": "test-type", "ItemId": "test-item-id", "Date": "` + date.Format("2006-01-02T15:04:05.9999999Z") + `", "UserId": "test-user-id", "Severity": "test-severity"}]}`),
	)

	if got := log.TotalRecordCount(); got != 1 {
		t.Errorf("TotalRecordCount() = %d; want 1", got)
	}

	if got := log.StartIndex(); got != 2 {
		t.Errorf("StartIndex() = %d; want 2", got)
	}

	if got := len(log.Items()); got != 1 {
		t.Errorf("len(Items()) = %d; want 1", got)
	}

	item := log.Items()[0]

	if got := item.Id(); got != 1 {
		t.Errorf("Items()[0]Id() = %d; want 1", got)
	}

	if got := item.Name(); got != "test-name" {
		t.Errorf("Items()[0]Name() = %s; want test-name", got)
	}

	if got := item.Overview(); got != "test-overview" {
		t.Errorf("Items()[0]Overview() = %s; want test-overview", got)
	}

	if got := item.ShortOverview(); got != "test-short-overview" {
		t.Errorf("Items()[0]ShortOverview() = %s; want test-short-overview", got)
	}

	if got := item.Type(); got != "test-type" {
		t.Errorf("Items()[0]Type() = %s; want test-type", got)
	}

	if got := item.ItemId(); got != "test-item-id" {
		t.Errorf("Items()[0]ItemId() = %s; want test-item-id", got)
	}

	if got := item.Date(); got.Format(time.DateTime) != date.Format(time.DateTime) {
		t.Errorf("Items()[0]Date() = %v; want %v", got, date)
	}

	if got := item.UserId(); got != "test-user-id" {
		t.Errorf("Items()[0]UserId() = %s; want test-user-id", got)
	}

	if got := item.Severity(); got != "test-severity" {
		t.Errorf("Items()[0]Severity() = %s; want test-severity", got)
	}
}

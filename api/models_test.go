package api_test

import (
	"testing"
	"time"

	"codeberg.org/jfenske/jellyfin-cli/api"
)

func TestSession(t *testing.T) {
	lastActivityDate := time.Now().Add(30 * time.Second)

	session := api.Session{
		"UserName":         "tester",
		"DeviceName":       "test",
		"LastActivityDate": lastActivityDate,
	}

	if got := session.GetUserName(); got != "tester" {
		t.Errorf("GetUserName() = %s; want tester", got)
	}

	if got := session.GetDeviceName(); got != "test" {
		t.Errorf("GetDeviceName() = %s; want test", got)
	}

	if got := session.GetLastActivityDate(); got != lastActivityDate {
		t.Errorf("GetLastActivityDate() = %v; want %v", got, lastActivityDate)
	}

	session["LastActivityDate"] = "2024-02-18T16:31:11.9906841Z"

	if got := session.GetLastActivityDate().Format(time.DateTime); got != "2024-02-18 16:31:11" {
		t.Errorf("GetLastActivityDate() = %v; want 2024-02-18 16:31:11", got)
	}
}

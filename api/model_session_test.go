package api_test

import (
	"testing"
	"time"

	"codeberg.org/jfenske/jellyfin-cli/api"
)

func TestSession(t *testing.T) {
	lastActivityDate := time.Now().Add(30 * time.Second)

	session := api.NewSession(
		[]byte(`{"UserName": "tester", "DeviceName": "test", "LastActivityDate": "` + lastActivityDate.Format("2006-01-02T15:04:05.9999999Z") + `"}`),
	)

	if got := session.UserName(); got != "tester" {
		t.Errorf("UserName() = %s; want tester", got)
	}

	if got := session.DeviceName(); got != "test" {
		t.Errorf("DeviceName() = %s; want test", got)
	}

	if got := session.LastActivityDate(); got.Format(time.DateTime) != lastActivityDate.Format(time.DateTime) {
		t.Errorf("LastActivityDate() = %v; want %v", got, lastActivityDate)
	}

	session = api.NewSession([]byte(`{"LastActivityDate": "2024-02-18T16:31:11.9906841Z"}`))

	if got := session.LastActivityDate().Format(time.DateTime); got != "2024-02-18 16:31:11" {
		t.Errorf("LastActivityDate() = %v; want 2024-02-18 16:31:11", got)
	}
}

package api

import (
	"time"

	"github.com/tidwall/gjson"
)

// Session defines a Jellyfin user session
type Session interface {
	Model

	UserName() string

	DeviceName() string

	LastActivityDate() time.Time
}

type sessionImpl struct {
	*genericModel
}

func NewSession(jsonBytes []byte) Session {
	return &sessionImpl{
		genericModel: newGenericModel(jsonBytes),
	}
}

func (s *sessionImpl) UserName() string {
	return s.JsonPath("UserName").String()
}

func (s *sessionImpl) DeviceName() string {
	return s.JsonPath("DeviceName").String()
}

func (s *sessionImpl) LastActivityDate() time.Time {
	result := s.JsonPath("LastActivityDate")

	if result.Exists() {
		switch result.Type {
		case gjson.String:
			// 2024-02-18T16:31:11.9906841Z
			if timeResult, err := time.Parse("2006-01-02T15:04:05.9999999Z", result.String()); err == nil {
				return timeResult
			}

		case gjson.Number:
			return time.Unix(result.Int(), 0)
		}
	}

	return time.Time{}
}

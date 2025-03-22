package api

import (
	"time"
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
	return s.JsonPathAsDate("LastActivityDate")
}

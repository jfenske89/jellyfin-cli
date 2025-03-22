package api

import (
	"time"
)

type ActivityLogItem interface {
	Id() int64

	Name() string

	Overview() string

	ShortOverview() string

	Type() string

	ItemId() string

	Date() time.Time

	UserId() string

	Severity() string
}

type activityLogItemImpl struct {
	*genericModel
}

func NewActivityLogItem(jsonBytes []byte) ActivityLogItem {
	return &activityLogItemImpl{
		genericModel: newGenericModel(jsonBytes),
	}
}

func (s *activityLogItemImpl) Id() int64 {
	return s.JsonPath("Id").Int()
}

func (s *activityLogItemImpl) Name() string {
	return s.JsonPath("Name").String()
}

func (s *activityLogItemImpl) Overview() string {
	return s.JsonPath("Overview").String()
}

func (s *activityLogItemImpl) ShortOverview() string {
	return s.JsonPath("ShortOverview").String()
}

func (s *activityLogItemImpl) Type() string {
	return s.JsonPath("Type").String()
}

func (s *activityLogItemImpl) ItemId() string {
	return s.JsonPath("ItemId").String()
}

func (s *activityLogItemImpl) Date() time.Time {
	return s.JsonPathAsDate("Date")
}

func (s *activityLogItemImpl) UserId() string {
	return s.JsonPath("UserId").String()
}

func (s *activityLogItemImpl) Severity() string {
	return s.JsonPath("Severity").String()
}

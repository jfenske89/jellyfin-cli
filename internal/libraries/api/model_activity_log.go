package api

// ActivityLog defines a Jellyfin user activityLog
type ActivityLog interface {
	Model

	Items() []ActivityLogItem

	TotalRecordCount() int64

	StartIndex() int64
}

type activityLogImpl struct {
	*genericModel
}

func NewActivityLog(jsonBytes []byte) ActivityLog {
	return &activityLogImpl{
		genericModel: newGenericModel(jsonBytes),
	}
}

func (s *activityLogImpl) Items() []ActivityLogItem {
	array := s.JsonPath("Items").Array()
	items := make([]ActivityLogItem, len(array))

	for i := range array {
		items[i] = NewActivityLogItem([]byte(array[i].Raw))
	}

	return items
}

func (s *activityLogImpl) TotalRecordCount() int64 {
	return s.JsonPath("TotalRecordCount").Int()
}

func (s *activityLogImpl) StartIndex() int64 {
	return s.JsonPath("StartIndex").Int()
}

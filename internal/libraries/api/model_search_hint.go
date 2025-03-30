package api

// SearchHint defines a Jellyfin search hint
type SearchHint interface {
	Model

	GetId() string

	Name() string

	Type() string
}

type searchHintImpl struct {
	*genericModel
}

func NewSearchHint(jsonBytes []byte) SearchHint {
	return &searchHintImpl{
		genericModel: newGenericModel(jsonBytes),
	}
}

func (s *searchHintImpl) GetId() string {
	return s.JsonPath("Id").String()
}

func (s *searchHintImpl) Name() string {
	return s.JsonPath("Name").String()
}

func (s *searchHintImpl) Type() string {
	return s.JsonPath("Type").String()
}

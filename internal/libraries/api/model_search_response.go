package api

// SearchResponse defines a Jellyfin search response
type SearchResponse interface {
	Model

	SearchHints() []SearchHint
	TotalRecordCount() int64
}

type searchResponseImpl struct {
	*genericModel
}

func NewSearchResponse(jsonBytes []byte) SearchResponse {
	return &searchResponseImpl{
		genericModel: newGenericModel(jsonBytes),
	}
}

func (s *searchResponseImpl) SearchHints() []SearchHint {
	array := s.JsonPath("SearchHints").Array()
	items := make([]SearchHint, len(array))

	for i := range array {
		items[i] = NewSearchHint([]byte(array[i].Raw))
	}

	return items
}

func (s *searchResponseImpl) TotalRecordCount() int64 {
	return s.JsonPath("TotalRecordCount").Int()
}

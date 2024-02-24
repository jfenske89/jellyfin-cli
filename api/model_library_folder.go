package api

// LibraryFolder defines a Jellyfin virtual library folder
type LibraryFolder interface {
	Model

	ItemId() string

	Name() string

	CollectionType() string
}

type libraryFolderImpl struct {
	*genericModel
}

func NewLibraryFolder(jsonBytes []byte) LibraryFolder {
	return &libraryFolderImpl{
		genericModel: newGenericModel(jsonBytes),
	}
}

func (s *libraryFolderImpl) ItemId() string {
	return s.JsonPath("ItemId").String()
}

func (s *libraryFolderImpl) Name() string {
	return s.JsonPath("Name").String()
}

func (s *libraryFolderImpl) CollectionType() string {
	return s.JsonPath("CollectionType").String()
}

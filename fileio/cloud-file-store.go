package fileio

// TODO
type CloudFileStore struct{}

func (c CloudFileStore) Save(path string, content string) error {
	panic("implement me")
}

func (c CloudFileStore) Load(path string) (content string, err error) {
	panic("implement me")
}

func (c CloudFileStore) List(path string) (subPaths []string, err error) {
	panic("implement me")
}

func (c CloudFileStore) FullName(path string) (fullPath string, err error) {
	panic("implement me")
}

package fileio

import (
	"fmt"
	"sync"
	"time"
)

var (
	storage FileStore
	once    sync.Once
)

func CurrStorage() FileStore {
	once.Do(func() {
		storage = SimpleFileStore{BasePath: ""}
	})
	return storage
}

func SetStorage(fs FileStore) {
	once.Do(func() {}) // Disable default setup if it has not yet happened
	storage = fs
}

type FileInfo struct {
	Name    string
	Path    string
	ModTime time.Time
}

type FileData struct {
	Filename string
	Content  string
}

type FileStore interface {
	// Save a file at the given path (filename included) with some content
	Save(path string, content string) error
	// Load contents at the given path (filename included)
	Load(path string) (content string, err error)
	// Move the file at path to the target directory
	Move(path string, targetDir string) error
	// List the sub paths at the given path, include directories and normal files.
	// The path of each returned FileInfo object will have the same relativity as the given path
	List(path string) (subPaths []FileInfo, err error)
	// Get more info about the given path
	Info(path string) (info FileInfo, err error)
	// Get the absolute full pathname of the given path
	FullName(path string) (fullPath string, err error)
	// Split the path into directory and filename
	Split(path string) (directory string, filename string)
	// The default upload path for this file storage, if filename is blank, return directory only
	UploadPath(userCode string, filename string) string
	// The default download path for this file storage, return directory only if filename is blank
	DownloadPath(userCode string, filename string) string
}

type FileExistsError struct {
	FileName string
}

func (e FileExistsError) Error() string {
	return fmt.Sprintf("%v already exists!", e.FileName)
}

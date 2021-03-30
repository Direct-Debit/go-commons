package fileio

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type SimpleFileStore struct {
	BasePath string
}

func (s SimpleFileStore) fullPath(path string) string {
	return filepath.Join(s.BasePath, path)
}

func (s SimpleFileStore) Save(path string, content string) error {
	path = s.fullPath(path)

	if _, err := os.Stat(path); os.IsExist(err) {
		return FileExistsError{FileName: path}
	}

	dir, _ := s.Split(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
		panic(fmt.Sprintf("Could not close file %s: %v", file.Name(), err))
	}()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func (s SimpleFileStore) Load(path string) (content string, err error) {
	path = s.fullPath(path)

	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return
	}

	var contents strings.Builder
	_, err = io.Copy(&contents, file)
	content = contents.String()
	return
}

func (s SimpleFileStore) Move(path string, targetDir string) error {
	fTarget := s.fullPath(targetDir)

	tInfo, err := os.Stat(fTarget)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil && !tInfo.IsDir() {
		return errors.New(targetDir + " is not a directory")
	}
	log.Trace("Checked target for move")

	content, err := s.Load(path)
	if err != nil {
		return err
	}
	_, name := s.Split(path)
	newPath := filepath.Join(targetDir, name)
	if err := s.Save(newPath, content); err != nil {
		return err
	}
	return os.Remove(path)
}

func (s SimpleFileStore) List(path string) (subPaths []FileInfo, err error) {
	fPath := s.fullPath(path)

	var files []os.FileInfo
	files, err = ioutil.ReadDir(fPath)
	subPaths = make([]FileInfo, len(files))
	for i, val := range files {
		sPath := filepath.Join(path, val.Name())
		subPaths[i] = FileInfo{
			Name:    val.Name(),
			Path:    sPath,
			ModTime: val.ModTime(),
		}
	}
	return
}

func (s SimpleFileStore) Info(path string) (info FileInfo, err error) {
	fPath := s.fullPath(path)
	inf, err := os.Stat(fPath)
	if err != nil {
		return FileInfo{}, err
	}
	return FileInfo{
		Name:    inf.Name(),
		Path:    path,
		ModTime: inf.ModTime(),
	}, nil
}

func (s SimpleFileStore) FullName(path string) (fullPath string, err error) {
	return filepath.Abs(s.fullPath(path))
}

func (s SimpleFileStore) Split(path string) (string, string) {
	return filepath.Split(path)
}

func (s SimpleFileStore) UploadPath(userCode string, filename string) string {
	return filepath.Join("uploads", userCode, filename)
}

func (s SimpleFileStore) DownloadPath(userCode string, filename string) string {
	return filepath.Join("downloads", userCode, filename)
}

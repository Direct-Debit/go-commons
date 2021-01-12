package fileio

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
)

type S3Store struct {
	s3     *s3.S3
	Bucket string
}

func NewS3Store(bucket string) S3Store {
	sess := session.Must(session.NewSession())
	return S3Store{s3: s3.New(sess), Bucket: bucket}
}

func (s S3Store) Save(path string, content string) error {
	panic("implement me")
}

func (s S3Store) Load(path string) (content string, err error) {
	output, err := s.s3.GetObject(&s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &path,
	})
	if err != nil {
		return "", err
	}

	var fileContent []byte
	buffer := bytes.NewBuffer(fileContent)
	if _, err = io.Copy(buffer, output.Body); err != nil {
		return "", err
	}

	content = buffer.String()
	return content, nil
}

func (s S3Store) Move(path string, targetDir string) error {
	panic("implement me")
}

func (s S3Store) List(path string) (subPaths []FileInfo, err error) {
	panic("implement me")
}

func (s S3Store) Info(path string) (info FileInfo, err error) {
	panic("implement me")
}

func (s S3Store) FullName(path string) (fullPath string, err error) {
	panic("implement me")
}

func (s S3Store) Split(path string) (directory string, filename string) {
	panic("implement me")
}

func (s S3Store) UploadPath(userCode string, filename string) string {
	panic("implement me")
}

func (s S3Store) DownloadPath(userCode string, filename string) string {
	panic("implement me")
}

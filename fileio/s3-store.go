package fileio

import (
	"bytes"
	"fmt"
	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
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
	_, err := s.s3.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(content),
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(path),
	})
	if errlib.ErrorError(err, "Couldn't save object to "+path) {
		return err
	}
	return nil
}

func (s S3Store) Load(path string) (content string, err error) {
	log.Trace(fmt.Sprintf("Downloading s3://%s%s", s.Bucket, path))
	output, err := s.s3.GetObject(&s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &path,
	})
	if err != nil {
		return "", err
	}

	var fileContent []byte
	buffer := bytes.NewBuffer(fileContent)
	n, err := io.Copy(buffer, output.Body)
	if errlib.ErrorError(err, "Couldn't copy bytes downloaded from s3") {
		return "", err
	}

	content = buffer.String()
	log.Trace(fmt.Sprintf("Downloaded %d bytes from s3", n))
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

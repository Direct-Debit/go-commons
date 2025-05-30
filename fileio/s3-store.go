package fileio

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

type S3Store struct {
	s3     *s3.S3
	Bucket *string
}

func NewS3Store(bucket string) S3Store {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return S3Store{s3: s3.New(sess), Bucket: &bucket}
}

func (s S3Store) Save(path string, content string) error {
	_, err := s.s3.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(content),
		Bucket: s.Bucket,
		Key:    aws.String(path),
	})
	if errlib.ErrorError(err, "Couldn't save object to "+path) {
		return err
	}
	return nil
}

func (s S3Store) Load(path string) (content string, err error) {
	log.Trace(fmt.Sprintf("Downloading s3://%s/%s", *s.Bucket, path))
	output, err := s.s3.GetObject(&s3.GetObjectInput{
		Bucket: s.Bucket,
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
	dir, name := s.Split(path)
	if targetDir[len(targetDir)-1] != '/' {
		targetDir += "/"
	}

	if dir == targetDir {
		log.Infof("Skipping move because source and target dir are the same (%s)", dir)
		return nil
	}

	_, err := s.s3.CopyObject(&s3.CopyObjectInput{
		Bucket:     s.Bucket,
		CopySource: aws.String(*s.Bucket + "/" + path),
		Key:        aws.String(targetDir + name),
	})
	if errlib.ErrorError(err, "Couldn't copy s3 file") {
		return err
	}

	return s.Delete(path)
}

func (s S3Store) Delete(path string) error {
	_, err := s.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.Bucket,
		Key:    &path,
	})
	errlib.ErrorError(err, "Couldn't delete s3 file")
	return err
}

func (s S3Store) List(path string) (subPaths []FileInfo, err error) {
	params := &s3.ListObjectsV2Input{
		Bucket: s.Bucket,
		Prefix: &path,
	}

	var content []*s3.Object
	err = s.s3.ListObjectsV2Pages(params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			content = append(content, page.Contents...)
			return true
		},
	)
	if err != nil {
		return nil, err
	}

	subPaths = make([]FileInfo, 0, len(content))
	for _, sp := range content {
		_, name := s.Split(*sp.Key)
		if name == "" {
			continue
		}

		subPaths = append(subPaths, FileInfo{
			Name:    name,
			Path:    *sp.Key,
			ModTime: *sp.LastModified,
		})
	}

	return subPaths, err
}

func (s S3Store) Info(path string) (info FileInfo, err error) {
	panic("implement me")
}

func (s S3Store) FullName(path string) (fullPath string, err error) {
	panic("implement me")
}

func (s S3Store) Split(path string) (directory string, filename string) {
	parts := strings.Split(path, "/")
	filename = parts[len(parts)-1]
	directory = strings.TrimSuffix(path, filename)
	return directory, filename
}

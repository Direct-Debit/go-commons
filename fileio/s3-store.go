package fileio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/Direct-Debit/go-commons/stdext"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	log "github.com/sirupsen/logrus"
)

type S3Store struct {
	s3              *s3.Client    // AWS S3 client
	Bucket          *string       // Name of the S3 bucket
	PresignDuration time.Duration // Duration for which the presigned URL is valid
}

// NewS3Store creates a new S3Store with the specified bucket name.
// It initializes the AWS session and S3 client.
// The PresignDuration is set to 24 hours by default.
func NewS3Store(bucket string) S3Store {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	return S3Store{s3: s3.NewFromConfig(cfg), Bucket: &bucket, PresignDuration: 24 * time.Hour}
}

func (s S3Store) Save(path string, content string) error {
	return s.SaveStream(path, strings.NewReader(content))
}

func (s S3Store) SaveStream(path string, content io.ReadSeeker) error {
	_, err := s.s3.PutObject(context.TODO(), &s3.PutObjectInput{
		Body:   content,
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
	output, err := s.s3.GetObject(context.TODO(), &s3.GetObjectInput{
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

	_, err := s.s3.CopyObject(context.TODO(), &s3.CopyObjectInput{
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
	_, err := s.s3.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
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

	var content []s3types.Object
	paginator := s3.NewListObjectsV2Paginator(s.s3, params)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		content = append(content, page.Contents...)
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

func (s S3Store) GetInfo(path string) (info FileInfo, err error) {
	output, err := s.s3.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: s.Bucket,
		Key:    &path,
	})
	if err != nil {
		return FileInfo{}, err
	}
	info = FileInfo{
		Name:    path,
		ModTime: *output.LastModified,
		Size:    *output.ContentLength,
	}
	return info, nil
}

func (s S3Store) GetFullName(path string) (fullPath string, err error) {
	fullPath = fmt.Sprintf("s3://%s/%s", *s.Bucket, strings.TrimPrefix(path, "/"))
	return fullPath, nil
}

func (s S3Store) Split(path string) (directory string, filename string) {
	parts := strings.Split(path, "/")
	filename = parts[len(parts)-1]
	directory = strings.TrimSuffix(path, filename)
	return directory, filename
}

func (s S3Store) GenerateDownloadLink(filePath string) (string, error) {
	presignClient := s3.NewPresignClient(s.s3)
	req, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: s.Bucket,
		Key:    &filePath,
	}, s3.WithPresignExpires(s.PresignDuration))
	if err != nil {
		return "", stdext.WrapError(err, "Couldn't create S3 presigned URL")
	}

	return req.URL, nil
}

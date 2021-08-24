package fileio

import (
	"fmt"
	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

type SFTPStore struct {
	Address  string
	User     string
	Password string

	client *sftp.Client
}

func (S *SFTPStore) connect() error {
	conf := &ssh.ClientConfig{
		User:            S.User,
		Auth:            []ssh.AuthMethod{ssh.Password(S.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", S.Address, conf)
	if err != nil {
		return errors.Wrap(err, "failed to dial ssh")
	}
	S.client, err = sftp.NewClient(conn)
	return errors.Wrap(err, "failed to create sftp client")
}

func (S *SFTPStore) disconnect() {
	if S.client == (*sftp.Client)(nil) {
		return
	}
	errlib.WarnError(S.client.Close(), "Could not disconnect from SFTP")
}

func (S *SFTPStore) Save(path string, content string) error {
	if err := S.connect(); err != nil {
		return err
	}
	defer S.disconnect()

	file, err := S.client.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return errors.Wrap(err, "could not create SFTP file "+path)
	}
	defer func() {
		errlib.WarnError(file.Close(), "Couldn't close SFTP file")
	}()

	_, err = file.Write([]byte(content))
	return errors.Wrap(err, "could not write to SFTP file "+path)
}

func (S *SFTPStore) Load(path string) (content string, err error) {
	if err := S.connect(); err != nil {
		return "", err
	}
	defer S.disconnect()

	file, err := S.client.OpenFile(path, os.O_RDONLY)
	if err != nil {
		return "", errors.Wrap(err, "failed to open SFTP file")
	}
	defer func() {
		errlib.WarnError(file.Close(), "Couldn't close SFTP file")
	}()

	var strBuilder strings.Builder
	_, err = file.WriteTo(&strBuilder)
	if err != nil {
		return "", errors.Wrap(err, "failed to write content to string builder")
	}
	return strBuilder.String(), nil
}

func (S *SFTPStore) Move(path string, targetDir string) error {
	panic("implement me")
}

func (S *SFTPStore) Delete(path string) error {
	if err := S.connect(); err != nil {
		return err
	}
	defer S.disconnect()

	return S.client.Remove(path)
}

func (S *SFTPStore) List(path string) (subPaths []FileInfo, err error) {
	if err := S.connect(); err != nil {
		return nil, err
	}
	defer S.disconnect()

	inf, err := S.client.ReadDir(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read SFTP directory")
	}
	subPaths = make([]FileInfo, len(inf))
	for i, info := range inf {
		subPaths[i] = FileInfo{
			Name:    info.Name(),
			Path:    fmt.Sprintf("%s/%s", path, info.Name()),
			ModTime: info.ModTime(),
		}
	}
	return subPaths, nil
}

func (S *SFTPStore) Info(path string) (info FileInfo, err error) {
	panic("implement me")
}

func (S *SFTPStore) FullName(path string) (fullPath string, err error) {
	panic("implement me")
}

func (S *SFTPStore) Split(path string) (directory string, filename string) {
	panic("implement me")
}

func (S *SFTPStore) UploadPath(userCode string, filename string) string {
	panic("implement me")
}

func (S *SFTPStore) DownloadPath(userCode string, filename string) string {
	panic("implement me")
}

package fileio

import (
	"fmt"
	"github.com/Direct-Debit/go-commons/errlib"
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

func (S SFTPStore) connect() error {
	conf := &ssh.ClientConfig{
		User:            S.User,
		Auth:            []ssh.AuthMethod{ssh.Password(S.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", S.Address, conf)
	if errlib.ErrorError(err, "Failed to dial ssh") {
		return err
	}
	S.client, err = sftp.NewClient(conn)
	errlib.ErrorError(err, "Failed to create sftp client")
	return err
}

func (S SFTPStore) disconnect() {
	errlib.ErrorError(S.client.Close(), "Could not disconnect from SFTP")
}

func (S SFTPStore) Save(path string, content string) error {
	if err := S.connect(); err != nil {
		return err
	}
	defer S.disconnect()

	file, err := S.client.OpenFile(path, os.O_WRONLY)
	if err != nil {
		return err
	}
	defer errlib.ErrorError(file.Close(), "Couldn't close SFTP file")

	_, err = file.Write([]byte(content))
	return err
}

func (S SFTPStore) Load(path string) (content string, err error) {
	if err := S.connect(); err != nil {
		return "", err
	}
	defer S.disconnect()

	file, err := S.client.OpenFile(path, os.O_RDONLY)
	if err != nil {
		return "", err
	}
	defer errlib.ErrorError(file.Close(), "Couldn't close SFTP file")

	var strBuilder strings.Builder
	_, err = file.WriteTo(&strBuilder)
	if err != nil {
		return "", err
	}
	return strBuilder.String(), nil
}

func (S SFTPStore) Move(path string, targetDir string) error {
	panic("implement me")
}

func (S SFTPStore) List(path string) (subPaths []FileInfo, err error) {
	if err := S.connect(); err != nil {
		return nil, err
	}
	defer S.disconnect()

	inf, err := S.client.ReadDir(path)
	if err != nil {
		return nil, err
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

func (S SFTPStore) Info(path string) (info FileInfo, err error) {
	panic("implement me")
}

func (S SFTPStore) FullName(path string) (fullPath string, err error) {
	panic("implement me")
}

func (S SFTPStore) Split(path string) (directory string, filename string) {
	panic("implement me")
}

func (S SFTPStore) UploadPath(userCode string, filename string) string {
	panic("implement me")
}

func (S SFTPStore) DownloadPath(userCode string, filename string) string {
	panic("implement me")
}

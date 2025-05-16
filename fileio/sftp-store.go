package fileio

import (
	"fmt"
	"os"
	"strings"

	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type SFTPStore struct {
	Address        string
	User           string
	Password       string
	PrivateKeyPath string

	client *sftp.Client
}

type resetError struct {
	message string
}

func (e *resetError) Error() string {
	return e.message
}

func (r *resetError) Is(err error) bool {
	_, ok := err.(*resetError)
	return ok
}

func (S *SFTPStore) connect() error {
	conf := &ssh.ClientConfig{
		User:            S.User,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if S.PrivateKeyPath != "" {
		key, err := os.ReadFile(S.PrivateKeyPath)
		if err != nil {
			return errors.Wrap(err, "failed to read private key")
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return errors.Wrap(err, "failed to parse private key")
		}
		conf.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else if S.Password != "" {
		conf.Auth = append(conf.Auth, ssh.Password(S.Password))
	} else {
		return errors.New("SFTP Store no authentication method provided")
	}

	conn, err := ssh.Dial("tcp", S.Address, conf)
	if err != nil {
		err := errors.Wrap(err, "failed to dial ssh")
		if strings.Contains(err.Error(), "connection reset by peer") {
			logrus.WithField("address", S.Address).Warn(err)
			return &resetError{message: err.Error()}
		}
		return err
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
	err := S.connect()
	if err != nil && !errors.Is(err, &resetError{}) {
		return err
	}
	defer S.disconnect()

	if err != nil && errors.Is(err, &resetError{}) {
		return nil
	}

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
	err = S.connect()
	if err != nil && !errors.Is(err, &resetError{}) {
		return "", err
	}
	defer S.disconnect()

	if err != nil && errors.Is(err, &resetError{}) {
		return "", nil
	}

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
	err := S.connect()
	if err != nil && !errors.Is(err, &resetError{}) {
		return err
	}
	defer S.disconnect()

	if err != nil && errors.Is(err, &resetError{}) {
		return nil
	}

	return S.client.Remove(path)
}

func (S *SFTPStore) List(path string) (subPaths []FileInfo, err error) {
	err = S.connect()
	if err != nil && !errors.Is(err, &resetError{}) {
		return nil, err
	}
	defer S.disconnect()

	if err != nil && errors.Is(err, &resetError{}) {
		return nil, nil
	}

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

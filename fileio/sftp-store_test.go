package fileio

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSFTPPrivateKey(t *testing.T) {
	sftpStore := SFTPStore{
		Address:        os.Getenv("SFTP_ADDRESS"),
		User:           os.Getenv("SFTP_USER"),
		PrivateKeyPath: os.Getenv("SFTP_PRIVATE_KEY_PATH"),
	}

	err := sftpStore.connect()
	assert.NoError(t, err)

	files, err := sftpStore.List("/")
	assert.NotEmpty(t, files)
	assert.NoError(t, err)

	sftpStore.disconnect()
}

func TestSFTPPassword(t *testing.T) {
	sftpStore := SFTPStore{
		Address:  "localhost:22",
		User:     os.Getenv("SFTP_USER"),
		Password: os.Getenv("SFTP_PASSWORD"),
	}

	err := sftpStore.connect()
	assert.NoError(t, err)

	files, err := sftpStore.List("/")
	assert.NotEmpty(t, files)
	assert.NoError(t, err)

	sftpStore.disconnect()
}

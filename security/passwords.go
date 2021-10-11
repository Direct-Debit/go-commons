package security

import (
	"crypto/rand"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

func GenSalt(length int) ([]byte, error) {
	s := make([]byte, length)
	_, err := rand.Read(s)
	return s, errors.Wrapf(err, "could not generate password salt")
}

func HashPassword(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}

package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	p := HashPassword("blikkies", []byte("blikkies"))
	assert.Equal(t, 32, len(p))
	assert.NotEqual(t, []byte("blikkies")[0], p[0])
	assert.NotEqual(t, []byte("blikkies")[1], p[1])
	assert.NotEqual(t, []byte("blikkies")[2], p[2])
	assert.NotEqual(t, []byte("blikkies")[3], p[3])
}

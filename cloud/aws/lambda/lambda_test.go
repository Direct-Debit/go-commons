package lambda

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Validate(t *testing.T) {
	res, err := NewClient().Validate("123456789", "450905", "savings")
	assert.Nil(t, err)
	assert.Equal(t, res, make(map[string]string))

	res, err = NewClient().Validate("123456789", "123456", "savings")
	assert.Nil(t, err)
	assert.Greater(t, len(res), 0)
}

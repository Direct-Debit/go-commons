package errlib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorError(t *testing.T) {
	b := ErrorError(fmt.Errorf("%s", "err"), "%s", "message")
	assert.True(t, b)
}

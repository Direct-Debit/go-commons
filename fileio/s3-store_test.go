package fileio

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestS3Store_List(t *testing.T) {
	testStore := NewS3Store("dps-production")
	files, err := testStore.List("grobank/reports/")
	assert.NoError(t, err)
	assert.LessOrEqual(t, 1000, len(files))
	println(len(files))
}

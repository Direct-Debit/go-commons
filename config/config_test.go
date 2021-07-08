package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetStrList(t *testing.T) {
	l := GetStrList("test_str_list")
	assert.Equal(t, 2, len(l))
	assert.Equal(t, "element1", l[0])

	assert.Panics(t, func() { GetStrList("test_str_list_ints") })
}

func TestGetStrListDef(t *testing.T) {
	l := GetStrListDef("no-exist", []string{"not configured"})
	assert.Equal(t, 1, len(l))
	assert.Equal(t, "not configured", l[0])

	l = GetStrListDef("test_str_list", []string{})
	assert.Equal(t, 2, len(l))
	assert.Equal(t, "element1", l[0])

	assert.Panics(t, func() { GetStrListDef("test_str_list_ints", []string{}) })
}

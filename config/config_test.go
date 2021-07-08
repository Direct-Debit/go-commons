package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInt64(t *testing.T) {
	l := GetInt64("test_int")
	assert.Equal(t, int64(7), l)
}

func TestGetInt(t *testing.T) {
	l := GetInt("test_int")
	assert.Equal(t, 7, l)
}

func TestGetIntDef(t *testing.T) {
	l := GetIntDef("test_int", 12)
	assert.Equal(t, 7, l)

	l = GetIntDef("missing_test_int", 12)
	assert.Equal(t, 12, l)
}

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

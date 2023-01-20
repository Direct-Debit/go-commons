package stdext

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIsNumeric(t *testing.T) {
	tables := []struct {
		in  string
		out bool
	}{
		{"0", true},
		{"7", true},
		{"A", false},
		{"H", false},
		{"10", true},
		{"16", true},
		{"1D", false},
		{"20", true},
		{"CA", false},
		{"100", true},
		{"10Q", false},
		{"LFLS", false},
		{"1234567890", true},
	}

	for _, table := range tables {
		result := IsNumeric(table.in)
		if result != table.out {
			t.Errorf("For %s: got %v, wanted %v", table.in, result, table.out)
		}
	}
}

func TestAbs(t *testing.T) {
	tables := []struct {
		in  int
		out int
	}{
		{0, 0},
		{7, 7},
		{-7, 7},
	}

	for _, table := range tables {
		result := Abs(table.in)
		if result != table.out {
			t.Errorf("For %d: got %d, wanted %d", table.in, result, table.out)
		}
	}
}

func TestRoundTo(t *testing.T) {
	tables := []struct {
		x        float64
		decimals int
		out      float64
	}{
		{0, 0, 0},
		{0.3, 0, 0},
		{0.7, 0, 1},
		{2387.9699997, 2, 2387.97},
	}

	for _, table := range tables {
		result := RoundTo(table.x, table.decimals)
		if result != table.out {
			t.Errorf("For %f, %d: got %f, wanted %f", table.x, table.decimals, result, table.out)
		}
	}
}

func TestCentToRand(t *testing.T) {
	tables := []struct {
		cent int
		rand float64
	}{
		{0, 0},
		{30, 0.3},
		{70, 0.7},
		{238796, 2387.96},
		{238797, 2387.97},
		{1000, 10},
		{10000, 100},
	}

	for _, table := range tables {
		result := CentToRand(table.cent)
		if result != table.rand {
			t.Errorf("For %d: got %f, wanted %f", table.cent, result, table.rand)
		}
	}
}

func TestStrSet(t *testing.T) {
	s := make(StrSet)
	s.Add("hh")
	s.Add("ah")
	s.Add("bh")
	s.Remove("bh")
	s.Remove("yh")

	assert.True(t, s.Has("hh"))
	assert.True(t, s.Has("ah"))
	assert.False(t, s.Has("bh"))
	assert.False(t, s.Has("xh"))
	assert.False(t, s.Has("yh"))

	assert.Equal(t, len(s), 2)
	assert.Equal(t, len(s.ToSlice()), 2)
}

func TestPtr(t *testing.T) {
	assert.Equal(t, "str", *(Ptr("str")))
	assert.Equal(t, 0, *(Ptr(0)))
	assert.Equal(t, 0.5, *(Ptr(0.5)))
	assert.Equal(t, true, *(Ptr(true)))
}

func TestRoundUp(t *testing.T) {
	assert.Equal(t, 0, RoundUp(0, 500))
	assert.Equal(t, 10_500_00, RoundUp(10_345_10, 500_00))
	assert.Equal(t, 10_500_00, RoundUp(10_500_00, 500_00))
	assert.Equal(t, -10_000_00, RoundUp(-10_345_10, 500_00))
	assert.Equal(t, -10_500_00, RoundUp(-10_500_00, 500_00))
}

func TestCachedCall(t *testing.T) {
	f := func(i string) (string, error) {
		time.Sleep(time.Second)
		return i, nil
	}

	m := make(map[string]string)
	v, err := CachedCall(f, "i", m)
	assert.NoError(t, err)
	assert.Equal(t, "i", v)
	assert.Equal(t, 1, len(m))

	v, err = CachedCall(f, "i", m)
	assert.NoError(t, err)
	assert.Equal(t, "i", v)
}

func TestSafeSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	assert.Equal(t, 2, len(SafeSlice(slice, 2, 4)))
	assert.Equal(t, 0, len(SafeSlice(slice, 4, 3)))
	assert.Equal(t, 2, len(SafeSlice(slice, 3, 100)))
	assert.Equal(t, 2, len(SafeSlice(slice, -6, 2)))
	assert.Equal(t, 5, len(SafeSlice(slice, -6, 100)))
	assert.Equal(t, 0, len(SafeSlice(slice, 32, 100)))
	assert.Equal(t, 0, len(SafeSlice(slice, -7, -1)))
}

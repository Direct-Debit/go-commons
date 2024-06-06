package stdext

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

func TestJoinMaps(t *testing.T) {
	mt := map[string]any{"a": 1, "b": "2"}
	ms := map[string]any{"b": 2, "c": "3"}

	JoinMaps(mt, ms)
	assert.Equal(t, 3, len(mt))
	assert.Equal(t, 1, mt["a"])
	assert.Equal(t, 2, mt["b"])
	assert.Equal(t, "3", mt["c"])
}

func BenchmarkFlatten(b *testing.B) {
	lists := make([][]int, 10_000)
	for i := range lists {
		lists[i] = make([]int, i+10)
		for j := range lists[i] {
			lists[i][j] = i + j
		}
	}

	for i := 0; i < b.N; i++ {
		Flatten(lists)
	}
}

func TestChunkify(t *testing.T) {
	list := make([]int, 247)
	for i := range list {
		list[i] = i + 1
	}

	lists := ChunkifyBySize(list, 32)
	for _, l := range lists {
		assert.LessOrEqual(t, len(l), 32)
		fmt.Println(l)
	}
	assert.Equal(t, 8, len(lists))

	lists = ChunkifyByCount(list, 32)
	for _, l := range lists {
		assert.LessOrEqual(t, len(l), 8)
		fmt.Println(l)
	}
	assert.Equal(t, 32, len(lists))

	flattened := Flatten(lists)
	assert.ElementsMatch(t, list, flattened)
}

func TestEllipticalTruncate(t *testing.T) {
	// Test with a string that has no spaces
	test_string := "1234567890AJQK"
	assert.Equal(t, "...", EllipticalTruncate(test_string, 3))
	assert.Equal(t, "1...", EllipticalTruncate(test_string, 4))
	assert.Equal(t, "12...", EllipticalTruncate(test_string, 5))
	assert.Equal(t, "123456...", EllipticalTruncate(test_string, 9))
	assert.Equal(t, "1234567...", EllipticalTruncate(test_string, 10))
	assert.Equal(t, "12345678...", EllipticalTruncate(test_string, 11))
	assert.Equal(t, test_string, EllipticalTruncate(test_string, len(test_string)))

	// Test with a string that has spaces
	test_string = "1 2 3 45 6 7 8 9 0 A J Q K"
	assert.Equal(t, "1 2 3 45...", EllipticalTruncate(test_string, 11))
	test_string = "1 2 3 4 5 6 7 8 9 0 A J Q K"
	assert.Equal(t, "1 2 3 4 ...", EllipticalTruncate(test_string, 11))

	// Test with a long string
	long_string := strings.Repeat("a ", 2024)
	assert.Equal(t, long_string[:1021]+"...", EllipticalTruncate(long_string, 1024))

	long_string = "[UnitTests] This is a test alert. It has more than 1024 it should truncate chars(aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa)"
	assert.Equal(t, long_string[:1021]+"...", EllipticalTruncate(long_string, 1024))

	logrus.Infof("EllipticalTruncate test passed with long string %v", long_string)
	logrus.Infof("EllipticalTruncate test passed with long result %v", EllipticalTruncate(long_string, 1024))
	logrus.Info("EllipticalTruncate test passed")
}

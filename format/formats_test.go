package format

import "testing"

func TestIntToBase36(t *testing.T) {
	tables := []struct {
		in  int
		out string
	}{
		{0, "0"},
		{7, "7"},
		{10, "A"},
		{17, "H"},
		{36, "10"},
		{42, "16"},
		{49, "1D"},
		{72, "20"},
		{442, "CA"},
		{1296, "100"},
		{1322, "10Q"},
		{1000000, "LFLS"},
	}

	for _, table := range tables {
		result := IntToBase36(table.in)
		if result != table.out {
			t.Errorf("Got %s, wanted %s", result, table.out)
		}
	}
}

func TestBase36toInt(t *testing.T) {
	tables := []struct {
		in      string
		out     int
		invalid bool
	}{
		{"0`", 0, true},
		{"0", 0, false},
		{"7", 7, false},
		{"00007", 7, false},
		{"A", 10, false},
		{"H", 17, false},
		{"10", 36, false},
		{"16", 42, false},
		{"1D", 49, false},
		{"20", 72, false},
		{"CA", 442, false},
		{"100", 1296, false},
		{"10Q", 1322, false},
		{"0010Q", 1322, false},
		{"LFLS", 1000000, false},
		{"lfls", 1000000, false},
		{"0)(", 0, true},
		{"0(AB)", 0, true},
	}

	for _, table := range tables {
		result, err := Base36toInt(table.in)
		if table.invalid && err == nil {
			t.Errorf("Expected error, but didn't get one for %v", table.in)
		}
		if !table.invalid && err != nil {
			t.Errorf("Expected no error for %v, but got: %v", table.in, err)
		}
		if result != table.out {
			t.Errorf("Got %v, wanted %v", result, table.out)
		}
	}
}

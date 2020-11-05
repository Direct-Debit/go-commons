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

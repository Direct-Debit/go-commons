package stdext

import "testing"

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
			t.Errorf("For %s, got %v, wanted %v", table.in, result, table.out)
		}
	}
}

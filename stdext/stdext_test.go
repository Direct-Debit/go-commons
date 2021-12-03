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

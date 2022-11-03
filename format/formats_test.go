package format

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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

func TestCentToCommaRand(t *testing.T) {
	tables := []struct {
		in  int
		out string
	}{
		{0, "0,00"},
		{7, "0,07"},
		{10, "0,10"},
		{17, "0,17"},
		{36, "0,36"},
		{42, "0,42"},
		{49, "0,49"},
		{72, "0,72"},
		{442, "4,42"},
		{1296, "12,96"},
		{1322, "13,22"},
		{1000000, "10000,00"},
	}

	for _, table := range tables {
		result := CentToCommaRand(table.in)
		if result != table.out {
			t.Errorf("Got %s, wanted %s", result, table.out)
		}
	}
}

func TestAnyAmountToCent(t *testing.T) {
	tables := []struct {
		in  string
		out int
	}{
		{"0,00", 0},
		{"0,07", 7},
		{"0,10", 10},
		{"0,17", 17},
		{"0,36", 36},
		{"0,42", 42},
		{"0,49", 49},
		{"0,72", 72},
		{"4,42", 442},
		{"12,96", 1296},
		{"13,22", 1322},
		{"10000,00", 1000000},
		{"2,194.50", 219450},
		{"0.00", 0},
	}

	for _, table := range tables {
		result, err := AnyAmountToCent(table.in)
		log.Infof("%s -> %d", table.in, result)

		assert.Nil(t, err)
		assert.Equal(t, table.out, result)
	}
}

func TestDateFormats(t *testing.T) {
	testTime := time.Date(2021, 9, 9, 9, 47, 12, 34, time.UTC)
	tests := []struct {
		format  string
		timeStr string
	}{
		{DateShort6, "210909"},
		{DateShort6Slashes, "21/09/09"},
		{DateShort8, "20210909"},
		{DateShortSlashes, "2021/09/09"},
		{DateShortDashes, "2021-09-09"},
		{DateTimeCompact, "210909094712"},
		{DateTimeShort, "09/09/2021 09:47"},
		{DateTimeShortDashes, "2021-09-09 09:47:12"},
		{DDsMMsYYYY, "09/09/2021"},
		{MonthYY, "Sep21"},
	}

	for _, c := range tests {
		res := testTime.Format(c.format)
		assert.Equal(t, c.timeStr, res)

		_, err := time.Parse(c.format, c.timeStr)
		assert.Nil(t, err)
	}
}

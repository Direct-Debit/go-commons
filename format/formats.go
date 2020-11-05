package format

import (
	"errors"
	"fmt"
	"github.com/Direct-Debit/go-commons/errlib"
	"math"
	"strconv"
	"strings"
)

const (
	DateShort6          = "060102"
	DateShort6Slashes   = "06/01/02"
	DateShort8          = "20060102"
	DateShortSlashes    = "2006/01/02"
	DateShortDashes     = "2006-01-02"
	DateTimeShort       = "02/01/2006 15:04"
	DateTimeShortDashes = "2006-01-02 15:04:05"
)

func CentToCommaRand(cent int) string {
	r := cent / 100
	c := cent % 100
	return fmt.Sprintf("%d,%02d", r, c)
}

func intToBase36Digit(i int) (string, error) {
	if i < 0 || i >= 36 {
		return "", errors.New(fmt.Sprintf("Can't convert %v to base 36 digit", i))
	}
	if i < 10 {
		return strconv.Itoa(i), nil
	}
	alphaPos := i - 10
	iVal := 'A' + rune(alphaPos)
	return string(iVal), nil
}

// Convert an integer to base36 where A-Z represent digits with values 10-35
func IntToBase36(i int) string {
	if i == 0 {
		return "0"
	}

	digits := make([]string, 1)
	unconverted := i
	for exp := 0; unconverted > 0; exp++ {
		digitPos := int(math.Pow(36.0, float64(exp)))
		digitVal := (unconverted / digitPos) % 36

		digit, err := intToBase36Digit(digitVal)
		errlib.FatalError(err, "Our math is broken")
		digits = append(digits, digit)

		unconverted = unconverted - (digitVal * digitPos)
	}

	var result strings.Builder
	for idx := len(digits) - 1; idx >= 0; idx-- {
		result.WriteString(digits[idx])
	}
	return result.String()
}

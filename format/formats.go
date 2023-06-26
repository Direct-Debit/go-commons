package format

import (
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
	DateTimeCompact     = "060102150405"
	DateTimeShort       = "02/01/2006 15:04"
	DateTimeShortDashes = "2006-01-02 15:04:05"
	DDsMMsYYYY          = "02/01/2006"
	DDdMMdYYYY          = "02-01-2006"
	DDdMMMdYY           = "02-JAN-06"
	MonthYY             = "Jan06"
	MMYY                = "0106"
	YYYYdMM             = "2006-01"
	RFC3339NanoFixed    = "2006-01-02T15:04:05.000000000Z07:00"
)

func CentToCommaRand(cent int) string {
	r := cent / 100
	c := cent % 100
	return fmt.Sprintf("%d,%02d", r, c)
}

func AnyAmountToCent(amount string) (int, error) {
	replacements := []string{",", "", ".", "", "R", "", " ", ""}
	if len(replacements)%2 != 0 {
		return 0, fmt.Errorf("uneven arguments for replacer %d, this is a devloper error", len(replacements))
	}
	replacer := strings.NewReplacer(replacements...)

	res, err := strconv.Atoi(replacer.Replace(amount))
	if err != nil {
		return 0, fmt.Errorf("failed to convert %s to cents: %w", amount, err)
	}
	return res, nil
}

func intToBase36Digit(i int) (string, error) {
	if i < 0 || i >= 36 {
		return "", fmt.Errorf("can't convert %v to base 36 digit", i)
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

func base36DigitToInt(d rune) (int, error) {
	if d >= 'A' && d <= 'Z' {
		return int(d-'A') + 10, nil
	}
	if d >= '0' && d <= '9' {
		return int(d - '0'), nil
	}
	return 0, fmt.Errorf("can't read %c as base 36 digit", d)
}

// Convert a base36 to int where A-Z represent digits with values 10-35
func Base36toInt(s string) (int, error) {
	s = strings.ToUpper(s)
	result := 0
	for i, d := range s {
		v, err := base36DigitToInt(d)
		if err != nil {
			return 0, err
		}
		digitPos := int(math.Pow(36.0, float64(len(s)-i-1)))
		result += v * digitPos
	}
	return result, nil
}

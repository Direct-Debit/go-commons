package stdext

import (
	"bufio"
	"math"
	"strconv"
	"strings"
)

// SplitParts splits the string into parts of max length n
func SplitParts(s string, n int) []string {
	if len(s) == 0 {
		return []string{}
	}

	count := int(math.Ceil(float64(len(s)) / float64(n)))
	result := make([]string, count)

	for i := 0; i < count; i++ {
		start := Min(i*n, len(s)-1)
		end := Min((i+1)*n, len(s))
		result[i] = s[start:end]
	}

	return result
}

func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

func IsNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func ContainsOnly(s, set string) bool {
	for _, c := range s {
		if !strings.ContainsRune(set, c) {
			return false
		}
	}
	return true
}

// ContainsI checks whether s contains substr in a case-insensitive way.
func ContainsI(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// EqualI does a case-insensitive comparison of strings.
// It literally wraps strings.EqualFold, which has a unintuitive name IMO
func EqualI(s1, s2 string) bool {
	return strings.EqualFold(s1, s2)
}

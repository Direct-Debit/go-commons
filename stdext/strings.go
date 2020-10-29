package stdext

import "math"

// Split the string into parts of max length n
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

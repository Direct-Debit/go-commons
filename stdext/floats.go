package stdext

import "math"

// Round to a certain number of digits. Might be prone to rounding errors
func RoundTo(x float64, decimals int) float64 {
	m := math.Pow10(decimals)
	return math.Round(x*m) / m
}

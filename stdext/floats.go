package stdext

import "math"

// Round to a certain number of digits. Might be prone to rounding errors
func RoundTo[T Float](x T, decimals int) T {
	m := math.Pow10(decimals)
	return T(math.Round(float64(x)*m) / m)
}

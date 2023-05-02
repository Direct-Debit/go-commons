package stdext

import "math"

func CentToRand(cent int) float64 {
	return RoundTo(float64(cent)/100, 2)
}

func RandToCent(rand float64) int {
	return int(math.Round(rand * 100))
}

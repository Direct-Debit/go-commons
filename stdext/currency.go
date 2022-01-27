package stdext

func CentToRand(cent int) float64 {
	return RoundTo(float64(cent)/100, 2)
}

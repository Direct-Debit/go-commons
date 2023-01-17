package stdext

func RoundUp[T Integer](x, increment T) T {
	part := x % increment
	if part == T(0) {
		return x
	}
	if part < 0 {
		return x - part
	}
	return x + (increment - part)
}

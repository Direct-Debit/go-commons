package stdext

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func Abs[T Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

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

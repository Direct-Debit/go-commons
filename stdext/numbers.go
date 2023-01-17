package stdext

type Complex interface {
	~complex64 | ~complex128
}

type Float interface {
	~float32 | ~float64
}

type Integer interface {
	Signed | Unsigned
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type RealNumber interface {
	Float | Integer
}

type Number interface {
	Complex | RealNumber
}

// Max returns the larger of x or y.
func Max[T RealNumber](x, y T) T {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min[T RealNumber](x, y T) T {
	if x > y {
		return y
	}
	return x
}

func Abs[T RealNumber](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

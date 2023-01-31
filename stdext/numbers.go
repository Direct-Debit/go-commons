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

// Max returns the largest number in nums, or the zero value of T if there are no nums
func Max[T RealNumber](nums ...T) T {
	if len(nums) == 0 {
		var nothing T
		return nothing
	}

	max := nums[0]
	for _, y := range nums {
		if y > max {
			max = y
		}
	}
	return max
}

// Min returns the smallest number in nums, or the zero value of T if there are no nums
func Min[T RealNumber](nums ...T) T {
	if len(nums) == 0 {
		var nothing T
		return nothing
	}

	min := nums[0]
	for _, y := range nums {
		if y < min {
			min = y
		}
	}
	return min
}

func Abs[T RealNumber](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

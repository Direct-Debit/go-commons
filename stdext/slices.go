package stdext

func SafeIdx[T any](idx int, arr []T) T {
	if 0 <= idx && idx < len(arr) {
		return arr[idx]
	}
	var nothing T
	return nothing
}

func SafeSlice[T any](slice []T, start, end int) []T {
	end = Max(0, Min(len(slice), end))
	start = Min(Max(0, start), end)
	return slice[start:end]
}

func FindInSlice[T comparable](slice []T, val T) (int, bool) {
	for k, item := range slice {
		if item == val {
			return k, true
		}
	}
	return -1, false
}

func Filter[T any](slice []T, keep func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, t := range slice {
		if keep(t) {
			result = append(result, t)
		}
	}
	return result
}

func Map[T any, U any](slice []T, conv func(T) U) []U {
	result := make([]U, len(slice))
	for i, t := range slice {
		result[i] = conv(t)
	}
	return result
}

func SafeIdxStr(idx int, arr []string) string {
	return SafeIdx(idx, arr)
}

func FindInStrSlice(slice []string, val string) (int, bool) {
	return FindInSlice(slice, val)
}

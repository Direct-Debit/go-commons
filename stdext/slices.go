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

func SafeIdxStr(idx int, arr []string) string {
	return SafeIdx(idx, arr)
}

func FindInStrSlice(slice []string, val string) (int, bool) {
	return FindInSlice(slice, val)
}

package stdext

func SafeIdxStr(idx int, arr []string) string {
	if idx < len(arr) {
		return arr[idx]
	}
	return ""
}

func SafeIdx[T any](idx int, arr []T) T {
	if 0 <= idx && idx < len(arr) {
		return arr[idx]
	}
	var nothing T
	return nothing
}

func FindInStrSlice(slice []string, val string) (int, bool) {
	for k, item := range slice {
		if item == val {
			return k, true
		}
	}
	return -1, false
}

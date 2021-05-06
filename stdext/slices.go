package stdext

func SafeIdxStr(idx int, arr []string) string {
	if idx < len(arr) {
		return arr[idx]
	}
	return ""
}

func Contains(slice []string, val string) (int, bool) {
	for k, item := range slice {
		if item == val {
			return k, true
		}
	}
	return -1, false
}

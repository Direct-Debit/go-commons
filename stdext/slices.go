package stdext

func SafeIdxStr(idx int, arr []string) string {
	if idx < len(arr) {
		return arr[idx]
	}
	return ""
}

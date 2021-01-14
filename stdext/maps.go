package stdext

func MapStrStrGetDef(m map[string]string, k string, d string) string {
	if v, ok := m[k]; ok {
		return v
	}
	return d
}

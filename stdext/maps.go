package stdext

func MapStrStrGetDef(m map[string]string, k string, d string) string {
	if v, ok := m[k]; ok {
		return v
	}
	return d
}

// MapGetDef returns the value at m[k] if k is in m, otherwise it returns d
func MapGetDef[K comparable, V any](m map[K]V, k K, d V) V {
	if v, ok := m[k]; ok {
		return v
	}
	return d
}

func JoinMaps(target map[string]interface{}, source map[string]interface{}) {
	for k, v := range source {
		target[k] = v
	}
}

// CachedCall tries to get a cached value from the given caching map before calling f and caching the result.
// If f returns an error, CachedCall will return f's output without caching it.
// CachedCall does not generate any errors of its own, but panics if cache is nil.
func CachedCall[K comparable, V any](f func(K) (V, error), input K, cache map[K]V) (V, error) {
	// Check cache and return if found
	if v, ok := cache[input]; ok {
		return v, nil
	}

	// call f to determine value to cache
	v, err := f(input)
	if err != nil {
		return v, err
	}

	// cache value
	cache[input] = v
	return v, nil
}

func Group[K comparable, V any](slice []V, key func(V) K) map[K][]V {
	m := make(map[K][]V)
	for _, v := range slice {
		k := key(v)
		if s, ok := m[k]; ok {
			m[k] = append(s, v)
		} else {
			m[k] = []V{v}
		}
	}
	return m
}

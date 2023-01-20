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
// CachedCall does not generate any errors of its own.
func CachedCall[K comparable, V any](f func(K) (V, error), cache map[K]V, key K) (V, error) {
	// Check cache and return if found
	if v, ok := cache[key]; ok {
		return v, nil
	}

	// call f to determine value to cache
	v, err := f(key)
	if err != nil {
		return v, err
	}

	// cache value
	cache[key] = v
	return v, nil
}

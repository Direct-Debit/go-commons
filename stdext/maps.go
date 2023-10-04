package stdext

// MapGetDef returns the value at m[k] if k is in m, otherwise it returns d
func MapGetDef[K comparable, V any](m map[K]V, k K, d V) V {
	if v, ok := m[k]; ok {
		return v
	}
	return d
}

// JoinMaps copies the values of the source map into the target map,
// overriding the target map if there is a key conflict.
// The target map is modified in place.
func JoinMaps[K comparable, V any](target map[K]V, source map[K]V) {
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

func MapStrStrGetDef(m map[string]string, k string, d string) string {
	return MapGetDef(m, k, d)
}

func MapToSlice[K comparable, V any, T any](m map[K]V, flatten func(K, V) T) []T {
	result := make([]T, 0, len(m))
	for k, v := range m {
		result = append(result, flatten(k, v))
	}
	return result
}

func SliceToMap[K comparable, V any](slice []V, key func(V) K) map[K]V {
	result := make(map[K]V, len(slice))
	for _, v := range slice {
		result[key(v)] = v
	}
	return result
}

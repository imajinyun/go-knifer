package maps

// This file provides common collection helpers aligned with hutool-core
// CollUtil, MapUtil, and ListUtil.

// MapIsEmpty reports whether the map is empty.
func MapIsEmpty[K comparable, V any](m map[K]V) bool { return len(m) == 0 }

// MapIsNotEmpty reports whether the map is not empty.
func MapIsNotEmpty[K comparable, V any](m map[K]V) bool { return len(m) > 0 }

// MapKeys returns all keys from the map. The order follows Go map iteration and is not stable.
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MapValues returns all values from the map. The order follows Go map iteration and is not stable.
func MapValues[K comparable, V any](m map[K]V) []V {
	vs := make([]V, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

// MapInverse swaps keys and values. V must be comparable, and later duplicate values overwrite earlier keys.
func MapInverse[K, V comparable](m map[K]V) map[V]K {
	out := make(map[V]K, len(m))
	for k, v := range m {
		out[v] = k
	}
	return out
}

// MapMerge merges maps into a new map; later maps overwrite earlier values for the same key.
func MapMerge[K comparable, V any](maps ...map[K]V) map[K]V {
	out := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

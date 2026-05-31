package maps

// This file provides common collection helpers aligned with hutool-core
// CollUtil, MapUtil, and ListUtil.

// IsEmpty reports whether the map is empty.
func IsEmpty[K comparable, V any](m map[K]V) bool { return len(m) == 0 }

// IsNotEmpty reports whether the map is not empty.
func IsNotEmpty[K comparable, V any](m map[K]V) bool { return len(m) > 0 }

// Keys returns all keys from the map. The order follows Go map iteration and is not stable.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all values from the map. The order follows Go map iteration and is not stable.
func Values[K comparable, V any](m map[K]V) []V {
	vs := make([]V, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

// Inverse swaps keys and values. V must be comparable, and later duplicate values overwrite earlier keys.
func Inverse[K, V comparable](m map[K]V) map[V]K {
	out := make(map[V]K, len(m))
	for k, v := range m {
		out[v] = k
	}
	return out
}

// Merge merges maps into a new map; later maps overwrite earlier values for the same key.
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	out := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

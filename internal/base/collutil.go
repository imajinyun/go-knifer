package base

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

// Union returns the deduplicated union of a and b.
func Union[T comparable](a, b []T) []T {
	return SliceDistinct(append(append([]T{}, a...), b...))
}

// Intersection returns the deduplicated intersection of a and b.
func Intersection[T comparable](a, b []T) []T {
	set := make(map[T]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}
	seen := make(map[T]struct{}, len(a))
	out := make([]T, 0)
	for _, v := range a {
		if _, ok := set[v]; !ok {
			continue
		}
		if _, dup := seen[v]; dup {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

// Subtract returns the set difference a - b while preserving the order from a.
func Subtract[T comparable](a, b []T) []T {
	bset := make(map[T]struct{}, len(b))
	for _, v := range b {
		bset[v] = struct{}{}
	}
	out := make([]T, 0, len(a))
	for _, v := range a {
		if _, ok := bset[v]; !ok {
			out = append(out, v)
		}
	}
	return out
}

// Page returns one page from a slice. pageNo starts from 1; invalid or out-of-range pages return an empty slice.
func Page[T any](a []T, pageNo, pageSize int) []T {
	if pageNo < 1 || pageSize <= 0 || len(a) == 0 {
		return []T{}
	}
	start := (pageNo - 1) * pageSize
	if start >= len(a) {
		return []T{}
	}
	end := start + pageSize
	if end > len(a) {
		end = len(a)
	}
	out := make([]T, end-start)
	copy(out, a[start:end])
	return out
}

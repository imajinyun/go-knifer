package maps

import stdmaps "maps"

// Merge returns the union of the given maps. When the same key appears in
// multiple maps, the value from a later map overrides earlier ones.
//
// Guarantees:
//   - Never returns nil; an empty input yields an empty (non-nil) map.
//   - Input maps are not modified.
func Merge[K comparable, V any](ms ...map[K]V) map[K]V {
	switch len(ms) {
	case 0:
		return make(map[K]V)
	case 1:
		if ms[0] == nil {
			return make(map[K]V)
		}
		return stdmaps.Clone(ms[0])
	}

	total := 0
	for _, m := range ms {
		total += len(m)
	}
	out := make(map[K]V, total)
	for _, m := range ms {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

// MergeFunc is like Merge but resolves conflicts via the supplied function.
// resolve(old, new) is invoked only when both old and new exist for a key.
// A nil resolve falls back to "last-write-wins" semantics.
func MergeFunc[K comparable, V any](resolve func(old, new V) V, ms ...map[K]V) map[K]V {
	if resolve == nil {
		return Merge(ms...)
	}
	total := 0
	for _, m := range ms {
		total += len(m)
	}
	out := make(map[K]V, total)
	for _, m := range ms {
		for k, v := range m {
			if old, ok := out[k]; ok {
				out[k] = resolve(old, v)
			} else {
				out[k] = v
			}
		}
	}
	return out
}

// MergeWithOverwrite merges srcMaps into dstMap in place.
// If a key exists in both dstMap and srcMaps, its value is overwritten;
// when multiple srcMaps share a key, the value from the last one wins.
func MergeWithOverwrite[K comparable, V any](dstMap map[K]V, srcMaps ...map[K]V) {
	merge(dstMap, true, srcMaps...)
}

// MergeWithoutOverwrite merges srcMaps into dstMap in place,
// keeping the existing value whenever a key already exists in dstMap.
// Among srcMaps, the first occurrence of a new key wins.
func MergeWithoutOverwrite[K comparable, V any](dstMap map[K]V, srcMaps ...map[K]V) {
	merge(dstMap, false, srcMaps...)
}

// merge copies all key-value pairs from srcMaps into dst.
// When overwrite is false, keys already present in dst are left unchanged.
func merge[K comparable, V any](dst map[K]V, overwrite bool, srcMaps ...map[K]V) {
	for _, src := range srcMaps {
		for k, v := range src {
			if !overwrite {
				if _, exists := dst[k]; exists {
					continue
				}
			}
			dst[k] = v
		}
	}
}

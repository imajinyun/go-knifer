package maps

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

// MergeCopyWithOverwrite merges srcMaps into a newly allocated map and returns it,
// leaving all inputs untouched.
// When multiple srcMaps share a key, the value from the last one wins.
// It returns nil if srcMaps is empty.
func MergeCopyWithOverwrite[K comparable, V any](srcMaps ...map[K]V) map[K]V {
	return mergeCopy(true, srcMaps...)
}

// MergeCopyWithoutOverwrite merges srcMaps into a newly allocated map and returns it,
// leaving all inputs untouched.
// Among srcMaps, the first occurrence of a key wins.
// It returns nil if srcMaps is empty.
func MergeCopyWithoutOverwrite[K comparable, V any](srcMaps ...map[K]V) map[K]V {
	return mergeCopy(false, srcMaps...)
}

// mergeCopy allocates a new map and merges srcMaps into it.
// The initial capacity is a heuristic: the size of the largest source map.
func mergeCopy[K comparable, V any](overwrite bool, srcMaps ...map[K]V) map[K]V {
	if len(srcMaps) == 0 {
		return nil
	}
	capHint := 0
	for _, src := range srcMaps {
		if len(src) > capHint {
			capHint = len(src)
		}
	}
	result := make(map[K]V, capHint)
	merge(result, overwrite, srcMaps...)
	return result
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

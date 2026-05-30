package slice

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

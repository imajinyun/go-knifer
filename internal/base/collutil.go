package base

// 对应 hutool-core CollUtil / MapUtil / ListUtil 的常用部分。

// MapIsEmpty Map 是否为空。
func MapIsEmpty[K comparable, V any](m map[K]V) bool { return len(m) == 0 }

// MapIsNotEmpty Map 是否非空。
func MapIsNotEmpty[K comparable, V any](m map[K]V) bool { return len(m) > 0 }

// MapKeys 取出 Map 的所有 Key。
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MapValues 取出 Map 的所有 Value。
func MapValues[K comparable, V any](m map[K]V) []V {
	vs := make([]V, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

// MapInverse 反转 Map（K/V 互换；要求 V 也 comparable）。
func MapInverse[K, V comparable](m map[K]V) map[V]K {
	out := make(map[V]K, len(m))
	for k, v := range m {
		out[v] = k
	}
	return out
}

// MapMerge 合并多个 Map，后者覆盖前者。
func MapMerge[K comparable, V any](maps ...map[K]V) map[K]V {
	out := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

// Union 并集（去重）。
func Union[T comparable](a, b []T) []T {
	return SliceDistinct(append(append([]T{}, a...), b...))
}

// Intersection 交集（去重）。
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

// Subtract 差集：a - b。
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

// Page 切片分页：pageNo 从 1 开始。返回该页内容（越界返回空切片）。
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

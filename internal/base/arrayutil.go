package base

import (
	"fmt"
	"strings"
)

// 对应 hutool-core ArrayUtil（基于 Go 泛型）。

// SliceIsEmpty 切片是否为空。
func SliceIsEmpty[T any](a []T) bool { return len(a) == 0 }

// SliceIsNotEmpty 切片是否非空。
func SliceIsNotEmpty[T any](a []T) bool { return len(a) > 0 }

// SliceContains 是否包含元素（要求 comparable）。
func SliceContains[T comparable](a []T, v T) bool {
	for _, x := range a {
		if x == v {
			return true
		}
	}
	return false
}

// SliceIndexOf 元素首次出现位置；不存在返回 -1。
func SliceIndexOf[T comparable](a []T, v T) int {
	for i, x := range a {
		if x == v {
			return i
		}
	}
	return -1
}

// SliceLastIndexOf 元素最后出现位置；不存在返回 -1。
func SliceLastIndexOf[T comparable](a []T, v T) int {
	for i := len(a) - 1; i >= 0; i-- {
		if a[i] == v {
			return i
		}
	}
	return -1
}

// SliceReverse 原地反转。
func SliceReverse[T any](a []T) []T {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}

// SliceDistinct 去重，保持顺序。
func SliceDistinct[T comparable](a []T) []T {
	seen := make(map[T]struct{}, len(a))
	out := make([]T, 0, len(a))
	for _, v := range a {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

// SliceJoin 字符串拼接。
func SliceJoin[T any](a []T, sep string) string {
	parts := make([]string, len(a))
	for i, v := range a {
		parts[i] = fmt.Sprint(v)
	}
	return strings.Join(parts, sep)
}

// SliceFilter 过滤。
func SliceFilter[T any](a []T, pred func(T) bool) []T {
	out := make([]T, 0, len(a))
	for _, v := range a {
		if pred(v) {
			out = append(out, v)
		}
	}
	return out
}

// SliceMap 映射。
func SliceMap[T, R any](a []T, fn func(T) R) []R {
	out := make([]R, len(a))
	for i, v := range a {
		out[i] = fn(v)
	}
	return out
}

// SliceSub 截取子切片，支持负数索引。
func SliceSub[T any](a []T, fromIndex, toIndex int) []T {
	n := len(a)
	if n == 0 {
		return []T{}
	}
	if fromIndex < 0 {
		fromIndex += n
	}
	if toIndex < 0 {
		toIndex += n
	}
	if fromIndex < 0 {
		fromIndex = 0
	}
	if toIndex > n {
		toIndex = n
	}
	if fromIndex > toIndex {
		fromIndex, toIndex = toIndex, fromIndex
	}
	if fromIndex >= n {
		return []T{}
	}
	out := make([]T, toIndex-fromIndex)
	copy(out, a[fromIndex:toIndex])
	return out
}

// SliceConcat 连接多个切片。
func SliceConcat[T any](slices ...[]T) []T {
	total := 0
	for _, s := range slices {
		total += len(s)
	}
	out := make([]T, 0, total)
	for _, s := range slices {
		out = append(out, s...)
	}
	return out
}

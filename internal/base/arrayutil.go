package base

import (
	"fmt"
	"strings"
)

// This file provides generic slice helpers aligned with hutool-core ArrayUtil.
// Functions return new slices where mutation would be surprising, while
// SliceReverse intentionally reverses the input slice in place for efficiency.

// SliceIsEmpty reports whether the slice is empty.
func SliceIsEmpty[T any](a []T) bool { return len(a) == 0 }

// SliceIsNotEmpty reports whether the slice is not empty.
func SliceIsNotEmpty[T any](a []T) bool { return len(a) > 0 }

// SliceContains reports whether the slice contains v. T must be comparable.
func SliceContains[T comparable](a []T, v T) bool {
	for _, x := range a {
		if x == v {
			return true
		}
	}
	return false
}

// SliceIndexOf returns the first index of v, or -1 when v is absent.
func SliceIndexOf[T comparable](a []T, v T) int {
	for i, x := range a {
		if x == v {
			return i
		}
	}
	return -1
}

// SliceLastIndexOf returns the last index of v, or -1 when v is absent.
func SliceLastIndexOf[T comparable](a []T, v T) int {
	for i := len(a) - 1; i >= 0; i-- {
		if a[i] == v {
			return i
		}
	}
	return -1
}

// SliceReverse reverses the input slice in place and returns the same slice.
func SliceReverse[T any](a []T) []T {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}

// SliceDistinct removes duplicates while preserving the first occurrence order.
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

// SliceJoin converts elements with fmt.Sprint and joins them with sep.
func SliceJoin[T any](a []T, sep string) string {
	parts := make([]string, len(a))
	for i, v := range a {
		parts[i] = fmt.Sprint(v)
	}
	return strings.Join(parts, sep)
}

// SliceFilter returns elements for which pred returns true.
func SliceFilter[T any](a []T, pred func(T) bool) []T {
	out := make([]T, 0, len(a))
	for _, v := range a {
		if pred(v) {
			out = append(out, v)
		}
	}
	return out
}

// SliceMap maps each element to another value while preserving order.
func SliceMap[T, R any](a []T, fn func(T) R) []R {
	out := make([]R, len(a))
	for i, v := range a {
		out[i] = fn(v)
	}
	return out
}

// SliceSub returns a copied sub-slice and supports negative indexes.
// Negative indexes are resolved from the end of the slice, and reversed ranges
// are normalized by swapping fromIndex and toIndex, following hutool behavior.
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

// SliceConcat concatenates multiple slices into a new slice.
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

package vslice

import sliceimpl "github.com/imajinyun/go-knifer/internal/slice"

func IsEmpty[T any](a []T) bool                    { return sliceimpl.SliceIsEmpty(a) }
func IsNotEmpty[T any](a []T) bool                 { return sliceimpl.SliceIsNotEmpty(a) }
func Contains[T comparable](a []T, v T) bool       { return sliceimpl.SliceContains(a, v) }
func IndexOf[T comparable](a []T, v T) int         { return sliceimpl.SliceIndexOf(a, v) }
func LastIndexOf[T comparable](a []T, v T) int     { return sliceimpl.SliceLastIndexOf(a, v) }
func Reverse[T any](a []T) []T                     { return sliceimpl.SliceReverse(a) }
func Distinct[T comparable](a []T) []T             { return sliceimpl.SliceDistinct(a) }
func Join[T any](a []T, sep string) string         { return sliceimpl.SliceJoin(a, sep) }
func Filter[T any](a []T, pred func(T) bool) []T   { return sliceimpl.SliceFilter(a, pred) }
func Map[T, R any](a []T, fn func(T) R) []R        { return sliceimpl.SliceMap(a, fn) }
func Sub[T any](a []T, fromIndex, toIndex int) []T { return sliceimpl.SliceSub(a, fromIndex, toIndex) }
func Concat[T any](slices ...[]T) []T              { return sliceimpl.SliceConcat(slices...) }
func Union[T comparable](a, b []T) []T             { return sliceimpl.Union(a, b) }
func Intersection[T comparable](a, b []T) []T      { return sliceimpl.Intersection(a, b) }
func Subtract[T comparable](a, b []T) []T          { return sliceimpl.Subtract(a, b) }
func Page[T any](a []T, pageNo, pageSize int) []T  { return sliceimpl.Page(a, pageNo, pageSize) }

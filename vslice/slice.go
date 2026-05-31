package vslice

import sliceimpl "github.com/imajinyun/go-knifer/internal/slice"

func IsEmpty[T any](a []T) bool                    { return sliceimpl.IsEmpty(a) }
func IsNotEmpty[T any](a []T) bool                 { return sliceimpl.IsNotEmpty(a) }
func Contains[T comparable](a []T, v T) bool       { return sliceimpl.Contains(a, v) }
func IndexOf[T comparable](a []T, v T) int         { return sliceimpl.IndexOf(a, v) }
func LastIndexOf[T comparable](a []T, v T) int     { return sliceimpl.LastIndexOf(a, v) }
func Reverse[T any](a []T) []T                     { return sliceimpl.Reverse(a) }
func Distinct[T comparable](a []T) []T             { return sliceimpl.Distinct(a) }
func Join[T any](a []T, sep string) string         { return sliceimpl.Join(a, sep) }
func Filter[T any](a []T, pred func(T) bool) []T   { return sliceimpl.Filter(a, pred) }
func Map[T, R any](a []T, fn func(T) R) []R        { return sliceimpl.Map(a, fn) }
func Sub[T any](a []T, fromIndex, toIndex int) []T { return sliceimpl.Sub(a, fromIndex, toIndex) }
func Concat[T any](slices ...[]T) []T              { return sliceimpl.Concat(slices...) }
func Union[T comparable](a, b []T) []T             { return sliceimpl.Union(a, b) }
func Intersection[T comparable](a, b []T) []T      { return sliceimpl.Intersection(a, b) }
func Subtract[T comparable](a, b []T) []T          { return sliceimpl.Subtract(a, b) }
func Page[T any](a []T, pageNo, pageSize int) []T  { return sliceimpl.Page(a, pageNo, pageSize) }

package vmap

import mapsimpl "github.com/imajinyun/go-knifer/internal/maps"

func IsEmpty[K comparable, V any](m map[K]V) bool        { return mapsimpl.IsEmpty(m) }
func IsNotEmpty[K comparable, V any](m map[K]V) bool     { return mapsimpl.IsNotEmpty(m) }
func Keys[K comparable, V any](m map[K]V) []K            { return mapsimpl.Keys(m) }
func Values[K comparable, V any](m map[K]V) []V          { return mapsimpl.Values(m) }
func Inverse[K, V comparable](m map[K]V) map[V]K         { return mapsimpl.Inverse(m) }
func Merge[K comparable, V any](maps ...map[K]V) map[K]V { return mapsimpl.Merge(maps...) }

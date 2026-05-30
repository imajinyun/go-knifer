package vconv

import convimpl "github.com/imajinyun/go-knifer/internal/conv"

func ToString(v any) string                       { return convimpl.ToString(v) }
func ToStringDefault(v any, def string) string    { return convimpl.ToStringDefault(v, def) }
func ToInt(v any) int                             { return convimpl.ToInt(v) }
func ToIntDefault(v any, def int) int             { return convimpl.ToIntDefault(v, def) }
func ToInt64(v any) int64                         { return convimpl.ToInt64(v) }
func ToInt64Default(v any, def int64) int64       { return convimpl.ToInt64Default(v, def) }
func ToFloat64(v any) float64                     { return convimpl.ToFloat64(v) }
func ToFloat64Default(v any, def float64) float64 { return convimpl.ToFloat64Default(v, def) }
func ToBool(v any) bool                           { return convimpl.ToBool(v) }
func ToBoolDefault(v any, def bool) bool          { return convimpl.ToBoolDefault(v, def) }
func ToBytes(v any) []byte                        { return convimpl.ToBytes(v) }

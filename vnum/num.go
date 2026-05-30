package vnum

import numimpl "github.com/imajinyun/go-knifer/internal/num"

type (
	Number  = numimpl.Number
	Ordered = numimpl.Ordered
)

func Add(a, b float64) float64                      { return numimpl.NumberAdd(a, b) }
func Sub(a, b float64) float64                      { return numimpl.NumberSub(a, b) }
func Mul(a, b float64) float64                      { return numimpl.NumberMul(a, b) }
func Div(a, b float64, scale int) float64           { return numimpl.NumberDiv(a, b, scale) }
func Round(v float64, scale int) float64            { return numimpl.Round(v, scale) }
func IsNumber(s string) bool                        { return numimpl.IsNumber(s) }
func IsInteger(s string) bool                       { return numimpl.IsInteger(s) }
func IsDigits(s string) bool                        { return numimpl.IsDigits(s) }
func Min[T Ordered](nums ...T) T                    { return numimpl.Min(nums...) }
func Max[T Ordered](nums ...T) T                    { return numimpl.Max(nums...) }
func Sum[T Number](nums ...T) T                     { return numimpl.Sum(nums...) }
func Avg[T Number](nums ...T) float64               { return numimpl.Avg(nums...) }
func Range(start, end, step int) []int              { return numimpl.Range(start, end, step) }
func Equals(a, b float64) bool                      { return numimpl.Equals(a, b) }
func DecimalFormat(format string, v float64) string { return numimpl.DecimalFormat(format, v) }

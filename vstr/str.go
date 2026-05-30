package vstr

import strimpl "github.com/imajinyun/go-knifer/internal/str"

func IsEmpty(s string) bool                       { return strimpl.IsEmpty(s) }
func IsNotEmpty(s string) bool                    { return strimpl.IsNotEmpty(s) }
func IsBlank(s string) bool                       { return strimpl.IsBlank(s) }
func IsNotBlank(s string) bool                    { return strimpl.IsNotBlank(s) }
func HasEmpty(strs ...string) bool                { return strimpl.HasEmpty(strs...) }
func HasBlank(strs ...string) bool                { return strimpl.HasBlank(strs...) }
func IsAllEmpty(strs ...string) bool              { return strimpl.IsAllEmpty(strs...) }
func IsAllBlank(strs ...string) bool              { return strimpl.IsAllBlank(strs...) }
func Trim(s string) string                        { return strimpl.Trim(s) }
func TrimToEmpty(s string) string                 { return strimpl.TrimToEmpty(s) }
func TrimStart(s string) string                   { return strimpl.TrimStart(s) }
func TrimEnd(s string) string                     { return strimpl.TrimEnd(s) }
func Sub(s string, fromIndex, toIndex int) string { return strimpl.Sub(s, fromIndex, toIndex) }
func SubBefore(s, sep string, isLastSeparator bool) string {
	return strimpl.SubBefore(s, sep, isLastSeparator)
}

func SubAfter(s, sep string, isLastSeparator bool) string {
	return strimpl.SubAfter(s, sep, isLastSeparator)
}
func Split(s, sep string) []string                   { return strimpl.Split(s, sep) }
func SplitTrim(s, sep string) []string               { return strimpl.SplitTrim(s, sep) }
func Repeat(s string, n int) string                  { return strimpl.Repeat(s, n) }
func PadLeft(s string, length int, pad rune) string  { return strimpl.PadLeft(s, length, pad) }
func PadRight(s string, length int, pad rune) string { return strimpl.PadRight(s, length, pad) }
func Contains(s, sub string) bool                    { return strimpl.Contains(s, sub) }
func ContainsAny(s string, subs ...string) bool      { return strimpl.ContainsAny(s, subs...) }
func ContainsAll(s string, subs ...string) bool      { return strimpl.ContainsAll(s, subs...) }
func ContainsIgnoreCase(s, sub string) bool          { return strimpl.ContainsIgnoreCase(s, sub) }
func StartsWith(s, prefix string) bool               { return strimpl.StartsWith(s, prefix) }
func EndsWith(s, suffix string) bool                 { return strimpl.EndsWith(s, suffix) }
func EqualsIgnoreCase(a, b string) bool              { return strimpl.EqualsIgnoreCase(a, b) }
func Reverse(s string) string                        { return strimpl.Reverse(s) }
func Format(template string, args ...any) string     { return strimpl.Format(template, args...) }
func RemovePrefix(s, prefix string) string           { return strimpl.RemovePrefix(s, prefix) }
func RemoveSuffix(s, suffix string) string           { return strimpl.RemoveSuffix(s, suffix) }
func AddPrefixIfNot(s, prefix string) string         { return strimpl.AddPrefixIfNot(s, prefix) }
func AddSuffixIfNot(s, suffix string) string         { return strimpl.AddSuffixIfNot(s, suffix) }
func Length(s string) int                            { return strimpl.Length(s) }
func DefaultIfNil[T any](v *T, def T) T              { return strimpl.DefaultIfNil(v, def) }
func DefaultIfEmpty(s, def string) string            { return strimpl.DefaultIfEmpty(s, def) }
func DefaultIfBlank(s, def string) string            { return strimpl.DefaultIfBlank(s, def) }
func EscapeHTML(s string) string                     { return strimpl.EscapeHTML(s) }
func UnescapeHTML(s string) string                   { return strimpl.UnescapeHTML(s) }
func ToCamelCase(s string) string                    { return strimpl.ToCamelCase(s) }
func ToPascalCase(s string) string                   { return strimpl.ToPascalCase(s) }
func ToUnderlineCase(s string) string                { return strimpl.ToUnderlineCase(s) }
func ToKebabCase(s string) string                    { return strimpl.ToKebabCase(s) }

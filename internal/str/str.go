package str

import (
	"fmt"
	"strings"
	"unicode"
)

// This file provides string helpers aligned with hutool-core StrUtil and CharSequenceUtil.

// IsEmpty reports whether s has zero length.
func IsEmpty(s string) bool { return len(s) == 0 }

// IsNotEmpty reports whether s is not empty.
func IsNotEmpty(s string) bool { return !IsEmpty(s) }

// IsBlank reports whether s is empty or contains only Unicode white space.
func IsBlank(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// IsNotBlank reports whether s is not blank.
func IsNotBlank(s string) bool { return !IsBlank(s) }

// HasEmpty reports whether any string is empty.
func HasEmpty(strs ...string) bool {
	for _, s := range strs {
		if IsEmpty(s) {
			return true
		}
	}
	return false
}

// HasBlank reports whether any string is blank.
func HasBlank(strs ...string) bool {
	for _, s := range strs {
		if IsBlank(s) {
			return true
		}
	}
	return false
}

// IsAllEmpty reports whether all strings are empty.
func IsAllEmpty(strs ...string) bool {
	for _, s := range strs {
		if IsNotEmpty(s) {
			return false
		}
	}
	return true
}

// IsAllBlank reports whether all strings are blank.
func IsAllBlank(strs ...string) bool {
	for _, s := range strs {
		if IsNotBlank(s) {
			return false
		}
	}
	return true
}

// Trim removes leading and trailing white space.
func Trim(s string) string { return strings.TrimSpace(s) }

// TrimToEmpty removes leading and trailing white space.
func TrimToEmpty(s string) string { return strings.TrimSpace(s) }

// TrimStart removes leading white space.
func TrimStart(s string) string { return strings.TrimLeftFunc(s, unicode.IsSpace) }

// TrimEnd removes trailing white space.
func TrimEnd(s string) string { return strings.TrimRightFunc(s, unicode.IsSpace) }

// Sub returns a substring by rune indexes and supports negative indexes from the end.
// fromIndex is inclusive and toIndex is exclusive; reversed ranges are normalized.
func Sub(s string, fromIndex, toIndex int) string {
	rs := []rune(s)
	n := len(rs)
	if n == 0 {
		return ""
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
	if fromIndex == toIndex {
		return ""
	}
	return string(rs[fromIndex:toIndex])
}

// SubBefore returns the text before sep. When isLastSeparator is true, the last sep is used.
func SubBefore(s, sep string, isLastSeparator bool) string {
	if s == "" || sep == "" {
		return s
	}
	var idx int
	if isLastSeparator {
		idx = strings.LastIndex(s, sep)
	} else {
		idx = strings.Index(s, sep)
	}
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// SubAfter returns the text after sep. When isLastSeparator is true, the last sep is used.
func SubAfter(s, sep string, isLastSeparator bool) string {
	if s == "" {
		return s
	}
	if sep == "" {
		return ""
	}
	var idx int
	if isLastSeparator {
		idx = strings.LastIndex(s, sep)
	} else {
		idx = strings.Index(s, sep)
	}
	if idx == -1 {
		return ""
	}
	return s[idx+len(sep):]
}

// Split splits s by sep and returns an empty slice for an empty input string.
func Split(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}

// SplitTrim splits s, trims each part, and drops blank parts.
func SplitTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// Repeat repeats s n times. Non-positive n returns an empty string.
func Repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(s, n)
}

// PadLeft pads s on the left to the requested rune length.
func PadLeft(s string, length int, pad rune) string {
	rs := []rune(s)
	if len(rs) >= length {
		return s
	}
	padding := make([]rune, length-len(rs))
	for i := range padding {
		padding[i] = pad
	}
	return string(padding) + s
}

// PadRight pads s on the right to the requested rune length.
func PadRight(s string, length int, pad rune) string {
	rs := []rune(s)
	if len(rs) >= length {
		return s
	}
	padding := make([]rune, length-len(rs))
	for i := range padding {
		padding[i] = pad
	}
	return s + string(padding)
}

// Contains reports whether s contains sub.
func Contains(s, sub string) bool { return strings.Contains(s, sub) }

// ContainsAny reports whether s contains any candidate substring.
func ContainsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// ContainsAll reports whether s contains all candidate substrings.
func ContainsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

// ContainsIgnoreCase reports whether s contains sub case-insensitively.
func ContainsIgnoreCase(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

// StartsWith reports whether s starts with prefix.
func StartsWith(s, prefix string) bool { return strings.HasPrefix(s, prefix) }

// EndsWith reports whether s ends with suffix.
func EndsWith(s, suffix string) bool { return strings.HasSuffix(s, suffix) }

// EqualsIgnoreCase compares strings case-insensitively.
func EqualsIgnoreCase(a, b string) bool { return strings.EqualFold(a, b) }

// Reverse reverses a string by rune, preserving multi-byte characters.
func Reverse(s string) string {
	rs := []rune(s)
	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		rs[i], rs[j] = rs[j], rs[i]
	}
	return string(rs)
}

// Format mimics hutool StrUtil.format by replacing {} placeholders in order.
//
//	Format("name={}, age={}", "tom", 12) -> "name=tom, age=12"
//
// Use \\{ to escape a literal opening brace.
func Format(template string, args ...any) string {
	if template == "" || len(args) == 0 {
		return template
	}
	var b strings.Builder
	b.Grow(len(template))
	idx := 0
	for i := 0; i < len(template); i++ {
		c := template[i]
		if c == '\\' && i+1 < len(template) && template[i+1] == '{' {
			b.WriteByte('{')
			i++
			continue
		}
		if c == '{' && i+1 < len(template) && template[i+1] == '}' {
			if idx < len(args) {
				fmt.Fprint(&b, args[idx])
				idx++
			} else {
				b.WriteString("{}")
			}
			i++
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}

// RemovePrefix removes prefix when present.
func RemovePrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

// RemoveSuffix removes suffix when present.
func RemoveSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

// AddPrefixIfNot adds prefix when it is not already present.
func AddPrefixIfNot(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s
	}
	return prefix + s
}

// AddSuffixIfNot adds suffix when it is not already present.
func AddSuffixIfNot(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s
	}
	return s + suffix
}

// Length returns the number of runes in s.
func Length(s string) int { return len([]rune(s)) }

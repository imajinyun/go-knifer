package base

import (
	"fmt"
	"strings"
	"unicode"
)

// 对应 hutool-core StrUtil / CharSequenceUtil。

// IsEmpty 字符串是否为空（长度为 0）。
func IsEmpty(s string) bool { return len(s) == 0 }

// IsNotEmpty 字符串是否非空。
func IsNotEmpty(s string) bool { return !IsEmpty(s) }

// IsBlank 字符串是否为空白（空字符串或仅含空白字符）。
func IsBlank(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// IsNotBlank 字符串是否非空白。
func IsNotBlank(s string) bool { return !IsBlank(s) }

// HasEmpty 多个字符串中是否有空字符串。
func HasEmpty(strs ...string) bool {
	for _, s := range strs {
		if IsEmpty(s) {
			return true
		}
	}
	return false
}

// HasBlank 多个字符串中是否有空白字符串。
func HasBlank(strs ...string) bool {
	for _, s := range strs {
		if IsBlank(s) {
			return true
		}
	}
	return false
}

// IsAllEmpty 多个字符串是否都为空。
func IsAllEmpty(strs ...string) bool {
	for _, s := range strs {
		if IsNotEmpty(s) {
			return false
		}
	}
	return true
}

// IsAllBlank 多个字符串是否都为空白。
func IsAllBlank(strs ...string) bool {
	for _, s := range strs {
		if IsNotBlank(s) {
			return false
		}
	}
	return true
}

// Trim 去除首尾空白。
func Trim(s string) string { return strings.TrimSpace(s) }

// TrimToEmpty 去除首尾空白，nil 视作空串。
func TrimToEmpty(s string) string { return strings.TrimSpace(s) }

// TrimStart 去除开头空白。
func TrimStart(s string) string { return strings.TrimLeftFunc(s, unicode.IsSpace) }

// TrimEnd 去除末尾空白。
func TrimEnd(s string) string { return strings.TrimRightFunc(s, unicode.IsSpace) }

// Sub 截取字符串，按 rune 计数，支持负数索引（从尾部开始）。
// fromIndex 包含，toIndex 不包含。
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

// SubBefore 截取分隔符 before 之前的字符串；若 isLastSeparator=true，按最后一次出现位置截取。
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

// SubAfter 截取分隔符 after 之后的字符串。
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

// Split 切分字符串。
func Split(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}

// SplitTrim 切分后修剪每段并丢弃空白段。
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

// Repeat 重复字符串 n 次。
func Repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(s, n)
}

// PadLeft 左侧填充至指定长度（按 rune 计算）。
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

// PadRight 右侧填充至指定长度。
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

// Contains 是否包含子串。
func Contains(s, sub string) bool { return strings.Contains(s, sub) }

// ContainsAny 是否包含任意一个子串。
func ContainsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// ContainsAll 是否包含全部子串。
func ContainsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

// ContainsIgnoreCase 忽略大小写包含。
func ContainsIgnoreCase(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

// StartsWith 是否以 prefix 开始。
func StartsWith(s, prefix string) bool { return strings.HasPrefix(s, prefix) }

// EndsWith 是否以 suffix 结束。
func EndsWith(s, suffix string) bool { return strings.HasSuffix(s, suffix) }

// EqualsIgnoreCase 忽略大小写比较。
func EqualsIgnoreCase(a, b string) bool { return strings.EqualFold(a, b) }

// Reverse 字符串反转（按 rune）。
func Reverse(s string) string {
	rs := []rune(s)
	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		rs[i], rs[j] = rs[j], rs[i]
	}
	return string(rs)
}

// Format 仿 hutool StrUtil.format：使用 {} 占位符，按顺序替换。
//
//	Format("name={}, age={}", "tom", 12) -> "name=tom, age=12"
//
// 使用 \\{ 转义。
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

// RemovePrefix 移除前缀。
func RemovePrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

// RemoveSuffix 移除后缀。
func RemoveSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

// AddPrefixIfNot 不存在前缀则添加。
func AddPrefixIfNot(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s
	}
	return prefix + s
}

// AddSuffixIfNot 不存在后缀则添加。
func AddSuffixIfNot(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s
	}
	return s + suffix
}

// Length 按 rune 计算字符串长度。
func Length(s string) int { return len([]rune(s)) }

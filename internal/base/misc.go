package base

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash/fnv"
	"regexp"
	"strings"
	"unicode"
)

// 对应 hutool-core CharUtil。

// IsBlankChar 是否为空白字符（包括 \u00A0 等）。
func IsBlankChar(r rune) bool {
	return unicode.IsSpace(r) || r == '\u00A0' || r == '\u2007' || r == '\u202F' || r == '\uFEFF'
}

// IsLetter 是否为字母。
func IsLetter(r rune) bool { return unicode.IsLetter(r) }

// IsDigit 是否为数字。
func IsDigit(r rune) bool { return unicode.IsDigit(r) }

// IsAscii 是否为 ASCII。
func IsAscii(r rune) bool { return r < 128 }

// IsLetterOrDigit 是否为字母或数字。
func IsLetterOrDigit(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) }

// 对应 hutool-core BooleanUtil。

// BoolNegate 取反。
func BoolNegate(b bool) bool { return !b }

// BoolToInt true=1, false=0。
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// BoolAnd 全真为真。
func BoolAnd(bs ...bool) bool {
	for _, b := range bs {
		if !b {
			return false
		}
	}
	return true
}

// BoolOr 任一为真即真。
func BoolOr(bs ...bool) bool {
	for _, b := range bs {
		if b {
			return true
		}
	}
	return false
}

// 对应 hutool-core HashUtil（轻量子集）。

// AdditiveHash 累加 hash。
func AdditiveHash(s string, prime int) int {
	if prime <= 0 {
		prime = 31
	}
	h := len(s)
	for _, r := range s {
		h += int(r)
	}
	return h % prime
}

// FnvHash FNV-1 32 位哈希。
func FnvHash(s string) uint32 {
	h := fnv.New32()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

// MD5Hex 计算 MD5（hex）。
func MD5Hex(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

// SHA1Hex 计算 SHA1（hex）。
func SHA1Hex(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

// SHA256Hex 计算 SHA256（hex）。
func SHA256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// 对应 hutool-core ReUtil / Validator。

// 常用正则。
var (
	rxEmail   = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)
	rxMobile  = regexp.MustCompile(`^1[3-9]\d{9}$`)
	rxURL     = regexp.MustCompile(`(?i)^https?://[^\s]+$`)
	rxIPv4    = regexp.MustCompile(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`)
	rxChinese = regexp.MustCompile(`^[\p{Han}]+$`)
	rxNumber  = regexp.MustCompile(`^-?\d+(\.\d+)?$`)
)

// IsEmail 是否为邮箱。
func IsEmail(s string) bool { return rxEmail.MatchString(s) }

// IsMobile 是否为中国大陆手机号。
func IsMobile(s string) bool { return rxMobile.MatchString(s) }

// IsURL 是否为 http/https URL。
func IsURL(s string) bool { return rxURL.MatchString(s) }

// IsIPv4 是否为 IPv4 地址。
func IsIPv4(s string) bool { return rxIPv4.MatchString(s) }

// IsChinese 是否全为中文字符。
func IsChinese(s string) bool { return s != "" && rxChinese.MatchString(s) }

// IsNumberStr 是否为数字字符串（支持小数和负号）。
func IsNumberStr(s string) bool { return rxNumber.MatchString(s) }

// ReMatch 是否匹配正则。
func ReMatch(pattern, s string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(s)
}

// ReFind 返回首个匹配，未命中返回空。
func ReFind(pattern, s string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}
	return re.FindString(s)
}

// ReFindAll 返回全部匹配。
func ReFindAll(pattern, s string) []string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	return re.FindAllString(s, -1)
}

// ReReplace 正则替换。
func ReReplace(pattern, s, replacement string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return s
	}
	return re.ReplaceAllString(s, replacement)
}

// 对应 hutool-core ObjUtil 部分常用方法。

// DefaultIfNil 为 nil 时返回默认值。
func DefaultIfNil[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}

// DefaultIfEmpty 字符串为空时返回默认值。
func DefaultIfEmpty(s, def string) string {
	if IsEmpty(s) {
		return def
	}
	return s
}

// DefaultIfBlank 字符串为空白时返回默认值。
func DefaultIfBlank(s, def string) string {
	if IsBlank(s) {
		return def
	}
	return s
}

// 对应 hutool-core EscapeUtil 简化版。

// HtmlEscape 等价 httpx.HTMLEscape，但 base 不依赖 httpx，提供独立简版。
func EscapeHTML(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	)
	return r.Replace(s)
}

// UnescapeHTML 解 HTML 实体（仅常用）。
func UnescapeHTML(s string) string {
	r := strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", "\"",
		"&#39;", "'",
		"&apos;", "'",
		"&nbsp;", "\u00A0",
	)
	return r.Replace(s)
}

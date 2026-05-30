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

// This section provides character helpers aligned with hutool-core CharUtil.

// IsBlankChar reports whether r is a blank character, including non-breaking spaces.
func IsBlankChar(r rune) bool {
	return unicode.IsSpace(r) || r == '\u00A0' || r == '\u2007' || r == '\u202F' || r == '\uFEFF'
}

// IsLetter reports whether r is a Unicode letter.
func IsLetter(r rune) bool { return unicode.IsLetter(r) }

// IsDigit reports whether r is a Unicode digit.
func IsDigit(r rune) bool { return unicode.IsDigit(r) }

// IsAscii reports whether r is an ASCII character.
func IsAscii(r rune) bool { return r < 128 }

// IsLetterOrDigit reports whether r is a Unicode letter or digit.
func IsLetterOrDigit(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) }

// This section provides boolean helpers aligned with hutool-core BooleanUtil.

// BoolNegate returns the logical negation of b.
func BoolNegate(b bool) bool { return !b }

// BoolToInt returns 1 for true and 0 for false.
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// BoolAnd returns true only when all inputs are true.
func BoolAnd(bs ...bool) bool {
	for _, b := range bs {
		if !b {
			return false
		}
	}
	return true
}

// BoolOr returns true when any input is true.
func BoolOr(bs ...bool) bool {
	for _, b := range bs {
		if b {
			return true
		}
	}
	return false
}

// This section provides a lightweight subset of hutool-core HashUtil.

// AdditiveHash calculates an additive hash modulo prime. Non-positive prime falls back to 31.
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

// FnvHash calculates a 32-bit FNV-1 hash.
func FnvHash(s string) uint32 {
	h := fnv.New32()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

// MD5Hex calculates the MD5 digest and returns lowercase hex text.
func MD5Hex(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

// SHA1Hex calculates the SHA-1 digest and returns lowercase hex text.
func SHA1Hex(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

// SHA256Hex calculates the SHA-256 digest and returns lowercase hex text.
func SHA256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// This section provides regular-expression and validation helpers aligned with hutool-core ReUtil and Validator.

// Common regular expressions used by validators.
var (
	rxEmail   = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)
	rxMobile  = regexp.MustCompile(`^1[3-9]\d{9}$`)
	rxURL     = regexp.MustCompile(`(?i)^https?://[^\s]+$`)
	rxIPv4    = regexp.MustCompile(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`)
	rxChinese = regexp.MustCompile(`^[\p{Han}]+$`)
	rxNumber  = regexp.MustCompile(`^-?\d+(\.\d+)?$`)
)

// IsEmail reports whether s is an email address.
func IsEmail(s string) bool { return rxEmail.MatchString(s) }

// IsMobile reports whether s is a mainland China mobile phone number.
func IsMobile(s string) bool { return rxMobile.MatchString(s) }

// IsURL reports whether s is an http or https URL.
func IsURL(s string) bool { return rxURL.MatchString(s) }

// IsIPv4 reports whether s is an IPv4 address.
func IsIPv4(s string) bool { return rxIPv4.MatchString(s) }

// IsChinese reports whether s consists only of Chinese Han characters.
func IsChinese(s string) bool { return s != "" && rxChinese.MatchString(s) }

// IsNumberStr reports whether s is a number string, including decimals and a leading minus sign.
func IsNumberStr(s string) bool { return rxNumber.MatchString(s) }

// ReMatch reports whether s matches pattern. Invalid patterns return false.
func ReMatch(pattern, s string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(s)
}

// ReFind returns the first match, or an empty string when there is no match or the pattern is invalid.
func ReFind(pattern, s string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}
	return re.FindString(s)
}

// ReFindAll returns all matches, or nil when the pattern is invalid.
func ReFindAll(pattern, s string) []string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	return re.FindAllString(s, -1)
}

// ReReplace replaces matches of pattern with replacement. Invalid patterns return the original string.
func ReReplace(pattern, s, replacement string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return s
	}
	return re.ReplaceAllString(s, replacement)
}

// This section provides common object helpers aligned with hutool-core ObjUtil.

// DefaultIfNil returns def when v is nil; otherwise it returns *v.
func DefaultIfNil[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}

// DefaultIfEmpty returns def when s is empty.
func DefaultIfEmpty(s, def string) string {
	if IsEmpty(s) {
		return def
	}
	return s
}

// DefaultIfBlank returns def when s is blank.
func DefaultIfBlank(s, def string) string {
	if IsBlank(s) {
		return def
	}
	return s
}

// This section provides a simplified subset of hutool-core EscapeUtil.

// EscapeHTML escapes common HTML-sensitive characters without depending on the HTTP package.
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

// UnescapeHTML unescapes common HTML entities.
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

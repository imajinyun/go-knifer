package net

import (
	"net/url"
	"strings"
)

// Decode unescapes percent-encoded text and converts plus signs to spaces.
func Decode(s string) (string, error) { return DecodePlus(s, true) }

// DecodeForPath unescapes percent-encoded path text without converting plus signs to spaces.
func DecodeForPath(s string) (string, error) { return DecodePlus(s, false) }

// DecodePlus unescapes percent-encoded text and controls whether plus signs become spaces.
func DecodePlus(s string, plusToSpace bool) (string, error) {
	if plusToSpace {
		return url.QueryUnescape(s)
	}
	return url.PathUnescape(strings.ReplaceAll(s, "+", "%2B"))
}

// EncodeAll percent-encodes every non-unreserved character.
func EncodeAll(s string) string { return encodeWith(s, isUnreserved, false) }

// Encode escapes text for URI path usage while keeping slash separators.
func Encode(s string) string { return EncodePath(s) }

// EncodeQuery escapes text for query/form usage. Spaces are encoded as '+'.
func EncodeQuery(s string) string { return url.QueryEscape(s) }

// EncodePathSegment escapes one path segment, including slash characters.
func EncodePathSegment(s string) string { return url.PathEscape(s) }

// EncodePath escapes each path segment and keeps slash separators.
func EncodePath(s string) string {
	parts := strings.Split(s, "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}

// EncodeFragment escapes URL fragment text.
func EncodeFragment(s string) string { return encodeWith(s, isFragmentSafe, false) }

// FormURLEncode escapes text for application/x-www-form-urlencoded usage.
func FormURLEncode(s string) string { return url.QueryEscape(s) }

func encodeWith(s string, safe func(byte) bool, spaceAsPlus bool) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' && spaceAsPlus {
			b.WriteByte('+')
			continue
		}
		if safe(c) {
			b.WriteByte(c)
			continue
		}
		b.WriteByte('%')
		const hex = "0123456789ABCDEF"
		b.WriteByte(hex[c>>4])
		b.WriteByte(hex[c&0x0f])
	}
	return b.String()
}

func isUnreserved(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '.' || c == '_' || c == '~'
}

func isFragmentSafe(c byte) bool {
	return isUnreserved(c) || strings.ContainsRune("!$&'()*+,;=:@/?", rune(c))
}

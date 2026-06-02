package url

import (
	neturl "net/url"
	"strings"
)

// DecodeForPath unescapes percent-encoded path text without converting plus signs to spaces.
func DecodeForPath(s string) (string, error) { return DecodePlus(s, false) }

// EncodeAll percent-encodes every non-unreserved character.
func EncodeAll(s string) string { return encodeWith(s, isUnreserved, false) }

// EncodeQuery escapes text for query/form usage. Spaces are encoded as '+'.
func EncodeQuery(s string) string { return neturl.QueryEscape(s) }

// EncodePathSegment escapes one path segment, including slash characters.
func EncodePathSegment(s string) string { return neturl.PathEscape(s) }

// EncodePath escapes each path segment and keeps slash separators.
func EncodePath(s string) string { return encodePathKeepSlash(s) }

// EncodeFragment escapes URL fragment text.
func EncodeFragment(s string) string { return encodeWith(s, isFragmentSafe, false) }

// FormURLEncode escapes text for application/x-www-form-urlencoded usage.
func FormURLEncode(s string) string { return neturl.QueryEscape(s) }

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

package str

import "strings"

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

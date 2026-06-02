package http

import (
	"html"
	"regexp"
	"strings"
)

// HTMLEscape escapes HTML, aligned with HtmlUtil.escape.
func HTMLEscape(s string) string { return html.EscapeString(s) }

// HTMLUnescape unescapes HTML, aligned with HtmlUtil.unescape.
func HTMLUnescape(s string) string { return html.UnescapeString(s) }

var (
	tagRegex     = regexp.MustCompile(`(?is)<[^>]+>`)
	commentRegex = regexp.MustCompile(`(?is)<!--.*?-->`)
)

// CleanHTML removes HTML tags and keeps plain text only, aligned with HtmlUtil.cleanHtmlTag.
func CleanHTML(s string) string {
	s = commentRegex.ReplaceAllString(s, "")
	s = tagRegex.ReplaceAllString(s, "")
	return strings.TrimSpace(s)
}

// FilterHTMLTag removes the specified HTML tags.
func FilterHTMLTag(s string, tagNames ...string) string {
	for _, tag := range tagNames {
		t := regexp.QuoteMeta(tag)
		// Tags with content: <tag ...>...</tag>.
		re := regexp.MustCompile(`(?is)<` + t + `(\s[^>]*)?>.*?</` + t + `\s*>`)
		s = re.ReplaceAllString(s, "")
		// Self-closing or single tags: <tag ... /> or <tag>.
		re2 := regexp.MustCompile(`(?is)<` + t + `(\s[^>]*)?/?>`)
		s = re2.ReplaceAllString(s, "")
	}
	return s
}

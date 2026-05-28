package http

import (
	"html"
	"regexp"
	"strings"
)

// HTMLEscape HTML 转义（对应 HtmlUtil.escape）。
func HTMLEscape(s string) string { return html.EscapeString(s) }

// HTMLUnescape HTML 反转义（对应 HtmlUtil.unescape）。
func HTMLUnescape(s string) string { return html.UnescapeString(s) }

var (
	tagRegex     = regexp.MustCompile(`(?is)<[^>]+>`)
	commentRegex = regexp.MustCompile(`(?is)<!--.*?-->`)
)

// CleanHTML 清理 HTML 标签，仅保留纯文本（对应 HtmlUtil.cleanHtmlTag）。
func CleanHTML(s string) string {
	s = commentRegex.ReplaceAllString(s, "")
	s = tagRegex.ReplaceAllString(s, "")
	return strings.TrimSpace(s)
}

// FilterHTMLTag 过滤指定 HTML 标签。
func FilterHTMLTag(s string, tagNames ...string) string {
	for _, tag := range tagNames {
		t := regexp.QuoteMeta(tag)
		// 包含内容标签：<tag ...>...</tag>
		re := regexp.MustCompile(`(?is)<` + t + `(\s[^>]*)?>.*?</` + t + `\s*>`)
		s = re.ReplaceAllString(s, "")
		// 自闭合或单标签：<tag ... /> 或 <tag>
		re2 := regexp.MustCompile(`(?is)<` + t + `(\s[^>]*)?/?>`)
		s = re2.ReplaceAllString(s, "")
	}
	return s
}

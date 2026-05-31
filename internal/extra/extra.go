package extra

import (
	"bytes"
	"html/template"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"

	urlimpl "github.com/imajinyun/go-knifer/internal/url"
	zipimpl "github.com/imajinyun/go-knifer/internal/zip"
)

var emojiPattern = regexp.MustCompile(`(?:[\x{1F1E6}-\x{1F1FF}]{2}|[#*0-9]\x{FE0F}?\x{20E3}|[\x{1F300}-\x{1FAFF}\x{2600}-\x{27BF}])(?:\x{FE0F}|\x{200D}[\x{1F300}-\x{1FAFF}\x{2600}-\x{27BF}]\x{FE0F}?)*`)

// Gzip compresses data using gzip.
func Gzip(data []byte) ([]byte, error) { return zipimpl.Gzip(data) }

// Gunzip decompresses gzip data.
func Gunzip(data []byte) ([]byte, error) { return zipimpl.Gunzip(data) }

// Zlib compresses data using zlib.
func Zlib(data []byte) ([]byte, error) { return zipimpl.Zlib(data) }

// Unzlib decompresses zlib data.
func Unzlib(data []byte) ([]byte, error) { return zipimpl.Unzlib(data) }

// ZipFiles creates a zip archive from files and directories.
//
// Directory entries are preserved so empty directories survive a round trip, and
// file metadata such as permissions is copied into the zip headers.
func ZipFiles(dest string, files ...string) (err error) {
	return zipimpl.ZipFiles(dest, true, files...)
}

// Unzip extracts a zip archive into destDir.
//
// Archive paths are validated before writing to prevent zip-slip path traversal.
func Unzip(src, destDir string) error { return zipimpl.UnzipTo(src, destDir) }

// ContainsEmoji reports whether s contains emoji-like runes.
func ContainsEmoji(s string) bool { return emojiPattern.MatchString(s) }

// RemoveEmoji removes emoji-like runes from s, including variation-selector and
// zero-width-joiner based emoji sequences.
func RemoveEmoji(s string) string { return emojiPattern.ReplaceAllString(s, "") }

// RenderTemplate renders a Go html/template string with data.
func RenderTemplate(tpl string, data any) (string, error) {
	t, err := template.New("go-knifer-extra").Parse(tpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// IsEmail reports whether s is a syntactically valid email address.
//
// Display-name forms such as "Alice <alice@example.com>" are rejected because
// this helper validates an address string rather than a mail header value.
func IsEmail(s string) bool {
	s = strings.TrimSpace(s)
	addr, err := mail.ParseAddress(s)
	return err == nil && addr.Name == "" && addr.Address == s
}

// IsURL reports whether s is an absolute URL with scheme and host.
func IsURL(s string) bool { return urlimpl.IsAbsoluteURL(s) }

// IsBlank reports whether s is empty or whitespace.
func IsBlank(s string) bool { return strings.TrimSpace(s) == "" }

// RuneLen returns the UTF-8 rune count.
func RuneLen(s string) int { return utf8.RuneCountInString(s) }

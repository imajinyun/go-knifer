package vextra

import extraimpl "github.com/imajinyun/go-knifer/internal/extra"

// Gzip compresses data using gzip.
func Gzip(data []byte) ([]byte, error) { return extraimpl.Gzip(data) }

// Gunzip decompresses gzip data.
func Gunzip(data []byte) ([]byte, error) { return extraimpl.Gunzip(data) }

// Zlib compresses data using zlib.
func Zlib(data []byte) ([]byte, error) { return extraimpl.Zlib(data) }

// Unzlib decompresses zlib data.
func Unzlib(data []byte) ([]byte, error) { return extraimpl.Unzlib(data) }

// ZipFiles creates a zip archive from files and directories.
func ZipFiles(dest string, files ...string) error { return extraimpl.ZipFiles(dest, files...) }

// Unzip extracts a zip archive into destDir.
func Unzip(src, destDir string) error { return extraimpl.Unzip(src, destDir) }

// ContainsEmoji reports whether s contains emoji-like runes.
func ContainsEmoji(s string) bool { return extraimpl.ContainsEmoji(s) }

// RemoveEmoji removes emoji-like runes from s.
func RemoveEmoji(s string) string { return extraimpl.RemoveEmoji(s) }

// RenderTemplate renders a Go html/template string with data.
func RenderTemplate(tpl string, data any) (string, error) { return extraimpl.RenderTemplate(tpl, data) }

// IsEmail reports whether s is a syntactically valid email address.
func IsEmail(s string) bool { return extraimpl.IsEmail(s) }

// IsURL reports whether s is an absolute URL with scheme and host.
func IsURL(s string) bool { return extraimpl.IsURL(s) }

// IsBlank reports whether s is empty or whitespace.
func IsBlank(s string) bool { return extraimpl.IsBlank(s) }

// RuneLen returns the UTF-8 rune count.
func RuneLen(s string) int { return extraimpl.RuneLen(s) }

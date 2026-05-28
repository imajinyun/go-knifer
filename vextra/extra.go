package vextra

import extraimpl "github.com/imajinyun/go-knifer/internal/extra"

// Gzip 使用 gzip 压缩数据，对应 Hutool CompressUtil 的 gzip 能力。Gzip compresses data using gzip.
func Gzip(data []byte) ([]byte, error) { return extraimpl.Gzip(data) }

// Gunzip 解压 gzip 数据。Gunzip decompresses gzip data.
func Gunzip(data []byte) ([]byte, error) { return extraimpl.Gunzip(data) }

// Zlib 使用 zlib 压缩数据。Zlib compresses data using zlib.
func Zlib(data []byte) ([]byte, error) { return extraimpl.Zlib(data) }

// Unzlib 解压 zlib 数据。Unzlib decompresses zlib data.
func Unzlib(data []byte) ([]byte, error) { return extraimpl.Unzlib(data) }

// ZipFiles 将文件或目录打包为 zip 归档。ZipFiles creates a zip archive from files.
func ZipFiles(dest string, files ...string) error { return extraimpl.ZipFiles(dest, files...) }

// Unzip 将 zip 归档解压到目标目录。Unzip extracts a zip archive into destDir.
func Unzip(src, destDir string) error { return extraimpl.Unzip(src, destDir) }

// ContainsEmoji 判断字符串是否包含 Emoji 字符，对应 Hutool EmojiUtil.containsEmoji。ContainsEmoji reports whether s contains emoji-like runes.
func ContainsEmoji(s string) bool { return extraimpl.ContainsEmoji(s) }

// RemoveEmoji 移除字符串中的 Emoji 字符。RemoveEmoji removes emoji-like runes from s.
func RemoveEmoji(s string) string { return extraimpl.RemoveEmoji(s) }

// RenderTemplate 使用 Go html/template 渲染模板字符串，对应 Hutool TemplateUtil 的基础模板能力。RenderTemplate renders a Go html/template string with data.
func RenderTemplate(tpl string, data any) (string, error) { return extraimpl.RenderTemplate(tpl, data) }

// IsEmail 判断字符串是否为合法邮箱地址，对应 Hutool ValidationUtil 的常见校验。IsEmail reports whether s is a syntactically valid email address.
func IsEmail(s string) bool { return extraimpl.IsEmail(s) }

// IsURL 判断字符串是否为包含 scheme 和 host 的绝对 URL。IsURL reports whether s is an absolute URL with scheme and host.
func IsURL(s string) bool { return extraimpl.IsURL(s) }

// IsBlank 判断字符串是否为空白。IsBlank reports whether s is empty or whitespace.
func IsBlank(s string) bool { return extraimpl.IsBlank(s) }

// RuneLen 返回 UTF-8 字符数量。RuneLen returns UTF-8 rune count.
func RuneLen(s string) int { return extraimpl.RuneLen(s) }

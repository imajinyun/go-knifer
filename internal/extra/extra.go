package extra

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"html/template"
	"io"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

var emojiPattern = regexp.MustCompile(`[\x{1F300}-\x{1FAFF}\x{2600}-\x{27BF}]`)

// Gzip 使用 gzip 压缩数据，对应 Hutool CompressUtil 的 gzip 能力。Gzip compresses data using gzip.
func Gzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		_ = w.Close()
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Gunzip 解压 gzip 数据。Gunzip decompresses gzip data.
func Gunzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

// Zlib 使用 zlib 压缩数据。Zlib compresses data using zlib.
func Zlib(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		_ = w.Close()
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unzlib 解压 zlib 数据。Unzlib decompresses zlib data.
func Unzlib(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

// ZipFiles 将文件或目录打包为 zip 归档。ZipFiles creates a zip archive from files.
func ZipFiles(dest string, files ...string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	zw := zip.NewWriter(out)
	defer zw.Close()
	for _, file := range files {
		if err := addFileToZip(zw, file, filepath.Base(file)); err != nil {
			return err
		}
	}
	return nil
}

// Unzip 将 zip 归档解压到目标目录。Unzip extracts a zip archive into destDir.
func Unzip(src, destDir string) error {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, f := range zr.File {
		if err := extractZipFile(f, destDir); err != nil {
			return err
		}
	}
	return nil
}

// ContainsEmoji 判断字符串是否包含 Emoji 字符，对应 Hutool EmojiUtil.containsEmoji。ContainsEmoji reports whether s contains emoji-like runes.
func ContainsEmoji(s string) bool { return emojiPattern.MatchString(s) }

// RemoveEmoji 移除字符串中的 Emoji 字符。RemoveEmoji removes emoji-like runes from s.
func RemoveEmoji(s string) string { return emojiPattern.ReplaceAllString(s, "") }

// RenderTemplate 使用 Go html/template 渲染模板字符串，对应 Hutool TemplateUtil 的基础模板能力。RenderTemplate renders a Go html/template string with data.
func RenderTemplate(tpl string, data any) (string, error) {
	t, err := template.New("hutool-extra").Parse(tpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// IsEmail 判断字符串是否为合法邮箱地址，对应 Hutool ValidationUtil 的常见校验。IsEmail reports whether s is a syntactically valid email address.
func IsEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

// IsURL 判断字符串是否为包含 scheme 和 host 的绝对 URL。IsURL reports whether s is an absolute URL with scheme and host.
func IsURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

// IsBlank 判断字符串是否为空白。IsBlank reports whether s is empty or whitespace.
func IsBlank(s string) bool { return strings.TrimSpace(s) == "" }

// RuneLen 返回 UTF-8 字符数量。RuneLen returns UTF-8 rune count.
func RuneLen(s string) int { return utf8.RuneCountInString(s) }

func addFileToZip(zw *zip.Writer, path, name string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := addFileToZip(zw, filepath.Join(path, entry.Name()), filepath.Join(name, entry.Name())); err != nil {
				return err
			}
		}
		return nil
	}
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := zw.Create(filepath.ToSlash(name))
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r)
	return err
}

func extractZipFile(f *zip.File, destDir string) error {
	target := filepath.Join(destDir, f.Name)
	cleanDest := filepath.Clean(destDir) + string(os.PathSeparator)
	cleanTarget := filepath.Clean(target)
	if !strings.HasPrefix(cleanTarget, cleanDest) && cleanTarget != filepath.Clean(destDir) {
		return os.ErrPermission
	}
	if f.FileInfo().IsDir() {
		return os.MkdirAll(target, f.Mode())
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	r, err := f.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, r)
	return err
}

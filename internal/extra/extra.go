package extra

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
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

var emojiPattern = regexp.MustCompile(`(?:[\x{1F1E6}-\x{1F1FF}]{2}|[#*0-9]\x{FE0F}?\x{20E3}|[\x{1F300}-\x{1FAFF}\x{2600}-\x{27BF}])(?:\x{FE0F}|\x{200D}[\x{1F300}-\x{1FAFF}\x{2600}-\x{27BF}]\x{FE0F}?)*`)

// Gzip compresses data using gzip.
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

// Gunzip decompresses gzip data.
func Gunzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()
	return io.ReadAll(r)
}

// Zlib compresses data using zlib.
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

// Unzlib decompresses zlib data.
func Unzlib(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()
	return io.ReadAll(r)
}

// ZipFiles creates a zip archive from files and directories.
//
// Directory entries are preserved so empty directories survive a round trip, and
// file metadata such as permissions is copied into the zip headers.
func ZipFiles(dest string, files ...string) (err error) {
	if dir := filepath.Dir(dest); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := out.Close(); err == nil {
			err = closeErr
		}
	}()
	zw := zip.NewWriter(out)
	defer func() {
		if closeErr := zw.Close(); err == nil {
			err = closeErr
		}
	}()
	for _, file := range files {
		if err := addFileToZip(zw, file, filepath.Base(file)); err != nil {
			return err
		}
	}
	return nil
}

// Unzip extracts a zip archive into destDir.
//
// Archive paths are validated before writing to prevent zip-slip path traversal.
func Unzip(src, destDir string) error {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() { _ = zr.Close() }()
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}
	for _, f := range zr.File {
		if err := extractZipFile(f, destDir); err != nil {
			return err
		}
	}
	return nil
}

// ContainsEmoji reports whether s contains emoji-like runes.
func ContainsEmoji(s string) bool { return emojiPattern.MatchString(s) }

// RemoveEmoji removes emoji-like runes from s, including variation-selector and
// zero-width-joiner based emoji sequences.
func RemoveEmoji(s string) string { return emojiPattern.ReplaceAllString(s, "") }

// RenderTemplate renders a Go html/template string with data.
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
func IsURL(s string) bool {
	if s == "" || strings.TrimSpace(s) != s {
		return false
	}
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	return u.IsAbs() && u.Host != ""
}

// IsBlank reports whether s is empty or whitespace.
func IsBlank(s string) bool { return strings.TrimSpace(s) == "" }

// RuneLen returns the UTF-8 rune count.
func RuneLen(s string) int { return utf8.RuneCountInString(s) }

func addFileToZip(zw *zip.Writer, path, name string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	zipName := filepath.ToSlash(filepath.Clean(name))
	if zipName == "." || strings.HasPrefix(zipName, "../") || strings.HasPrefix(zipName, "/") {
		return fmt.Errorf("invalid zip entry name %q", name)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = zipName
	header.Method = zip.Deflate
	header.SetMode(info.Mode())

	if info.IsDir() {
		header.Name += "/"
		if _, err := zw.CreateHeader(header); err != nil {
			return err
		}
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

	w, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		linkTarget, err := os.Readlink(path)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, linkTarget)
		return err
	}

	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()
	_, err = io.Copy(w, r)
	return err
}

func extractZipFile(f *zip.File, destDir string) error {
	target, err := safeZipTarget(destDir, f.Name)
	if err != nil {
		return err
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
	defer func() { _ = r.Close() }()
	w, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, r); err != nil {
		_ = w.Close()
		return err
	}
	return w.Close()
}

func safeZipTarget(destDir, name string) (string, error) {
	if name == "" || filepath.IsAbs(name) {
		return "", fmt.Errorf("illegal file path in zip: %q", name)
	}
	target := filepath.Join(destDir, filepath.FromSlash(name))
	destAbs, err := filepath.Abs(destDir)
	if err != nil {
		return "", err
	}
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(destAbs, targetAbs)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("illegal file path in zip: %q", name)
	}
	return target, nil
}

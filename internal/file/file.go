package file

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// This section provides IO helpers aligned with hutool-core IoUtil.

// ReadAll reads all data from r.
func ReadAll(r io.Reader) ([]byte, error) { return io.ReadAll(r) }

// ReadString reads all data from r and returns it as a string.
func ReadString(r io.Reader) (string, error) {
	b, err := ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadLines reads all lines from r. The scanner buffer is enlarged to support lines up to 1 MiB.
func ReadLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// IoCopy copies from src to dst and returns the number of bytes written.
func IoCopy(dst io.Writer, src io.Reader) (int64, error) { return io.Copy(dst, src) }

// CloseQuietly closes c and ignores the returned error.
func CloseQuietly(c io.Closer) {
	if c == nil {
		return
	}
	_ = c.Close()
}

// This section provides file and filename helpers aligned with hutool-core FileUtil and FileNameUtil.

// FileExists reports whether a file or directory exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsFile reports whether path exists and is a regular file.
func IsFile(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

// IsDirectory reports whether path exists and is a directory.
func IsDirectory(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}

// FileReadString reads the whole file as a string.
func FileReadString(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FileReadBytes reads all bytes from a file.
func FileReadBytes(path string) ([]byte, error) { return os.ReadFile(path) }

// FileReadLines reads all lines from a file.
func FileReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer CloseQuietly(f)
	return ReadLines(f)
}

// FileWriteString writes content to a file, overwriting existing data and creating parent directories.
func FileWriteString(path, content string) error {
	if err := Mkdir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// FileWriteBytes writes bytes to a file, overwriting existing data and creating parent directories.
func FileWriteBytes(path string, data []byte) error {
	if err := Mkdir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// FileAppendString appends content to a file and creates parent directories when needed.
func FileAppendString(path, content string) error {
	if err := Mkdir(filepath.Dir(path)); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer CloseQuietly(f)
	_, err = f.WriteString(content)
	return err
}

// Mkdir creates a directory tree. Empty and current-directory paths are treated as no-ops.
func Mkdir(dir string) error {
	if dir == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

// Touch creates an empty file when it does not exist.
func Touch(path string) error {
	if FileExists(path) {
		return nil
	}
	if err := Mkdir(filepath.Dir(path)); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}

// Del removes a file or directory recursively. Missing paths are treated as success.
func Del(path string) error {
	if !FileExists(path) {
		return nil
	}
	return os.RemoveAll(path)
}

// FileCopy copies a file and overwrites the destination when it already exists.
func FileCopy(src, dst string) error {
	if !IsFile(src) {
		return errors.New("source is not a file: " + src)
	}
	if err := Mkdir(filepath.Dir(dst)); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer CloseQuietly(in)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer CloseQuietly(out)
	_, err = io.Copy(out, in)
	return err
}

// MainName returns the file name without its extension; parent directories are ignored.
func MainName(path string) string {
	name := filepath.Base(path)
	ext := filepath.Ext(name)
	if ext == "" {
		return name
	}
	return strings.TrimSuffix(name, ext)
}

// Extension returns the file extension without the leading dot, or an empty string when absent.
func Extension(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return ""
	}
	return ext[1:]
}

// FileSize returns the file size in bytes, or -1 when the path is missing or not a file.
func FileSize(path string) int64 {
	st, err := os.Stat(path)
	if err != nil || st.IsDir() {
		return -1
	}
	return st.Size()
}

// ReaderFromString converts a string to an io.Reader.
func ReaderFromString(s string) io.Reader { return bytes.NewBufferString(s) }

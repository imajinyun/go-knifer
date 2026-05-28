package base

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 对应 hutool-core IoUtil。

// ReadAll 读取 Reader 全部内容。
func ReadAll(r io.Reader) ([]byte, error) { return io.ReadAll(r) }

// ReadString 读取 Reader 全部内容为字符串。
func ReadString(r io.Reader) (string, error) {
	b, err := ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadLines 按行读取 Reader 所有行。
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

// IoCopy 拷贝 Reader 到 Writer，返回已写入字节数。
func IoCopy(dst io.Writer, src io.Reader) (int64, error) { return io.Copy(dst, src) }

// CloseQuietly 安静关闭。
func CloseQuietly(c io.Closer) {
	if c == nil {
		return
	}
	_ = c.Close()
}

// 对应 hutool-core FileUtil / FileNameUtil。

// FileExists 文件或目录是否存在。
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsFile 是否为文件。
func IsFile(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

// IsDirectory 是否为目录。
func IsDirectory(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}

// FileReadString 读取文件全部内容为字符串。
func FileReadString(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FileReadBytes 读取文件全部字节。
func FileReadBytes(path string) ([]byte, error) { return os.ReadFile(path) }

// FileReadLines 按行读取文件。
func FileReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer CloseQuietly(f)
	return ReadLines(f)
}

// FileWriteString 将字符串写入文件（覆盖，自动创建目录）。
func FileWriteString(path, content string) error {
	if err := Mkdir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// FileWriteBytes 写入字节（覆盖）。
func FileWriteBytes(path string, data []byte) error {
	if err := Mkdir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// FileAppendString 追加字符串。
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

// Mkdir 递归创建目录。
func Mkdir(dir string) error {
	if dir == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

// Touch 不存在则创建空文件。
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

// Del 删除文件或目录（递归）。
func Del(path string) error {
	if !FileExists(path) {
		return nil
	}
	return os.RemoveAll(path)
}

// FileCopy 复制文件；目标存在时覆盖。
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

// MainName 不含扩展名的文件名（含路径中的目录会被忽略）。
func MainName(path string) string {
	name := filepath.Base(path)
	ext := filepath.Ext(name)
	if ext == "" {
		return name
	}
	return strings.TrimSuffix(name, ext)
}

// Extension 文件扩展名（不含 '.'，无扩展时返回空）。
func Extension(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return ""
	}
	return ext[1:]
}

// FileSize 文件大小（字节），不存在或非文件返回 -1。
func FileSize(path string) int64 {
	st, err := os.Stat(path)
	if err != nil || st.IsDir() {
		return -1
	}
	return st.Size()
}

// ReaderFromString 字符串转 Reader。
func ReaderFromString(s string) io.Reader { return bytes.NewBufferString(s) }

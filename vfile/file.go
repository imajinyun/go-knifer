package vfile

import (
	"io"

	fileimpl "github.com/imajinyun/go-knifer/internal/file"
)

func ReadAll(r io.Reader) ([]byte, error)              { return fileimpl.ReadAll(r) }
func ReadString(r io.Reader) (string, error)           { return fileimpl.ReadString(r) }
func ReadLines(r io.Reader) ([]string, error)          { return fileimpl.ReadLines(r) }
func Copy(dst io.Writer, src io.Reader) (int64, error) { return fileimpl.IoCopy(dst, src) }
func CloseQuietly(c io.Closer)                         { fileimpl.CloseQuietly(c) }
func Exists(path string) bool                          { return fileimpl.FileExists(path) }
func IsFile(path string) bool                          { return fileimpl.IsFile(path) }
func IsDirectory(path string) bool                     { return fileimpl.IsDirectory(path) }
func ReadFileString(path string) (string, error)       { return fileimpl.FileReadString(path) }
func ReadFileBytes(path string) ([]byte, error)        { return fileimpl.FileReadBytes(path) }
func ReadFileLines(path string) ([]string, error)      { return fileimpl.FileReadLines(path) }
func WriteFileString(path, content string) error       { return fileimpl.FileWriteString(path, content) }
func WriteFileBytes(path string, data []byte) error    { return fileimpl.FileWriteBytes(path, data) }
func AppendFileString(path, content string) error      { return fileimpl.FileAppendString(path, content) }
func Mkdir(dir string) error                           { return fileimpl.Mkdir(dir) }
func Touch(path string) error                          { return fileimpl.Touch(path) }
func Del(path string) error                            { return fileimpl.Del(path) }
func CopyFile(src, dst string) error                   { return fileimpl.FileCopy(src, dst) }
func MainName(path string) string                      { return fileimpl.MainName(path) }
func Extension(path string) string                     { return fileimpl.Extension(path) }
func Size(path string) int64                           { return fileimpl.FileSize(path) }
func ReaderFromString(s string) io.Reader              { return fileimpl.ReaderFromString(s) }

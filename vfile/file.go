package vfile

import (
	"io"
	"io/fs"

	fileimpl "github.com/imajinyun/go-knifer/internal/file"
)

// WriteOption customizes file write helpers.
type WriteOption = fileimpl.WriteOption

// DirOption customizes directory helpers.
type DirOption = fileimpl.DirOption

// WithFilePerm sets the file permission used when creating files.
func WithFilePerm(perm fs.FileMode) WriteOption { return fileimpl.WithFilePerm(perm) }

// WithDirPerm sets the parent-directory permission used when creating directories.
func WithDirPerm(perm fs.FileMode) WriteOption { return fileimpl.WithDirPerm(perm) }

// WithOverwrite controls whether an existing destination file may be replaced.
func WithOverwrite(overwrite bool) WriteOption { return fileimpl.WithOverwrite(overwrite) }

// WithCreateParents controls whether parent directories are created automatically.
func WithCreateParents(create bool) WriteOption { return fileimpl.WithCreateParents(create) }

// WithMkdirPerm sets the directory permission used by Mkdir.
func WithMkdirPerm(perm fs.FileMode) DirOption { return fileimpl.WithMkdirPerm(perm) }

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
func WriteFileString(path, content string, opts ...WriteOption) error {
	return fileimpl.FileWriteString(path, content, opts...)
}

func WriteFileBytes(path string, data []byte, opts ...WriteOption) error {
	return fileimpl.FileWriteBytes(path, data, opts...)
}

func AppendFileString(path, content string, opts ...WriteOption) error {
	return fileimpl.FileAppendString(path, content, opts...)
}
func Mkdir(dir string, opts ...DirOption) error { return fileimpl.Mkdir(dir, opts...) }
func Touch(path string, opts ...WriteOption) error {
	return fileimpl.Touch(path, opts...)
}
func Del(path string) error { return fileimpl.Del(path) }
func CopyFile(src, dst string, opts ...WriteOption) error {
	return fileimpl.FileCopy(src, dst, opts...)
}
func MainName(path string) string         { return fileimpl.MainName(path) }
func Extension(path string) string        { return fileimpl.Extension(path) }
func Size(path string) int64              { return fileimpl.FileSize(path) }
func ReaderFromString(s string) io.Reader { return fileimpl.ReaderFromString(s) }

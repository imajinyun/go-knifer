package file

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	knifer "github.com/imajinyun/go-knifer"
)

// Tests cover the utility toolkit-core IoUtilTest, FileUtilTest, and FileNameUtilTest.

func TestReadString(t *testing.T) {
	r := ReaderFromString("hello world")
	got, err := ReadString(r)
	if err != nil || got != "hello world" {
		t.Fatalf("ReadString: %v %q", err, got)
	}
}

func TestReadLines(t *testing.T) {
	r := ReaderFromString("a\nb\nc")
	lines, err := ReadLines(r)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(lines) != 3 || lines[0] != "a" || lines[2] != "c" {
		t.Fatalf("ReadLines: %v", lines)
	}
}

func TestReadOptions(t *testing.T) {
	if got, err := ReadStringWithOptions(ReaderFromString("abc"), WithMaxBytes(3)); err != nil || got != "abc" {
		t.Fatalf("ReadStringWithOptions exact limit = %q, %v", got, err)
	}
	if _, err := ReadStringWithOptions(ReaderFromString("abcd"), WithMaxBytes(3)); err == nil {
		t.Fatal("ReadStringWithOptions over limit error = nil")
	}

	lines, err := ReadLinesWithOptions(ReaderFromString("abc"), WithMaxBytes(3), WithInitialLineBuffer(1), WithMaxLineBytes(4))
	if err != nil {
		t.Fatalf("ReadLinesWithOptions exact limit: %v", err)
	}
	if len(lines) != 1 || lines[0] != "abc" {
		t.Fatalf("ReadLinesWithOptions lines = %v", lines)
	}
	if _, err := ReadLinesWithOptions(ReaderFromString("abcd"), WithMaxBytes(3), WithMaxLineBytes(4)); err == nil {
		t.Fatal("ReadLinesWithOptions over limit error = nil")
	}
}

func TestFileWriteRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "a.txt")
	if err := FileWriteString(path, "你好"); err != nil {
		t.Fatalf("write: %v", err)
	}
	if !FileExists(path) || !IsFile(path) {
		t.Fatalf("FileExists/IsFile failed")
	}
	if FileSize(path) <= 0 {
		t.Fatalf("FileSize failed")
	}
	got, err := FileReadString(path)
	if err != nil || got != "你好" {
		t.Fatalf("read: %v %q", err, got)
	}
	if err := FileAppendString(path, "X"); err != nil {
		t.Fatalf("append: %v", err)
	}
	got, _ = FileReadString(path)
	if !strings.HasSuffix(got, "X") {
		t.Fatalf("append result: %q", got)
	}
}

func TestFileLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lines.txt")
	if err := FileWriteString(path, "x\ny\nz"); err != nil {
		t.Fatalf("write: %v", err)
	}
	lines, err := FileReadLines(path)
	if err != nil {
		t.Fatalf("read lines: %v", err)
	}
	if len(lines) != 3 || lines[2] != "z" {
		t.Fatalf("lines: %v", lines)
	}
	if _, err := FileReadLinesWithOptions(path, WithMaxBytes(2)); err == nil {
		t.Fatal("FileReadLinesWithOptions over limit error = nil")
	}
}

func TestMkdirTouchDel(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "x", "y")
	if err := Mkdir(sub); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if !IsDirectory(sub) {
		t.Fatalf("IsDirectory failed")
	}
	f := filepath.Join(sub, "a.txt")
	if err := Touch(f); err != nil {
		t.Fatalf("touch: %v", err)
	}
	if !IsFile(f) {
		t.Fatalf("touch did not create file")
	}
	if err := Del(f); err != nil {
		t.Fatalf("del: %v", err)
	}
	if FileExists(f) {
		t.Fatalf("Del did not remove")
	}
}

func TestFileCopy(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.txt")
	dst := filepath.Join(dir, "sub", "b.txt")
	if err := FileWriteString(src, "hello"); err != nil {
		t.Fatalf("write src: %v", err)
	}
	if err := FileCopy(src, dst); err != nil {
		t.Fatalf("copy: %v", err)
	}
	got, _ := FileReadString(dst)
	if got != "hello" {
		t.Fatalf("copy content: %q", got)
	}
	missingErr := FileCopy(filepath.Join(dir, "missing.txt"), filepath.Join(dir, "unused.txt"))
	assertFileCode(t, missingErr, knifer.ErrCodeInvalidInput)
}

func TestWriteOptions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.txt")
	if err := FileWriteString(path, "first", WithFilePerm(0o600)); err != nil {
		t.Fatalf("write with options: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("perm = %o, want 600", got)
	}
	if err := FileWriteString(path, "second", WithOverwrite(false)); !errors.Is(err, os.ErrExist) {
		t.Fatalf("overwrite false error = %v, want exists", err)
	}
	if err := FileWriteString(path, "second", WithOverwrite(false)); !errors.Is(err, knifer.ErrCodeInternal) {
		t.Fatalf("overwrite false code = %v, want internal", err)
	}
	missingParent := filepath.Join(dir, "missing", "b.txt")
	assertFileCode(t, FileWriteString(missingParent, "x", WithCreateParents(false)), knifer.ErrCodeNotFound)
}

func TestMkdirAndCopyOptions(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "private")
	if err := Mkdir(sub, WithMkdirPerm(0o700)); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	info, err := os.Stat(sub)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o700 {
		t.Fatalf("dir perm = %o, want 700", got)
	}
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := FileWriteString(src, "src"); err != nil {
		t.Fatalf("write src: %v", err)
	}
	if err := FileWriteString(dst, "dst"); err != nil {
		t.Fatalf("write dst: %v", err)
	}
	if err := FileCopy(src, dst, WithOverwrite(false)); !errors.Is(err, os.ErrExist) {
		t.Fatalf("copy overwrite false error = %v, want exists", err)
	}
	if err := FileCopy(src, dst, WithOverwrite(false)); !errors.Is(err, knifer.ErrCodeInternal) {
		t.Fatalf("copy overwrite false code = %v, want internal", err)
	}
}

func TestFileReadProviderOptions(t *testing.T) {
	opened := ""
	open := func(path string) (io.ReadCloser, error) {
		opened = path
		return io.NopCloser(strings.NewReader("a\nb")), nil
	}
	got, err := FileReadStringWithOptions("virtual.txt", WithOpen(open), WithMaxBytes(3))
	if err != nil || got != "a\nb" {
		t.Fatalf("FileReadStringWithOptions() = %q, %v", got, err)
	}
	if opened != "virtual.txt" {
		t.Fatalf("open path = %q, want virtual.txt", opened)
	}
	lines, err := FileReadLinesWithOptions("virtual-lines.txt", WithOpen(open))
	if err != nil || len(lines) != 2 || lines[1] != "b" {
		t.Fatalf("FileReadLinesWithOptions() = %v, %v", lines, err)
	}
	if _, err := FileReadBytesWithOptions("too-large.txt", WithOpen(open), WithMaxBytes(2)); err == nil {
		t.Fatal("FileReadBytesWithOptions() over limit error = nil")
	}
	if _, err := FileReadStringWithOptions("ignored", WithOpen(func(string) (io.ReadCloser, error) {
		return nil, os.ErrNotExist
	})); !errors.Is(err, knifer.ErrCodeNotFound) {
		t.Fatalf("FileReadStringWithOptions() missing code = %v, want not found", err)
	}
}

func TestFileWriteProviderOptions(t *testing.T) {
	var mkdirPath string
	var mkdirPerm fs.FileMode
	var openPath string
	var openFlag int
	var openPerm fs.FileMode
	buf := &bytes.Buffer{}
	closer := &bufferWriteCloser{Buffer: buf}
	err := FileWriteBytes("parent/out.txt", []byte("payload"),
		WithDirPerm(0o700),
		WithFilePerm(0o600),
		WithMkdirAll(func(path string, perm fs.FileMode) error {
			mkdirPath = path
			mkdirPerm = perm
			return nil
		}),
		WithOpenFile(func(path string, flag int, perm fs.FileMode) (io.WriteCloser, error) {
			openPath = path
			openFlag = flag
			openPerm = perm
			return closer, nil
		}),
	)
	if err != nil {
		t.Fatalf("FileWriteBytes() error = %v", err)
	}
	if mkdirPath != "parent" || mkdirPerm != 0o700 {
		t.Fatalf("mkdir = %q/%o, want parent/700", mkdirPath, mkdirPerm)
	}
	if openPath != "parent/out.txt" || openPerm != 0o600 || openFlag&os.O_TRUNC == 0 {
		t.Fatalf("open = %q/%o/%#x", openPath, openPerm, openFlag)
	}
	if got := buf.String(); got != "payload" || !closer.closed {
		t.Fatalf("written = %q closed=%v, want payload/true", got, closer.closed)
	}
	if err := FileAppendString("append.txt", "!", WithCreateParents(false), WithOpenFile(func(string, int, fs.FileMode) (io.WriteCloser, error) {
		return &bufferWriteCloser{Buffer: buf}, nil
	})); err != nil {
		t.Fatalf("FileAppendString() with provider error = %v", err)
	}
}

func TestFileStatAndDeleteProviderOptions(t *testing.T) {
	regular := fakeFileInfo{name: "file.txt", size: 42}
	dir := fakeFileInfo{name: "dir", dir: true}
	stat := func(path string) (fs.FileInfo, error) {
		switch path {
		case "file.txt":
			return regular, nil
		case "dir":
			return dir, nil
		default:
			return nil, os.ErrNotExist
		}
	}
	if !FileExistsWithOptions("file.txt", WithStat(stat)) || !IsFileWithOptions("file.txt", WithStat(stat)) {
		t.Fatal("file stat providers failed")
	}
	if !IsDirectoryWithOptions("dir", WithStat(stat)) {
		t.Fatal("directory stat provider failed")
	}
	if got := FileSizeWithOptions("file.txt", WithStat(stat)); got != 42 {
		t.Fatalf("FileSizeWithOptions() = %d, want 42", got)
	}
	removed := ""
	if err := DelWithOptions("file.txt", WithStat(stat), WithRemoveAll(func(path string) error {
		removed = path
		return nil
	})); err != nil {
		t.Fatalf("DelWithOptions() error = %v", err)
	}
	if removed != "file.txt" {
		t.Fatalf("removed path = %q, want file.txt", removed)
	}
	removed = ""
	if err := DelWithOptions("missing.txt", WithStat(stat), WithRemoveAll(func(path string) error {
		removed = path
		return nil
	})); err != nil || removed != "" {
		t.Fatalf("DelWithOptions() missing = %v, removed %q", err, removed)
	}
}

func TestFileCopyProviderOptions(t *testing.T) {
	var copied bytes.Buffer
	if err := FileCopy("src.txt", "dst/out.txt",
		WithStat(func(path string) (fs.FileInfo, error) {
			if path == "src.txt" {
				return fakeFileInfo{name: path, size: 3}, nil
			}
			return nil, os.ErrNotExist
		}),
		WithMkdirAll(func(string, fs.FileMode) error { return nil }),
		WithOpen(func(path string) (io.ReadCloser, error) {
			if path != "src.txt" {
				t.Fatalf("open path = %q, want src.txt", path)
			}
			return io.NopCloser(strings.NewReader("abc")), nil
		}),
		WithOpenFile(func(path string, flag int, perm fs.FileMode) (io.WriteCloser, error) {
			if path != "dst/out.txt" {
				t.Fatalf("open destination = %q, want dst/out.txt", path)
			}
			return &bufferWriteCloser{Buffer: &copied}, nil
		}),
	); err != nil {
		t.Fatalf("FileCopy() with providers error = %v", err)
	}
	if got := copied.String(); got != "abc" {
		t.Fatalf("copied content = %q, want abc", got)
	}
}

type bufferWriteCloser struct {
	*bytes.Buffer
	closed bool
}

func (w *bufferWriteCloser) Close() error {
	w.closed = true
	return nil
}

type fakeFileInfo struct {
	name string
	size int64
	dir  bool
}

func (f fakeFileInfo) Name() string { return f.name }
func (f fakeFileInfo) Size() int64  { return f.size }
func (f fakeFileInfo) Mode() fs.FileMode {
	if f.dir {
		return fs.ModeDir | 0o755
	}
	return 0o644
}
func (f fakeFileInfo) ModTime() time.Time { return time.Unix(0, 0) }
func (f fakeFileInfo) IsDir() bool        { return f.dir }
func (f fakeFileInfo) Sys() any           { return nil }

func assertFileCode(t *testing.T, err error, code knifer.ErrCode) {
	t.Helper()
	if err == nil {
		t.Fatalf("err = nil, want %s", code)
	}
	if !errors.Is(err, code) {
		t.Fatalf("errors.Is(%v, %s) = false", err, code)
	}
	got, ok := knifer.CodeOf(err)
	if !ok || got != code {
		t.Fatalf("CodeOf(%v) = %q, %v; want %q, true", err, got, ok, code)
	}
}

func TestMainNameAndExtension(t *testing.T) {
	if MainName("/x/y/foo.txt") != "foo" {
		t.Fatalf("MainName failed")
	}
	if Extension("/x/y/foo.txt") != "txt" {
		t.Fatalf("Extension failed")
	}
	if MainName("foo") != "foo" || Extension("foo") != "" {
		t.Fatalf("no-ext failed")
	}
}

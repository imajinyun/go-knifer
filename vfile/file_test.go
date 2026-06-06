package vfile

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"strings"
	"testing"
	"time"

	knifer "github.com/imajinyun/go-knifer"
)

func TestFileFacade(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/a/b.txt"
	if err := WriteFileString(path, "line1\nline2"); err != nil {
		t.Fatal(err)
	}
	if !Exists(path) || !IsFile(path) || !IsDirectory(dir) {
		t.Fatal("file predicates failed")
	}
	if got, err := ReadFileString(path); err != nil || got != "line1\nline2" {
		t.Fatalf("ReadFileString = %q, %v", got, err)
	}
	if got, err := ReadFileBytes(path); err != nil || string(got) != "line1\nline2" {
		t.Fatalf("ReadFileBytes = %q, %v", got, err)
	}
	if lines, err := ReadFileLines(path); err != nil || len(lines) != 2 {
		t.Fatalf("ReadFileLines = %v, %v", lines, err)
	}
	if _, err := Copy(&strings.Builder{}, ReaderFromString("x")); err != nil {
		t.Fatal(err)
	}
	if MainName(path) != "b" || Extension(path) != "txt" || Size(path) <= 0 {
		t.Fatal("file name/size helpers failed")
	}
	copyPath := dir + "/copy.txt"
	if err := CopyFile(path, copyPath); err != nil || !Exists(copyPath) {
		t.Fatalf("CopyFile failed: %v", err)
	}
	if err := AppendFileString(copyPath, "!"); err != nil {
		t.Fatal(err)
	}
	if err := Touch(dir + "/touch.txt"); err != nil {
		t.Fatal(err)
	}
	if err := Del(dir + "/a"); err != nil || Exists(path) {
		t.Fatalf("Del failed: %v", err)
	}
	CloseQuietly(nil)
}

func TestFacadeFileErrorContract(t *testing.T) {
	err := CopyFile(t.TempDir()+"/missing.txt", t.TempDir()+"/out.txt")
	if err == nil {
		t.Fatal("CopyFile() error = nil, want invalid input")
	}
	if !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("errors.Is(err, ErrCodeInvalidInput) = false: %v", err)
	}
	code, ok := knifer.CodeOf(err)
	if !ok || code != knifer.ErrCodeInvalidInput {
		t.Fatalf("CodeOf(err) = %q, %v; want invalid input", code, ok)
	}
	var fileErr *Error
	if !errors.As(err, &fileErr) {
		t.Fatalf("errors.As(err, *vfile.Error) = false: %v", err)
	}
}

func TestFacadeReadOptions(t *testing.T) {
	if got, err := ReadStringWithOptions(ReaderFromString("abc"), WithMaxBytes(3)); err != nil || got != "abc" {
		t.Fatalf("ReadStringWithOptions() = %q, %v; want abc, nil", got, err)
	}
	if _, err := ReadStringWithOptions(ReaderFromString("abcd"), WithMaxBytes(3)); err == nil {
		t.Fatal("ReadStringWithOptions() over limit error = nil")
	}
	lines, err := ReadLinesWithOptions(ReaderFromString("abc"), WithInitialLineBuffer(1), WithMaxLineBytes(4))
	if err != nil {
		t.Fatalf("ReadLinesWithOptions() error = %v", err)
	}
	if len(lines) != 1 || lines[0] != "abc" {
		t.Fatalf("ReadLinesWithOptions() = %v, want [abc]", lines)
	}
}

func TestFacadeWriteAndDirOptions(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/private/a.txt"
	if err := WriteFileString(path, "first", WithCreateParents(false)); err == nil {
		t.Fatal("WriteFileString() should fail when parent directory is missing and WithCreateParents(false) is set")
	}
	if err := Mkdir(dir+"/private", WithMkdirPerm(0o700)); err != nil {
		t.Fatalf("Mkdir() with options error = %v", err)
	}
	dirInfo, err := os.Stat(dir + "/private")
	if err != nil {
		t.Fatalf("stat private dir: %v", err)
	}
	if got := dirInfo.Mode().Perm(); got != 0o700 {
		t.Fatalf("private dir perm = %o, want 700", got)
	}
	if err := WriteFileString(path, "first", WithFilePerm(0o600)); err != nil {
		t.Fatalf("WriteFileString() with options error = %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat file: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("file perm = %o, want 600", got)
	}
	if err := WriteFileBytes(path, []byte("second"), WithOverwrite(false)); err == nil {
		t.Fatal("WriteFileBytes() should reject overwrite=false for existing file")
	}
	if err := AppendFileString(path, "!"); err != nil {
		t.Fatalf("AppendFileString() error = %v", err)
	}
	copyPath := dir + "/copy.txt"
	if err := CopyFile(path, copyPath, WithFilePerm(0o600)); err != nil {
		t.Fatalf("CopyFile() with options error = %v", err)
	}
	if err := CopyFile(path, copyPath, WithOverwrite(false)); err == nil {
		t.Fatal("CopyFile() should reject overwrite=false for existing destination")
	}
	if err := Touch(dir+"/touch.txt", WithFilePerm(0o600)); err != nil {
		t.Fatalf("Touch() with options error = %v", err)
	}
}

func TestFacadeProviderOptions(t *testing.T) {
	if got, err := ReadFileStringWithOptions("virtual.txt", WithOpen(func(path string) (io.ReadCloser, error) {
		if path != "virtual.txt" {
			t.Fatalf("read path = %q, want virtual.txt", path)
		}
		return io.NopCloser(strings.NewReader("virtual")), nil
	})); err != nil || got != "virtual" {
		t.Fatalf("ReadFileStringWithOptions() = %q, %v", got, err)
	}
	if !ExistsWithOptions("x", WithStat(func(path string) (fs.FileInfo, error) {
		return fakeFacadeFileInfo{name: path}, nil
	})) {
		t.Fatal("ExistsWithOptions() = false, want true")
	}
	removed := false
	if err := DelWithOptions("x",
		WithStat(func(string) (fs.FileInfo, error) { return fakeFacadeFileInfo{name: "x"}, nil }),
		WithRemoveAll(func(string) error { removed = true; return nil }),
	); err != nil || !removed {
		t.Fatalf("DelWithOptions() = %v, removed=%v", err, removed)
	}
}

type fakeFacadeFileInfo struct {
	name string
}

func (f fakeFacadeFileInfo) Name() string       { return f.name }
func (f fakeFacadeFileInfo) Size() int64        { return 1 }
func (f fakeFacadeFileInfo) Mode() fs.FileMode  { return 0o644 }
func (f fakeFacadeFileInfo) ModTime() time.Time { return time.Unix(0, 0) }
func (f fakeFacadeFileInfo) IsDir() bool        { return false }
func (f fakeFacadeFileInfo) Sys() any           { return nil }

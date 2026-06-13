package vzip_test

import (
	archivezip "archive/zip"
	"bytes"
	"compress/flate"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vzip"
)

func TestFacadeZipAndCompression(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "data.zip")
	if err := vzip.ZipEntries(archive, vzip.EntryData{Name: "hello.txt", Data: []byte("hello")}); err != nil {
		t.Fatalf("ZipEntries: %v", err)
	}
	data, err := vzip.GetBytes(archive, "hello.txt")
	if err != nil || string(data) != "hello" {
		t.Fatalf("GetBytes: %q %v", data, err)
	}
	dest := filepath.Join(tmp, "dest")
	if err := vzip.UnzipTo(archive, dest); err != nil {
		t.Fatalf("UnzipTo: %v", err)
	}
	if got, err := os.ReadFile(filepath.Join(dest, "hello.txt")); err != nil || string(got) != "hello" {
		t.Fatalf("unzipped: %q %v", got, err)
	}
	gz, err := vzip.GzipString("hello")
	if err != nil {
		t.Fatalf("GzipString: %v", err)
	}
	text, err := vzip.UnGzipString(gz)
	if err != nil || text != "hello" {
		t.Fatalf("UnGzipString: %q %v", text, err)
	}
	dataBytes := []byte("hello the utility toolkit zip facade")
	gzipBytes, err := vzip.Gzip(dataBytes)
	if err != nil {
		t.Fatalf("Gzip: %v", err)
	}
	out, err := vzip.Gunzip(gzipBytes)
	if err != nil || !bytes.Equal(out, dataBytes) {
		t.Fatalf("Gunzip: %q %v", out, err)
	}
	zlibBytes, err := vzip.Zlib(dataBytes)
	if err != nil {
		t.Fatalf("Zlib: %v", err)
	}
	out, err = vzip.Unzlib(zlibBytes)
	if err != nil || !bytes.Equal(out, dataBytes) {
		t.Fatalf("Unzlib: %q %v", out, err)
	}
}

func TestFacadeZipOptions(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "options.zip")
	if err := vzip.ZipEntriesWithOptions(
		archive,
		[]vzip.EntryData{{Name: "hello.txt", Data: []byte("hello")}},
		vzip.WithFilePerm(0o600),
		vzip.WithCompressionLevel(1),
	); err != nil {
		t.Fatalf("ZipEntriesWithOptions: %v", err)
	}
	if err := vzip.ZipEntriesWithOptions(
		archive,
		[]vzip.EntryData{{Name: "hello.txt", Data: []byte("hello")}},
		vzip.WithOverwrite(false),
	); err == nil {
		t.Fatal("ZipEntriesWithOptions should reject overwrite=false for existing archive")
	}

	data, err := vzip.GetBytesWithOptions(archive, "hello.txt", vzip.WithMaxBytes(5))
	if err != nil || string(data) != "hello" {
		t.Fatalf("GetBytesWithOptions = %q, %v", data, err)
	}
	if _, err := vzip.GetBytesWithOptions(archive, "hello.txt", vzip.WithMaxBytes(4)); err == nil {
		t.Fatal("GetBytesWithOptions should reject content larger than max bytes")
	}

	dest := filepath.Join(tmp, "dest")
	if err := vzip.UnzipToWithOptions(archive, dest, vzip.WithDirPerm(0o700), vzip.WithFilePerm(0o600)); err != nil {
		t.Fatalf("UnzipToWithOptions: %v", err)
	}
	if got, err := os.ReadFile(filepath.Join(dest, "hello.txt")); err != nil || string(got) != "hello" {
		t.Fatalf("unzipped via options: %q %v", got, err)
	}

	gz, err := vzip.GzipString("hello")
	if err != nil {
		t.Fatalf("GzipString: %v", err)
	}
	if _, err := vzip.UnGzipWithOptions(gz, vzip.WithMaxBytes(4)); err == nil {
		t.Fatal("UnGzipWithOptions should reject content larger than max bytes")
	}
	zlibBytes, err := vzip.Zlib([]byte("hello"))
	if err != nil {
		t.Fatalf("Zlib: %v", err)
	}
	if _, err := vzip.UnZlibWithOptions(zlibBytes, vzip.WithMaxBytes(4)); err == nil {
		t.Fatal("UnZlibWithOptions should reject content larger than max bytes")
	}
}

func TestFacadeZipAppendAndGzipOptions(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "append.zip")
	if err := vzip.ZipEntries(archive, vzip.EntryData{Name: "a.txt", Data: []byte("a")}); err != nil {
		t.Fatalf("ZipEntries: %v", err)
	}
	keepFile := filepath.Join(tmp, "keep.txt")
	if err := os.WriteFile(keepFile, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	filter := func(path string, info os.FileInfo) bool {
		return info.IsDir() || filepath.Ext(path) == ".txt"
	}
	if err := vzip.AppendWithOptions(archive, keepFile, vzip.WithFileFilter(filter), vzip.WithCompressionLevel(flate.BestSpeed)); err != nil {
		t.Fatalf("AppendWithOptions: %v", err)
	}
	data, err := vzip.GetBytes(archive, "keep.txt")
	if err != nil || string(data) != "keep" {
		t.Fatalf("appended keep.txt = %q, %v", data, err)
	}

	payload := []byte("hello gzip options")
	gz, err := vzip.GzipWithOptions(payload, vzip.WithCompressionLevel(flate.BestSpeed))
	if err != nil {
		t.Fatalf("GzipWithOptions: %v", err)
	}
	out, err := vzip.UnGzip(gz)
	if err != nil || !bytes.Equal(out, payload) {
		t.Fatalf("UnGzip = %q, %v", out, err)
	}
	gz, err = vzip.GzipReaderWithOptions(bytes.NewReader(payload), 0, vzip.WithCompressionLevel(flate.NoCompression))
	if err != nil {
		t.Fatalf("GzipReaderWithOptions: %v", err)
	}
	out, err = vzip.UnGzip(gz)
	if err != nil || !bytes.Equal(out, payload) {
		t.Fatalf("UnGzip reader = %q, %v", out, err)
	}
}

func TestFacadeZipCreationUsingOptions(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src")
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "keep.txt"), []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "skip.log"), []byte("skip"), 0o644); err != nil {
		t.Fatal(err)
	}

	archive := filepath.Join(tmp, "filtered.zip")
	filter := func(path string, info os.FileInfo) bool {
		return info.IsDir() || filepath.Ext(path) == ".txt"
	}
	if err := vzip.ZipFilesUsingOptions(archive, []string{src}, vzip.WithSourceDir(true), vzip.WithFileFilter(filter)); err != nil {
		t.Fatalf("ZipFilesUsingOptions: %v", err)
	}
	data, err := vzip.GetBytes(archive, "src/keep.txt")
	if err != nil || string(data) != "keep" {
		t.Fatalf("GetBytes keep = %q, %v", data, err)
	}
	if _, err := vzip.GetBytes(archive, "src/skip.log"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("skip.log err = %v, want not exist", err)
	}

	var buf bytes.Buffer
	if err := vzip.ZipToWriterUsingOptions(&buf, []string{src}, vzip.WithFileFilter(filter)); err != nil {
		t.Fatalf("ZipToWriterUsingOptions: %v", err)
	}
	bufReader, err := archivezip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}
	if len(bufReader.File) != 1 || bufReader.File[0].Name != "keep.txt" {
		t.Fatalf("writer archive entries = %#v", bufReader.File)
	}
	entry, err := vzip.GetStream(bufReader.File[0])
	if err != nil {
		t.Fatalf("GetStream: %v", err)
	}
	defer func() { _ = entry.Close() }()
	if _, err := io.ReadAll(entry); err != nil {
		t.Fatalf("read entry: %v", err)
	}
}

func TestFacadeZipDefaultFileAndWriterHelpers(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src")
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "keep.txt"), []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "skip.log"), []byte("skip"), 0o644); err != nil {
		t.Fatal(err)
	}

	autoArchive, err := vzip.Zip(filepath.Join(src, "keep.txt"))
	if err != nil {
		t.Fatalf("Zip: %v", err)
	}
	if got, err := vzip.GetBytes(autoArchive, "keep.txt"); err != nil || string(got) != "keep" {
		t.Fatalf("Zip content = %q, %v", got, err)
	}

	toArchive := filepath.Join(tmp, "to.zip")
	if err := vzip.ZipTo(src, toArchive, true); err != nil {
		t.Fatalf("ZipTo: %v", err)
	}
	if got, err := vzip.GetBytes(toArchive, "src/keep.txt"); err != nil || string(got) != "keep" {
		t.Fatalf("ZipTo content = %q, %v", got, err)
	}

	filesArchive := filepath.Join(tmp, "files.zip")
	if err := vzip.ZipFiles(filesArchive, false, filepath.Join(src, "keep.txt")); err != nil {
		t.Fatalf("ZipFiles: %v", err)
	}
	if got, err := vzip.GetBytes(filesArchive, "keep.txt"); err != nil || string(got) != "keep" {
		t.Fatalf("ZipFiles content = %q, %v", got, err)
	}

	filterArchive := filepath.Join(tmp, "filter.zip")
	filter := func(path string, info os.FileInfo) bool {
		return info.IsDir() || filepath.Ext(path) == ".txt"
	}
	if err := vzip.ZipFilesFilter(filterArchive, false, filter, src); err != nil {
		t.Fatalf("ZipFilesFilter: %v", err)
	}
	if _, err := vzip.GetBytes(filterArchive, "skip.log"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ZipFilesFilter skip err = %v, want not exist", err)
	}

	var buf bytes.Buffer
	if err := vzip.ZipToWriter(&buf, false, filter, src); err != nil {
		t.Fatalf("ZipToWriter: %v", err)
	}
	zr, err := archivezip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil || len(zr.File) != 1 || zr.File[0].Name != "keep.txt" {
		t.Fatalf("ZipToWriter archive = %#v, %v", zr, err)
	}
}

func TestFacadeZipAppendUnzipAndCompressionOptions(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "append-default.zip")
	if err := vzip.ZipEntries(archive, vzip.EntryData{Name: "base.txt", Data: []byte("base")}); err != nil {
		t.Fatalf("ZipEntries: %v", err)
	}
	extra := filepath.Join(tmp, "extra.txt")
	if err := os.WriteFile(extra, []byte("extra"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := vzip.Append(archive, extra); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if got, err := vzip.GetBytes(archive, "extra.txt"); err != nil || string(got) != "extra" {
		t.Fatalf("Append content = %q, %v", got, err)
	}

	defaultDest, err := vzip.Unzip(archive)
	if err != nil {
		t.Fatalf("Unzip: %v", err)
	}
	if got, err := os.ReadFile(filepath.Join(defaultDest, "base.txt")); err != nil || string(got) != "base" {
		t.Fatalf("Unzip default output = %q, %v", got, err)
	}
	if err := vzip.UnzipToLimit(archive, filepath.Join(tmp, "limit"), 1); err == nil {
		t.Fatal("UnzipToLimit should reject content larger than limit")
	}

	payload := []byte("compression option payload")
	source := filepath.Join(tmp, "payload.txt")
	if err := os.WriteFile(source, payload, 0o644); err != nil {
		t.Fatal(err)
	}
	gz, err := vzip.GzipFileWithOptions(source, vzip.WithMaxBytes(int64(len(payload))))
	if err != nil {
		t.Fatalf("GzipFileWithOptions: %v", err)
	}
	if out, err := vzip.UnGzipReaderWithOptions(bytes.NewReader(gz), len(payload), vzip.WithMaxBytes(int64(len(payload)))); err != nil || !bytes.Equal(out, payload) {
		t.Fatalf("UnGzipReaderWithOptions = %q, %v", out, err)
	}
	if _, err := vzip.GzipFileWithOptions(source, vzip.WithMaxBytes(1)); err == nil {
		t.Fatal("GzipFileWithOptions max bytes error = nil")
	}

	zlibBytes, err := vzip.ZlibFileWithOptions(source, flate.BestSpeed, vzip.WithMaxBytes(int64(len(payload))))
	if err != nil {
		t.Fatalf("ZlibFileWithOptions: %v", err)
	}
	if out, err := vzip.UnZlibReaderWithOptions(bytes.NewReader(zlibBytes), len(payload), vzip.WithMaxBytes(int64(len(payload)))); err != nil || !bytes.Equal(out, payload) {
		t.Fatalf("UnZlibReaderWithOptions = %q, %v", out, err)
	}
	if got, err := vzip.ZlibLevelWithOptions(payload, flate.NoCompression, vzip.WithMaxBytes(int64(len(payload)))); err != nil || len(got) == 0 {
		t.Fatalf("ZlibLevelWithOptions len=%d err=%v", len(got), err)
	}
	if got, err := vzip.ZlibReaderWithOptions(bytes.NewReader(payload), flate.BestSpeed, len(payload), vzip.WithMaxBytes(int64(len(payload)))); err != nil || len(got) == 0 {
		t.Fatalf("ZlibReaderWithOptions len=%d err=%v", len(got), err)
	}
}

func TestFacadeZipErrorContract(t *testing.T) {
	_, err := vzip.GetStream(nil)
	if err == nil {
		t.Fatal("GetStream(nil) error = nil, want invalid input")
	}
	if !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("errors.Is(err, ErrCodeInvalidInput) = false: %v", err)
	}
	code, ok := knifer.CodeOf(err)
	if !ok || code != knifer.ErrCodeInvalidInput {
		t.Fatalf("CodeOf(err) = %q, %v; want invalid input", code, ok)
	}
	var zipErr *vzip.Error
	if !errors.As(err, &zipErr) {
		t.Fatalf("errors.As(err, *vzip.Error) = false: %v", err)
	}
}

func TestFacadeZipArchiveReaderAndWriterHelpers(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "helpers.zip")

	if err := vzip.ZipData(archive, "text/data.txt", "hello"); err != nil {
		t.Fatalf("ZipData: %v", err)
	}
	rc, err := vzip.Open(archive)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if len(rc.File) != 1 || rc.File[0].Name != "text/data.txt" {
		_ = rc.Close()
		t.Fatalf("Open files = %#v", rc.File)
	}
	_ = rc.Close()

	names, err := vzip.ListFileNames(archive, "text")
	if err != nil || len(names) != 1 || names[0] != "data.txt" {
		t.Fatalf("ListFileNames = %#v, %v", names, err)
	}
	seen := false
	if err := vzip.Read(archive, func(f *archivezip.File) error {
		seen = f.Name == "text/data.txt"
		return nil
	}); err != nil || !seen {
		t.Fatalf("Read seen=%v err=%v", seen, err)
	}
	reader, err := vzip.Get(archive, "text/data.txt")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	data, err := io.ReadAll(reader)
	_ = reader.Close()
	if err != nil || string(data) != "hello" {
		t.Fatalf("Get data = %q, %v", data, err)
	}

	bytesArchive := filepath.Join(tmp, "bytes.zip")
	if err := vzip.ZipBytes(bytesArchive, "bytes.bin", []byte{1, 2, 3}); err != nil {
		t.Fatalf("ZipBytes: %v", err)
	}
	if got, err := vzip.GetBytes(bytesArchive, "bytes.bin"); err != nil || !bytes.Equal(got, []byte{1, 2, 3}) {
		t.Fatalf("ZipBytes content = %v, %v", got, err)
	}

	var buf bytes.Buffer
	zw := vzip.NewWriter(&buf)
	w, err := zw.Create("manual.txt")
	if err != nil {
		t.Fatalf("Create manual entry: %v", err)
	}
	if _, err := w.Write([]byte("manual")); err != nil {
		t.Fatalf("write manual entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close manual writer: %v", err)
	}
	zr, err := archivezip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil || len(zr.File) != 1 || zr.File[0].Name != "manual.txt" {
		t.Fatalf("NewWriter archive = %#v, %v", zr, err)
	}
}

func TestFacadeZipStreamAndReaderHelpers(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "streams.zip")

	if err := vzip.ZipStreams(archive, vzip.StreamEntry{Name: "stream.txt", Reader: strings.NewReader("stream")}); err != nil {
		t.Fatalf("ZipStreams: %v", err)
	}
	if got, err := vzip.GetBytes(archive, "stream.txt"); err != nil || string(got) != "stream" {
		t.Fatalf("ZipStreams content = %q, %v", got, err)
	}

	var buf bytes.Buffer
	if err := vzip.ZipEntriesToWriter(&buf, vzip.EntryData{Name: "entry.txt", Data: []byte("entry")}); err != nil {
		t.Fatalf("ZipEntriesToWriter: %v", err)
	}
	var streamBuf bytes.Buffer
	if err := vzip.ZipStreamsToWriter(&streamBuf, vzip.StreamEntry{Name: "stream2.txt", Reader: strings.NewReader("stream2")}); err != nil {
		t.Fatalf("ZipStreamsToWriter: %v", err)
	}
	streamReader, err := archivezip.NewReader(bytes.NewReader(streamBuf.Bytes()), int64(streamBuf.Len()))
	if err != nil || len(streamReader.File) != 1 || streamReader.File[0].Name != "stream2.txt" {
		t.Fatalf("ZipStreamsToWriter archive = %#v, %v", streamReader, err)
	}

	dest := filepath.Join(tmp, "dest")
	zr, err := archivezip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}
	if err := vzip.UnzipReaderTo(zr, dest); err != nil {
		t.Fatalf("UnzipReaderTo: %v", err)
	}
	if got, err := os.ReadFile(filepath.Join(dest, "entry.txt")); err != nil || string(got) != "entry" {
		t.Fatalf("UnzipReaderTo output = %q, %v", got, err)
	}
	if err := vzip.UnzipReaderToLimit(zr, filepath.Join(tmp, "limited"), 1); err == nil {
		t.Fatal("UnzipReaderToLimit should reject archive larger than limit")
	}
}

func TestFacadeGzipZlibAndReadFileHelpers(t *testing.T) {
	tmp := t.TempDir()
	source := filepath.Join(tmp, "payload.txt")
	payload := []byte("payload for compression helpers")
	if err := os.WriteFile(source, payload, 0o600); err != nil {
		t.Fatal(err)
	}

	gz, err := vzip.GzipFile(source)
	if err != nil {
		t.Fatalf("GzipFile: %v", err)
	}
	gunzip, err := vzip.UnGzipReader(bytes.NewReader(gz), len(payload))
	if err != nil || !bytes.Equal(gunzip, payload) {
		t.Fatalf("UnGzipReader = %q, %v", gunzip, err)
	}

	zlibBytes, err := vzip.ZlibString("hello zlib", flate.BestSpeed)
	if err != nil {
		t.Fatalf("ZlibString: %v", err)
	}
	zlibText, err := vzip.UnZlibString(zlibBytes)
	if err != nil || zlibText != "hello zlib" {
		t.Fatalf("UnZlibString = %q, %v", zlibText, err)
	}
	levelBytes, err := vzip.ZlibLevel(payload, flate.BestCompression)
	if err != nil {
		t.Fatalf("ZlibLevel: %v", err)
	}
	levelOut, err := vzip.UnZlibReader(bytes.NewReader(levelBytes), len(payload))
	if err != nil || !bytes.Equal(levelOut, payload) {
		t.Fatalf("UnZlibReader = %q, %v", levelOut, err)
	}
	readerBytes, err := vzip.ZlibReader(bytes.NewReader(payload), flate.NoCompression, len(payload))
	if err != nil {
		t.Fatalf("ZlibReader: %v", err)
	}
	if out, err := vzip.UnZlib(readerBytes); err != nil || !bytes.Equal(out, payload) {
		t.Fatalf("UnZlib reader bytes = %q, %v", out, err)
	}

	read, err := vzip.ReadFile(source)
	if err != nil || !bytes.Equal(read, payload) {
		t.Fatalf("ReadFile = %q, %v", read, err)
	}
	read, err = vzip.ReadFileWithOptions("/virtual/payload.txt", vzip.WithReadFile(func(path string) ([]byte, error) {
		if path != "/virtual/payload.txt" {
			return nil, os.ErrNotExist
		}
		return []byte("virtual"), nil
	}))
	if err != nil || string(read) != "virtual" {
		t.Fatalf("ReadFileWithOptions = %q, %v", read, err)
	}
}

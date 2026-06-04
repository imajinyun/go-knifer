package zip

import (
	archivezip "archive/zip"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultBufferSize = 32
	maxInt64          = int64(1<<63 - 1)
)

// FileFilter decides whether a source path should be added to an archive.
type FileFilter func(path string, info os.FileInfo) bool

// Entry describes an archive entry.
type Entry = archivezip.File

// Writer is a ZIP archive writer.
type Writer = archivezip.Writer

// Reader is a ZIP archive reader.
type Reader = archivezip.ReadCloser

// EntryData represents in-memory content to add into a ZIP archive.
type EntryData struct {
	Name string
	Data []byte
}

// StreamEntry represents stream content to add into a ZIP archive.
type StreamEntry struct {
	Name   string
	Reader io.Reader
}

// Open opens a ZIP file for reading.
func Open(path string) (*archivezip.ReadCloser, error) { return archivezip.OpenReader(path) }

// NewWriter returns a ZIP writer for out.
func NewWriter(out io.Writer) *archivezip.Writer { return archivezip.NewWriter(out) }

// GetStream returns a reader for entry.
func GetStream(entry *archivezip.File) (io.ReadCloser, error) {
	if entry == nil {
		return nil, invalidInputf("zip entry is nil")
	}
	return entry.Open()
}

// Append appends srcPath into zipPath by rewriting the archive.
func Append(zipPath, srcPath string) error {
	return appendWithFilter(zipPath, srcPath, nil)
}

// Zip creates an archive next to srcPath and returns the archive path.
func Zip(srcPath string) (string, error) {
	dest := strings.TrimSuffix(srcPath, filepath.Ext(srcPath)) + ".zip"
	return dest, ZipTo(srcPath, dest, false)
}

// ZipTo creates an archive at zipPath from srcPath.
func ZipTo(srcPath, zipPath string, withSrcDir bool) error {
	return ZipFiles(zipPath, withSrcDir, srcPath)
}

// ZipFiles creates a ZIP archive from source files or directories.
func ZipFiles(dest string, withSrcDir bool, srcFiles ...string) (err error) {
	return ZipFilesFilter(dest, withSrcDir, nil, srcFiles...)
}

// ZipFilesFilter creates a ZIP archive and filters source paths.
func ZipFilesFilter(dest string, withSrcDir bool, filter FileFilter, srcFiles ...string) (err error) {
	if err := validateZipTarget(dest, srcFiles...); err != nil {
		return err
	}
	if dir := filepath.Dir(dest); dir != "." {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return err
		}
	}
	out, err := os.Create(dest) // #nosec G304 -- destination path is an explicit caller-provided archive output.
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := out.Close(); err == nil {
			err = closeErr
		}
	}()
	return ZipToWriter(out, withSrcDir, filter, srcFiles...)
}

// ZipToWriter writes source files or directories into out as a ZIP archive.
func ZipToWriter(out io.Writer, withSrcDir bool, filter FileFilter, srcFiles ...string) (err error) {
	zw := archivezip.NewWriter(out)
	defer func() {
		if closeErr := zw.Close(); err == nil {
			err = closeErr
		}
	}()
	for _, src := range srcFiles {
		if src == "" {
			continue
		}
		info, err := os.Lstat(src)
		if err != nil {
			return err
		}
		base := filepath.Dir(src)
		name := filepath.Base(src)
		if info.IsDir() && !withSrcDir {
			base = src
			name = ""
		}
		if err := addPath(zw, src, base, name, filter); err != nil {
			return err
		}
	}
	return nil
}

// ZipData creates or overwrites zipFile and adds one text entry.
func ZipData(zipFile, path, data string) error {
	return ZipBytes(zipFile, path, []byte(data))
}

// ZipBytes creates or overwrites zipFile and adds one byte entry.
func ZipBytes(zipFile, path string, data []byte) error {
	return ZipEntries(zipFile, EntryData{Name: path, Data: data})
}

// ZipEntries creates or overwrites zipFile and adds in-memory entries.
func ZipEntries(zipFile string, entries ...EntryData) (err error) {
	if dir := filepath.Dir(zipFile); dir != "." {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return err
		}
	}
	out, err := os.Create(zipFile) // #nosec G304 -- destination path is an explicit caller-provided archive output.
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := out.Close(); err == nil {
			err = closeErr
		}
	}()
	return ZipEntriesToWriter(out, entries...)
}

// ZipEntriesToWriter writes in-memory entries into out as a ZIP archive.
func ZipEntriesToWriter(out io.Writer, entries ...EntryData) (err error) {
	streams := make([]StreamEntry, 0, len(entries))
	for _, entry := range entries {
		streams = append(streams, StreamEntry{Name: entry.Name, Reader: bytes.NewReader(entry.Data)})
	}
	return ZipStreamsToWriter(out, streams...)
}

// ZipStreams creates or overwrites zipFile and adds stream entries.
func ZipStreams(zipFile string, entries ...StreamEntry) (err error) {
	if dir := filepath.Dir(zipFile); dir != "." {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return err
		}
	}
	out, err := os.Create(zipFile) // #nosec G304 -- destination path is an explicit caller-provided archive output.
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := out.Close(); err == nil {
			err = closeErr
		}
	}()
	return ZipStreamsToWriter(out, entries...)
}

// ZipStreamsToWriter writes stream entries into out as a ZIP archive.
func ZipStreamsToWriter(out io.Writer, entries ...StreamEntry) (err error) {
	zw := archivezip.NewWriter(out)
	defer func() {
		if closeErr := zw.Close(); err == nil {
			err = closeErr
		}
	}()
	for _, entry := range entries {
		name, err := cleanEntryName(entry.Name)
		if err != nil {
			return err
		}
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, entry.Reader); err != nil {
			return err
		}
	}
	return nil
}

// Unzip extracts zipFile into a sibling directory named after the archive.
func Unzip(zipFile string) (string, error) {
	dest := strings.TrimSuffix(zipFile, filepath.Ext(zipFile))
	return dest, UnzipTo(zipFile, dest)
}

// UnzipTo extracts zipFile into destDir.
func UnzipTo(zipFile, destDir string) error { return UnzipToLimit(zipFile, destDir, -1) }

// UnzipToLimit extracts zipFile into destDir and optionally limits total uncompressed size.
func UnzipToLimit(zipFile, destDir string, limit int64) error {
	r, err := archivezip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()
	return UnzipReaderToLimit(&r.Reader, destDir, limit)
}

// UnzipReaderTo extracts archive reader contents into destDir.
func UnzipReaderTo(r *archivezip.Reader, destDir string) error {
	return UnzipReaderToLimit(r, destDir, -1)
}

// UnzipReaderToLimit extracts archive reader contents into destDir and optionally limits total size.
func UnzipReaderToLimit(r *archivezip.Reader, destDir string, limit int64) error {
	if r == nil {
		return invalidInputf("zip reader is nil")
	}
	if err := os.MkdirAll(destDir, 0o750); err != nil {
		return err
	}
	var total int64
	for _, f := range r.File {
		if limit > 0 {
			if f.UncompressedSize64 > uint64(maxInt64) {
				return invalidInputf("uncompressed size exceeds int64 limit")
			}
			size := int64(f.UncompressedSize64) // #nosec G115 -- guarded by the maxInt64 check above.
			if total > maxInt64-size {
				return invalidInputf("uncompressed size exceeds int64 limit")
			}
			total += size
			if total > limit {
				return invalidInputf("uncompressed size exceeds limit")
			}
		}
		if err := extractFile(f, destDir); err != nil {
			return err
		}
	}
	return nil
}

// Get returns a reader for the named entry in zipFile.
func Get(zipFile, name string) (io.ReadCloser, error) {
	r, err := archivezip.OpenReader(zipFile)
	if err != nil {
		return nil, err
	}
	for _, f := range r.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				_ = r.Close()
				return nil, err
			}
			return &readCloserWithClose{ReadCloser: rc, close: r.Close}, nil
		}
	}
	_ = r.Close()
	return nil, notFound("zip entry not found: "+name, os.ErrNotExist)
}

// GetBytes returns the content of the named entry in zipFile.
func GetBytes(zipFile, name string) ([]byte, error) {
	rc, err := Get(zipFile, name)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rc.Close() }()
	return io.ReadAll(rc)
}

// Read walks every archive entry and calls consumer.
func Read(zipFile string, consumer func(*archivezip.File) error) error {
	r, err := archivezip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()
	for _, f := range r.File {
		if err := consumer(f); err != nil {
			return err
		}
	}
	return nil
}

// ListFileNames returns direct file names under dir inside zipFile.
func ListFileNames(zipFile, dir string) ([]string, error) {
	r, err := archivezip.OpenReader(zipFile)
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()
	if strings.TrimSpace(dir) != "" && !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	names := make([]string, 0)
	for _, f := range r.File {
		name := f.Name
		if dir != "" {
			if !strings.HasPrefix(name, dir) {
				continue
			}
			name = strings.TrimPrefix(name, dir)
		}
		if name != "" && !strings.Contains(name, "/") && !f.FileInfo().IsDir() {
			names = append(names, name)
		}
	}
	return names, nil
}

// Gzip compresses data using gzip.
func Gzip(data []byte) ([]byte, error) { return GzipReader(bytes.NewReader(data), len(data)) }

// GzipString compresses text using gzip.
func GzipString(content string) ([]byte, error) { return Gzip([]byte(content)) }

// GzipFile compresses a file using gzip and returns compressed bytes.
func GzipFile(path string) ([]byte, error) {
	// #nosec G304 -- SDK file helper intentionally opens the caller-provided path.
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return GzipReader(f, int(info.Size()))
}

// GzipReader compresses all bytes from r using gzip.
func GzipReader(r io.Reader, estimatedLength int) ([]byte, error) {
	if estimatedLength <= 0 {
		estimatedLength = defaultBufferSize
	}
	var buf bytes.Buffer
	buf.Grow(estimatedLength)
	w := gzip.NewWriter(&buf)
	if _, err := io.Copy(w, r); err != nil {
		_ = w.Close()
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnGzip decompresses gzip data.
func UnGzip(data []byte) ([]byte, error) { return UnGzipReader(bytes.NewReader(data), len(data)) }

// Gunzip decompresses gzip data.
func Gunzip(data []byte) ([]byte, error) { return UnGzip(data) }

// UnGzipString decompresses gzip data and returns text.
func UnGzipString(data []byte) (string, error) {
	out, err := UnGzip(data)
	return string(out), err
}

// UnGzipReader decompresses all gzip bytes from r.
func UnGzipReader(r io.Reader, estimatedLength int) ([]byte, error) {
	if estimatedLength <= 0 {
		estimatedLength = defaultBufferSize
	}
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer func() { _ = zr.Close() }()
	var buf bytes.Buffer
	buf.Grow(estimatedLength)
	_, err = io.Copy(&buf, zr) // #nosec G110 -- this low-level helper intentionally decompresses caller-provided gzip data.
	return buf.Bytes(), err
}

// Zlib compresses data using zlib with the default compression level.
func Zlib(data []byte) ([]byte, error) { return ZlibLevel(data, flate.DefaultCompression) }

// ZlibString compresses text using zlib with the specified compression level.
func ZlibString(content string, level int) ([]byte, error) { return ZlibLevel([]byte(content), level) }

// ZlibFile compresses a file using zlib with the specified compression level.
func ZlibFile(path string, level int) ([]byte, error) {
	// #nosec G304 -- SDK file helper intentionally opens the caller-provided path.
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return ZlibReader(f, level, int(info.Size()))
}

// ZlibLevel compresses data using zlib with the specified compression level.
func ZlibLevel(data []byte, level int) ([]byte, error) {
	return ZlibReader(bytes.NewReader(data), level, len(data))
}

// ZlibReader compresses all bytes from r using zlib with the specified compression level.
func ZlibReader(r io.Reader, level, estimatedLength int) ([]byte, error) {
	if estimatedLength <= 0 {
		estimatedLength = defaultBufferSize
	}
	var buf bytes.Buffer
	buf.Grow(estimatedLength)
	w, err := zlib.NewWriterLevel(&buf, level)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(w, r); err != nil {
		_ = w.Close()
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnZlib decompresses zlib data.
func UnZlib(data []byte) ([]byte, error) { return UnZlibReader(bytes.NewReader(data), len(data)) }

// Unzlib decompresses zlib data.
func Unzlib(data []byte) ([]byte, error) { return UnZlib(data) }

// UnZlibString decompresses zlib data and returns text.
func UnZlibString(data []byte) (string, error) {
	out, err := UnZlib(data)
	return string(out), err
}

// UnZlibReader decompresses all zlib bytes from r.
func UnZlibReader(r io.Reader, estimatedLength int) ([]byte, error) {
	if estimatedLength <= 0 {
		estimatedLength = defaultBufferSize
	}
	zr, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer func() { _ = zr.Close() }()
	var buf bytes.Buffer
	buf.Grow(estimatedLength)
	_, err = io.Copy(&buf, zr) // #nosec G110 -- this low-level helper intentionally decompresses caller-provided zlib data.
	return buf.Bytes(), err
}

func appendWithFilter(zipPath, srcPath string, filter FileFilter) error {
	tmp, err := os.CreateTemp(filepath.Dir(zipPath), ".zip-append-*.zip")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	zw := archivezip.NewWriter(tmp)
	if _, err := os.Stat(zipPath); err == nil {
		r, err := archivezip.OpenReader(zipPath)
		if err != nil {
			_ = zw.Close()
			_ = tmp.Close()
			_ = os.Remove(tmpPath)
			return err
		}
		for _, f := range r.File {
			if err := copyExistingEntry(zw, f); err != nil {
				_ = r.Close()
				_ = zw.Close()
				_ = tmp.Close()
				_ = os.Remove(tmpPath)
				return err
			}
		}
		_ = r.Close()
	}
	info, err := os.Lstat(srcPath)
	if err != nil {
		_ = zw.Close()
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	base := filepath.Dir(srcPath)
	name := filepath.Base(srcPath)
	if info.IsDir() && filepath.Dir(srcPath) == srcPath {
		base = srcPath
		name = ""
	}
	if err := addPath(zw, srcPath, base, name, filter); err != nil {
		_ = zw.Close()
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	if err := zw.Close(); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return os.Rename(tmpPath, zipPath)
}

func addPath(zw *archivezip.Writer, path, base, name string, filter FileFilter) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if filter != nil && !filter(path, info) {
		return nil
	}
	zipName := name
	if zipName == "" {
		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}
		if rel == "." {
			zipName = ""
		} else {
			zipName = rel
		}
	}
	zipName = filepath.ToSlash(filepath.Clean(zipName))
	if zipName == "." {
		zipName = ""
	}
	if zipName != "" {
		if _, err := cleanEntryName(zipName); err != nil {
			return err
		}
	}
	if info.IsDir() {
		if zipName != "" {
			header, err := archivezip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name = strings.TrimSuffix(zipName, "/") + "/"
			header.SetMode(info.Mode())
			if _, err := zw.CreateHeader(header); err != nil {
				return err
			}
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			child := filepath.Join(path, entry.Name())
			childName := entry.Name()
			if zipName != "" {
				childName = filepath.Join(zipName, entry.Name())
			}
			if err := addPath(zw, child, base, childName, filter); err != nil {
				return err
			}
		}
		return nil
	}
	return addFile(zw, path, zipName, info)
}

func addFile(zw *archivezip.Writer, path, zipName string, info os.FileInfo) error {
	header, err := archivezip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = zipName
	header.Method = archivezip.Deflate
	header.SetMode(info.Mode())
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
	r, err := os.Open(path) // #nosec G304 -- archive creation intentionally reads caller-provided source paths.
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()
	_, err = io.Copy(w, r) // #nosec G110 -- copying existing archive entries preserves caller-provided archive contents.
	return err
}

func extractFile(f *archivezip.File, destDir string) error {
	target, err := safeZipTarget(destDir, f.Name)
	if err != nil {
		return err
	}
	if f.FileInfo().IsDir() {
		return os.MkdirAll(target, f.Mode())
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o750); err != nil {
		return err
	}
	r, err := f.Open()
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()
	w, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode()) // #nosec G304 -- target is validated by safeZipTarget before extraction.
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, r); err != nil { // #nosec G110 -- unzip extraction is guarded by safeZipTarget and optional UnzipToLimit size checks.
		_ = w.Close()
		return err
	}
	return w.Close()
}

func copyExistingEntry(zw *archivezip.Writer, f *archivezip.File) error {
	r, err := f.Open()
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()
	header := f.FileHeader
	w, err := zw.CreateHeader(&header)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r) // #nosec G110 -- copying existing archive entries preserves caller-provided archive contents.
	return err
}

func cleanEntryName(name string) (string, error) {
	if name == "" || filepath.IsAbs(name) {
		return "", invalidInputf("invalid zip entry name %q", name)
	}
	cleaned := filepath.ToSlash(filepath.Clean(name))
	if cleaned == "." || strings.HasPrefix(cleaned, "../") || cleaned == ".." || strings.HasPrefix(cleaned, "/") {
		return "", invalidInputf("invalid zip entry name %q", name)
	}
	return cleaned, nil
}

func safeZipTarget(destDir, name string) (string, error) {
	cleaned, err := cleanEntryName(name)
	if err != nil {
		return "", err
	}
	target := filepath.Join(destDir, filepath.FromSlash(cleaned))
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
		return "", invalidInputf("invalid zip entry name %q", name)
	}
	return target, nil
}

func validateZipTarget(zipFile string, srcFiles ...string) error {
	info, err := os.Stat(zipFile)
	if err == nil && info.IsDir() {
		return invalidInputf("zip file %q must not be a directory", zipFile)
	}
	zipAbs, err := filepath.Abs(zipFile)
	if err != nil {
		return err
	}
	zipDir := filepath.Dir(zipAbs)
	for _, src := range srcFiles {
		if src == "" {
			continue
		}
		info, err := os.Stat(src)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			continue
		}
		srcAbs, err := filepath.Abs(src)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcAbs, zipDir)
		if err != nil {
			return err
		}
		if rel == "." || (!strings.HasPrefix(rel, "..") && rel != "") {
			return invalidInputf("zip file path %q must not be inside source directory %q", zipFile, src)
		}
	}
	return nil
}

type readCloserWithClose struct {
	io.ReadCloser
	close func() error
}

func (r *readCloserWithClose) Close() error {
	err1 := r.ReadCloser.Close()
	err2 := r.close()
	if err1 != nil {
		return err1
	}
	return err2
}

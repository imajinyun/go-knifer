package net

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// UploadSetting configures multipart form parsing and file saving.
type UploadSetting struct {
	MaxFileSize     int64
	MemoryThreshold int64
	TmpUploadPath   string
	FileExts        []string
	AllowFileExts   bool
}

// NewUploadSetting returns a default upload setting.
func NewUploadSetting() UploadSetting {
	return UploadSetting{MaxFileSize: 32 << 20, MemoryThreshold: 32 << 20, AllowFileExts: true}
}

// MultipartFormData wraps a parsed multipart form.
type MultipartFormData struct {
	Form   *multipart.Form
	loaded bool
}

// ParseMultipartForm parses multipart/form-data from an HTTP request.
func ParseMultipartForm(r *http.Request, setting UploadSetting) (*MultipartFormData, error) {
	if setting.MemoryThreshold <= 0 {
		setting.MemoryThreshold = 32 << 20
	}
	r.Body = http.MaxBytesReader(nil, r.Body, setting.MaxFileSize)        //nolint:bodyclose // request body lifecycle is owned by caller.
	if err := r.ParseMultipartForm(setting.MemoryThreshold); err != nil { //nolint:gosec // request body is bounded by MaxBytesReader above.
		return nil, err
	}
	return &MultipartFormData{Form: r.MultipartForm, loaded: true}, nil
}

// IsLoaded reports whether the form has been parsed.
func (m *MultipartFormData) IsLoaded() bool { return m != nil && m.loaded }

// GetParam returns the first value for name.
func (m *MultipartFormData) GetParam(name string) string {
	values := m.GetListParam(name)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// GetParamNames returns parameter names.
func (m *MultipartFormData) GetParamNames() []string {
	if m == nil || m.Form == nil {
		return nil
	}
	out := make([]string, 0, len(m.Form.Value))
	for k := range m.Form.Value {
		out = append(out, k)
	}
	return out
}

// GetArrayParam returns all values for name.
func (m *MultipartFormData) GetArrayParam(name string) []string { return m.GetListParam(name) }

// GetListParam returns all values for name.
func (m *MultipartFormData) GetListParam(name string) []string {
	if m == nil || m.Form == nil {
		return nil
	}
	return m.Form.Value[name]
}

// GetParamMap returns first parameter values.
func (m *MultipartFormData) GetParamMap() map[string]string {
	out := map[string]string{}
	if m == nil || m.Form == nil {
		return out
	}
	for k, values := range m.Form.Value {
		if len(values) > 0 {
			out[k] = values[0]
		}
	}
	return out
}

// GetParamListMap returns all parameter values.
func (m *MultipartFormData) GetParamListMap() map[string][]string {
	if m == nil || m.Form == nil {
		return map[string][]string{}
	}
	return m.Form.Value
}

// GetFile returns the first file for name.
func (m *MultipartFormData) GetFile(name string) *multipart.FileHeader {
	files := m.GetFileList(name)
	if len(files) == 0 {
		return nil
	}
	return files[0]
}

// GetFiles returns all files for name.
func (m *MultipartFormData) GetFiles(name string) []*multipart.FileHeader { return m.GetFileList(name) }

// GetFileList returns all files for name.
func (m *MultipartFormData) GetFileList(name string) []*multipart.FileHeader {
	if m == nil || m.Form == nil {
		return nil
	}
	return m.Form.File[name]
}

// GetFileParamNames returns file parameter names.
func (m *MultipartFormData) GetFileParamNames() []string {
	if m == nil || m.Form == nil {
		return nil
	}
	out := make([]string, 0, len(m.Form.File))
	for k := range m.Form.File {
		out = append(out, k)
	}
	return out
}

// GetFileMap returns first file values.
func (m *MultipartFormData) GetFileMap() map[string]*multipart.FileHeader {
	out := map[string]*multipart.FileHeader{}
	if m == nil || m.Form == nil {
		return out
	}
	for k, files := range m.Form.File {
		if len(files) > 0 {
			out[k] = files[0]
		}
	}
	return out
}

// GetFileListValueMap returns all file values.
func (m *MultipartFormData) GetFileListValueMap() map[string][]*multipart.FileHeader {
	if m == nil || m.Form == nil {
		return map[string][]*multipart.FileHeader{}
	}
	return m.Form.File
}

// SaveUploadedFile saves file to destPath.
func SaveUploadedFile(file *multipart.FileHeader, destPath string) error {
	if file == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(destPath), 0o750); err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()
	dst, err := os.Create(destPath) // #nosec G304 -- caller controls destination path.
	if err != nil {
		return err
	}
	defer func() { _ = dst.Close() }()
	_, err = io.Copy(dst, src)
	return err
}

// UploadFileName returns the uploaded file name.
func UploadFileName(file *multipart.FileHeader) string {
	if file == nil {
		return ""
	}
	return file.Filename
}

// UploadFileSize returns the uploaded file size.
func UploadFileSize(file *multipart.FileHeader) int64 {
	if file == nil {
		return 0
	}
	return file.Size
}

// UploadFileContentType returns the uploaded file content type header.
func UploadFileContentType(file *multipart.FileHeader) string {
	if file == nil {
		return ""
	}
	return file.Header.Get("Content-Type")
}

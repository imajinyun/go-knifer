package http

import (
	"bytes"
	"io"
	"mime/multipart"
)

// buildMultipartBody builds a multipart request body from form fields and files.
func buildMultipartBody(form map[string]any, files []*formFile, newWriter MultipartWriterFactory) (io.Reader, string, error) {
	buf := &bytes.Buffer{}
	if newWriter == nil {
		newWriter = newMultipartWriter
	}
	w := newWriter(buf)

	for k, v := range form {
		if err := w.WriteField(k, toString(v)); err != nil {
			return nil, "", NewHTTPError("write multipart field failed", err)
		}
	}
	for _, f := range files {
		fw, err := w.CreateFormFile(f.field, f.fileName)
		if err != nil {
			return nil, "", NewHTTPError("create multipart file failed", err)
		}
		if f.reader != nil {
			if _, err := io.Copy(fw, f.reader); err != nil {
				return nil, "", NewHTTPError("copy multipart file failed", err)
			}
		} else if _, err := fw.Write(f.data); err != nil {
			return nil, "", NewHTTPError("write multipart file failed", err)
		}
	}
	if err := w.Close(); err != nil {
		return nil, "", NewHTTPError("close multipart writer failed", err)
	}
	return buf, w.FormDataContentType(), nil
}

func newMultipartWriter(w io.Writer) MultipartWriter { return multipart.NewWriter(w) }

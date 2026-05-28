package http

import (
	"bytes"
	"io"
	"mime/multipart"
)

// buildMultipartBody 根据表单字段和文件构造 multipart 请求体。
func buildMultipartBody(form map[string]any, files []*formFile) (io.Reader, string, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)

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

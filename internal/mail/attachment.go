package mail

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"path/filepath"
	"strings"
)

// Attachment is a MIME file part.
type Attachment struct {
	Name        string
	ContentType ContentType
	ContentID   string
	Encoding    Encoding
	Size        int64
	Open        func() (io.ReadCloser, error)
	isInline    bool
}

// NewAttachment creates an attachment from bytes.
func NewAttachment(name string, content []byte, contentType ContentType) (Attachment, error) {
	return newBytesFile(name, content, contentType, "", false)
}

// NewInline creates an inline attachment from bytes with a Content-ID.
func NewInline(name, contentID string, content []byte, contentType ContentType) (Attachment, error) {
	return newBytesFile(name, content, contentType, contentID, true)
}

func newBytesFile(name string, content []byte, contentType ContentType, contentID string, inline bool) (Attachment, error) {
	name = strings.TrimSpace(name)
	if name == "" || hasCRLF(name) || hasCRLF(contentID) {
		return Attachment{}, ErrInvalidHeader
	}
	if contentType == "" {
		contentType = detectContentType(name)
	}
	data := append([]byte(nil), content...)
	return Attachment{
		Name:        filepath.Base(name),
		ContentType: contentType,
		ContentID:   strings.Trim(contentID, "<>"),
		Encoding:    EncodingBase64,
		Size:        int64(len(data)),
		Open: func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(data)), nil
		},
		isInline: inline,
	}, nil
}

func detectContentType(name string) ContentType {
	if ext := filepath.Ext(name); ext != "" {
		if contentType := mime.TypeByExtension(ext); contentType != "" {
			if mediaType, _, err := mime.ParseMediaType(contentType); err == nil {
				return ContentType(mediaType)
			}
			return ContentType(contentType)
		}
	}
	return TypeApplicationOctetStream
}

func (a Attachment) validate(maxBytes int64) error {
	if strings.TrimSpace(a.Name) == "" || hasCRLF(a.Name) || hasCRLF(a.ContentID) {
		return ErrInvalidHeader
	}
	if a.Open == nil {
		return errors.New("mail: missing attachment opener")
	}
	if maxBytes > 0 && a.Size > maxBytes {
		return ErrAttachmentTooLarge
	}
	return nil
}

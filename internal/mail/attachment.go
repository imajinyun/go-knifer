package mail

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"os"
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

// NewAttachmentReader creates an attachment from a reader opener.
func NewAttachmentReader(name string, size int64, contentType ContentType, open func() (io.ReadCloser, error)) (Attachment, error) {
	return newReaderFile(name, size, contentType, "", false, open)
}

// NewInlineReader creates an inline attachment from a reader opener with a Content-ID.
func NewInlineReader(
	name string,
	contentID string,
	size int64,
	contentType ContentType,
	open func() (io.ReadCloser, error),
) (Attachment, error) {
	return newReaderFile(name, size, contentType, contentID, true, open)
}

// NewAttachmentFile creates an attachment loaded lazily from path.
func NewAttachmentFile(path string) (Attachment, error) {
	return newFileAttachment(path, "", false)
}

// NewInlineFile creates an inline attachment loaded lazily from path with a Content-ID.
func NewInlineFile(path, contentID string) (Attachment, error) {
	return newFileAttachment(path, contentID, true)
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

func newReaderFile(
	name string,
	size int64,
	contentType ContentType,
	contentID string,
	inline bool,
	open func() (io.ReadCloser, error),
) (Attachment, error) {
	name = strings.TrimSpace(name)
	if name == "" || hasCRLF(name) || hasCRLF(contentID) {
		return Attachment{}, ErrInvalidHeader
	}
	if open == nil {
		return Attachment{}, errors.New("mail: missing attachment opener")
	}
	if contentType == "" {
		contentType = detectContentType(name)
	}
	return Attachment{
		Name:        filepath.Base(name),
		ContentType: contentType,
		ContentID:   strings.Trim(contentID, "<>"),
		Encoding:    EncodingBase64,
		Size:        size,
		Open:        open,
		isInline:    inline,
	}, nil
}

func newFileAttachment(path, contentID string, inline bool) (Attachment, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Attachment{}, err
	}
	return newReaderFile(filepath.Base(path), info.Size(), detectContentType(path), contentID, inline, func() (io.ReadCloser, error) {
		return os.Open(path) // #nosec G304 -- caller controls attachment path.
	})
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

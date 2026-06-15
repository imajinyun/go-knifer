package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"strings"
	"time"
)

func (m *Message) write(w io.Writer) error {
	if err := m.Validate(); err != nil {
		return err
	}
	header := strings.Builder{}
	writeHeaderLine(&header, "From", m.From.String())
	writeHeaderLine(&header, "To", formatAddressList(m.To))
	if len(m.Cc) > 0 {
		writeHeaderLine(&header, "Cc", formatAddressList(m.Cc))
	}
	if len(m.ReplyTo) > 0 {
		writeHeaderLine(&header, "Reply-To", formatAddressList(m.ReplyTo))
	}
	if m.Subject != "" {
		writeHeaderLine(&header, "Subject", mime.BEncoding.Encode(m.Charset.String(), m.Subject))
	}
	if !m.Date.IsZero() {
		writeHeaderLine(&header, "Date", m.Date.Format(time.RFC1123Z))
	}
	if m.MessageID != "" {
		writeHeaderLine(&header, "Message-ID", "<"+m.MessageID+">")
	}
	writeHeaderLine(&header, "MIME-Version", "1.0")
	for _, field := range m.Headers.fields {
		writeHeaderLine(&header, field.Name, field.Values...)
	}
	body, contentType, err := m.renderBody()
	if err != nil {
		return err
	}
	writeHeaderLine(&header, "Content-Type", contentType)
	header.WriteString("\r\n")
	if _, err := io.WriteString(w, header.String()); err != nil {
		return fmt.Errorf("write message headers: %w", err)
	}
	if _, err := body.WriteTo(w); err != nil {
		return fmt.Errorf("write message body: %w", err)
	}
	return nil
}

func (m *Message) renderBody() (*bytes.Buffer, string, error) {
	var body bytes.Buffer
	hasAttachments := len(m.Attachments) > 0
	hasInlines := len(m.Inlines) > 0
	if hasAttachments {
		mw, contentType, err := newMultipartWriter(&body, "mixed", m.boundaryGenerator)
		if err != nil {
			return nil, "", err
		}
		if err := m.writeRelatedOrAlternative(mw); err != nil {
			return nil, "", err
		}
		for _, attachment := range m.Attachments {
			if err := writeFilePart(mw, attachment, false, m.Charset); err != nil {
				return nil, "", err
			}
		}
		if err := mw.Close(); err != nil {
			return nil, "", fmt.Errorf("close mixed multipart: %w", err)
		}
		return &body, contentType, nil
	}
	if hasInlines {
		mw, contentType, err := newMultipartWriter(&body, "related", m.boundaryGenerator)
		if err != nil {
			return nil, "", err
		}
		if err := m.writeAlternative(mw); err != nil {
			return nil, "", err
		}
		for _, inline := range m.Inlines {
			if err := writeFilePart(mw, inline, true, m.Charset); err != nil {
				return nil, "", err
			}
		}
		if err := mw.Close(); err != nil {
			return nil, "", fmt.Errorf("close related multipart: %w", err)
		}
		return &body, contentType, nil
	}
	if m.Text != "" && m.HTML != "" {
		mw, contentType, err := newMultipartWriter(&body, "alternative", m.boundaryGenerator)
		if err != nil {
			return nil, "", err
		}
		if err := m.writeTextPart(mw, TypeTextPlain, m.Text); err != nil {
			return nil, "", err
		}
		if err := m.writeTextPart(mw, TypeTextHTML, m.HTML); err != nil {
			return nil, "", err
		}
		if err := mw.Close(); err != nil {
			return nil, "", fmt.Errorf("close alternative multipart: %w", err)
		}
		return &body, contentType, nil
	}
	contentType := TypeTextPlain
	content := m.Text
	if m.HTML != "" {
		contentType = TypeTextHTML
		content = m.HTML
	}
	encoded, err := encodeBytes([]byte(content), m.Encoding)
	if err != nil {
		return nil, "", err
	}
	body.Write(encoded)
	return &body, fmt.Sprintf(`%s; charset=%s`, contentType, m.Charset), nil
}

func (m *Message) writeRelatedOrAlternative(mw *multipart.Writer) error {
	if len(m.Inlines) == 0 {
		return m.writeAlternative(mw)
	}
	var body bytes.Buffer
	related, contentType, err := newMultipartWriter(&body, "related", m.boundaryGenerator)
	if err != nil {
		return err
	}
	if err := m.writeAlternative(related); err != nil {
		return err
	}
	for _, inline := range m.Inlines {
		if err := writeFilePart(related, inline, true, m.Charset); err != nil {
			return err
		}
	}
	if err := related.Close(); err != nil {
		return fmt.Errorf("close related multipart: %w", err)
	}
	part, err := mw.CreatePart(textproto.MIMEHeader{"Content-Type": {contentType}})
	if err != nil {
		return fmt.Errorf("create related part: %w", err)
	}
	_, err = body.WriteTo(part)
	return err
}

func (m *Message) writeAlternative(mw *multipart.Writer) error {
	if m.Text != "" && m.HTML != "" {
		var body bytes.Buffer
		alt, contentType, err := newMultipartWriter(&body, "alternative", m.boundaryGenerator)
		if err != nil {
			return err
		}
		if err := m.writeTextPart(alt, TypeTextPlain, m.Text); err != nil {
			return err
		}
		if err := m.writeTextPart(alt, TypeTextHTML, m.HTML); err != nil {
			return err
		}
		if err := alt.Close(); err != nil {
			return fmt.Errorf("close alternative multipart: %w", err)
		}
		part, err := mw.CreatePart(textproto.MIMEHeader{"Content-Type": {contentType}})
		if err != nil {
			return fmt.Errorf("create alternative part: %w", err)
		}
		_, err = body.WriteTo(part)
		return err
	}
	if m.HTML != "" {
		return m.writeTextPart(mw, TypeTextHTML, m.HTML)
	}
	return m.writeTextPart(mw, TypeTextPlain, m.Text)
}

func (m *Message) writeTextPart(mw *multipart.Writer, contentType ContentType, content string) error {
	header := textproto.MIMEHeader{}
	header.Set("Content-Type", fmt.Sprintf(`%s; charset=%s`, contentType, m.Charset))
	header.Set("Content-Transfer-Encoding", m.Encoding.String())
	part, err := mw.CreatePart(header)
	if err != nil {
		return fmt.Errorf("create text part: %w", err)
	}
	encoded, err := encodeBytes([]byte(content), m.Encoding)
	if err != nil {
		return err
	}
	_, err = part.Write(encoded)
	return err
}

func newMultipartWriter(w io.Writer, subtype string, generator BoundaryGenerator) (*multipart.Writer, string, error) {
	mw := multipart.NewWriter(w)
	boundary, err := generator()
	if err != nil {
		return nil, "", err
	}
	if err := mw.SetBoundary(boundary); err != nil {
		return nil, "", fmt.Errorf("set mime boundary: %w", err)
	}
	return mw, fmt.Sprintf(`multipart/%s; boundary=%q`, subtype, boundary), nil
}

func writeFilePart(mw *multipart.Writer, file Attachment, inline bool, charset Charset) error {
	header := textproto.MIMEHeader{}
	name := mime.BEncoding.Encode(charset.String(), file.Name)
	contentType := file.ContentType
	if contentType == "" {
		contentType = detectContentType(file.Name)
	}
	header.Set("Content-Type", fmt.Sprintf(`%s; name=%q`, contentType, name))
	header.Set("Content-Transfer-Encoding", EncodingBase64.String())
	disposition := "attachment"
	if inline {
		disposition = "inline"
	}
	header.Set("Content-Disposition", fmt.Sprintf(`%s; filename=%q`, disposition, name))
	if inline {
		contentID := file.ContentID
		if contentID == "" {
			contentID = file.Name
		}
		header.Set("Content-ID", "<"+strings.Trim(contentID, "<>")+">")
	}
	part, err := mw.CreatePart(header)
	if err != nil {
		return fmt.Errorf("create file part: %w", err)
	}
	r, err := file.Open()
	if err != nil {
		return fmt.Errorf("open file part: %w", err)
	}
	defer r.Close()
	encoder := base64.NewEncoder(base64.StdEncoding, newBase64LineWriter(part))
	if _, err := io.Copy(encoder, r); err != nil {
		return fmt.Errorf("encode file part: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("close file encoder: %w", err)
	}
	return nil
}

func encodeBytes(data []byte, encoding Encoding) ([]byte, error) {
	var buf bytes.Buffer
	switch encoding {
	case EncodingQuotedPrintable:
		w := quotedprintable.NewWriter(&buf)
		if _, err := w.Write(data); err != nil {
			return nil, fmt.Errorf("encode quoted-printable body: %w", err)
		}
		if err := w.Close(); err != nil {
			return nil, fmt.Errorf("close quoted-printable body: %w", err)
		}
	case EncodingBase64:
		w := base64.NewEncoder(base64.StdEncoding, newBase64LineWriter(&buf))
		if _, err := w.Write(data); err != nil {
			return nil, fmt.Errorf("encode base64 body: %w", err)
		}
		if err := w.Close(); err != nil {
			return nil, fmt.Errorf("close base64 body: %w", err)
		}
	case Encoding7Bit, Encoding8Bit:
		buf.Write(data)
	default:
		return nil, fmt.Errorf("%w: unsupported encoding %q", ErrInvalidHeader, encoding)
	}
	return buf.Bytes(), nil
}

type base64LineWriter struct {
	out  io.Writer
	line int
}

func newBase64LineWriter(out io.Writer) *base64LineWriter { return &base64LineWriter{out: out} }

func (w *base64LineWriter) Write(p []byte) (int, error) {
	written := 0
	for _, b := range p {
		if w.line == 76 {
			if _, err := w.out.Write([]byte("\r\n")); err != nil {
				return written, err
			}
			w.line = 0
		}
		if _, err := w.out.Write([]byte{b}); err != nil {
			return written, err
		}
		w.line++
		written++
	}
	return written, nil
}

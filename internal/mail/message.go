package mail

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const defaultMaxAttachmentBytes int64 = 25 << 20

// BoundaryGenerator creates MIME multipart boundaries.
type BoundaryGenerator func() (string, error)

// MessageOption customizes message construction.
type MessageOption func(*Message) error

// Message is an email message with text, HTML, inline files, and attachments.
type Message struct {
	From        *Address
	To          []*Address
	Cc          []*Address
	Bcc         []*Address
	ReplyTo     []*Address
	Subject     string
	Text        string
	HTML        string
	Headers     Header
	Attachments []Attachment
	Inlines     []Attachment
	Date        time.Time
	MessageID   string
	Charset     Charset
	Encoding    Encoding

	maxAttachmentBytes int64
	boundaryGenerator  BoundaryGenerator
}

// NewMessage creates and validates an email message.
func NewMessage(opts ...MessageOption) (*Message, error) {
	m := &Message{
		Date:               time.Now(),
		Charset:            CharsetUTF8,
		Encoding:           EncodingQuotedPrintable,
		maxAttachmentBytes: defaultMaxAttachmentBytes,
		boundaryGenerator:  randomBoundary,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(m); err != nil {
			return nil, err
		}
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}

// Validate checks message invariants before rendering or sending.
func (m *Message) Validate() error {
	if m == nil {
		return ErrMissingBody
	}
	if m.From == nil {
		return ErrMissingFrom
	}
	if len(m.To)+len(m.Cc)+len(m.Bcc) == 0 {
		return ErrMissingRecipient
	}
	if m.Text == "" && m.HTML == "" {
		return ErrMissingBody
	}
	for _, addr := range append(append(cloneAddresses(m.To), m.Cc...), m.Bcc...) {
		if addr == nil || validateEmail(addr.Email) != nil || hasCRLF(addr.Name) {
			return ErrInvalidAddress
		}
	}
	if validateEmail(m.From.Email) != nil || hasCRLF(m.From.Name) || hasCRLF(m.Subject) || hasCRLF(m.MessageID) {
		return ErrInvalidHeader
	}
	for _, attachment := range m.Attachments {
		if err := attachment.validate(m.maxAttachmentBytes); err != nil {
			return err
		}
	}
	for _, inline := range m.Inlines {
		if err := inline.validate(m.maxAttachmentBytes); err != nil {
			return err
		}
	}
	return nil
}

// Recipients returns all SMTP envelope recipients, including Bcc.
func (m *Message) Recipients() []string {
	addresses := append(append(cloneAddresses(m.To), m.Cc...), m.Bcc...)
	recipients := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		if addr != nil {
			recipients = append(recipients, addr.Email)
		}
	}
	return recipients
}

// WriteTo renders the message to w.
func (m *Message) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer
	if err := m.write(&buf); err != nil {
		return 0, err
	}
	n, err := w.Write(buf.Bytes())
	return int64(n), err
}

// Bytes renders the message to a byte slice.
func (m *Message) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	if _, err := m.WriteTo(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// WithFrom sets the From address.
func WithFrom(address string) MessageOption {
	return func(m *Message) error {
		addr, err := ParseAddress(address)
		if err != nil {
			return err
		}
		m.From = addr
		return nil
	}
}

// WithFromAddress sets the From address.
func WithFromAddress(addr *Address) MessageOption {
	return func(m *Message) error {
		if addr == nil {
			return ErrInvalidAddress
		}
		copyAddr := *addr
		m.From = &copyAddr
		return nil
	}
}

// WithTo appends To recipients.
func WithTo(addresses ...string) MessageOption {
	return appendAddressOption(func(m *Message, a []*Address) { m.To = append(m.To, a...) }, addresses...)
}

// WithCc appends Cc recipients.
func WithCc(addresses ...string) MessageOption {
	return appendAddressOption(func(m *Message, a []*Address) { m.Cc = append(m.Cc, a...) }, addresses...)
}

// WithBcc appends Bcc recipients.
func WithBcc(addresses ...string) MessageOption {
	return appendAddressOption(func(m *Message, a []*Address) { m.Bcc = append(m.Bcc, a...) }, addresses...)
}

// WithReplyTo appends Reply-To addresses.
func WithReplyTo(addresses ...string) MessageOption {
	return appendAddressOption(func(m *Message, a []*Address) { m.ReplyTo = append(m.ReplyTo, a...) }, addresses...)
}

// WithSubject sets the Subject header.
func WithSubject(subject string) MessageOption {
	return func(m *Message) error {
		if hasCRLF(subject) {
			return ErrInvalidHeader
		}
		m.Subject = subject
		return nil
	}
}

// WithText sets the text/plain body.
func WithText(text string) MessageOption { return func(m *Message) error { m.Text = text; return nil } }

// WithHTML sets the text/html body.
func WithHTML(html string) MessageOption { return func(m *Message) error { m.HTML = html; return nil } }

// WithHeader appends an additional header.
func WithHeader(name string, values ...string) MessageOption {
	return func(m *Message) error { return m.Headers.Add(name, values...) }
}

// WithAttachment appends an attachment from bytes.
func WithAttachment(name string, content []byte, contentType ContentType) MessageOption {
	return func(m *Message) error {
		attachment, err := NewAttachment(name, content, contentType)
		if err != nil {
			return err
		}
		m.Attachments = append(m.Attachments, attachment)
		return nil
	}
}

// WithInline appends an inline file from bytes.
func WithInline(name, contentID string, content []byte, contentType ContentType) MessageOption {
	return func(m *Message) error {
		inline, err := NewInline(name, contentID, content, contentType)
		if err != nil {
			return err
		}
		m.Inlines = append(m.Inlines, inline)
		return nil
	}
}

// WithAttachmentFile appends an attachment loaded lazily from path.
func WithAttachmentFile(path string) MessageOption {
	return func(m *Message) error {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat attachment: %w", err)
		}
		attachment := Attachment{
			Name:        filepath.Base(path),
			ContentType: detectContentType(path),
			Encoding:    EncodingBase64,
			Size:        info.Size(),
			Open: func() (io.ReadCloser, error) {
				return os.Open(path) // #nosec G304 -- caller controls attachment path.
			},
		}
		m.Attachments = append(m.Attachments, attachment)
		return nil
	}
}

// WithDate sets the Date header value.
func WithDate(t time.Time) MessageOption { return func(m *Message) error { m.Date = t; return nil } }

// WithMessageID sets the Message-ID header without angle brackets.
func WithMessageID(id string) MessageOption {
	return func(m *Message) error {
		if hasCRLF(id) {
			return ErrInvalidHeader
		}
		m.MessageID = strings.Trim(id, "<>")
		return nil
	}
}

// WithCharset sets the charset used for message text parts and encoded headers.
func WithCharset(charset Charset) MessageOption {
	return func(m *Message) error {
		if strings.TrimSpace(charset.String()) == "" || hasCRLF(charset.String()) {
			return ErrInvalidHeader
		}
		m.Charset = charset
		return nil
	}
}

// WithEncoding sets the transfer encoding used for text parts.
func WithEncoding(encoding Encoding) MessageOption {
	return func(m *Message) error {
		if strings.TrimSpace(encoding.String()) == "" || hasCRLF(encoding.String()) {
			return ErrInvalidHeader
		}
		m.Encoding = encoding
		return nil
	}
}

// WithMaxAttachmentBytes sets the per-attachment size limit. A non-positive value disables the limit.
func WithMaxAttachmentBytes(maxBytes int64) MessageOption {
	return func(m *Message) error { m.maxAttachmentBytes = maxBytes; return nil }
}

// WithBoundaryGenerator injects the MIME boundary generator.
func WithBoundaryGenerator(generator BoundaryGenerator) MessageOption {
	return func(m *Message) error {
		if generator == nil {
			return errors.New("mail: nil boundary generator")
		}
		m.boundaryGenerator = generator
		return nil
	}
}

func appendAddressOption(set func(*Message, []*Address), addresses ...string) MessageOption {
	return func(m *Message) error {
		for _, raw := range addresses {
			parsed, err := ParseAddressList(raw)
			if err != nil {
				return err
			}
			set(m, parsed)
		}
		return nil
	}
}

func randomBoundary() (string, error) {
	var buf [24]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", fmt.Errorf("generate mime boundary: %w", err)
	}
	return hex.EncodeToString(buf[:]), nil
}

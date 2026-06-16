package mail

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewMessageRendersMixedRelatedAlternative(t *testing.T) {
	msg, err := NewMessage(
		WithFrom("Sender <sender@example.com>"),
		WithTo("Receiver <receiver@example.com>"),
		WithCc("copy@example.com"),
		WithBcc("hidden@example.com"),
		WithSubject("hello 世界"),
		WithText("plain body"),
		WithHTML(`<p><img src="cid:logo">html body</p>`),
		WithInline("logo.png", "logo", []byte("inline-data"), "image/png"),
		WithAttachment("report.txt", []byte("attachment-data"), TypeTextPlain),
		WithDate(time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)),
		WithMessageID("message@example.com"),
		WithBoundaryGenerator(sequenceBoundary("mixed-boundary", "related-boundary", "alternative-boundary")),
	)
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}

	var buf bytes.Buffer
	if _, err := msg.WriteTo(&buf); err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}
	raw := buf.String()

	assertContains(t, raw, "From: \"Sender\" <sender@example.com>\r\n")
	assertContains(t, raw, "To: \"Receiver\" <receiver@example.com>\r\n")
	assertContains(t, raw, "Cc: <copy@example.com>\r\n")
	assertContains(t, raw, "Subject: =?UTF-8?b?")
	assertContains(t, raw, "Message-ID: <message@example.com>\r\n")
	assertContains(t, raw, `Content-Type: multipart/mixed; boundary="mixed-boundary"`)
	assertContains(t, raw, `Content-Type: multipart/related; boundary="related-boundary"`)
	assertContains(t, raw, `Content-Type: multipart/alternative; boundary="alternative-boundary"`)
	assertContains(t, raw, "Content-Id: <logo>")
	assertContains(t, raw, "Content-Disposition: attachment;")
	if strings.Contains(raw, "Bcc:") {
		t.Fatalf("rendered message leaked Bcc header:\n%s", raw)
	}

	recipients := strings.Join(msg.Recipients(), ",")
	assertContains(t, recipients, "receiver@example.com")
	assertContains(t, recipients, "copy@example.com")
	assertContains(t, recipients, "hidden@example.com")
}

func TestNewMessageValidation(t *testing.T) {
	tests := []struct {
		name     string
		opts     []MessageOption
		expected error
	}{
		{
			name:     "missing from",
			opts:     []MessageOption{WithTo("to@example.com"), WithText("body")},
			expected: ErrMissingFrom,
		},
		{
			name:     "missing recipient",
			opts:     []MessageOption{WithFrom("from@example.com"), WithText("body")},
			expected: ErrMissingRecipient,
		},
		{
			name:     "missing body",
			opts:     []MessageOption{WithFrom("from@example.com"), WithTo("to@example.com")},
			expected: ErrMissingBody,
		},
		{
			name: "header injection",
			opts: []MessageOption{
				WithFrom("from@example.com"),
				WithTo("to@example.com"),
				WithSubject("safe\r\nBcc: attacker@example.com"),
				WithText("body"),
			},
			expected: ErrInvalidHeader,
		},
		{
			name: "attachment too large",
			opts: []MessageOption{
				WithFrom("from@example.com"),
				WithTo("to@example.com"),
				WithText("body"),
				WithMaxAttachmentBytes(1),
				WithAttachment("big.txt", []byte("too big"), TypeTextPlain),
			},
			expected: ErrAttachmentTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMessage(tt.opts...)
			if !errors.Is(err, tt.expected) {
				t.Fatalf("NewMessage() error = %v, want %v", err, tt.expected)
			}
		})
	}
}

func TestParseAddressListRejectsCRLF(t *testing.T) {
	_, err := ParseAddressList("to@example.com\nBcc: attacker@example.com")
	if !errors.Is(err, ErrInvalidAddress) {
		t.Fatalf("ParseAddressList() error = %v, want %v", err, ErrInvalidAddress)
	}
}

func TestAddressAndHeaderHelpers(t *testing.T) {
	addr, err := NewAddress("Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("NewAddress() error = %v", err)
	}
	if addr.String() != `"Alice" <alice@example.com>` {
		t.Fatalf("Address.String() = %q", addr.String())
	}
	if empty := (*Address)(nil).String(); empty != "" {
		t.Fatalf("nil Address.String() = %q, want empty", empty)
	}
	if _, err := NewAddress("bad\nname", "alice@example.com"); !errors.Is(err, ErrInvalidAddress) {
		t.Fatalf("NewAddress(CRLF) error = %v, want %v", err, ErrInvalidAddress)
	}
	if _, err := NewAddress("Alice", "Alice <alice@example.com>"); !errors.Is(err, ErrInvalidAddress) {
		t.Fatalf("NewAddress(display email) error = %v, want %v", err, ErrInvalidAddress)
	}

	var header Header
	if err := header.Add("X-Test", "one"); err != nil {
		t.Fatalf("Header.Add() error = %v", err)
	}
	if err := header.Set("X-Test", "two", "three"); err != nil {
		t.Fatalf("Header.Set(existing) error = %v", err)
	}
	if err := header.Set("X-Other", "value"); err != nil {
		t.Fatalf("Header.Set(new) error = %v", err)
	}
	values := header.Values("x-test")
	if strings.Join(values, ",") != "two,three" {
		t.Fatalf("Header.Values() = %v", values)
	}
	values[0] = "mutated"
	if got := header.Values("X-Test")[0]; got != "two" {
		t.Fatalf("Header.Values() returned mutable slice, got %q", got)
	}
	if got := header.Values("missing"); got != nil {
		t.Fatalf("Header.Values(missing) = %v, want nil", got)
	}
	if err := header.Add("Bad:Name", "value"); !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("Header.Add(bad name) error = %v, want %v", err, ErrInvalidHeader)
	}
	if err := header.Add("X-Bad", "value\r\nInjected: true"); !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("Header.Add(bad value) error = %v, want %v", err, ErrInvalidHeader)
	}
}

func TestMessageOptionsAndRenderingPaths(t *testing.T) {
	from := &Address{Name: "Sender", Email: "sender@example.com"}
	msg, err := NewMessage(
		WithFromAddress(from),
		WithTo("to@example.com"),
		WithReplyTo("reply@example.com"),
		WithSubject("html"),
		WithHTML("<strong>Hello</strong>"),
		WithHeader("X-Custom", "a", "b"),
		WithCharset(CharsetASCII),
		WithEncoding(EncodingBase64),
	)
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	from.Email = "changed@example.com"
	raw, err := msg.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	text := string(raw)
	assertContains(t, text, `From: "Sender" <sender@example.com>`)
	assertContains(t, text, "Reply-To: <reply@example.com>\r\n")
	assertContains(t, text, "X-Custom: a, b\r\n")
	assertContains(t, text, "Content-Type: text/html; charset=US-ASCII\r\n")
	assertContains(t, text, "PHN0cm9uZz5IZWxsbzwvc3Ryb25nPg==")
}

func TestMessageAttachmentFileAndInlineOnlyRendering(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.txt")
	if err := os.WriteFile(path, []byte("file-body"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	msg, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("plain"),
		WithAttachmentFile(path),
		WithBoundaryGenerator(sequenceBoundary("mixed-file")),
	)
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	raw, err := msg.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	text := string(raw)
	assertContains(t, text, `Content-Type: multipart/mixed; boundary="mixed-file"`)
	assertContains(t, text, `Content-Type: text/plain; name=report.txt`)
	assertContains(t, text, "ZmlsZS1ib2R5")

	inline, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithHTML(`<img src="cid:logo.png">`),
		WithInline("logo.png", "", []byte("inline"), ""),
		WithBoundaryGenerator(sequenceBoundary("related-only")),
	)
	if err != nil {
		t.Fatalf("NewMessage(inline) error = %v", err)
	}
	raw, err = inline.Bytes()
	if err != nil {
		t.Fatalf("Bytes(inline) error = %v", err)
	}
	text = string(raw)
	assertContains(t, text, `Content-Type: multipart/related; boundary="related-only"`)
	assertContains(t, text, "Content-Id: <logo.png>")
	assertContains(t, text, `Content-Type: image/png; name=logo.png`)
}

func TestAttachmentHeaderParametersUseRFC2231Encoding(t *testing.T) {
	msg, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("plain"),
		WithAttachment("报告 2026(最终).txt", []byte("report"), TypeTextPlain),
		WithInline("内嵌 logo.png", "logo", []byte("inline"), "image/png"),
		WithBoundaryGenerator(sequenceBoundary("mixed-rfc2231", "related-rfc2231")),
	)
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	raw, err := msg.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	text := string(raw)
	assertContains(t, text, `Content-Type: text/plain; name*=utf-8''`)
	assertContains(t, text, `Content-Disposition: attachment; filename*=utf-8''`)
	assertContains(t, text, `Content-Type: image/png; name*=utf-8''`)
	assertContains(t, text, `Content-Disposition: inline; filename*=utf-8''`)
	if strings.Contains(text, `name="=?UTF-8?`) || strings.Contains(text, `filename="=?UTF-8?`) {
		t.Fatalf("file parameters used encoded-word syntax instead of RFC 2231 encoding:\n%s", text)
	}
}

func TestAttachmentReaderAndInlineFileRendering(t *testing.T) {
	dir := t.TempDir()
	logoPath := filepath.Join(dir, "logo.png")
	if err := os.WriteFile(logoPath, []byte("logo-data"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	readerAttachment, err := NewAttachmentReader("reader.txt", 11, TypeTextPlain, func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("reader-data")), nil
	})
	if err != nil {
		t.Fatalf("NewAttachmentReader() error = %v", err)
	}
	if readerAttachment.Name != "reader.txt" || readerAttachment.Size != 11 {
		t.Fatalf("reader attachment = %#v", readerAttachment)
	}
	if _, err := NewAttachmentReader("reader.txt", 1, TypeTextPlain, nil); err == nil {
		t.Fatal("NewAttachmentReader(nil opener) error = nil, want error")
	}
	inlineFile, err := NewInlineFile(logoPath, "logo-id")
	if err != nil {
		t.Fatalf("NewInlineFile() error = %v", err)
	}
	if inlineFile.Name != "logo.png" || inlineFile.ContentID != "logo-id" {
		t.Fatalf("inline file = %#v", inlineFile)
	}

	msg, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithHTML(`<img src="cid:logo-id">`),
		WithAttachmentReader("reader.txt", 11, TypeTextPlain, func() (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("reader-data")), nil
		}),
		WithInlineFile(logoPath, "logo-id"),
		WithBoundaryGenerator(sequenceBoundary("mixed-reader", "related-reader")),
	)
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	raw, err := msg.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	text := string(raw)
	assertContains(t, text, `Content-Type: multipart/mixed; boundary="mixed-reader"`)
	assertContains(t, text, `Content-Type: multipart/related; boundary="related-reader"`)
	assertContains(t, text, "cmVhZGVyLWRhdGE=")
	assertContains(t, text, "Content-Id: <logo-id>")
	assertContains(t, text, "bG9nby1kYXRh")
}

func TestMessageOptionValidationErrors(t *testing.T) {
	if _, err := NewMessage(WithFromAddress(nil)); !errors.Is(err, ErrInvalidAddress) {
		t.Fatalf("WithFromAddress(nil) error = %v, want %v", err, ErrInvalidAddress)
	}
	if _, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("body"),
		WithHeader("X-Bad", "bad\nvalue"),
	); !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("WithHeader(CRLF) error = %v, want %v", err, ErrInvalidHeader)
	}
	if _, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("body"),
		WithMessageID("id\nnext"),
	); !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("WithMessageID(CRLF) error = %v, want %v", err, ErrInvalidHeader)
	}
	if _, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("body"),
		WithCharset(Charset("bad\ncharset")),
	); !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("WithCharset(CRLF) error = %v, want %v", err, ErrInvalidHeader)
	}
	if _, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("body"),
		WithEncoding(Encoding("bad\nencoding")),
	); !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("WithEncoding(CRLF) error = %v, want %v", err, ErrInvalidHeader)
	}
	if _, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("body"),
		WithBoundaryGenerator(nil),
	); err == nil {
		t.Fatal("WithBoundaryGenerator(nil) error = nil, want error")
	}
	if _, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("body"),
		WithEnvelopeFrom("bad address"),
	); !errors.Is(err, ErrInvalidAddress) {
		t.Fatalf("WithEnvelopeFrom(invalid) error = %v, want %v", err, ErrInvalidAddress)
	}
	if _, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("body"),
		WithAttachmentFile(filepath.Join(t.TempDir(), "missing.txt")),
	); err == nil {
		t.Fatal("WithAttachmentFile(missing) error = nil, want error")
	}
	if _, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("body"),
		WithInlineFile(filepath.Join(t.TempDir(), "missing.png"), "logo"),
	); err == nil {
		t.Fatal("WithInlineFile(missing) error = nil, want error")
	}
	if _, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("body"),
		WithAttachment("bad\nname.txt", []byte("x"), TypeTextPlain),
	); !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("WithAttachment(CRLF) error = %v, want %v", err, ErrInvalidHeader)
	}
}

func TestMessageEncodingAndBoundaryErrors(t *testing.T) {
	for _, tt := range []struct {
		name     string
		encoding Encoding
		body     string
	}{
		{name: "seven bit", encoding: Encoding7Bit, body: "plain"},
		{name: "eight bit", encoding: Encoding8Bit, body: "héllo"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := NewMessage(
				WithFrom("from@example.com"),
				WithTo("to@example.com"),
				WithText(tt.body),
				WithEncoding(tt.encoding),
			)
			if err != nil {
				t.Fatalf("NewMessage() error = %v", err)
			}
			raw, err := msg.Bytes()
			if err != nil {
				t.Fatalf("Bytes() error = %v", err)
			}
			assertContains(t, string(raw), tt.body)
		})
	}

	msg, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("plain"),
		WithHTML("<p>html</p>"),
	)
	if err != nil {
		t.Fatalf("NewMessage(alternative) error = %v", err)
	}
	raw, err := msg.Bytes()
	if err != nil {
		t.Fatalf("Bytes(alternative) error = %v", err)
	}
	assertContains(t, string(raw), "Content-Type: multipart/alternative;")

	msg.Encoding = Encoding("x-bad")
	if _, err := msg.Bytes(); !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("Bytes(unsupported encoding) error = %v, want %v", err, ErrInvalidHeader)
	}

	badBoundary, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("plain"),
		WithHTML("<p>html</p>"),
		WithBoundaryGenerator(func() (string, error) { return "bad\r\nboundary", nil }),
	)
	if err != nil {
		t.Fatalf("NewMessage(bad boundary) error = %v", err)
	}
	if _, err := badBoundary.Bytes(); err == nil {
		t.Fatal("Bytes(bad boundary) error = nil, want error")
	}
}

func sequenceBoundary(values ...string) BoundaryGenerator {
	idx := 0
	return func() (string, error) {
		if idx >= len(values) {
			return "extra-boundary", nil
		}
		value := values[idx]
		idx++
		return value, nil
	}
}

func assertContains(t *testing.T, got, expected string) {
	t.Helper()
	if !strings.Contains(got, expected) {
		t.Fatalf("expected %q to contain %q", got, expected)
	}
}

package mail

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

func TestNewInlineReader(t *testing.T) {
	inline, err := NewInlineReader("inline.txt", "cid-123", 11, TypeTextPlain, func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("inline-data")), nil
	})
	if err != nil {
		t.Fatalf("NewInlineReader() error = %v", err)
	}
	if inline.Name != "inline.txt" || inline.ContentID != "cid-123" {
		t.Fatalf("inline = %#v", inline)
	}
}

func TestWithInlineReader(t *testing.T) {
	msg, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithHTML(`<img src="cid:cid-123">`),
		WithInlineReader("inline.txt", "cid-123", 11, TypeTextPlain, func() (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("inline-data")), nil
		}),
		WithBoundaryGenerator(sequenceBoundary("with-inline-reader")),
	)
	if err != nil {
		t.Fatalf("NewMessage() with WithInlineReader error = %v", err)
	}
	raw, err := msg.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	text := string(raw)
	assertContains(t, text, "Content-Id: <cid-123>")
	assertContains(t, text, "aW5saW5lLWRhdGE=")
}

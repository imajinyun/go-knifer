package mail

import (
	"bytes"
	"errors"
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

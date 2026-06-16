package mail

import (
	"errors"
	"path/filepath"
	"testing"
)

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

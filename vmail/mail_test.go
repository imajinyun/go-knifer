package vmail

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewMessageFacade(t *testing.T) {
	message, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithSubject("hello"),
		WithText("body"),
	)
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	raw, err := message.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	if !strings.Contains(string(raw), "Subject: hello") {
		t.Fatalf("message did not contain encoded subject: %s", raw)
	}
}

func TestSendTextUsesInjectedProvider(t *testing.T) {
	var got *Message
	err := SendText(
		context.Background(),
		"smtp.example.com",
		587,
		"from@example.com",
		[]string{"to@example.com"},
		"subject",
		"body",
		WithSenderProvider(func(config Config) (Sender, error) {
			if config.Host != "smtp.example.com" || config.Port != 587 {
				t.Fatalf("Config = %#v", config)
			}
			return SenderFunc(func(ctx context.Context, message *Message) error {
				got = message
				return nil
			}), nil
		}),
	)
	if err != nil {
		t.Fatalf("SendText() error = %v", err)
	}
	if got == nil || got.Subject != "subject" || got.Text != "body" {
		t.Fatalf("sent message = %#v", got)
	}
}

func TestFacadeExportsSentinelErrors(t *testing.T) {
	_, err := NewMessage(WithTo("to@example.com"), WithText("body"))
	if !errors.Is(err, ErrMissingFrom) {
		t.Fatalf("NewMessage() error = %v, want %v", err, ErrMissingFrom)
	}
}

func TestFacadeConstructorsAndMessageOptions(t *testing.T) {
	addr, err := NewAddress("Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("NewAddress() error = %v", err)
	}
	if _, err := ParseAddress(addr.String()); err != nil {
		t.Fatalf("ParseAddress() error = %v", err)
	}
	list, err := ParseAddressList("bob@example.com, carol@example.com")
	if err != nil {
		t.Fatalf("ParseAddressList() error = %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("ParseAddressList() len = %d, want 2", len(list))
	}
	attachment, err := NewAttachment("report.txt", []byte("report"), TypeTextPlain)
	if err != nil {
		t.Fatalf("NewAttachment() error = %v", err)
	}
	inline, err := NewInline("logo.png", "logo", []byte("inline"), "")
	if err != nil {
		t.Fatalf("NewInline() error = %v", err)
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "extra.txt")
	if err := os.WriteFile(path, []byte("extra"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	message, err := NewMessage(
		WithFromAddress(addr),
		WithTo("to@example.com"),
		WithCc("cc@example.com"),
		WithBcc("bcc@example.com"),
		WithReplyTo("reply@example.com"),
		WithSubject("facade"),
		WithText("plain"),
		WithHTML("<b>html</b>"),
		WithHeader("X-Facade", "yes"),
		WithAttachment(attachment.Name, []byte("report"), attachment.ContentType),
		WithInline(inline.Name, inline.ContentID, []byte("inline"), inline.ContentType),
		WithAttachmentFile(path),
		WithDate(time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)),
		WithMessageID("facade@example.com"),
		WithCharset(CharsetUTF8),
		WithEncoding(EncodingQuotedPrintable),
		WithMaxAttachmentBytes(1024),
		WithBoundaryGenerator(sequenceBoundary("mixed", "related", "alternative")),
	)
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	raw, err := message.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}
	text := string(raw)
	assertContains(t, text, "X-Facade: yes")
	assertContains(t, text, "Message-ID: <facade@example.com>")
	assertContains(t, text, `Content-Type: multipart/mixed; boundary="mixed"`)
}

func TestFacadeSendAndClientOptions(t *testing.T) {
	message, err := NewMessage(WithFrom("from@example.com"), WithTo("to@example.com"), WithText("body"))
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	var got *Message
	auth := facadeSMTPAuth{mechanism: "CUSTOM"}
	provider := func(config Config) (Sender, error) {
		if config.Username != "user" || config.Password != "pass" || !config.AllowPlainAuth {
			t.Fatalf("Config auth = %#v", config)
		}
		if config.Auth == nil {
			t.Fatal("Config Auth is nil")
		}
		if config.TLSPolicy != TLSNone || config.LocalName != "mail.local" || config.Timeout != time.Second {
			t.Fatalf("Config transport = %#v", config)
		}
		if config.TLSConfig == nil || config.TLSConfig.ServerName != "smtp.example.com" {
			t.Fatalf("Config TLS = %#v", config.TLSConfig)
		}
		return SenderFunc(func(ctx context.Context, message *Message) error {
			got = message
			return nil
		}), nil
	}
	if err := Send(context.Background(), "smtp.example.com", 587, message,
		WithAuth("user", "pass"),
		WithSMTPAuth(auth),
		WithTLSConfig(&tls.Config{ServerName: "smtp.example.com", MinVersion: tls.VersionTLS12}),
		WithTLSPolicy(TLSNone),
		WithAllowPlainAuth(true),
		WithTimeout(time.Second),
		WithLocalName("mail.local"),
		WithDialContext(func(context.Context, string, string) (net.Conn, error) { return nil, errors.New("unused") }),
		WithSenderProvider(provider),
	); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if got != message {
		t.Fatalf("sent message = %p, want %p", got, message)
	}

	got = nil
	if err := SendHTML(context.Background(), "smtp.example.com", 587, "from@example.com", []string{"to@example.com"}, "subject", "<p>html</p>", WithSenderProvider(func(Config) (Sender, error) {
		return SenderFunc(func(ctx context.Context, message *Message) error {
			got = message
			return nil
		}), nil
	})); err != nil {
		t.Fatalf("SendHTML() error = %v", err)
	}
	if got == nil || got.HTML != "<p>html</p>" {
		t.Fatalf("SendHTML message = %#v", got)
	}

	client, err := NewClient("smtp.example.com", 587, WithSenderProvider(provider))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}
}

func sequenceBoundary(values ...string) BoundaryGenerator {
	idx := 0
	return func() (string, error) {
		value := values[idx]
		idx++
		return value, nil
	}
}

type facadeSMTPAuth struct{ mechanism string }

func (a facadeSMTPAuth) Start(*smtp.ServerInfo) (string, []byte, error) { return a.mechanism, nil, nil }

func (a facadeSMTPAuth) Next([]byte, bool) ([]byte, error) { return nil, nil }

func assertContains(t *testing.T, got, expected string) {
	t.Helper()
	if !strings.Contains(got, expected) {
		t.Fatalf("expected %q to contain %q", got, expected)
	}
}

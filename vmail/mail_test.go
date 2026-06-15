package vmail

import (
	"context"
	"errors"
	"strings"
	"testing"
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

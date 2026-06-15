package mail

import (
	"context"
	"errors"
	"testing"
)

func TestClientSendUsesInjectedSender(t *testing.T) {
	msg, err := NewMessage(
		WithFrom("from@example.com"),
		WithTo("to@example.com"),
		WithText("body"),
	)
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}

	var called bool
	client, err := NewClient("smtp.example.com", 587, WithSenderProvider(func(config Config) (Sender, error) {
		if config.Host != "smtp.example.com" || config.Port != 587 {
			t.Fatalf("Config = %#v", config)
		}
		return SenderFunc(func(ctx context.Context, got *Message) error {
			called = true
			if got != msg {
				t.Fatalf("Send message = %p, want %p", got, msg)
			}
			return nil
		}), nil
	}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if !called {
		t.Fatal("injected sender was not called")
	}
}

func TestClientSendValidatesMessageBeforeProvider(t *testing.T) {
	client, err := NewClient("smtp.example.com", 587, WithSenderProvider(func(Config) (Sender, error) {
		t.Fatal("provider should not be called for invalid message")
		return nil, nil
	}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	err = client.Send(context.Background(), &Message{})
	if !errors.Is(err, ErrMissingFrom) {
		t.Fatalf("Send() error = %v, want %v", err, ErrMissingFrom)
	}
}

func TestNewClientRejectsBadConfig(t *testing.T) {
	tests := []struct {
		name string
		host string
		port int
		opts []ClientOption
	}{
		{name: "empty host", host: "", port: 587},
		{name: "empty port", host: "smtp.example.com", port: 0},
		{name: "bad local name", host: "smtp.example.com", port: 587, opts: []ClientOption{WithLocalName("ok\nBAD")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.host, tt.port, tt.opts...)
			if err == nil {
				t.Fatal("NewClient() error = nil, want error")
			}
		})
	}
}

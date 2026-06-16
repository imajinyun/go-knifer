package mail

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/smtp"
	"strings"
	"sync"
	"testing"
	"time"
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

func TestSendTextAndHTMLConvenienceFunctions(t *testing.T) {
	for _, tt := range []struct {
		name     string
		send     func(context.Context, SenderProvider) error
		wantBody string
	}{
		{
			name: "text",
			send: func(ctx context.Context, provider SenderProvider) error {
				return SendText(ctx, "smtp.example.com", 587, "from@example.com", []string{"to@example.com"}, "subject", "plain", WithSenderProvider(provider))
			},
			wantBody: "plain",
		},
		{
			name: "html",
			send: func(ctx context.Context, provider SenderProvider) error {
				return SendHTML(ctx, "smtp.example.com", 587, "from@example.com", []string{"to@example.com"}, "subject", "<b>html</b>", WithSenderProvider(provider))
			},
			wantBody: "<b>html</b>",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var got *Message
			provider := func(config Config) (Sender, error) {
				if config.Host != "smtp.example.com" || config.Port != 587 {
					t.Fatalf("Config = %#v", config)
				}
				return SenderFunc(func(ctx context.Context, message *Message) error {
					got = message
					return nil
				}), nil
			}
			if err := tt.send(context.Background(), provider); err != nil {
				t.Fatalf("send() error = %v", err)
			}
			if got == nil || !strings.Contains(got.Text+got.HTML, tt.wantBody) {
				t.Fatalf("sent message = %#v, want body %q", got, tt.wantBody)
			}
		})
	}
}

func TestAccountQuickSendUsesAccountDefaults(t *testing.T) {
	auth := testSMTPAuth{mechanism: "CUSTOM"}
	tlsConfig := &tls.Config{ServerName: "smtp.example.com", MinVersion: tls.VersionTLS12}
	account := Account{
		Host:           "smtp.example.com",
		Port:           587,
		Username:       "user@example.com",
		Password:       "secret",
		Auth:           auth,
		From:           "from@example.com",
		FromName:       "Sender",
		TLSConfig:      tlsConfig,
		TLSPolicy:      TLSNone,
		AllowPlainAuth: true,
		Timeout:        time.Second,
		LocalName:      "mail.local",
	}

	var got *Message
	provider := func(config Config) (Sender, error) {
		if config.Host != account.Host || config.Port != account.Port {
			t.Fatalf("Config address = %#v", config)
		}
		if config.Username != account.Username || config.Password != account.Password || config.Auth == nil {
			t.Fatalf("Config auth = %#v", config)
		}
		if config.TLSPolicy != TLSNone || !config.AllowPlainAuth || config.Timeout != time.Second || config.LocalName != "mail.local" {
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
	err := SendAccountText(
		context.Background(),
		account,
		[]string{"to@example.com"},
		"subject",
		"plain",
		WithQuickMessageOptions(WithHeader("X-Quick", "yes")),
		WithQuickClientOptions(WithSenderProvider(provider)),
	)
	if err != nil {
		t.Fatalf("SendAccountText() error = %v", err)
	}
	if got == nil || got.From.Email != "from@example.com" || got.From.Name != "Sender" {
		t.Fatalf("sent From = %#v", got)
	}
	if got.Subject != "subject" || got.Text != "plain" || got.Headers.Values("X-Quick")[0] != "yes" {
		t.Fatalf("sent message = %#v", got)
	}

	tlsConfig.ServerName = "mutated.example.com"
	if account.TLSConfig.ServerName != "mutated.example.com" {
		t.Fatalf("test setup did not mutate original TLSConfig")
	}
}

func TestQuickSendAndAccountValidation(t *testing.T) {
	provider := func(config Config) (Sender, error) {
		return SenderFunc(func(ctx context.Context, message *Message) error { return nil }), nil
	}
	account := Account{Host: "smtp.example.com", Port: 587, Username: "user@example.com"}
	if err := QuickSend(
		context.Background(),
		account,
		WithQuickMessageOptions(WithTo("to@example.com"), WithSubject("subject"), WithHTML("<p>html</p>")),
		WithQuickClientOptions(WithSenderProvider(provider), WithTLSPolicy(TLSNone)),
	); err != nil {
		t.Fatalf("QuickSend() error = %v", err)
	}

	err := QuickSend(
		context.Background(),
		Account{Host: "smtp.example.com", Port: 587},
		WithQuickMessageOptions(WithTo("to@example.com"), WithText("body")),
		WithQuickClientOptions(WithSenderProvider(provider)),
	)
	if !errors.Is(err, ErrMissingFrom) {
		t.Fatalf("QuickSend() error = %v, want %v", err, ErrMissingFrom)
	}

	quickErr := errors.New("quick option failed")
	err = QuickSend(context.Background(), account, func(*quickConfig) error { return quickErr })
	if !errors.Is(err, quickErr) {
		t.Fatalf("QuickSend(option error) = %v, want %v", err, quickErr)
	}
}

func TestClientOptionsAndProviderErrors(t *testing.T) {
	auth := testSMTPAuth{mechanism: "CUSTOM"}
	tlsConfig := &tls.Config{ServerName: "custom.example.com", MinVersion: tls.VersionTLS13}
	client, err := NewClient("smtp.example.com", 587,
		WithAuth("user", "pass"),
		WithSMTPAuth(auth),
		WithTLSConfig(tlsConfig),
		WithTLSPolicy(TLSPolicyUnknown),
		WithAllowPlainAuth(true),
		WithTimeout(time.Second),
		WithLocalName("mail.local"),
		WithDialContext(func(context.Context, string, string) (net.Conn, error) { return nil, errors.New("dial blocked") }),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	tlsConfig.ServerName = "mutated.example.com"
	if client.config.Username != "user" || client.config.Password != "pass" || !client.config.AllowPlainAuth {
		t.Fatalf("auth/plain config = %#v", client.config)
	}
	if client.config.Auth == nil {
		t.Fatal("custom auth was not configured")
	}
	if client.config.TLSPolicy != TLSMandatoryStartTLS {
		t.Fatalf("TLSPolicy = %v, want %v", client.config.TLSPolicy, TLSMandatoryStartTLS)
	}
	if client.config.TLSConfig.ServerName != "custom.example.com" || client.config.TLSConfig.MinVersion != tls.VersionTLS13 {
		t.Fatalf("TLSConfig was not cloned: %#v", client.config.TLSConfig)
	}

	if _, err := NewClient("smtp.example.com", 587, WithTLSConfig(nil)); err != nil {
		t.Fatalf("NewClient(WithTLSConfig(nil)) error = %v", err)
	}
	if _, err := NewClient("smtp.example.com", 587, WithDialContext(nil)); err == nil {
		t.Fatal("NewClient(WithDialContext(nil)) error = nil, want error")
	}
	if _, err := NewClient("smtp.example.com", 587, WithSenderProvider(nil)); err == nil {
		t.Fatal("NewClient(WithSenderProvider(nil)) error = nil, want error")
	}

	msg, err := NewMessage(WithFrom("from@example.com"), WithTo("to@example.com"), WithText("body"))
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	providerErr := errors.New("provider failed")
	client, err = NewClient("smtp.example.com", 587, WithSenderProvider(func(Config) (Sender, error) {
		return nil, providerErr
	}))
	if err != nil {
		t.Fatalf("NewClient(provider) error = %v", err)
	}
	var nilCtx context.Context
	if err := client.Send(nilCtx, msg); !errors.Is(err, providerErr) {
		t.Fatalf("Send(provider error) = %v, want %v", err, providerErr)
	}
}

func TestSMTPClientSendAgainstFakeServer(t *testing.T) {
	server, err := newFakeSMTPServer(t)
	if err != nil {
		t.Fatalf("newFakeSMTPServer() error = %v", err)
	}
	defer server.Close()

	msg, err := NewMessage(WithFrom("from@example.com"), WithTo("to@example.com"), WithSubject("hello"), WithText("body"))
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	client, err := NewClient(server.Host(), server.Port(), WithTLSPolicy(TLSNone), WithTimeout(time.Second))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if err := server.Wait(); err != nil {
		t.Fatalf("fake SMTP server error = %v", err)
	}
	if !strings.Contains(server.Data(), "Subject: hello") || !strings.Contains(server.Data(), "body") {
		t.Fatalf("SMTP DATA = %q", server.Data())
	}
}

func TestSMTPClientUsesEnvelopeSenderAndDedupedRecipients(t *testing.T) {
	server, err := newFakeSMTPServer(t)
	if err != nil {
		t.Fatalf("newFakeSMTPServer() error = %v", err)
	}
	defer server.Close()

	msg, err := NewMessage(
		WithFrom("header@example.com"),
		WithEnvelopeFrom("bounce@example.com"),
		WithTo("to@example.com", "TO@example.com"),
		WithCc("cc@example.com", "to@example.com"),
		WithBcc("hidden@example.com", "cc@example.com"),
		WithText("body"),
	)
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	client, err := NewClient(server.Host(), server.Port(), WithTLSPolicy(TLSNone), WithTimeout(time.Second))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if err := server.Wait(); err != nil {
		t.Fatalf("fake SMTP server error = %v", err)
	}
	if got := server.MailFrom(); got != "<bounce@example.com>" {
		t.Fatalf("MAIL FROM = %q, want bounce envelope", got)
	}
	wantRecipients := []string{"<to@example.com>", "<cc@example.com>", "<hidden@example.com>"}
	if got := server.RcptTo(); strings.Join(got, ",") != strings.Join(wantRecipients, ",") {
		t.Fatalf("RCPT TO = %v, want %v", got, wantRecipients)
	}
	if recipients := msg.Recipients(); strings.Join(recipients, ",") != "to@example.com,cc@example.com,hidden@example.com" {
		t.Fatalf("Message.Recipients() = %v", recipients)
	}
}

func TestSMTPClientRequiresStartTLS(t *testing.T) {
	server, err := newFakeSMTPServer(t)
	if err != nil {
		t.Fatalf("newFakeSMTPServer() error = %v", err)
	}
	defer server.Close()

	msg, err := NewMessage(WithFrom("from@example.com"), WithTo("to@example.com"), WithText("body"))
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	client, err := NewClient(server.Host(), server.Port(), WithTimeout(time.Second))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.Send(context.Background(), msg); !errors.Is(err, ErrTLSRequired) {
		t.Fatalf("Send() error = %v, want %v", err, ErrTLSRequired)
	}
	if err := server.Wait(); err != nil {
		t.Fatalf("fake SMTP server error = %v", err)
	}
}

func TestSMTPClientRejectsPlainAuthWithoutTLS(t *testing.T) {
	server, err := newFakeSMTPServer(t)
	if err != nil {
		t.Fatalf("newFakeSMTPServer() error = %v", err)
	}
	defer server.Close()

	msg, err := NewMessage(WithFrom("from@example.com"), WithTo("to@example.com"), WithText("body"))
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	client, err := NewClient(server.Host(), server.Port(), WithTLSPolicy(TLSNone), WithAuth("user", "pass"), WithTimeout(time.Second))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.Send(context.Background(), msg); !errors.Is(err, ErrPlainAuth) {
		t.Fatalf("Send() error = %v, want %v", err, ErrPlainAuth)
	}
	if err := server.Wait(); err != nil {
		t.Fatalf("fake SMTP server error = %v", err)
	}
}

func TestSMTPClientUsesCustomAuth(t *testing.T) {
	server, err := newFakeSMTPServer(t, withFakeSMTPAuth("CUSTOM", "token", true))
	if err != nil {
		t.Fatalf("newFakeSMTPServer() error = %v", err)
	}
	defer server.Close()

	msg, err := NewMessage(WithFrom("from@example.com"), WithTo("to@example.com"), WithText("body"))
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	client, err := NewClient(server.Host(), server.Port(),
		WithTLSPolicy(TLSNone),
		WithAllowPlainAuth(true),
		WithSMTPAuth(testSMTPAuth{mechanism: "CUSTOM", initial: []byte("token")}),
		WithTimeout(time.Second),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if err := server.Wait(); err != nil {
		t.Fatalf("fake SMTP server error = %v", err)
	}
	if !server.Authenticated() {
		t.Fatal("server did not receive successful custom AUTH")
	}
}

func TestSMTPClientReturnsAuthFailure(t *testing.T) {
	server, err := newFakeSMTPServer(t, withFakeSMTPAuth("CUSTOM", "token", false))
	if err != nil {
		t.Fatalf("newFakeSMTPServer() error = %v", err)
	}
	defer server.Close()

	msg, err := NewMessage(WithFrom("from@example.com"), WithTo("to@example.com"), WithText("body"))
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	client, err := NewClient(server.Host(), server.Port(),
		WithTLSPolicy(TLSNone),
		WithAllowPlainAuth(true),
		WithSMTPAuth(testSMTPAuth{mechanism: "CUSTOM", initial: []byte("token")}),
		WithTimeout(time.Second),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.Send(context.Background(), msg); err == nil || !strings.Contains(err.Error(), "smtp auth") {
		t.Fatalf("Send() error = %v, want smtp auth error", err)
	}
	if err := server.Wait(); err != nil {
		t.Fatalf("fake SMTP server error = %v", err)
	}
}

func TestSMTPClientStartTLS(t *testing.T) {
	cert := newTestCertificate(t)
	server, err := newFakeSMTPServer(t, withFakeSMTPStartTLS(cert))
	if err != nil {
		t.Fatalf("newFakeSMTPServer() error = %v", err)
	}
	defer server.Close()

	msg, err := NewMessage(WithFrom("from@example.com"), WithTo("to@example.com"), WithText("secure body"))
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	client, err := NewClient(server.Host(), server.Port(),
		WithTLSConfig(&tls.Config{RootCAs: cert.pool, ServerName: "localhost", MinVersion: tls.VersionTLS12}),
		WithTimeout(time.Second),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if err := server.Wait(); err != nil {
		t.Fatalf("fake SMTP server error = %v", err)
	}
	if !server.TLSActive() || !strings.Contains(server.Data(), "secure body") {
		t.Fatalf("TLSActive=%v DATA=%q", server.TLSActive(), server.Data())
	}
}

func TestSMTPClientImplicitTLS(t *testing.T) {
	cert := newTestCertificate(t)
	server, err := newFakeSMTPServer(t, withFakeSMTPImplicitTLS(cert))
	if err != nil {
		t.Fatalf("newFakeSMTPServer() error = %v", err)
	}
	defer server.Close()

	msg, err := NewMessage(WithFrom("from@example.com"), WithTo("to@example.com"), WithText("implicit body"))
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	client, err := NewClient(server.Host(), server.Port(),
		WithTLSPolicy(TLSImplicit),
		WithTLSConfig(&tls.Config{RootCAs: cert.pool, ServerName: "localhost", MinVersion: tls.VersionTLS12}),
		WithTimeout(time.Second),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if err := server.Wait(); err != nil {
		t.Fatalf("fake SMTP server error = %v", err)
	}
	if !server.TLSActive() || !strings.Contains(server.Data(), "implicit body") {
		t.Fatalf("TLSActive=%v DATA=%q", server.TLSActive(), server.Data())
	}
}

func TestSMTPClientContextCancelClosesConnection(t *testing.T) {
	server, err := newFakeSMTPServer(t, withFakeSMTPHangOnData())
	if err != nil {
		t.Fatalf("newFakeSMTPServer() error = %v", err)
	}
	defer server.Close()

	msg, err := NewMessage(WithFrom("from@example.com"), WithTo("to@example.com"), WithText("body"))
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}
	client, err := NewClient(server.Host(), server.Port(), WithTLSPolicy(TLSNone), WithTimeout(0))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	canceled := make(chan error, 1)
	go func() { canceled <- client.Send(ctx, msg) }()
	server.WaitForDataCommand(t)
	cancel()
	select {
	case err := <-canceled:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Send() error = %v, want %v", err, context.Canceled)
		}
	case <-time.After(time.Second):
		t.Fatal("Send() did not return after context cancellation")
	}
	if err := server.Wait(); err != nil {
		t.Fatalf("fake SMTP server error = %v", err)
	}
}

func TestClientDialReusesConnectionWithReset(t *testing.T) {
	server, err := newFakeSMTPServer(t)
	if err != nil {
		t.Fatalf("newFakeSMTPServer() error = %v", err)
	}
	defer server.Close()

	client, err := NewClient(server.Host(), server.Port(), WithTLSPolicy(TLSNone))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	sendCloser, err := client.Dial(context.Background())
	if err != nil {
		t.Fatalf("Dial() error = %v", err)
	}

	first, err := NewMessage(WithFrom("from@example.com"), WithTo("first@example.com"), WithSubject("first"), WithText("first body"))
	if err != nil {
		t.Fatalf("NewMessage(first) error = %v", err)
	}
	second, err := NewMessage(WithFrom("from@example.com"), WithTo("second@example.com"), WithSubject("second"), WithText("second body"))
	if err != nil {
		t.Fatalf("NewMessage(second) error = %v", err)
	}
	if err := sendCloser.Send(context.Background(), first); err != nil {
		t.Fatalf("Send(first) error = %v", err)
	}
	if got := server.RSETCount(); got != 0 {
		t.Fatalf("RSET count after first send = %d, want 0", got)
	}
	if err := sendCloser.Send(context.Background(), second); err != nil {
		t.Fatalf("Send(second) error = %v", err)
	}
	if got := server.RSETCount(); got != 1 {
		t.Fatalf("RSET count after second send = %d, want 1", got)
	}
	if err := sendCloser.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if err := sendCloser.Close(); err != nil {
		t.Fatalf("second Close() error = %v", err)
	}
	if err := sendCloser.Send(context.Background(), second); err == nil {
		t.Fatal("Send() after Close() error = nil")
	}
	if err := server.Wait(); err != nil {
		t.Fatalf("fake SMTP server error = %v", err)
	}
}

func TestSMTPHelpers(t *testing.T) {
	ctx, cancel := withClientTimeout(context.Background(), 0)
	if _, ok := ctx.Deadline(); ok {
		t.Fatal("withClientTimeout(0) unexpectedly set a deadline")
	}
	cancel()

	deadlineCtx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	ctx, cancel = withClientTimeout(deadlineCtx, time.Second)
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("withClientTimeout(existing deadline) removed deadline")
	}
	existingDeadline, _ := deadlineCtx.Deadline()
	if !deadline.Equal(existingDeadline) {
		t.Fatalf("deadline = %v, want %v", deadline, existingDeadline)
	}

	sender := smtpSender{config: Config{Host: "smtp.example.com"}}
	config := sender.tlsConfig()
	if config.ServerName != "smtp.example.com" || config.MinVersion != tls.VersionTLS12 {
		t.Fatalf("default tlsConfig = %#v", config)
	}

	original := &tls.Config{MinVersion: tls.VersionTLS13}
	sender.config.TLSConfig = original
	config = sender.tlsConfig()
	if config.ServerName != "smtp.example.com" || config.MinVersion != tls.VersionTLS13 {
		t.Fatalf("cloned tlsConfig = %#v", config)
	}
	config.ServerName = "mutated.example.com"
	if original.ServerName != "" {
		t.Fatalf("tlsConfig mutated original: %#v", original)
	}

	dialErr := errors.New("dial failed")
	sender.config.DialContext = func(context.Context, string, string) (net.Conn, error) { return nil, dialErr }
	if _, err := sender.dial(context.Background(), "smtp.example.com:587", config); !errors.Is(err, dialErr) {
		t.Fatalf("dial() error = %v, want %v", err, dialErr)
	}
}

type fakeSMTPServer struct {
	listener      net.Listener
	done          chan error
	dataStarted   chan struct{}
	mu            sync.Mutex
	data          string
	mailFrom      string
	rcptTo        []string
	rsetCount     int
	cert          *testCertificate
	startTLS      bool
	implicitTLS   bool
	tlsActive     bool
	authMechanism string
	authInitial   string
	authOK        bool
	authenticated bool
	hangOnData    bool
	once          sync.Once
}

type fakeSMTPOption func(*fakeSMTPServer)

func newFakeSMTPServer(t *testing.T, opts ...fakeSMTPOption) (*fakeSMTPServer, error) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	server := &fakeSMTPServer{
		listener:    listener,
		done:        make(chan error, 1),
		dataStarted: make(chan struct{}),
		authOK:      true,
	}
	for _, opt := range opts {
		opt(server)
	}
	go server.serve()
	return server, nil
}

func withFakeSMTPStartTLS(cert *testCertificate) fakeSMTPOption {
	return func(s *fakeSMTPServer) {
		s.cert = cert
		s.startTLS = true
	}
}

func withFakeSMTPImplicitTLS(cert *testCertificate) fakeSMTPOption {
	return func(s *fakeSMTPServer) {
		s.cert = cert
		s.implicitTLS = true
	}
}

func withFakeSMTPAuth(mechanism, initial string, ok bool) fakeSMTPOption {
	return func(s *fakeSMTPServer) {
		s.authMechanism = mechanism
		s.authInitial = initial
		s.authOK = ok
	}
}

func withFakeSMTPHangOnData() fakeSMTPOption {
	return func(s *fakeSMTPServer) { s.hangOnData = true }
}

func (s *fakeSMTPServer) Host() string {
	host, _, _ := net.SplitHostPort(s.listener.Addr().String())
	return host
}

func (s *fakeSMTPServer) Port() int {
	_, port, _ := net.SplitHostPort(s.listener.Addr().String())
	n, _ := strconvAtoi(port)
	return n
}

func (s *fakeSMTPServer) Close() {
	_ = s.listener.Close()
}

func (s *fakeSMTPServer) Data() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data
}

func (s *fakeSMTPServer) MailFrom() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.mailFrom
}

func (s *fakeSMTPServer) RcptTo() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.rcptTo...)
}

func (s *fakeSMTPServer) RSETCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rsetCount
}

func (s *fakeSMTPServer) TLSActive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tlsActive
}

func (s *fakeSMTPServer) Authenticated() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.authenticated
}

func (s *fakeSMTPServer) WaitForDataCommand(t *testing.T) {
	t.Helper()
	select {
	case <-s.dataStarted:
	case <-time.After(time.Second):
		t.Fatal("fake smtp server did not receive DATA command")
	}
}

func (s *fakeSMTPServer) Wait() error {
	select {
	case err := <-s.done:
		return err
	case <-time.After(2 * time.Second):
		return errors.New("fake smtp server timed out")
	}
}

func (s *fakeSMTPServer) serve() {
	conn, err := s.listener.Accept()
	if err != nil {
		s.done <- err
		return
	}
	defer func() { _ = conn.Close() }()
	if s.implicitTLS {
		conn = tls.Server(conn, s.cert.serverConfig())
		if err := conn.(*tls.Conn).Handshake(); err != nil {
			s.done <- err
			return
		}
		s.setTLSActive()
	}
	reader := bufio.NewReader(conn)
	if _, err := io.WriteString(conn, "220 fake.smtp ESMTP\r\n"); err != nil {
		s.done <- err
		return
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				s.done <- nil
				return
			}
			s.done <- err
			return
		}
		line = strings.TrimRight(line, "\r\n")
		switch {
		case strings.HasPrefix(line, "EHLO") || strings.HasPrefix(line, "HELO"):
			if _, err := io.WriteString(conn, s.ehloResponse()); err != nil {
				s.done <- err
				return
			}
		case line == "STARTTLS" && s.startTLS:
			if _, err := io.WriteString(conn, "220 Ready to start TLS\r\n"); err != nil {
				s.done <- err
				return
			}
			conn = tls.Server(conn, s.cert.serverConfig())
			if err := conn.(*tls.Conn).Handshake(); err != nil {
				s.done <- err
				return
			}
			s.setTLSActive()
			reader = bufio.NewReader(conn)
		case strings.HasPrefix(line, "AUTH "):
			if err := s.handleAuth(conn, line); err != nil {
				s.done <- err
				return
			}
		case line == "*":
			s.done <- nil
			return
		case strings.HasPrefix(line, "MAIL FROM:"):
			s.mu.Lock()
			s.mailFrom = strings.TrimPrefix(line, "MAIL FROM:")
			s.mu.Unlock()
			if _, err := io.WriteString(conn, "250 OK\r\n"); err != nil {
				s.done <- err
				return
			}
		case strings.HasPrefix(line, "RCPT TO:"):
			s.mu.Lock()
			s.rcptTo = append(s.rcptTo, strings.TrimPrefix(line, "RCPT TO:"))
			s.mu.Unlock()
			if _, err := io.WriteString(conn, "250 OK\r\n"); err != nil {
				s.done <- err
				return
			}
		case line == "RSET":
			s.mu.Lock()
			s.rsetCount++
			s.mu.Unlock()
			if _, err := io.WriteString(conn, "250 OK\r\n"); err != nil {
				s.done <- err
				return
			}
		case line == "DATA":
			s.once.Do(func() { close(s.dataStarted) })
			if s.hangOnData {
				_, err := reader.ReadString('\n')
				if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
					s.done <- nil
					return
				}
				s.done <- err
				return
			}
			if _, err := io.WriteString(conn, "354 End data with <CR><LF>.<CR><LF>\r\n"); err != nil {
				s.done <- err
				return
			}
			var data strings.Builder
			for {
				dataLine, err := reader.ReadString('\n')
				if err != nil {
					s.done <- err
					return
				}
				if strings.TrimRight(dataLine, "\r\n") == "." {
					break
				}
				data.WriteString(dataLine)
			}
			s.mu.Lock()
			s.data = data.String()
			s.mu.Unlock()
			if _, err := io.WriteString(conn, "250 OK queued\r\n"); err != nil {
				s.done <- err
				return
			}
		case line == "QUIT":
			_, err := io.WriteString(conn, "221 Bye\r\n")
			s.done <- err
			return
		default:
			s.done <- fmt.Errorf("unexpected SMTP command %q", line)
			return
		}
	}
}

func (s *fakeSMTPServer) ehloResponse() string {
	var builder strings.Builder
	builder.WriteString("250-fake.smtp\r\n")
	if s.startTLS && !s.TLSActive() {
		builder.WriteString("250-STARTTLS\r\n")
	}
	if s.authMechanism != "" {
		builder.WriteString("250-AUTH " + s.authMechanism + "\r\n")
	}
	builder.WriteString("250 OK\r\n")
	return builder.String()
}

func (s *fakeSMTPServer) handleAuth(conn net.Conn, line string) error {
	fields := strings.Fields(line)
	if len(fields) < 2 || fields[1] != s.authMechanism {
		_, err := io.WriteString(conn, "504 unsupported auth\r\n")
		return err
	}
	if !s.authOK {
		_, err := io.WriteString(conn, "535 auth failed\r\n")
		return err
	}
	initial := ""
	if len(fields) > 2 {
		decoded, err := base64.StdEncoding.DecodeString(fields[2])
		if err != nil {
			return err
		}
		initial = string(decoded)
	}
	if initial != s.authInitial {
		_, err := io.WriteString(conn, "535 auth failed\r\n")
		return err
	}
	s.mu.Lock()
	s.authenticated = true
	s.mu.Unlock()
	_, err := io.WriteString(conn, "235 authenticated\r\n")
	return err
}

func (s *fakeSMTPServer) setTLSActive() {
	s.mu.Lock()
	s.tlsActive = true
	s.mu.Unlock()
}

type testSMTPAuth struct {
	mechanism string
	initial   []byte
}

func (a testSMTPAuth) Start(*smtp.ServerInfo) (string, []byte, error) {
	return a.mechanism, a.initial, nil
}

func (a testSMTPAuth) Next([]byte, bool) ([]byte, error) { return nil, nil }

type testCertificate struct {
	cert tls.Certificate
	pool *x509.CertPool
}

func (c *testCertificate) serverConfig() *tls.Config {
	return &tls.Config{Certificates: []tls.Certificate{c.cert}, MinVersion: tls.VersionTLS12}
}

func newTestCertificate(t *testing.T) *testCertificate {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"localhost"},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("CreateCertificate() error = %v", err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatalf("MarshalECPrivateKey() error = %v", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		t.Fatalf("X509KeyPair() error = %v", err)
	}
	parsed, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("ParseCertificate() error = %v", err)
	}
	pool := x509.NewCertPool()
	pool.AddCert(parsed)
	return &testCertificate{cert: cert, pool: pool}
}

func strconvAtoi(value string) (int, error) {
	var n int
	for _, r := range value {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("invalid digit %q", r)
		}
		n = n*10 + int(r-'0')
	}
	return n, nil
}

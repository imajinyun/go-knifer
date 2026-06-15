package mail

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
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

func TestClientOptionsAndProviderErrors(t *testing.T) {
	tlsConfig := &tls.Config{ServerName: "custom.example.com", MinVersion: tls.VersionTLS13}
	client, err := NewClient("smtp.example.com", 587,
		WithAuth("user", "pass"),
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
	listener net.Listener
	done     chan error
	mu       sync.Mutex
	data     string
}

func newFakeSMTPServer(t *testing.T) (*fakeSMTPServer, error) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	server := &fakeSMTPServer{
		listener: listener,
		done:     make(chan error, 1),
	}
	go server.serve()
	return server, nil
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
			if _, err := io.WriteString(conn, "250-fake.smtp\r\n250 OK\r\n"); err != nil {
				s.done <- err
				return
			}
		case strings.HasPrefix(line, "MAIL FROM:"):
			if _, err := io.WriteString(conn, "250 OK\r\n"); err != nil {
				s.done <- err
				return
			}
		case strings.HasPrefix(line, "RCPT TO:"):
			if _, err := io.WriteString(conn, "250 OK\r\n"); err != nil {
				s.done <- err
				return
			}
		case line == "DATA":
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

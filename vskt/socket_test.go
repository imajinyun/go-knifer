package vskt_test

import (
	"context"
	"errors"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/imajinyun/go-knifer/vskt"
)

type facadeFakeDialer struct {
	calls   atomic.Int32
	network string
	address string
	server  net.Conn
}

type facadeFakeAddr string

func (a facadeFakeAddr) Network() string { return "tcp" }
func (a facadeFakeAddr) String() string  { return string(a) }

type facadeFakeListener struct {
	addr net.Addr
}

func (l *facadeFakeListener) Accept() (net.Conn, error) { return nil, net.ErrClosed }
func (l *facadeFakeListener) Close() error              { return nil }
func (l *facadeFakeListener) Addr() net.Addr            { return l.addr }

func (d *facadeFakeDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	d.calls.Add(1)
	d.network = network
	d.address = address
	client, server := net.Pipe()
	d.server = server
	return client, nil
}

func TestFacadeSocketConfig(t *testing.T) {
	cfg := vskt.NewSocketConfig()
	if cfg == nil {
		t.Fatal("expected non-nil socket config")
	}
}

func TestFacadeSocketConfigWithOptions(t *testing.T) {
	listener := &facadeFakeListener{addr: facadeFakeAddr("listener")}
	client, server := net.Pipe()
	defer func() { _ = server.Close() }()
	cfg := vskt.NewSocketConfigWithOptions(
		vskt.WithThreadPoolSize(2),
		vskt.WithReadTimeout(100),
		vskt.WithWriteTimeout(200),
		vskt.WithReadBufferSize(64),
		vskt.WithWriteBufferSize(128),
		vskt.WithListenerFactory(func(*net.TCPAddr) (net.Listener, error) { return listener, nil }),
		vskt.WithConnFactory(func(*net.TCPAddr) (net.Conn, error) { return client, nil }),
	)
	if cfg.ThreadPoolSize != 2 || cfg.ReadTimeout != 100 || cfg.WriteTimeout != 200 ||
		cfg.ReadBufferSize != 64 || cfg.WriteBufferSize != 128 {
		t.Fatalf("NewSocketConfigWithOptions not applied: %+v", cfg)
	}
	if cfg.ListenerFactory == nil || cfg.ConnFactory == nil {
		t.Fatal("expected listener and connection factories")
	}
}

func TestFacadeSocketConfigThreadPoolSizeFunc(t *testing.T) {
	calls := 0
	cfg := vskt.NewSocketConfigWithOptions(vskt.WithThreadPoolSizeFunc(func() int {
		calls++
		return 9
	}))
	if calls != 1 || cfg.ThreadPoolSize != 9 {
		t.Fatalf("WithThreadPoolSizeFunc calls=%d size=%d, want 1/9", calls, cfg.ThreadPoolSize)
	}
}

func TestFacadeSocketFactories(t *testing.T) {
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9999}
	listener := &facadeFakeListener{addr: facadeFakeAddr("facade-aio")}
	aio, err := vskt.NewAioServerAddrWithOptions(addr, nil, vskt.WithListenerFactory(func(got *net.TCPAddr) (net.Listener, error) {
		if got != addr {
			return nil, errors.New("unexpected aio addr")
		}
		return listener, nil
	}))
	if err != nil {
		t.Fatalf("NewAioServerAddrWithOptions: %v", err)
	}
	if aio.Listener() != listener || aio.LocalAddr().String() != "facade-aio" {
		t.Fatalf("aio listener = %#v addr=%v", aio.Listener(), aio.LocalAddr())
	}
	_ = aio.Close()

	listener = &facadeFakeListener{addr: facadeFakeAddr("facade-nio")}
	nio, err := vskt.NewNioServerAddrWithOptions(addr, nil, vskt.WithListenerFactory(func(got *net.TCPAddr) (net.Listener, error) {
		if got != addr {
			return nil, errors.New("unexpected nio addr")
		}
		return listener, nil
	}))
	if err != nil {
		t.Fatalf("NewNioServerAddrWithOptions: %v", err)
	}
	if nio.Listener() != listener || nio.LocalAddr().String() != "facade-nio" {
		t.Fatalf("nio listener = %#v addr=%v", nio.Listener(), nio.LocalAddr())
	}
	_ = nio.Close()

	client, server := net.Pipe()
	defer func() { _ = server.Close() }()
	nioClient, err := vskt.NewNioClientAddrWithOptions(addr, vskt.WithConnFactory(func(got *net.TCPAddr) (net.Conn, error) {
		if got != addr {
			return nil, errors.New("unexpected client addr")
		}
		return client, nil
	}))
	if err != nil {
		t.Fatalf("NewNioClientAddrWithOptions: %v", err)
	}
	if nioClient.Channel() != client {
		t.Fatalf("nio client channel = %#v", nioClient.Channel())
	}
	_ = nioClient.Close()
}

func TestFacadeSocketConnectWithOptions(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()

	go func() {
		conn, _ := ln.Accept()
		if conn != nil {
			_ = conn.Close()
		}
	}()

	addr := ln.Addr().(*net.TCPAddr)
	conn, err := vskt.SocketConnectWithOptions("127.0.0.1", addr.Port, vskt.WithConnectTimeout(time.Second), vskt.WithConnectNetwork("tcp"))
	if err != nil {
		t.Fatalf("SocketConnectWithOptions failed: %v", err)
	}
	defer func() { _ = conn.Close() }()
	if !vskt.SocketIsConnected(conn) {
		t.Fatal("SocketConnectWithOptions should return a connected socket")
	}
}

func TestFacadeSocketConnectOptionVariants(t *testing.T) {
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1234}
	dialer := &facadeFakeDialer{}
	conn, err := vskt.SocketConnectAddrWithOptions(addr, vskt.WithConnectDialer(dialer), vskt.WithConnectNetwork("tcp4"))
	if err != nil {
		t.Fatalf("SocketConnectAddrWithOptions failed: %v", err)
	}
	_ = conn.Close()
	_ = dialer.server.Close()
	if dialer.calls.Load() != 1 || dialer.network != "tcp4" || dialer.address != "127.0.0.1:1234" {
		t.Fatalf("dialer call = (%d, %q, %q)", dialer.calls.Load(), dialer.network, dialer.address)
	}

	dialer = &facadeFakeDialer{}
	conn, err = vskt.ChannelDialWithOptions(addr, vskt.WithConnectDialer(dialer))
	if err != nil {
		t.Fatalf("ChannelDialWithOptions failed: %v", err)
	}
	_ = conn.Close()
	_ = dialer.server.Close()

	dialer = &facadeFakeDialer{}
	client, err := vskt.NewAioClientWithOptions(addr, &vskt.SimpleIoAction{}, vskt.WithConnectDialer(dialer))
	if err != nil {
		t.Fatalf("NewAioClientWithOptions failed: %v", err)
	}
	_ = client.Close()
	_ = dialer.server.Close()
}

func TestFacadeSocketIsConnected(t *testing.T) {
	// nil conn should not be connected
	if vskt.SocketIsConnected(nil) {
		t.Fatal("expected nil conn to be disconnected")
	}
}

func TestFacadeSocketError(t *testing.T) {
	err := vskt.NewSocketErrorMsg("test error")
	if err == nil {
		t.Fatal("expected non-nil socket error")
	}
	if err.Error() != "test error" {
		t.Fatalf("expected 'test error', got %q", err.Error())
	}
}

func TestFacadeOperations(t *testing.T) {
	// verify operation constants are accessible
	_ = vskt.OpRead
	_ = vskt.OpWrite
	_ = vskt.OpConnect
	_ = vskt.OpAccept
}

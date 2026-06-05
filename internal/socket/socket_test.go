package socket

import (
	"context"
	"errors"
	"net"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	knifer "github.com/imajinyun/go-knifer"
)

func closeAndReport(t *testing.T, closeFn func() error) {
	t.Helper()
	if err := closeFn(); err != nil {
		t.Errorf("close failed: %v", err)
	}
}

func waitForInt32(t *testing.T, get func() int32, want int32) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if get() >= want {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("value = %d, want %d", get(), want)
}

type fakeDialer struct {
	calls   atomic.Int32
	network string
	address string
	server  net.Conn
}

func (d *fakeDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	d.calls.Add(1)
	d.network = network
	d.address = address
	client, server := net.Pipe()
	d.server = server
	return client, nil
}

func TestSocketConfigDefaults(t *testing.T) {
	cfg := NewSocketConfig()
	if cfg.ThreadPoolSize != runtime.NumCPU() {
		t.Errorf("ThreadPoolSize 默认应为 CPU 核数，实际 %d", cfg.ThreadPoolSize)
	}
	if cfg.ReadBufferSize != DefaultBufferSize || cfg.WriteBufferSize != DefaultBufferSize {
		t.Errorf("默认缓冲区大小不正确：%d / %d", cfg.ReadBufferSize, cfg.WriteBufferSize)
	}

	cfg.SetThreadPoolSize(8).SetReadTimeout(100).SetWriteTimeout(200).
		SetReadBufferSize(1024).SetWriteBufferSize(2048)
	if cfg.ThreadPoolSize != 8 || cfg.ReadTimeout != 100 || cfg.WriteTimeout != 200 ||
		cfg.ReadBufferSize != 1024 || cfg.WriteBufferSize != 2048 {
		t.Errorf("链式 setter 未生效: %+v", cfg)
	}
}

func TestSocketConfigOptions(t *testing.T) {
	cfg := NewSocketConfigWithOptions(
		WithThreadPoolSize(2),
		WithReadTimeout(100),
		WithWriteTimeout(200),
		WithReadBufferSize(64),
		WithWriteBufferSize(128),
	)
	if cfg.ThreadPoolSize != 2 || cfg.ReadTimeout != 100 || cfg.WriteTimeout != 200 ||
		cfg.ReadBufferSize != 64 || cfg.WriteBufferSize != 128 {
		t.Fatalf("NewSocketConfigWithOptions not applied: %+v", cfg)
	}
}

func TestSocketUtilConnect(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer closeAndReport(t, ln.Close)

	go func() {
		c, _ := ln.Accept()
		if c != nil {
			_ = c.Close()
		}
	}()

	addr := ln.Addr().(*net.TCPAddr)
	conn, err := Connect("127.0.0.1", addr.Port, time.Second)
	if err != nil {
		t.Fatalf("Connect 失败: %v", err)
	}
	defer closeAndReport(t, conn.Close)

	if !IsConnected(conn) {
		t.Errorf("IsConnected 应返回 true")
	}
	if GetRemoteAddress(conn) == nil {
		t.Errorf("GetRemoteAddress 不应返回 nil")
	}
}

func TestSocketConnectWithOptions(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer closeAndReport(t, ln.Close)

	go func() {
		c, _ := ln.Accept()
		if c != nil {
			_ = c.Close()
		}
	}()

	addr := ln.Addr().(*net.TCPAddr)
	conn, err := ConnectWithOptions("127.0.0.1", addr.Port, WithConnectTimeout(time.Second), WithConnectNetwork("tcp"))
	if err != nil {
		t.Fatalf("ConnectWithOptions failed: %v", err)
	}
	defer closeAndReport(t, conn.Close)
}

func TestSocketConnectOptionsUseDialerAndContext(t *testing.T) {
	dialer := &fakeDialer{}
	conn, err := ConnectWithOptions("example.com", 443, WithConnectNetwork("tcp4"), WithConnectDialer(dialer))
	if err != nil {
		t.Fatalf("ConnectWithOptions with fake dialer failed: %v", err)
	}
	defer closeAndReport(t, conn.Close)
	defer closeAndReport(t, dialer.server.Close)
	if dialer.calls.Load() != 1 || dialer.network != "tcp4" || dialer.address != "example.com:443" {
		t.Fatalf("dialer call = (%d, %q, %q), want (1, tcp4, example.com:443)", dialer.calls.Load(), dialer.network, dialer.address)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := ConnectWithOptions("127.0.0.1", 1, WithConnectContext(ctx)); err == nil {
		t.Fatal("ConnectWithOptions with canceled context error = nil")
	}
}

func TestSocketAddrChannelAndAioClientOptions(t *testing.T) {
	dialer := &fakeDialer{}
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1234}
	conn, err := ConnectAddrWithOptions(addr, WithConnectDialer(dialer))
	if err != nil {
		t.Fatalf("ConnectAddrWithOptions failed: %v", err)
	}
	closeAndReport(t, conn.Close)
	closeAndReport(t, dialer.server.Close)

	dialer = &fakeDialer{}
	conn, err = ChannelUtilDialWithOptions(addr, WithConnectDialer(dialer))
	if err != nil {
		t.Fatalf("ChannelUtilDialWithOptions failed: %v", err)
	}
	closeAndReport(t, conn.Close)
	closeAndReport(t, dialer.server.Close)

	dialer = &fakeDialer{}
	client, err := NewAioClientWithOptions(addr, &echoIoAction{}, WithConnectDialer(dialer))
	if err != nil {
		t.Fatalf("NewAioClientWithOptions failed: %v", err)
	}
	closeAndReport(t, client.Close)
	closeAndReport(t, dialer.server.Close)
}

func TestSocketRuntimeError(t *testing.T) {
	base := net.ErrClosed
	e := WrapSocketError(base, "wrapped")
	if e == nil {
		t.Fatal("WrapSocketError 不应返回 nil")
	}
	if e.Unwrap() != base {
		t.Errorf("Unwrap 失败")
	}
	if e.Error() == "" {
		t.Errorf("Error 不应为空")
	}
	if !errors.Is(e, knifer.ErrCodeInternal) {
		t.Errorf("SocketRuntimeError 应匹配 ErrCodeInternal")
	}
	if !errors.Is(e, base) {
		t.Errorf("SocketRuntimeError 应保留 cause 链")
	}
	if WrapSocketError(nil, "x") != nil {
		t.Errorf("nil err 应返回 nil")
	}
	if NewSocketErrorf("hello %s", "world").Error() != "hello world" {
		t.Errorf("格式化失败")
	}
}

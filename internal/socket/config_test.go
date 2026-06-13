package socket

import (
	"net"
	"runtime"
	"testing"
)

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
	listener := &fakeListener{addr: factoryFakeAddr("listener")}
	client, server := net.Pipe()
	defer closeAndReport(t, server.Close)
	runnerCalled := false
	cfg := NewSocketConfigWithOptions(
		WithThreadPoolSize(2),
		WithReadTimeout(100),
		WithWriteTimeout(200),
		WithReadBufferSize(64),
		WithWriteBufferSize(128),
		WithRunner(func(fn func()) { runnerCalled = true; fn() }),
		WithListenerFactory(func(*net.TCPAddr) (net.Listener, error) { return listener, nil }),
		WithConnFactory(func(*net.TCPAddr) (net.Conn, error) { return client, nil }),
	)
	if cfg.ThreadPoolSize != 2 || cfg.ReadTimeout != 100 || cfg.WriteTimeout != 200 ||
		cfg.ReadBufferSize != 64 || cfg.WriteBufferSize != 128 {
		t.Fatalf("NewSocketConfigWithOptions not applied: %+v", cfg)
	}
	if cfg.ListenerFactory == nil || cfg.ConnFactory == nil || cfg.Runner == nil {
		t.Fatal("expected listener, connection, and runner providers")
	}
	cfg.Runner(func() {})
	if !runnerCalled {
		t.Fatal("custom runner was not called")
	}
}

func TestSocketConfigThreadPoolSizeFunc(t *testing.T) {
	calls := 0
	cfg := NewSocketConfigWithOptions(WithThreadPoolSizeFunc(func() int {
		calls++
		return 7
	}))
	if calls != 1 || cfg.ThreadPoolSize != 7 {
		t.Fatalf("WithThreadPoolSizeFunc calls=%d size=%d, want 1/7", calls, cfg.ThreadPoolSize)
	}
}

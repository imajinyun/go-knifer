package socket

import (
	"net"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
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

func TestSocketConfigProviderFallbacks(t *testing.T) {
	cfg := NewSocketConfigWithOptions(
		nil,
		WithThreadPoolSizeFunc(nil),
		WithClock(nil),
		WithRunner(nil),
		WithListenerFactory(nil),
		WithConnFactory(nil),
		WithSocketIPParser(nil),
	)
	if cfg.ThreadPoolSize != runtime.NumCPU() || cfg.Clock != nil || cfg.Runner != nil || cfg.ListenerFactory != nil || cfg.ConnFactory != nil || cfg.IPParser != nil {
		t.Fatalf("nil provider options should preserve defaults: %+v", cfg)
	}
	if got := parseIPWithConfig(cfg, "127.0.0.1"); !got.Equal(net.ParseIP("127.0.0.1")) {
		t.Fatalf("parseIPWithConfig default = %v", got)
	}
	customIP := net.ParseIP("127.0.0.9")
	cfg.SetSocketIPParser(func(string) net.IP { return customIP })
	if got := parseIPWithConfig(cfg, "example.test"); !got.Equal(customIP) {
		t.Fatalf("parseIPWithConfig custom = %v", got)
	}

	var ran atomic.Bool
	runWithConfig(nil, func() { ran.Store(true) })
	waitForBool(t, ran.Load)
	ran.Store(false)
	runWithConfig(&SocketConfig{Runner: func(fn func()) { fn() }}, func() { ran.Store(true) })
	if !ran.Load() {
		t.Fatal("runWithConfig custom runner did not run")
	}
}

func TestConcurrencyLimiterBoundaries(t *testing.T) {
	if limiter := newConcurrencyLimiter(nil); limiter != nil {
		t.Fatalf("newConcurrencyLimiter(nil) = %#v, want nil", limiter)
	}
	if limiter := newConcurrencyLimiter(&SocketConfig{ThreadPoolSize: 0}); limiter != nil {
		t.Fatalf("newConcurrencyLimiter(0) = %#v, want nil", limiter)
	}
	if !acquireConcurrencySlot(nil, make(chan struct{})) {
		t.Fatal("nil limiter should always acquire")
	}
	releaseConcurrencySlot(nil)

	limiter := newConcurrencyLimiter(&SocketConfig{ThreadPoolSize: 1})
	done := make(chan struct{})
	if !acquireConcurrencySlot(limiter, done) {
		t.Fatal("first acquire should succeed")
	}
	close(done)
	if acquireConcurrencySlot(limiter, done) {
		t.Fatal("acquire should fail when limiter is full and done is closed")
	}
	releaseConcurrencySlot(limiter)
	select {
	case limiter <- struct{}{}:
		releaseConcurrencySlot(limiter)
	case <-time.After(time.Second):
		t.Fatal("releaseConcurrencySlot did not release limiter capacity")
	}
}

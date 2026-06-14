package socket

import (
	"net"
	"sync/atomic"
	"testing"
	"time"
)

type blockingChannelHandler struct {
	started atomic.Int32
	release chan struct{}
}

func (h *blockingChannelHandler) Handle(conn net.Conn) error {
	h.started.Add(1)
	<-h.release
	return nil
}

func TestNioServerThreadPoolSizeLimitsHandlers(t *testing.T) {
	handler := &blockingChannelHandler{release: make(chan struct{})}
	server, err := NewNioServerAddrWithConfig(
		&net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0},
		NewSocketConfigWithOptions(WithThreadPoolSize(1)),
	)
	if err != nil {
		t.Fatal(err)
	}
	server.SetChannelHandler(handler)
	defer closeAndReport(t, server.Close)
	server.ListenAsync()

	addr := server.LocalAddr().String()
	first, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	waitForInt32(t, func() int32 { return handler.started.Load() }, 1)

	second, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer closeAndReport(t, second.Close)
	time.Sleep(50 * time.Millisecond)
	if got := handler.started.Load(); got != 1 {
		t.Fatalf("started handlers = %d, want 1 while first handler occupies the only slot", got)
	}

	close(handler.release)
	closeAndReport(t, first.Close)
	waitForInt32(t, func() int32 { return handler.started.Load() }, 2)
}

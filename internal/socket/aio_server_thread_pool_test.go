package socket

import (
	"bytes"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

type blockingIoAction struct {
	accepted atomic.Int32
	release  chan struct{}
}

func (a *blockingIoAction) Accept(session *AioSession) {
	a.accepted.Add(1)
	<-a.release
}

func (a *blockingIoAction) DoAction(session *AioSession, data *bytes.Buffer) {}

func (a *blockingIoAction) Failed(err error, session *AioSession) {}

func TestAioServerThreadPoolSizeLimitsHandlers(t *testing.T) {
	action := &blockingIoAction{release: make(chan struct{})}
	server, err := NewAioServerAddr(
		&net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0},
		NewSocketConfigWithOptions(WithThreadPoolSize(1)),
	)
	if err != nil {
		t.Fatal(err)
	}
	server.SetIoAction(action)
	defer closeAndReport(t, server.Close)
	server.Start(false)

	addr := server.LocalAddr().String()
	first, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	waitForInt32(t, func() int32 { return action.accepted.Load() }, 1)

	second, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer closeAndReport(t, second.Close)
	time.Sleep(50 * time.Millisecond)
	if got := action.accepted.Load(); got != 1 {
		t.Fatalf("accepted = %d, want 1 while first handler occupies the only slot", got)
	}

	close(action.release)
	closeAndReport(t, first.Close)
	waitForInt32(t, func() int32 { return action.accepted.Load() }, 2)
}

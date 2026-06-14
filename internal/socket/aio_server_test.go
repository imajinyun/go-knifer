package socket

import (
	"bytes"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

// echoIoAction implements IoAction and writes received data back as-is.
type echoIoAction struct {
	accepted atomic.Int32
	failed   atomic.Int32
}

func (a *echoIoAction) Accept(session *AioSession) {
	a.accepted.Add(1)
}

func (a *echoIoAction) DoAction(session *AioSession, data *bytes.Buffer) {
	if data == nil || data.Len() == 0 {
		return
	}
	_, _ = session.Write(data.Bytes())
}

func (a *echoIoAction) Failed(err error, session *AioSession) {
	a.failed.Add(1)
}

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

func TestAioServerEcho(t *testing.T) {
	action := &echoIoAction{}
	server, err := NewAioServerAddr(&net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0}, NewSocketConfig())
	if err != nil {
		t.Fatal(err)
	}
	server.SetIoAction(action)
	defer closeAndReport(t, server.Close)

	server.Start(false)

	addr := server.LocalAddr().(*net.TCPAddr)
	conn, err := net.DialTimeout("tcp", addr.String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer closeAndReport(t, conn.Close)

	want := []byte("hello-aio")
	if _, err := conn.Write(want); err != nil {
		t.Fatal(err)
	}
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	got := make([]byte, len(want))
	if _, err := conn.Read(got); err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Errorf("AioServer 回显数据不一致: got=%q want=%q", got, want)
	}
	if action.accepted.Load() != 1 {
		t.Errorf("Accept 应被回调 1 次，实际 %d", action.accepted.Load())
	}
}

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

func TestAioServerCloseClosesActiveConnections(t *testing.T) {
	client, serverConn := net.Pipe()
	defer closeAndReport(t, client.Close)
	listener := &queuedListener{addr: factoryFakeAddr("aio"), conns: make(chan net.Conn, 1)}
	listener.conns <- serverConn
	accepted := make(chan struct{})
	aio, err := NewAioServerAddrWithOptions(&net.TCPAddr{Port: 1}, nil, WithListenerFactory(func(*net.TCPAddr) (net.Listener, error) {
		return listener, nil
	}))
	if err != nil {
		t.Fatalf("NewAioServerAddrWithOptions: %v", err)
	}
	aio.SetIoAction(&SimpleIoAction{OnAccept: func(*AioSession) { close(accepted) }})
	aio.Start(false)
	select {
	case <-accepted:
	case <-time.After(time.Second):
		t.Fatal("aio server did not accept connection")
	}
	done := make(chan error, 1)
	go func() { done <- aio.Close() }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Close error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("AioServer.Close blocked with active connection")
	}
}

package socket

import (
	"bytes"
	"net"
	"sync"
	"testing"
	"time"
)

// clientIoAction notifies done after receiving one message.
type clientIoAction struct {
	mu      sync.Mutex
	message []byte
	done    chan struct{}
}

func (a *clientIoAction) Accept(session *AioSession) {}

func (a *clientIoAction) DoAction(session *AioSession, data *bytes.Buffer) {
	a.mu.Lock()
	a.message = append(a.message, data.Bytes()...)
	a.mu.Unlock()
	select {
	case a.done <- struct{}{}:
	default:
	}
}

func (a *clientIoAction) Failed(err error, session *AioSession) {}

func TestAioClient(t *testing.T) {
	server, err := NewAioServerAddr(&net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0}, NewSocketConfig())
	if err != nil {
		t.Fatal(err)
	}
	server.SetIoAction(&echoIoAction{})
	defer closeAndReport(t, server.Close)
	server.Start(false)

	addr := server.LocalAddr().(*net.TCPAddr)
	clientAction := &clientIoAction{done: make(chan struct{}, 1)}

	client, err := NewAioClient(addr, clientAction)
	if err != nil {
		t.Fatal(err)
	}
	defer closeAndReport(t, client.Close)

	if _, err := client.Write([]byte("ping")); err != nil {
		t.Fatal(err)
	}
	client.Read()

	select {
	case <-clientAction.done:
	case <-time.After(2 * time.Second):
		t.Fatal("AioClient 未在超时内收到回显")
	}

	clientAction.mu.Lock()
	got := string(clientAction.message)
	clientAction.mu.Unlock()
	if got != "ping" {
		t.Errorf("AioClient 收到错误数据: %q", got)
	}
}

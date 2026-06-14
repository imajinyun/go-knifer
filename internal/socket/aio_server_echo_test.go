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

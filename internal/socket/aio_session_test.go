package socket

import (
	"bytes"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type deadlineConn struct {
	readDeadline  time.Time
	writeDeadline time.Time
}

func (c *deadlineConn) Read([]byte) (int, error)          { return 0, io.EOF }
func (c *deadlineConn) Write(p []byte) (int, error)       { return len(p), nil }
func (c *deadlineConn) Close() error                      { return nil }
func (c *deadlineConn) LocalAddr() net.Addr               { return factoryFakeAddr("local") }
func (c *deadlineConn) RemoteAddr() net.Addr              { return factoryFakeAddr("remote") }
func (c *deadlineConn) SetDeadline(time.Time) error       { return nil }
func (c *deadlineConn) SetReadDeadline(t time.Time) error { c.readDeadline = t; return nil }
func (c *deadlineConn) SetWriteDeadline(t time.Time) error {
	c.writeDeadline = t
	return nil
}

func TestAioSessionClockControlsDeadlines(t *testing.T) {
	fixed := time.Date(2026, 6, 6, 12, 0, 0, 0, time.UTC)
	cfg := NewSocketConfigWithOptions(
		WithReadTimeout(1000),
		WithWriteTimeout(2000),
		WithClock(func() time.Time { return fixed }),
	)

	writeConn := &deadlineConn{}
	writeSession := NewAioSession(writeConn, nil, cfg)
	if _, err := writeSession.Write([]byte("x")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if want := fixed.Add(2 * time.Second); !writeConn.writeDeadline.Equal(want) {
		t.Fatalf("write deadline = %v, want %v", writeConn.writeDeadline, want)
	}

	readConn := &deadlineConn{}
	readSession := NewAioSession(readConn, nil, cfg)
	if readSession.doRead() {
		t.Fatal("doRead should fail when fake connection returns EOF")
	}
	if want := fixed.Add(time.Second); !readConn.readDeadline.Equal(want) {
		t.Fatalf("read deadline = %v, want %v", readConn.readDeadline, want)
	}
}

func TestAioSessionCloseDuringReadCallbackKeepsBuffers(t *testing.T) {
	client, server := net.Pipe()
	defer closeAndReport(t, server.Close)

	entered := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)
	session := NewAioSession(client, &SimpleIoAction{OnDoAction: func(session *AioSession, data *bytes.Buffer) {
		if data == nil || data.String() != "hello" {
			done <- errors.New("unexpected read data")
			return
		}
		close(entered)
		<-release
		if session.ReadBuffer() == nil || session.WriteBuffer() == nil {
			done <- errors.New("session buffers were cleared while callback was active")
			return
		}
		done <- nil
	}}, nil)

	session.Read()
	if _, err := server.Write([]byte("hello")); err != nil {
		t.Fatalf("server write: %v", err)
	}
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("read callback was not entered")
	}
	closeAndReport(t, session.Close)
	close(release)
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("read callback did not finish")
	}
}

func TestAioSessionReadSerializesSharedBuffer(t *testing.T) {
	client, server := net.Pipe()
	defer closeAndReport(t, server.Close)
	defer closeAndReport(t, client.Close)

	var active int32
	var maxActive atomic.Int32
	var callbacks atomic.Int32
	release := make(chan struct{})
	session := NewAioSession(client, &SimpleIoAction{OnDoAction: func(session *AioSession, data *bytes.Buffer) {
		n := atomic.AddInt32(&active, 1)
		for {
			old := maxActive.Load()
			if n <= old || maxActive.CompareAndSwap(old, n) {
				break
			}
		}
		callbacks.Add(1)
		<-release
		atomic.AddInt32(&active, -1)
	}}, nil)

	session.Read().Read()
	var writes sync.WaitGroup
	writes.Add(2)
	go func() {
		defer writes.Done()
		_, _ = server.Write([]byte("one"))
	}()
	go func() {
		defer writes.Done()
		_, _ = server.Write([]byte("two"))
	}()
	waitForInt32(t, callbacks.Load, 1)
	time.Sleep(30 * time.Millisecond)
	if callbacks.Load() != 1 {
		t.Fatal("second callback overlapped the first read callback")
	}
	close(release)
	writes.Wait()
	waitForInt32(t, callbacks.Load, 2)
	if maxActive.Load() != 1 {
		t.Fatalf("max active callbacks = %d, want 1", maxActive.Load())
	}
}

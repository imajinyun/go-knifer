package socket

import (
	"io"
	"net"
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

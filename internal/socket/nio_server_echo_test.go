package socket

import (
	"net"
	"testing"
	"time"
)

// echoChannelHandler echoes data read from the connection.
type echoChannelHandler struct{}

func (h *echoChannelHandler) Handle(conn net.Conn) error {
	buf := make([]byte, 1024)
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}
	_, err = conn.Write(buf[:n])
	return err
}

func TestNioServerEcho(t *testing.T) {
	server, err := NewNioServerAddr(&net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})
	if err != nil {
		t.Fatal(err)
	}
	server.SetChannelHandler(&echoChannelHandler{})
	defer closeAndReport(t, server.Close)

	server.ListenAsync()

	addr := server.LocalAddr().(*net.TCPAddr)
	conn, err := net.DialTimeout("tcp", addr.String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer closeAndReport(t, conn.Close)

	want := []byte("hello-nio")
	if _, err := conn.Write(want); err != nil {
		t.Fatal(err)
	}
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	got := make([]byte, len(want))
	if _, err := conn.Read(got); err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Errorf("回显数据不一致: got=%q want=%q", got, want)
	}
}

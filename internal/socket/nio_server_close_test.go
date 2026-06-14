package socket

import (
	"net"
	"testing"
	"time"
)

func TestNioServerCloseClosesActiveConnections(t *testing.T) {
	client, serverConn := net.Pipe()
	defer closeAndReport(t, client.Close)
	listener := &queuedListener{addr: factoryFakeAddr("nio"), conns: make(chan net.Conn, 1)}
	listener.conns <- serverConn
	entered := make(chan struct{})
	nio, err := NewNioServerAddrWithOptions(&net.TCPAddr{Port: 1}, nil, WithListenerFactory(func(*net.TCPAddr) (net.Listener, error) {
		return listener, nil
	}))
	if err != nil {
		t.Fatalf("NewNioServerAddrWithOptions: %v", err)
	}
	nio.SetChannelHandler(ChannelHandlerFunc(func(conn net.Conn) error {
		select {
		case <-entered:
		default:
			close(entered)
		}
		buf := make([]byte, 1)
		_, err := conn.Read(buf)
		return err
	}))
	nio.ListenAsync()
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("nio server did not handle connection")
	}
	done := make(chan error, 1)
	go func() { done <- nio.Close() }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Close error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("NioServer.Close blocked with active connection")
	}
}

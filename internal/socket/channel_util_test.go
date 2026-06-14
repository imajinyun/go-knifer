package socket

import (
	"net"
	"testing"
)

func TestChannelUtilDialWithOptionsUsesDialer(t *testing.T) {
	dialer := &fakeDialer{}
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1234}
	conn, err := ChannelUtilDialWithOptions(addr, WithConnectDialer(dialer))
	if err != nil {
		t.Fatalf("ChannelUtilDialWithOptions failed: %v", err)
	}
	closeAndReport(t, conn.Close)
	closeAndReport(t, dialer.server.Close)
}

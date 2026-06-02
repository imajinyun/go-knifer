package socket

import (
	"net"
	"time"

	netimpl "github.com/imajinyun/go-knifer/internal/net"
)

// ChannelUtilDial opens a TCP connection; pool size should be handled by upper layers.
func ChannelUtilDial(addr *net.TCPAddr, timeout time.Duration) (net.Conn, error) {
	if addr == nil {
		return nil, NewSocketErrorMsg("address must not be nil")
	}
	conn, err := netimpl.Connect(addr.IP.String(), addr.Port, timeout)
	if err != nil {
		return nil, NewSocketError(err)
	}
	return conn, nil
}

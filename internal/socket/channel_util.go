package socket

import (
	"net"
	"time"
)

// ChannelUtilDial opens a TCP connection; pool size should be handled by upper layers.
func ChannelUtilDial(addr *net.TCPAddr, timeout time.Duration) (net.Conn, error) {
	return ChannelUtilDialWithOptions(addr, WithConnectTimeout(timeout))
}

// ChannelUtilDialWithOptions opens a TCP connection with custom dial options.
func ChannelUtilDialWithOptions(addr *net.TCPAddr, opts ...ConnectOption) (net.Conn, error) {
	if addr == nil {
		return nil, NewSocketErrorMsg("address must not be nil")
	}
	conn, err := ConnectAddrWithOptions(addr, opts...)
	if err != nil {
		return nil, NewSocketError(err)
	}
	return conn, nil
}

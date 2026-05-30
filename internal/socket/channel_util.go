package socket

import (
	"net"
	"time"
)

// ChannelUtilDial is aligned with hutool ChannelUtil.connect.
// In Go a synchronous Dial is enough; pool size should be handled by upper layers.
func ChannelUtilDial(addr *net.TCPAddr, timeout time.Duration) (net.Conn, error) {
	if addr == nil {
		return nil, NewSocketErrorMsg("address must not be nil")
	}
	if timeout <= 0 {
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			return nil, NewSocketError(err)
		}
		return conn, nil
	}
	conn, err := net.DialTimeout("tcp", addr.String(), timeout)
	if err != nil {
		return nil, NewSocketError(err)
	}
	return conn, nil
}

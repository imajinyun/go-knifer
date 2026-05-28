package socket

import (
	"net"
	"time"
)

// ChannelUtilDial 对应 hutool ChannelUtil.connect，
// 在 Go 中使用同步 Dial 即可，poolSize 用于上层并发限制。
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

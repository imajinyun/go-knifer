package socket

import (
	"net"
	"time"
)

// Connect 创建 Socket 并连接到指定地址，对应 hutool SocketUtil.connect。
// 当 timeout<=0 时，使用默认（无超时）连接。
func Connect(hostname string, port int, timeout time.Duration) (net.Conn, error) {
	address := net.JoinHostPort(hostname, itoa(port))
	if timeout <= 0 {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return nil, NewSocketError(err)
		}
		return conn, nil
	}
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, NewSocketError(err)
	}
	return conn, nil
}

// ConnectAddr 通过 net.TCPAddr 创建连接。
func ConnectAddr(addr *net.TCPAddr, timeout time.Duration) (net.Conn, error) {
	if addr == nil {
		return nil, NewSocketErrorMsg("address must not be nil")
	}
	return Connect(addr.IP.String(), addr.Port, timeout)
}

// GetRemoteAddress 获取远端地址，channel 为 nil 或未连接时返回 nil。
func GetRemoteAddress(conn net.Conn) net.Addr {
	if conn == nil {
		return nil
	}
	return conn.RemoteAddr()
}

// IsConnected 判断当前连接是否已建立（远端地址可获取）。
func IsConnected(conn net.Conn) bool {
	return GetRemoteAddress(conn) != nil
}

// itoa 简易整数转字符串，避免引入额外依赖。
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	if negative {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if negative {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

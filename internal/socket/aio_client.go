package socket

import (
	"bytes"
	"net"
	"time"
)

// AioClient 对应 hutool aio.AioClient，是基于 AIO 的 Socket 客户端。
type AioClient struct {
	session *AioSession
}

// NewAioClient 通过地址创建 AioClient，使用默认配置。
func NewAioClient(addr *net.TCPAddr, action IoAction[*bytes.Buffer]) (*AioClient, error) {
	return NewAioClientWithConfig(addr, action, NewSocketConfig())
}

// NewAioClientWithConfig 通过地址、IO 处理器、配置构造客户端。
func NewAioClientWithConfig(addr *net.TCPAddr, action IoAction[*bytes.Buffer], cfg *SocketConfig) (*AioClient, error) {
	if addr == nil {
		return nil, NewSocketErrorMsg("address must not be nil")
	}
	if action == nil {
		return nil, NewSocketErrorMsg("ioAction must not be nil")
	}
	if cfg == nil {
		cfg = NewSocketConfig()
	}
	conn, err := dialAio(addr, cfg)
	if err != nil {
		return nil, err
	}
	c := NewAioClientWithConn(conn, action, cfg)
	return c, nil
}

// NewAioClientWithConn 直接基于已建立的连接构造客户端。
func NewAioClientWithConn(conn net.Conn, action IoAction[*bytes.Buffer], cfg *SocketConfig) *AioClient {
	if cfg == nil {
		cfg = NewSocketConfig()
	}
	session := NewAioSession(conn, action, cfg)
	c := &AioClient{session: session}
	action.Accept(session)
	return c
}

// Session 返回底层会话。
func (c *AioClient) Session() *AioSession { return c.session }

// IoAction 返回当前 IO 处理器。
func (c *AioClient) IoAction() IoAction[*bytes.Buffer] {
	if c.session == nil {
		return nil
	}
	return c.session.IoAction()
}

// Read 触发一次异步读，读取结果会回调到 IoAction.DoAction。
func (c *AioClient) Read() *AioClient {
	if c.session != nil {
		c.session.Read()
	}
	return c
}

// Write 写出数据。
func (c *AioClient) Write(data []byte) (*AioClient, error) {
	if c.session == nil {
		return c, NewSocketErrorMsg("session is nil")
	}
	if _, err := c.session.Write(data); err != nil {
		return c, err
	}
	return c, nil
}

// Close 关闭客户端。
func (c *AioClient) Close() error {
	if c.session == nil {
		return nil
	}
	return c.session.Close()
}

// dialAio 内部建立连接，遵循配置中的写超时。
func dialAio(addr *net.TCPAddr, cfg *SocketConfig) (net.Conn, error) {
	timeout := time.Duration(cfg.WriteTimeout) * time.Millisecond
	return ChannelUtilDial(addr, timeout)
}

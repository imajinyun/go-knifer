package socket

import (
	"bytes"
	"net"
	"time"
)

// AioClient is an AIO-style socket client aligned with hutool aio.AioClient.
type AioClient struct {
	session *AioSession
}

// NewAioClient creates an AioClient with the default configuration.
func NewAioClient(addr *net.TCPAddr, action IoAction[*bytes.Buffer]) (*AioClient, error) {
	return NewAioClientWithConfig(addr, action, NewSocketConfig())
}

// NewAioClientWithConfig creates a client from an address, IO action, and configuration.
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

// NewAioClientWithConn creates a client from an established connection.
func NewAioClientWithConn(conn net.Conn, action IoAction[*bytes.Buffer], cfg *SocketConfig) *AioClient {
	if cfg == nil {
		cfg = NewSocketConfig()
	}
	session := NewAioSession(conn, action, cfg)
	c := &AioClient{session: session}
	action.Accept(session)
	return c
}

// Session returns the underlying session.
func (c *AioClient) Session() *AioSession { return c.session }

// IoAction returns the current IO action.
func (c *AioClient) IoAction() IoAction[*bytes.Buffer] {
	if c.session == nil {
		return nil
	}
	return c.session.IoAction()
}

// Read triggers one asynchronous read and passes the result to IoAction.DoAction.
func (c *AioClient) Read() *AioClient {
	if c.session != nil {
		c.session.Read()
	}
	return c
}

// Write writes data.
func (c *AioClient) Write(data []byte) (*AioClient, error) {
	if c.session == nil {
		return c, NewSocketErrorMsg("session is nil")
	}
	if _, err := c.session.Write(data); err != nil {
		return c, err
	}
	return c, nil
}

// Close closes the client.
func (c *AioClient) Close() error {
	if c.session == nil {
		return nil
	}
	return c.session.Close()
}

// dialAio creates the connection using the write timeout from the configuration.
func dialAio(addr *net.TCPAddr, cfg *SocketConfig) (net.Conn, error) {
	timeout := time.Duration(cfg.WriteTimeout) * time.Millisecond
	return ChannelUtilDial(addr, timeout)
}

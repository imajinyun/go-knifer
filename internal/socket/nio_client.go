package socket

import (
	"net"
	"sync"
	"sync/atomic"
)

// NioClient 对应 hutool nio.NioClient，是基于事件驱动的 TCP 客户端。
type NioClient struct {
	conn    net.Conn
	handler ChannelHandler

	closed atomic.Bool
	mu     sync.Mutex
	wg     sync.WaitGroup
}

// NewNioClient 创建并连接到指定 host:port。
func NewNioClient(host string, port int) (*NioClient, error) {
	return NewNioClientAddr(&net.TCPAddr{IP: net.ParseIP(host), Port: port})
}

// NewNioClientAddr 通过具体地址创建并连接。
func NewNioClientAddr(addr *net.TCPAddr) (*NioClient, error) {
	if addr == nil {
		return nil, NewSocketErrorMsg("address must not be nil")
	}
	c := &NioClient{}
	if err := c.init(addr); err != nil {
		return nil, err
	}
	return c, nil
}

// init 完成初始化连接。
func (c *NioClient) init(addr *net.TCPAddr) error {
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return NewSocketError(err)
	}
	c.conn = conn
	return nil
}

// SetChannelHandler 设置数据处理器。
func (c *NioClient) SetChannelHandler(handler ChannelHandler) *NioClient {
	c.handler = handler
	return c
}

// Channel 返回底层 net.Conn。
func (c *NioClient) Channel() net.Conn {
	return c.conn
}

// Listen 异步监听服务端推送的数据，等价于 hutool 的 listen()。
func (c *NioClient) Listen() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			if c.closed.Load() {
				return
			}
			if c.handler == nil {
				return
			}
			if err := c.handler.Handle(c.conn); err != nil {
				// 错误发生时直接退出监听（连接已关闭、读取失败等）
				return
			}
		}
	}()
}

// Write 写出多个数据片段。
func (c *NioClient) Write(datas ...[]byte) (*NioClient, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, d := range datas {
		if len(d) == 0 {
			continue
		}
		if _, err := c.conn.Write(d); err != nil {
			return c, NewSocketError(err)
		}
	}
	return c, nil
}

// Close 关闭客户端。
func (c *NioClient) Close() error {
	if c.closed.Swap(true) {
		return nil
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.wg.Wait()
	return nil
}

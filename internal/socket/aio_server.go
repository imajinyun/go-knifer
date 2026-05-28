package socket

import (
	"bytes"
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

// AioServer 对应 hutool aio.AioServer，基于异步 IO 风格的 Socket 服务端。
// Go 中通过 goroutine + 阻塞 Read 模拟 AIO 的回调语义。
type AioServer struct {
	listener net.Listener
	ioAction IoAction[*bytes.Buffer]
	config   *SocketConfig

	closed atomic.Bool
	wg     sync.WaitGroup
	mu     sync.Mutex
}

// NewAioServer 通过端口构造服务端。
func NewAioServer(port int) (*AioServer, error) {
	return NewAioServerAddr(&net.TCPAddr{Port: port}, NewSocketConfig())
}

// NewAioServerAddr 通过具体地址和配置构造服务端。
func NewAioServerAddr(addr *net.TCPAddr, cfg *SocketConfig) (*AioServer, error) {
	if addr == nil {
		return nil, NewSocketErrorMsg("address must not be nil")
	}
	if cfg == nil {
		cfg = NewSocketConfig()
	}
	s := &AioServer{config: cfg}
	if err := s.init(addr); err != nil {
		return nil, err
	}
	return s, nil
}

// init 初始化监听器。
func (s *AioServer) init(addr *net.TCPAddr) error {
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return NewSocketError(err)
	}
	s.listener = ln
	return nil
}

// SetIoAction 设置 IO 处理器。
func (s *AioServer) SetIoAction(action IoAction[*bytes.Buffer]) *AioServer {
	s.ioAction = action
	return s
}

// IoAction 返回 IO 处理器。
func (s *AioServer) IoAction() IoAction[*bytes.Buffer] {
	return s.ioAction
}

// Listener 返回底层监听器。
func (s *AioServer) Listener() net.Listener {
	return s.listener
}

// LocalAddr 返回本地监听地址。
func (s *AioServer) LocalAddr() net.Addr {
	if s.listener == nil {
		return nil
	}
	return s.listener.Addr()
}

// Config 返回配置。
func (s *AioServer) Config() *SocketConfig { return s.config }

// IsOpen 服务端是否仍在运行。
func (s *AioServer) IsOpen() bool { return !s.closed.Load() }

// Start 启动服务端，sync 表示是否阻塞当前 goroutine。
func (s *AioServer) Start(sync bool) {
	if sync {
		s.acceptLoop()
		return
	}
	go s.acceptLoop()
}

// acceptLoop 持续接收新连接。
func (s *AioServer) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() || errors.Is(err, net.ErrClosed) {
				return
			}
			if s.ioAction != nil {
				s.ioAction.Failed(NewSocketError(err), nil)
			}
			continue
		}
		s.handleAccept(conn)
	}
}

// handleAccept 为每个连接创建 AioSession 并触发回调。
func (s *AioServer) handleAccept(conn net.Conn) {
	if s.ioAction == nil {
		_ = conn.Close()
		return
	}
	session := NewAioSession(conn, s.ioAction, s.config)
	// 同步触发 Accept
	s.ioAction.Accept(session)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() { _ = session.Close() }()
		// 持续读取，模拟 AIO 链式回调
		for session.IsOpen() && !s.closed.Load() {
			if !session.doRead() {
				return
			}
		}
	}()
}

// Close 关闭服务端。
func (s *AioServer) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed.Swap(true) {
		return nil
	}
	if s.listener != nil {
		_ = s.listener.Close()
	}
	s.wg.Wait()
	return nil
}

package socket

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

// NioServer 对应 hutool nio.NioServer，是基于事件驱动的 TCP 服务端。
// Go 中以 goroutine + 阻塞 Accept/Read 实现等效语义。
type NioServer struct {
	listener net.Listener
	handler  ChannelHandler
	addr     *net.TCPAddr

	closed atomic.Bool
	wg     sync.WaitGroup
	mu     sync.Mutex
}

// NewNioServer 通过端口创建 NioServer 并完成初始化。
func NewNioServer(port int) (*NioServer, error) {
	return NewNioServerAddr(&net.TCPAddr{Port: port})
}

// NewNioServerAddr 通过具体地址创建 NioServer。
func NewNioServerAddr(addr *net.TCPAddr) (*NioServer, error) {
	if addr == nil {
		return nil, NewSocketErrorMsg("address must not be nil")
	}
	s := &NioServer{addr: addr}
	if err := s.init(addr); err != nil {
		return nil, err
	}
	return s, nil
}

// init 初始化监听器。
func (s *NioServer) init(addr *net.TCPAddr) error {
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return NewSocketError(err)
	}
	s.listener = ln
	return nil
}

// SetChannelHandler 设置数据处理器。
func (s *NioServer) SetChannelHandler(handler ChannelHandler) *NioServer {
	s.handler = handler
	return s
}

// Listener 暴露底层 net.Listener。
func (s *NioServer) Listener() net.Listener {
	return s.listener
}

// LocalAddr 返回本地监听地址，便于动态端口测试。
func (s *NioServer) LocalAddr() net.Addr {
	if s.listener == nil {
		return nil
	}
	return s.listener.Addr()
}

// Start 等价于 hutool 中的 start()，开始监听并阻塞当前 goroutine。
func (s *NioServer) Start() {
	s.Listen()
}

// Listen 同步阻塞监听，等同于 hutool 的 listen()。
func (s *NioServer) Listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			if errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}
		s.handleAccept(conn)
	}
}

// ListenAsync 异步启动监听，返回的 channel 在监听结束时关闭。
func (s *NioServer) ListenAsync() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		s.Listen()
	}()
	return done
}

// handleAccept 在新 goroutine 中持续处理来自该连接的"读事件"。
func (s *NioServer) handleAccept(conn net.Conn) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() { _ = conn.Close() }()
		if s.handler == nil {
			return
		}
		for {
			if s.closed.Load() {
				return
			}
			// 模拟 NIO 读事件：每当连接可读时回调 handler。
			// handler 内部通常会使用一次 conn.Read 来消费数据。
			if err := s.handler.Handle(conn); err != nil {
				return
			}
		}
	}()
}

// Close 关闭服务端。
func (s *NioServer) Close() error {
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

// IsOpen 服务端是否仍在运行。
func (s *NioServer) IsOpen() bool {
	return !s.closed.Load()
}

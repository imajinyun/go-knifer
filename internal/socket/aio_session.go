package socket

import (
	"bytes"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// AioSession 对应 hutool aio.AioSession，每个客户端连接对应一个会话对象。
// 会话内部维护读/写缓冲，并以异步回调的方式触发 IoAction.DoAction。
type AioSession struct {
	conn        net.Conn
	ioAction    IoAction[*bytes.Buffer]
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
	scratch     []byte // 用于一次 Read 的临时空间

	readTimeout  time.Duration
	writeTimeout time.Duration

	closed atomic.Bool
	mu     sync.Mutex
}

// NewAioSession 构造 AioSession。
func NewAioSession(conn net.Conn, ioAction IoAction[*bytes.Buffer], cfg *SocketConfig) *AioSession {
	if cfg == nil {
		cfg = NewSocketConfig()
	}
	readSize := cfg.ReadBufferSize
	if readSize <= 0 {
		readSize = DefaultBufferSize
	}
	writeSize := cfg.WriteBufferSize
	if writeSize <= 0 {
		writeSize = DefaultBufferSize
	}
	return &AioSession{
		conn:         conn,
		ioAction:     ioAction,
		readBuffer:   bytes.NewBuffer(make([]byte, 0, readSize)),
		writeBuffer:  bytes.NewBuffer(make([]byte, 0, writeSize)),
		scratch:      make([]byte, readSize),
		readTimeout:  time.Duration(cfg.ReadTimeout) * time.Millisecond,
		writeTimeout: time.Duration(cfg.WriteTimeout) * time.Millisecond,
	}
}

// Conn 返回底层连接。
func (s *AioSession) Conn() net.Conn { return s.conn }

// ReadBuffer 返回读缓冲。
func (s *AioSession) ReadBuffer() *bytes.Buffer { return s.readBuffer }

// WriteBuffer 返回写缓冲。
func (s *AioSession) WriteBuffer() *bytes.Buffer { return s.writeBuffer }

// IoAction 返回 IO 处理器。
func (s *AioSession) IoAction() IoAction[*bytes.Buffer] { return s.ioAction }

// RemoteAddress 返回远端地址。
func (s *AioSession) RemoteAddress() net.Addr {
	return GetRemoteAddress(s.conn)
}

// Read 异步读取一次数据，等价于 hutool 的 read()。
// 在 Go 中通过 goroutine 实现"异步" + 完成后回调 IoAction。
func (s *AioSession) Read() *AioSession {
	if !s.IsOpen() {
		return s
	}
	go s.doRead()
	return s
}

// doRead 执行一次读取并回调，返回 false 表示读取失败/连接已关闭。
func (s *AioSession) doRead() bool {
	if !s.IsOpen() {
		return false
	}
	if s.readTimeout > 0 {
		_ = s.conn.SetReadDeadline(time.Now().Add(s.readTimeout))
	} else {
		_ = s.conn.SetReadDeadline(time.Time{})
	}
	n, err := s.conn.Read(s.scratch)
	if err != nil {
		if s.ioAction != nil {
			s.ioAction.Failed(err, s)
		}
		_ = s.Close()
		return false
	}
	s.readBuffer.Reset()
	s.readBuffer.Write(s.scratch[:n])
	s.callbackRead()
	return true
}

// callbackRead 执行 IoAction 的 DoAction 回调。
func (s *AioSession) callbackRead() {
	if s.ioAction != nil {
		s.ioAction.DoAction(s, s.readBuffer)
	}
}

// Write 写出数据。
func (s *AioSession) Write(data []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.IsOpen() {
		return 0, NewSocketErrorMsg("session is closed")
	}
	if s.writeTimeout > 0 {
		_ = s.conn.SetWriteDeadline(time.Now().Add(s.writeTimeout))
	} else {
		_ = s.conn.SetWriteDeadline(time.Time{})
	}
	n, err := s.conn.Write(data)
	if err != nil {
		return n, NewSocketError(err)
	}
	return n, nil
}

// WriteAndClose 写出数据并关闭写方向。
func (s *AioSession) WriteAndClose(data []byte) error {
	if _, err := s.Write(data); err != nil {
		return err
	}
	return s.CloseOut()
}

// IsOpen 会话是否仍然打开。
func (s *AioSession) IsOpen() bool {
	return s.conn != nil && !s.closed.Load()
}

// CloseIn 关闭读方向。
func (s *AioSession) CloseIn() error {
	if tc, ok := s.conn.(*net.TCPConn); ok {
		if err := tc.CloseRead(); err != nil {
			return NewSocketError(err)
		}
	}
	return nil
}

// CloseOut 关闭写方向。
func (s *AioSession) CloseOut() error {
	if tc, ok := s.conn.(*net.TCPConn); ok {
		if err := tc.CloseWrite(); err != nil {
			return NewSocketError(err)
		}
	}
	return nil
}

// Close 关闭会话。
func (s *AioSession) Close() error {
	if s.closed.Swap(true) {
		return nil
	}
	if s.conn != nil {
		_ = s.conn.Close()
	}
	s.readBuffer = nil
	s.writeBuffer = nil
	return nil
}

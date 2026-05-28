package socket

import "net"

// ChannelHandler 对应 hutool nio.ChannelHandler，
// 用户实现该接口即可处理对应的 Conn（读、写）。
type ChannelHandler interface {
	Handle(conn net.Conn) error
}

// ChannelHandlerFunc 函数式 ChannelHandler。
type ChannelHandlerFunc func(conn net.Conn) error

// Handle 实现 ChannelHandler。
func (f ChannelHandlerFunc) Handle(conn net.Conn) error { return f(conn) }

// Operation 对应 hutool nio.Operation，
// 在 Go 中以常量形式标识感兴趣的事件类型，主要用于语义对齐。
type Operation int

const (
	// OpRead 读操作
	OpRead Operation = 1 << iota
	// OpWrite 写操作
	OpWrite
	// OpConnect 连接操作
	OpConnect
	// OpAccept 接受连接操作
	OpAccept
)

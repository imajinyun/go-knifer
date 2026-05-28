// Package gksocket 对应 hutool-socket，提供基于 NIO/AIO 风格的 Socket 通信封装。
//
// 本包包含以下主要类型：
//   - SocketConfig：通讯配置（线程池大小、超时、缓冲区等）。
//   - SocketRuntimeError：Socket 运行时错误。
//   - SocketUtil（Connect / GetRemoteAddress / IsConnected）：Socket 通用工具。
//   - ChannelUtilDial：建立连接的工具方法。
//   - Protocol/MsgEncoder/MsgDecoder：消息编解码协议接口。
//   - NioServer / NioClient / ChannelHandler：NIO 风格的服务端/客户端。
//   - AioServer / AioClient / AioSession / IoAction / SimpleIoAction：AIO 风格的服务端/客户端。
package socket

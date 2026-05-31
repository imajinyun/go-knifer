// Package socket provides NIO/AIO-style socket communication helpers.
//
// Main types in this package:
//   - SocketConfig: communication options such as pool size, timeout, and buffers.
//   - SocketRuntimeError: socket runtime errors.
//   - SocketUtil (Connect / GetRemoteAddress / IsConnected): common socket helpers.
//   - ChannelUtilDial: helper for establishing connections.
//   - Protocol/MsgEncoder/MsgDecoder: message encoding and decoding interfaces.
//   - NioServer / NioClient / ChannelHandler: NIO-style server/client helpers.
//   - AioServer / AioClient / AioSession / IoAction / SimpleIoAction: AIO-style server/client helpers.
package socket

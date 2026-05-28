package socket

import "bytes"

// IoAction 对应 hutool aio.IoAction。
// 用户实现该接口以处理 Socket 会话的各类事件。
// T 表示经过解码后的业务数据类型；如果不需要解码可直接使用 []byte 或 *bytes.Buffer。
type IoAction[T any] interface {
	// Accept 接收客户端连接（会话建立）事件。
	Accept(session *AioSession)
	// DoAction 执行数据处理（消息读取）。
	DoAction(session *AioSession, data T)
	// Failed 数据读取失败的回调。
	Failed(err error, session *AioSession)
}

// SimpleIoAction 对应 hutool SimpleIoAction，
// 提供 Accept 与 Failed 的默认实现，业务方仅需实现 DoAction。
type SimpleIoAction struct {
	OnAccept   func(session *AioSession)
	OnDoAction func(session *AioSession, data *bytes.Buffer)
	OnFailed   func(err error, session *AioSession)
}

// Accept 默认空实现。
func (s *SimpleIoAction) Accept(session *AioSession) {
	if s.OnAccept != nil {
		s.OnAccept(session)
	}
}

// DoAction 调用注册的回调。
func (s *SimpleIoAction) DoAction(session *AioSession, data *bytes.Buffer) {
	if s.OnDoAction != nil {
		s.OnDoAction(session, data)
	}
}

// Failed 默认仅记录错误，可通过回调覆盖。
func (s *SimpleIoAction) Failed(err error, session *AioSession) {
	if s.OnFailed != nil {
		s.OnFailed(err, session)
	}
}

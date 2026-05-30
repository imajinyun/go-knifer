package socket

import "bytes"

// IoAction is aligned with hutool aio.IoAction.
// Implement this interface to handle socket session events.
// T is the decoded business data type; use []byte or *bytes.Buffer directly when no decoding is needed.
type IoAction[T any] interface {
	// Accept handles a client connection event when the session is established.
	Accept(session *AioSession)
	// DoAction handles data after a message is read.
	DoAction(session *AioSession, data T)
	// Failed handles read failures.
	Failed(err error, session *AioSession)
}

// SimpleIoAction is aligned with hutool SimpleIoAction.
// It provides default Accept and Failed behavior, so callers only need to implement DoAction.
type SimpleIoAction struct {
	OnAccept   func(session *AioSession)
	OnDoAction func(session *AioSession, data *bytes.Buffer)
	OnFailed   func(err error, session *AioSession)
}

// Accept invokes the configured accept callback when present.
func (s *SimpleIoAction) Accept(session *AioSession) {
	if s.OnAccept != nil {
		s.OnAccept(session)
	}
}

// DoAction invokes the configured action callback.
func (s *SimpleIoAction) DoAction(session *AioSession, data *bytes.Buffer) {
	if s.OnDoAction != nil {
		s.OnDoAction(session, data)
	}
}

// Failed invokes the configured failure callback when present.
func (s *SimpleIoAction) Failed(err error, session *AioSession) {
	if s.OnFailed != nil {
		s.OnFailed(err, session)
	}
}

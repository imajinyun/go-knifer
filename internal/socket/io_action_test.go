package socket

import (
	"bytes"
	"testing"
)

func TestSimpleIoAction(t *testing.T) {
	var (
		acceptCalled bool
		failed       error
		received     []byte
	)
	action := &SimpleIoAction{
		OnAccept: func(session *AioSession) { acceptCalled = true },
		OnDoAction: func(session *AioSession, data *bytes.Buffer) {
			received = append(received, data.Bytes()...)
		},
		OnFailed: func(err error, session *AioSession) { failed = err },
	}

	action.Accept(nil)
	action.DoAction(nil, bytes.NewBufferString("hi"))
	action.Failed(NewSocketErrorMsg("oops"), nil)

	if !acceptCalled {
		t.Errorf("OnAccept 未被调用")
	}
	if string(received) != "hi" {
		t.Errorf("OnDoAction 数据错误: %q", received)
	}
	if failed == nil {
		t.Errorf("OnFailed 未传递错误")
	}
}

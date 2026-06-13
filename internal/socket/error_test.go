package socket

import (
	"errors"
	"net"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

func TestSocketRuntimeError(t *testing.T) {
	base := net.ErrClosed
	e := WrapSocketError(base, "wrapped")
	if e == nil {
		t.Fatal("WrapSocketError 不应返回 nil")
	}
	if e.Unwrap() != base {
		t.Errorf("Unwrap 失败")
	}
	if e.Error() == "" {
		t.Errorf("Error 不应为空")
	}
	if !errors.Is(e, knifer.ErrCodeInternal) {
		t.Errorf("SocketRuntimeError 应匹配 ErrCodeInternal")
	}
	if !errors.Is(e, base) {
		t.Errorf("SocketRuntimeError 应保留 cause 链")
	}
	if WrapSocketError(nil, "x") != nil {
		t.Errorf("nil err 应返回 nil")
	}
	if NewSocketErrorf("hello %s", "world").Error() != "hello world" {
		t.Errorf("格式化失败")
	}
}

package socket

import (
	"bytes"
	"errors"
	"net"
	"testing"
	"time"
)

func TestAioSessionCloseDuringReadCallbackKeepsBuffers(t *testing.T) {
	client, server := net.Pipe()
	defer closeAndReport(t, server.Close)

	entered := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)
	session := NewAioSession(client, &SimpleIoAction{OnDoAction: func(session *AioSession, data *bytes.Buffer) {
		if data == nil || data.String() != "hello" {
			done <- errors.New("unexpected read data")
			return
		}
		close(entered)
		<-release
		if session.ReadBuffer() == nil || session.WriteBuffer() == nil {
			done <- errors.New("session buffers were cleared while callback was active")
			return
		}
		done <- nil
	}}, nil)

	session.Read()
	if _, err := server.Write([]byte("hello")); err != nil {
		t.Fatalf("server write: %v", err)
	}
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("read callback was not entered")
	}
	closeAndReport(t, session.Close)
	close(release)
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("read callback did not finish")
	}
}

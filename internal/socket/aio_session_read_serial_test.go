package socket

import (
	"bytes"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAioSessionReadSerializesSharedBuffer(t *testing.T) {
	client, server := net.Pipe()
	defer closeAndReport(t, server.Close)
	defer closeAndReport(t, client.Close)

	var active int32
	var maxActive atomic.Int32
	var callbacks atomic.Int32
	release := make(chan struct{})
	session := NewAioSession(client, &SimpleIoAction{OnDoAction: func(session *AioSession, data *bytes.Buffer) {
		n := atomic.AddInt32(&active, 1)
		for {
			old := maxActive.Load()
			if n <= old || maxActive.CompareAndSwap(old, n) {
				break
			}
		}
		callbacks.Add(1)
		<-release
		atomic.AddInt32(&active, -1)
	}}, nil)

	session.Read().Read()
	var writes sync.WaitGroup
	writes.Add(2)
	go func() {
		defer writes.Done()
		_, _ = server.Write([]byte("one"))
	}()
	go func() {
		defer writes.Done()
		_, _ = server.Write([]byte("two"))
	}()
	waitForInt32(t, callbacks.Load, 1)
	time.Sleep(30 * time.Millisecond)
	if callbacks.Load() != 1 {
		t.Fatal("second callback overlapped the first read callback")
	}
	close(release)
	writes.Wait()
	waitForInt32(t, callbacks.Load, 2)
	if maxActive.Load() != 1 {
		t.Fatalf("max active callbacks = %d, want 1", maxActive.Load())
	}
}

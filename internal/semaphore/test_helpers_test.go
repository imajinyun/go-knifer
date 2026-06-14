package semaphore

import (
	"testing"
	"time"
)

func assertNoAcquire(t *testing.T, ch <-chan error) {
	t.Helper()
	select {
	case err := <-ch:
		t.Fatalf("unexpected acquire result: %v", err)
	case <-time.After(30 * time.Millisecond):
	}
}

func assertPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}

func waitUntil(t *testing.T, fn func() bool) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if fn() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("condition not met before deadline")
}

func queueLen(sem *Semaphore) int {
	sem.mux.Lock()
	defer sem.mux.Unlock()
	return len(sem.queues)
}

package semaphore

import (
	"context"
	"errors"
	"testing"
)

func TestAcquireWaitsInFIFOOrder(t *testing.T) {
	sem := New(2)
	if err := sem.Acquire(context.Background(), 2); err != nil {
		t.Fatal(err)
	}

	first := make(chan error, 1)
	second := make(chan error, 1)
	go func() { first <- sem.Acquire(context.Background(), 2) }()
	waitUntil(t, func() bool { return queueLen(sem) == 1 })
	go func() { second <- sem.Acquire(context.Background(), 1) }()
	waitUntil(t, func() bool { return queueLen(sem) == 2 })

	sem.Release(1)
	assertNoAcquire(t, first)
	assertNoAcquire(t, second)

	sem.Release(1)
	if err := <-first; err != nil {
		t.Fatalf("first acquire error = %v", err)
	}
	assertNoAcquire(t, second)

	sem.Release(2)
	if err := <-second; err != nil {
		t.Fatalf("second acquire error = %v", err)
	}
	sem.Release(1)
}

func TestAcquireContextCancelRemovesWaiterAndNotifiesNext(t *testing.T) {
	sem := New(2)
	if err := sem.Acquire(context.Background(), 2); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	first := make(chan error, 1)
	second := make(chan error, 1)
	go func() { first <- sem.Acquire(ctx, 2) }()
	waitUntil(t, func() bool { return queueLen(sem) == 1 })
	go func() { second <- sem.Acquire(context.Background(), 1) }()
	waitUntil(t, func() bool { return queueLen(sem) == 2 })

	cancel()
	if err := <-first; !errors.Is(err, context.Canceled) {
		t.Fatalf("first acquire error = %v, want context.Canceled", err)
	}

	sem.Release(1)
	if err := <-second; err != nil {
		t.Fatalf("second acquire error = %v", err)
	}
	sem.Release(1)
}

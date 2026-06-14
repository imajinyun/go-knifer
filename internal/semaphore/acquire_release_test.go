package semaphore

import (
	"context"
	"testing"
)

func TestAcquireTryAcquireRelease(t *testing.T) {
	sem := New(2)

	if sem.Cap() != 2 {
		t.Fatalf("Cap() = %d, want 2", sem.Cap())
	}
	if err := sem.Acquire(context.Background(), 2); err != nil {
		t.Fatalf("Acquire() error = %v", err)
	}
	if sem.Use() != 2 {
		t.Fatalf("Use() = %d, want 2", sem.Use())
	}
	if sem.TryAcquire(1) {
		t.Fatal("TryAcquire() should fail when capacity is full")
	}

	sem.Release(1)
	if sem.Use() != 1 {
		t.Fatalf("Use() after Release() = %d, want 1", sem.Use())
	}
	if !sem.TryAcquire(1) {
		t.Fatal("TryAcquire() should succeed after one permit is released")
	}
	sem.Release(2)
	if sem.Use() != 0 {
		t.Fatalf("Use() after final Release() = %d, want 0", sem.Use())
	}
}

func TestReleaseEReturnsErrorAndReleasesPermits(t *testing.T) {
	sem := New(2)
	if err := sem.Acquire(context.Background(), 2); err != nil {
		t.Fatal(err)
	}
	if err := sem.ReleaseE(1); err != nil {
		t.Fatalf("ReleaseE(1) error = %v", err)
	}
	if sem.Use() != 1 {
		t.Fatalf("Use() after ReleaseE(1) = %d, want 1", sem.Use())
	}
	if err := sem.ReleaseE(1); err != nil {
		t.Fatalf("ReleaseE(1) final error = %v", err)
	}
	if sem.Use() != 0 {
		t.Fatalf("Use() after final ReleaseE(1) = %d, want 0", sem.Use())
	}
}

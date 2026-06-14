package http

import (
	"testing"
	"time"
)

func TestGlobalTimeout(t *testing.T) {
	old := GetGlobalTimeout()
	defer SetGlobalTimeout(old)

	SetGlobalTimeout(7 * time.Second)
	if got := GetGlobalTimeout(); got != 7*time.Second {
		t.Fatalf("timeout: %v", got)
	}
}

func TestDefaultGlobalTimeoutIsBounded(t *testing.T) {
	previous := SnapshotGlobalConfig()
	defer ConfigureGlobalConfig(previous)

	ResetGlobalConfig()
	if got := GetGlobalTimeout(); got != defaultGlobalTimeout || got <= 0 {
		t.Fatalf("default timeout = %v, want positive %v", got, defaultGlobalTimeout)
	}
}

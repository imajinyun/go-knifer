package vsys_test

import (
	"runtime"
	"testing"

	"github.com/imajinyun/go-knifer/vsys"
)

func TestFacadePID(t *testing.T) {
	pid := vsys.CurrentPID()
	if pid <= 0 {
		t.Fatalf("expected positive pid, got %d", pid)
	}
	if got := vsys.CurrentPIDWithOptions(vsys.WithPIDFunc(func() int { return 99 })); got != 99 {
		t.Fatalf("CurrentPIDWithOptions = %d", got)
	}
}

func TestFacadeMemory(t *testing.T) {
	total := vsys.TotalMemory()
	free := vsys.FreeMemory()
	max := vsys.MaxMemory()
	if total == 0 && free == 0 && max == 0 {
		t.Fatal("expected at least one memory metric to be non-zero")
	}
	opt := vsys.WithReadMemStatsFunc(func(stats *runtime.MemStats) {
		stats.Sys = 300
		stats.HeapSys = 200
		stats.HeapIdle = 100
	})
	if vsys.MaxMemoryWithOptions(opt) != 300 || vsys.TotalMemoryWithOptions(opt) != 200 || vsys.FreeMemoryWithOptions(opt) != 100 {
		t.Fatal("expected memory option providers to be used")
	}
}

func TestFacadeGoroutineCount(t *testing.T) {
	count := vsys.TotalGoroutineCount()
	if count < 1 {
		t.Fatalf("expected at least 1 goroutine, got %d", count)
	}
	if got := vsys.TotalGoroutineCountWithOptions(vsys.WithProcessNumGoroutineFunc(func() int { return 12 })); got != 12 {
		t.Fatalf("TotalGoroutineCountWithOptions = %d", got)
	}
}

package system

import (
	"runtime"
	"strings"
)

// RuntimeInfo describes current Go runtime memory usage and goroutine statistics.
type RuntimeInfo struct {
	stats runtime.MemStats
}

// NewRuntimeInfo creates RuntimeInfo.
func NewRuntimeInfo() *RuntimeInfo {
	r := &RuntimeInfo{}
	runtime.ReadMemStats(&r.stats)
	return r
}

// Refresh updates memory statistics.
func (r *RuntimeInfo) Refresh() *RuntimeInfo {
	runtime.ReadMemStats(&r.stats)
	return r
}

// GetMaxMemory returns the Go process memory upper bound approximation.
func (r *RuntimeInfo) GetMaxMemory() uint64 { return r.stats.Sys }

// GetTotalMemory returns total memory currently requested from the OS.
func (r *RuntimeInfo) GetTotalMemory() uint64 { return r.stats.HeapSys }

// GetFreeMemory returns currently idle heap memory.
func (r *RuntimeInfo) GetFreeMemory() uint64 { return r.stats.HeapIdle }

// GetUsableMemory returns an approximation of usable memory as Sys - HeapInuse.
func (r *RuntimeInfo) GetUsableMemory() uint64 {
	if r.stats.Sys < r.stats.HeapInuse {
		return 0
	}
	return r.stats.Sys - r.stats.HeapInuse
}

// GetHeapInuse returns currently in-use heap memory.
func (r *RuntimeInfo) GetHeapInuse() uint64 { return r.stats.HeapInuse }

// GetGoroutineCount returns the current goroutine count.
func (r *RuntimeInfo) GetGoroutineCount() int { return runtime.NumGoroutine() }

// GetMemStats returns a copy of the underlying MemStats.
func (r *RuntimeInfo) GetMemStats() runtime.MemStats { return r.stats }

// String implements fmt.Stringer.
func (r *RuntimeInfo) String() string {
	var b strings.Builder
	appendLine(&b, "Max Memory:        ", readableSize(r.GetMaxMemory()))
	appendLine(&b, "Total Memory:      ", readableSize(r.GetTotalMemory()))
	appendLine(&b, "Free Memory:       ", readableSize(r.GetFreeMemory()))
	appendLine(&b, "Usable Memory:     ", readableSize(r.GetUsableMemory()))
	appendLine(&b, "Heap In Use:       ", readableSize(r.GetHeapInuse()))
	appendLine(&b, "Goroutine Count:   ", r.GetGoroutineCount())
	return b.String()
}

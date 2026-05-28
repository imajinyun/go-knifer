package system

import (
	"runtime"
	"strings"
)

// RuntimeInfo 对应 hutool RuntimeInfo + 部分 JVM 内存信息。
// 提供当前 Go 运行时的内存使用与协程数等信息。
type RuntimeInfo struct {
	stats runtime.MemStats
}

// NewRuntimeInfo 构造 RuntimeInfo。
func NewRuntimeInfo() *RuntimeInfo {
	r := &RuntimeInfo{}
	runtime.ReadMemStats(&r.stats)
	return r
}

// Refresh 刷新内存统计信息。
func (r *RuntimeInfo) Refresh() *RuntimeInfo {
	runtime.ReadMemStats(&r.stats)
	return r
}

// GetMaxMemory 返回 Go 进程的内存上限（HeapSys，近似最大）。
func (r *RuntimeInfo) GetMaxMemory() uint64 { return r.stats.Sys }

// GetTotalMemory 返回当前已经从 OS 申请的内存总量。
func (r *RuntimeInfo) GetTotalMemory() uint64 { return r.stats.HeapSys }

// GetFreeMemory 返回当前空闲内存（HeapIdle）。
func (r *RuntimeInfo) GetFreeMemory() uint64 { return r.stats.HeapIdle }

// GetUsableMemory 返回 Sys - HeapInuse 近似的可用内存。
func (r *RuntimeInfo) GetUsableMemory() uint64 {
	if r.stats.Sys < r.stats.HeapInuse {
		return 0
	}
	return r.stats.Sys - r.stats.HeapInuse
}

// GetHeapInuse 返回当前堆已使用内存。
func (r *RuntimeInfo) GetHeapInuse() uint64 { return r.stats.HeapInuse }

// GetGoroutineCount 返回当前协程数。
func (r *RuntimeInfo) GetGoroutineCount() int { return runtime.NumGoroutine() }

// GetMemStats 返回底层 MemStats（拷贝）。
func (r *RuntimeInfo) GetMemStats() runtime.MemStats { return r.stats }

// String 实现 fmt.Stringer。
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

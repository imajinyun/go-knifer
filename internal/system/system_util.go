package system

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
)

// 单例缓存，避免重复采集。
var (
	hostOnce    sync.Once
	hostInfo    *HostInfo
	osOnce      sync.Once
	osInfo      *OsInfo
	userOnce    sync.Once
	userInfo    *UserInfo
	goOnce      sync.Once
	goInfo      *GoInfo
	runtimeOnce sync.Once
	runtimeRef  *RuntimeInfo
)

// GetHostInfo 取得主机信息（单例）。
func GetHostInfo() *HostInfo {
	hostOnce.Do(func() { hostInfo = NewHostInfo() })
	return hostInfo
}

// GetOsInfo 取得 OS 信息（单例）。
func GetOsInfo() *OsInfo {
	osOnce.Do(func() { osInfo = NewOsInfo() })
	return osInfo
}

// GetUserInfo 取得用户信息（单例）。
func GetUserInfo() *UserInfo {
	userOnce.Do(func() { userInfo = NewUserInfo() })
	return userInfo
}

// GetGoInfo 取得 Go 运行时元信息（单例）。
func GetGoInfo() *GoInfo {
	goOnce.Do(func() { goInfo = NewGoInfo() })
	return goInfo
}

// GetRuntimeInfo 取得运行时内存信息。每次调用都会刷新。
func GetRuntimeInfo() *RuntimeInfo {
	runtimeOnce.Do(func() { runtimeRef = NewRuntimeInfo() })
	return runtimeRef.Refresh()
}

// GetCurrentPID 返回当前进程 PID。
func GetCurrentPID() int {
	return os.Getpid()
}

// GetTotalMemory 当前 Go 程序从 OS 申请的内存总量。
func GetTotalMemory() uint64 {
	return GetRuntimeInfo().GetTotalMemory()
}

// GetFreeMemory 当前 Go 程序的空闲内存。
func GetFreeMemory() uint64 {
	return GetRuntimeInfo().GetFreeMemory()
}

// GetMaxMemory 当前 Go 程序的最大内存（系统内存上限）。
func GetMaxMemory() uint64 {
	return GetRuntimeInfo().GetMaxMemory()
}

// GetTotalThreadCount 取得总协程数（对应 hutool 的总线程数）。
func GetTotalThreadCount() int {
	return runtime.NumGoroutine()
}

// Get 通过 key 取得环境变量；若 quiet=false 且变量缺失，将打印警告信息到 stderr。
// 与 hutool SystemUtil.get(key, quiet) 行为对应。
func Get(key string, quiet bool) string {
	v, ok := os.LookupEnv(key)
	if !ok && !quiet {
		fmt.Fprintf(os.Stderr, "[gksystem] env %q not found\n", key)
	}
	return v
}

// GetOrDefault 取得环境变量，若不存在或为空，返回 def。
func GetOrDefault(key, def string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def
	}
	return v
}

// GetInt 取得环境变量并转为 int，转换失败返回 def。
func GetInt(key string, def int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// GetBool 取得环境变量并转为 bool，转换失败返回 def。
func GetBool(key string, def bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

// DumpSystemInfo 将系统信息输出到 stdout，对应 hutool SystemUtil.dumpSystemInfo()。
func DumpSystemInfo() {
	DumpSystemInfoTo(os.Stdout)
}

// DumpSystemInfoTo 将系统信息输出到指定 Writer。
func DumpSystemInfoTo(w io.Writer) {
	const sep = "--------------\n"
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetGoInfo())
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetOsInfo())
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetUserInfo())
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetHostInfo())
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetRuntimeInfo())
	_, _ = fmt.Fprint(w, sep)
}

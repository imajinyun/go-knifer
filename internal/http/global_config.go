package http

import (
	"sync"
	"time"
)

// 全局默认配置（对应 hutool-http HttpGlobalConfig）。
var (
	globalMu               sync.RWMutex
	globalTimeout          = 0 * time.Second // 0 表示使用 HTTP 客户端默认值
	globalMaxRedirects     = 10
	globalIgnoreEOFError   = true
	globalDecodeURL        = false
	globalFollowRedirects  = true
	globalDefaultUserAgent = ""
	globalTrustAnyHost     = false
	globalBoundary         = "--------------------gokitFormBoundary"
)

// SetGlobalTimeout 设置全局默认超时。
func SetGlobalTimeout(d time.Duration) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalTimeout = d
}

// GetGlobalTimeout 获取全局默认超时。
func GetGlobalTimeout() time.Duration {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalTimeout
}

// SetGlobalMaxRedirects 设置全局最大重定向次数。
func SetGlobalMaxRedirects(n int) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalMaxRedirects = n
}

// GetGlobalMaxRedirects 获取全局最大重定向次数。
func GetGlobalMaxRedirects() int {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalMaxRedirects
}

// SetGlobalFollowRedirects 设置是否跟随重定向。
func SetGlobalFollowRedirects(b bool) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalFollowRedirects = b
}

// GetGlobalFollowRedirects 获取是否跟随重定向。
func GetGlobalFollowRedirects() bool {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalFollowRedirects
}

// SetGlobalUserAgent 设置全局默认 User-Agent。
func SetGlobalUserAgent(ua string) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalDefaultUserAgent = ua
}

// GetGlobalUserAgent 获取全局默认 User-Agent。
func GetGlobalUserAgent() string {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalDefaultUserAgent
}

// SetIgnoreEOFError 设置是否忽略 EOF 错误。
func SetIgnoreEOFError(b bool) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalIgnoreEOFError = b
}

// IsIgnoreEOFError 是否忽略 EOF 错误。
func IsIgnoreEOFError() bool {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalIgnoreEOFError
}

// SetTrustAnyHost 设置是否信任所有主机（HTTPS 跳过证书校验）。
func SetTrustAnyHost(b bool) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalTrustAnyHost = b
}

// IsTrustAnyHost 是否信任所有主机。
func IsTrustAnyHost() bool {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalTrustAnyHost
}

// SetGlobalBoundary 设置 multipart 默认 boundary。
func SetGlobalBoundary(b string) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalBoundary = b
}

// GetGlobalBoundary 获取 multipart 默认 boundary。
func GetGlobalBoundary() string {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalBoundary
}

// SetGlobalDecodeURL 设置是否对 URL 自动解码。
func SetGlobalDecodeURL(b bool) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalDecodeURL = b
}

// IsGlobalDecodeURL 是否对 URL 自动解码。
func IsGlobalDecodeURL() bool {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalDecodeURL
}

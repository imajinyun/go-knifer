package http

import (
	"net/http"
	"sync"
	"time"
)

// Global default configuration, aligned with the utility toolkit-http HttpGlobalConfig.
var (
	globalMu               sync.RWMutex
	globalTimeout          = 0 * time.Second // 0 means using the HTTP client's default timeout.
	globalMaxRedirects     = 10
	globalIgnoreEOFError   = true
	globalDecodeURL        = false
	globalFollowRedirects  = true
	globalDefaultUserAgent = ""
	globalTrustAnyHost     = false
	globalBoundary         = "--------------------gokitFormBoundary"
)

// GlobalConfig is an immutable snapshot of package-level HTTP defaults.
type GlobalConfig struct {
	Timeout          time.Duration
	MaxRedirects     int
	IgnoreEOFError   bool
	DecodeURL        bool
	FollowRedirects  bool
	DefaultUserAgent string
	TrustAnyHost     bool
	Boundary         string
	Headers          http.Header
	CookieJar        http.CookieJar
}

// SnapshotGlobalConfig returns a consistent copy of the current package-level HTTP defaults.
func SnapshotGlobalConfig() GlobalConfig {
	globalMu.RLock()
	cfg := GlobalConfig{
		Timeout:          globalTimeout,
		MaxRedirects:     globalMaxRedirects,
		IgnoreEOFError:   globalIgnoreEOFError,
		DecodeURL:        globalDecodeURL,
		FollowRedirects:  globalFollowRedirects,
		DefaultUserAgent: globalDefaultUserAgent,
		TrustAnyHost:     globalTrustAnyHost,
		Boundary:         globalBoundary,
	}
	globalMu.RUnlock()
	cfg.Headers = CloneGlobalHeaders()
	cfg.CookieJar = GetCookieJar()
	return cfg
}

// SetGlobalTimeout sets the global default timeout.
func SetGlobalTimeout(d time.Duration) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalTimeout = d
}

// GetGlobalTimeout returns the global default timeout.
func GetGlobalTimeout() time.Duration {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalTimeout
}

// SetGlobalMaxRedirects sets the global maximum redirect count.
func SetGlobalMaxRedirects(n int) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalMaxRedirects = n
}

// GetGlobalMaxRedirects returns the global maximum redirect count.
func GetGlobalMaxRedirects() int {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalMaxRedirects
}

// SetGlobalFollowRedirects sets whether redirects are followed.
func SetGlobalFollowRedirects(b bool) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalFollowRedirects = b
}

// GetGlobalFollowRedirects reports whether redirects are followed.
func GetGlobalFollowRedirects() bool {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalFollowRedirects
}

// SetGlobalUserAgent sets the global default User-Agent.
func SetGlobalUserAgent(ua string) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalDefaultUserAgent = ua
}

// GetGlobalUserAgent returns the global default User-Agent.
func GetGlobalUserAgent() string {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalDefaultUserAgent
}

// SetIgnoreEOFError sets whether EOF errors are ignored.
func SetIgnoreEOFError(b bool) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalIgnoreEOFError = b
}

// IsIgnoreEOFError reports whether EOF errors are ignored.
func IsIgnoreEOFError() bool {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalIgnoreEOFError
}

// SetTrustAnyHost sets whether all hosts are trusted, skipping HTTPS certificate verification.
func SetTrustAnyHost(b bool) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalTrustAnyHost = b
}

// IsTrustAnyHost reports whether all hosts are trusted.
func IsTrustAnyHost() bool {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalTrustAnyHost
}

// SetGlobalBoundary sets the default multipart boundary.
func SetGlobalBoundary(b string) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalBoundary = b
}

// GetGlobalBoundary returns the default multipart boundary.
func GetGlobalBoundary() string {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalBoundary
}

// SetGlobalDecodeURL sets whether URLs are decoded automatically.
func SetGlobalDecodeURL(b bool) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalDecodeURL = b
}

// IsGlobalDecodeURL reports whether URLs are decoded automatically.
func IsGlobalDecodeURL() bool {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalDecodeURL
}

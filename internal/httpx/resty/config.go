package resty

import (
	"sync"
	"time"
)

// HeaderValues stores HTTP header values without depending on net/http types.
type HeaderValues map[string][]string

var (
	globalMu               sync.RWMutex
	globalTimeout          time.Duration
	globalMaxRedirects     = 10
	globalFollowRedirects  = true
	globalDefaultUserAgent = ""
	globalTrustAnyHost     = false

	globalHeadersMu sync.RWMutex
	globalHeaders   = HeaderValues{}

	cookieMu       sync.RWMutex
	cookieDisabled bool
)

func init() {
	setHeader(globalHeaders, string(HeaderAccept), "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	setHeader(globalHeaders, string(HeaderAcceptEncoding), "gzip, deflate")
	setHeader(globalHeaders, string(HeaderAcceptLanguage), "zh-CN,zh;q=0.8")
	setHeader(globalHeaders, string(HeaderUserAgent),
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/72.0.3626.109 Safari/537.36")
}

// SetGlobalTimeout sets the global default timeout.
func SetGlobalTimeout(d time.Duration) { globalMu.Lock(); defer globalMu.Unlock(); globalTimeout = d }

// GetGlobalTimeout returns the global default timeout.
func GetGlobalTimeout() time.Duration {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalTimeout
}

// SetGlobalMaxRedirects sets the global maximum redirect count.
func SetGlobalMaxRedirects(n int) { globalMu.Lock(); defer globalMu.Unlock(); globalMaxRedirects = n }

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

// SetTrustAnyHost sets whether all hosts are trusted, skipping HTTPS certificate verification.
func SetTrustAnyHost(b bool) { globalMu.Lock(); defer globalMu.Unlock(); globalTrustAnyHost = b }

// IsTrustAnyHost reports whether all hosts are trusted.
func IsTrustAnyHost() bool { globalMu.RLock(); defer globalMu.RUnlock(); return globalTrustAnyHost }

// SetGlobalHeader sets a global default request header.
func SetGlobalHeader(name, value string) {
	globalHeadersMu.Lock()
	defer globalHeadersMu.Unlock()
	setHeader(globalHeaders, name, value)
}

// AddGlobalHeader appends a global default request header value.
func AddGlobalHeader(name, value string) {
	globalHeadersMu.Lock()
	defer globalHeadersMu.Unlock()
	globalHeaders[name] = append(globalHeaders[name], value)
}

// RemoveGlobalHeader removes a global default request header.
func RemoveGlobalHeader(name string) {
	globalHeadersMu.Lock()
	defer globalHeadersMu.Unlock()
	delete(globalHeaders, name)
}

// CloneGlobalHeaders returns a copy of global default request headers.
func CloneGlobalHeaders() HeaderValues {
	globalHeadersMu.RLock()
	defer globalHeadersMu.RUnlock()
	out := HeaderValues{}
	for k, v := range globalHeaders {
		out[k] = append([]string(nil), v...)
	}
	return out
}

// CloseCookie disables global cookie management.
func CloseCookie() {
	cookieMu.Lock()
	defer cookieMu.Unlock()
	cookieDisabled = true
}

func isCookieDisabled() bool {
	cookieMu.RLock()
	defer cookieMu.RUnlock()
	return cookieDisabled
}

func setHeader(headers HeaderValues, name, value string) {
	headers[name] = []string{value}
}

func getHeader(headers HeaderValues, name string) string {
	if values := headers[name]; len(values) > 0 {
		return values[0]
	}
	return ""
}

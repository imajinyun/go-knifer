package http

import (
	"net/http"
	"sync"
)

// GlobalHeaders maintains global default request headers, aligned with the utility toolkit-http GlobalHeaders.
var (
	globalHeadersMu sync.RWMutex
	globalHeaders   = defaultGlobalHeaders()
)

func defaultGlobalHeaders() http.Header {
	headers := http.Header{}
	// Align with the utility toolkit defaults, excluding unsupported encodings.
	headers.Set(string(HeaderAccept), "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	headers.Set(string(HeaderAcceptEncoding), "gzip, deflate")
	headers.Set(string(HeaderAcceptLanguage), "zh-CN,zh;q=0.8")
	headers.Set(string(HeaderUserAgent),
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/72.0.3626.109 Safari/537.36")
	return headers
}

// SetGlobalHeader sets a global default request header.
func SetGlobalHeader(name, value string) {
	globalHeadersMu.Lock()
	defer globalHeadersMu.Unlock()
	globalHeaders.Set(name, value)
}

// AddGlobalHeader appends a global default request header value.
func AddGlobalHeader(name, value string) {
	globalHeadersMu.Lock()
	defer globalHeadersMu.Unlock()
	globalHeaders.Add(name, value)
}

// RemoveGlobalHeader removes a global default request header.
func RemoveGlobalHeader(name string) {
	globalHeadersMu.Lock()
	defer globalHeadersMu.Unlock()
	globalHeaders.Del(name)
}

// CloneGlobalHeaders returns a copy of global default request headers.
func CloneGlobalHeaders() http.Header {
	globalHeadersMu.RLock()
	defer globalHeadersMu.RUnlock()
	return cloneHeader(globalHeaders)
}

func cloneHeader(headers http.Header) http.Header {
	out := http.Header{}
	for k, v := range headers {
		out[k] = append([]string(nil), v...)
	}
	return out
}

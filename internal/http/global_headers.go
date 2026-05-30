package http

import (
	"net/http"
	"sync"
)

// GlobalHeaders maintains global default request headers, aligned with hutool-http GlobalHeaders.
var (
	globalHeadersMu sync.RWMutex
	globalHeaders   = http.Header{}
)

func init() {
	// Align with hutool defaults, excluding unsupported encodings.
	globalHeaders.Set(string(HeaderAccept), "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	globalHeaders.Set(string(HeaderAcceptEncoding), "gzip, deflate")
	globalHeaders.Set(string(HeaderAcceptLanguage), "zh-CN,zh;q=0.8")
	globalHeaders.Set(string(HeaderUserAgent),
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/72.0.3626.109 Safari/537.36")
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
	out := http.Header{}
	for k, v := range globalHeaders {
		out[k] = append([]string(nil), v...)
	}
	return out
}

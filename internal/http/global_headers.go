package http

import (
	"net/http"
	"sync"
)

// GlobalHeaders 维护全局默认请求头（对应 hutool-http GlobalHeaders）。
var (
	globalHeadersMu sync.RWMutex
	globalHeaders   = http.Header{}
)

func init() {
	// 与 hutool 默认值对齐
	globalHeaders.Set(string(HeaderAccept), "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	globalHeaders.Set(string(HeaderAcceptEncoding), "gzip, deflate, br")
	globalHeaders.Set(string(HeaderAcceptLanguage), "zh-CN,zh;q=0.8")
	globalHeaders.Set(string(HeaderUserAgent),
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/72.0.3626.109 Safari/537.36")
}

// SetGlobalHeader 设置一个全局默认请求头。
func SetGlobalHeader(name, value string) {
	globalHeadersMu.Lock()
	defer globalHeadersMu.Unlock()
	globalHeaders.Set(name, value)
}

// AddGlobalHeader 追加一个全局默认请求头。
func AddGlobalHeader(name, value string) {
	globalHeadersMu.Lock()
	defer globalHeadersMu.Unlock()
	globalHeaders.Add(name, value)
}

// RemoveGlobalHeader 移除一个全局默认请求头。
func RemoveGlobalHeader(name string) {
	globalHeadersMu.Lock()
	defer globalHeadersMu.Unlock()
	globalHeaders.Del(name)
}

// CloneGlobalHeaders 拷贝一份全局默认请求头。
func CloneGlobalHeaders() http.Header {
	globalHeadersMu.RLock()
	defer globalHeadersMu.RUnlock()
	out := http.Header{}
	for k, v := range globalHeaders {
		out[k] = append([]string(nil), v...)
	}
	return out
}

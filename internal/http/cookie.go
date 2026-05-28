package http

import (
	"net/http"
	"net/http/cookiejar"
	"sync"
)

// GlobalCookie 提供全局 Cookie 管理（对应 hutool-http GlobalCookieManager）。
var (
	cookieMu  sync.RWMutex
	cookieJar http.CookieJar
)

func init() {
	jar, _ := cookiejar.New(nil)
	cookieJar = jar
}

// SetCookieJar 自定义全局 CookieJar，传入 nil 表示禁用 Cookie 管理。
func SetCookieJar(jar http.CookieJar) {
	cookieMu.Lock()
	defer cookieMu.Unlock()
	cookieJar = jar
}

// GetCookieJar 获取当前全局 CookieJar。
func GetCookieJar() http.CookieJar {
	cookieMu.RLock()
	defer cookieMu.RUnlock()
	return cookieJar
}

// CloseCookie 关闭全局 Cookie 管理。
func CloseCookie() { SetCookieJar(nil) }

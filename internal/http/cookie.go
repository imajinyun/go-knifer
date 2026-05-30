package http

import (
	"net/http"
	"net/http/cookiejar"
	"sync"
)

// GlobalCookie provides global cookie management, aligned with hutool-http GlobalCookieManager.
var (
	cookieMu  sync.RWMutex
	cookieJar http.CookieJar
)

func init() {
	jar, _ := cookiejar.New(nil)
	cookieJar = jar
}

// SetCookieJar customizes the global CookieJar; nil disables cookie management.
func SetCookieJar(jar http.CookieJar) {
	cookieMu.Lock()
	defer cookieMu.Unlock()
	cookieJar = jar
}

// GetCookieJar returns the current global CookieJar.
func GetCookieJar() http.CookieJar {
	cookieMu.RLock()
	defer cookieMu.RUnlock()
	return cookieJar
}

// CloseCookie disables global cookie management.
func CloseCookie() { SetCookieJar(nil) }

package http

import (
	"net/http"
	"net/http/cookiejar"
)

func newDefaultCookieJar() http.CookieJar {
	jar, _ := cookiejar.New(nil)
	return jar
}

// SetCookieJar customizes the global CookieJar; nil disables cookie management.
func SetCookieJar(jar http.CookieJar) {
	globalMu.Lock()
	defer globalMu.Unlock()
	cookieJar = jar
}

// GetCookieJar returns the current global CookieJar.
func GetCookieJar() http.CookieJar {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return cookieJar
}

// CloseCookie disables global cookie management.
func CloseCookie() { SetCookieJar(nil) }

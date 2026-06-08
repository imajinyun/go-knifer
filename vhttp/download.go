package vhttp

import httpx "github.com/imajinyun/go-knifer/internal/httpx/http"

// DownloadString delegates to the internal httpx implementation.
//
// Deprecated: use DownloadStringE to handle request and read errors explicitly.
func DownloadString(rawURL, customCharset string) string {
	return DownloadStringWithOptions(rawURL, customCharset)
}

// DownloadStringWithOptions downloads remote text with per-request options.
//
// Deprecated: use DownloadStringEWithOptions to handle request and read errors explicitly.
func DownloadStringWithOptions(rawURL, customCharset string, opts ...RequestOption) string {
	return httpx.DownloadStringWithOptions(rawURL, customCharset, opts...)
}

// DownloadStringE downloads remote text and returns an error on request or read failure.
func DownloadStringE(rawURL, customCharset string) (string, error) {
	return DownloadStringEWithOptions(rawURL, customCharset)
}

// DownloadStringEWithOptions downloads remote text with per-request options and returns an error on failure.
func DownloadStringEWithOptions(rawURL, customCharset string, opts ...RequestOption) (string, error) {
	return httpx.DownloadStringEWithOptions(rawURL, customCharset, opts...)
}

// DownloadStringSafeE downloads remote text with SSRF-oriented safety checks enabled.
func DownloadStringSafeE(rawURL, customCharset string, opts ...RequestOption) (string, error) {
	return httpx.DownloadStringSafeE(rawURL, customCharset, opts...)
}

// DownloadBytes delegates to the internal httpx implementation.
//
// Deprecated: use DownloadBytesE to handle request and read errors explicitly.
func DownloadBytes(rawURL string) []byte {
	return DownloadBytesWithOptions(rawURL)
}

// DownloadBytesWithOptions downloads and returns bytes with per-request options.
//
// Deprecated: use DownloadBytesEWithOptions to handle request and read errors explicitly.
func DownloadBytesWithOptions(rawURL string, opts ...RequestOption) []byte {
	return httpx.DownloadBytesWithOptions(rawURL, opts...)
}

// DownloadBytesE downloads and returns bytes or an error.
func DownloadBytesE(rawURL string) ([]byte, error) { return DownloadBytesEWithOptions(rawURL) }

// DownloadBytesEWithOptions downloads and returns bytes with per-request options or an error.
func DownloadBytesEWithOptions(rawURL string, opts ...RequestOption) ([]byte, error) {
	return httpx.DownloadBytesEWithOptions(rawURL, opts...)
}

// DownloadBytesSafeE downloads and returns bytes with SSRF-oriented safety checks enabled.
func DownloadBytesSafeE(rawURL string, opts ...RequestOption) ([]byte, error) {
	return httpx.DownloadBytesSafeE(rawURL, opts...)
}

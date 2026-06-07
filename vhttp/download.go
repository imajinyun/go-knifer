package vhttp

import httpx "github.com/imajinyun/go-knifer/internal/httpx/http"

// DownloadString delegates to the internal httpx implementation.
func DownloadString(rawURL, customCharset string) string {
	return DownloadStringWithOptions(rawURL, customCharset)
}

// DownloadStringWithOptions downloads remote text with per-request options.
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

// DownloadBytes delegates to the internal httpx implementation.
func DownloadBytes(rawURL string) []byte {
	return DownloadBytesWithOptions(rawURL)
}

// DownloadBytesWithOptions downloads and returns bytes with per-request options.
func DownloadBytesWithOptions(rawURL string, opts ...RequestOption) []byte {
	return httpx.DownloadBytesWithOptions(rawURL, opts...)
}

// DownloadBytesE downloads and returns bytes or an error.
func DownloadBytesE(rawURL string) ([]byte, error) { return DownloadBytesEWithOptions(rawURL) }

// DownloadBytesEWithOptions downloads and returns bytes with per-request options or an error.
func DownloadBytesEWithOptions(rawURL string, opts ...RequestOption) ([]byte, error) {
	return httpx.DownloadBytesEWithOptions(rawURL, opts...)
}

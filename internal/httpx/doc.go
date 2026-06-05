// Package httpx provides HTTP client and server implementations.
//
// It is organized into subpackages:
//   - http: standard-library-based HTTP request builder with chainable API
//   - resty: Resty v3 wrapper for advanced HTTP operations
//   - internal/shared: shared types and error definitions for both subpackages
//
// This package is exposed through the vhttp and vresty facades.
package httpx

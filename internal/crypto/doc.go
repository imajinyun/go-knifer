// Package crypto provides security-oriented digest, HMAC, AES, RSA, and PEM
// encoding helpers.
//
// Keep general-purpose hash helpers such as additive/FNV and lightweight digest
// shortcuts in internal/hash. This package owns APIs that are used in security,
// signing, encryption, and key material handling contexts.
package crypto

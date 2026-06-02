// Package crypto provides security-oriented digest, HMAC, AES, RSA, and PEM
// encoding helpers.
//
// Keep general-purpose, non-cryptographic hash helpers such as additive/FNV in
// internal/hash. This package owns all digest APIs (MD5/SHA family) plus the
// signing, encryption, and key material handling used in security contexts.
package crypto

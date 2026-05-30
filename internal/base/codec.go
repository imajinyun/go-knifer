package base

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
)

// This file provides encoding and decoding helpers aligned with hutool-core
// Base64, HexUtil, and URLUtil.

// Base64Encode encodes bytes with standard Base64 encoding.
func Base64Encode(data []byte) string { return base64.StdEncoding.EncodeToString(data) }

// Base64EncodeStr encodes a string with standard Base64 encoding.
func Base64EncodeStr(s string) string { return Base64Encode([]byte(s)) }

// Base64Decode decodes a standard Base64 string.
func Base64Decode(s string) ([]byte, error) { return base64.StdEncoding.DecodeString(s) }

// Base64DecodeStr decodes a standard Base64 string and returns text.
func Base64DecodeStr(s string) (string, error) {
	b, err := Base64Decode(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Base64URLEncode encodes bytes with URL-safe Base64 encoding.
func Base64URLEncode(data []byte) string { return base64.URLEncoding.EncodeToString(data) }

// Base64URLDecode decodes a URL-safe Base64 string.
func Base64URLDecode(s string) ([]byte, error) { return base64.URLEncoding.DecodeString(s) }

// HexEncode encodes bytes as a lowercase hexadecimal string.
func HexEncode(data []byte) string { return hex.EncodeToString(data) }

// HexEncodeStr encodes a string as lowercase hexadecimal text.
func HexEncodeStr(s string) string { return HexEncode([]byte(s)) }

// HexDecode decodes a hexadecimal string.
func HexDecode(s string) ([]byte, error) { return hex.DecodeString(s) }

// HexDecodeStr decodes a hexadecimal string and returns text.
func HexDecodeStr(s string) (string, error) {
	b, err := HexDecode(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// URLEncode escapes a string for use in URL query components.
func URLEncode(s string) string { return url.QueryEscape(s) }

// URLDecode unescapes a URL query component string.
func URLDecode(s string) (string, error) { return url.QueryUnescape(s) }

package base

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
)

// 对应 hutool-core Base64 / HexUtil / URLUtil 编解码部分。

// Base64Encode 标准 Base64 编码。
func Base64Encode(data []byte) string { return base64.StdEncoding.EncodeToString(data) }

// Base64EncodeStr 字符串的 Base64 编码。
func Base64EncodeStr(s string) string { return Base64Encode([]byte(s)) }

// Base64Decode 标准 Base64 解码。
func Base64Decode(s string) ([]byte, error) { return base64.StdEncoding.DecodeString(s) }

// Base64DecodeStr 解码并返回字符串。
func Base64DecodeStr(s string) (string, error) {
	b, err := Base64Decode(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Base64URLEncode URL 安全 Base64 编码。
func Base64URLEncode(data []byte) string { return base64.URLEncoding.EncodeToString(data) }

// Base64URLDecode URL 安全 Base64 解码。
func Base64URLDecode(s string) ([]byte, error) { return base64.URLEncoding.DecodeString(s) }

// HexEncode 十六进制编码。
func HexEncode(data []byte) string { return hex.EncodeToString(data) }

// HexEncodeStr 字符串的 Hex 编码。
func HexEncodeStr(s string) string { return HexEncode([]byte(s)) }

// HexDecode 十六进制解码。
func HexDecode(s string) ([]byte, error) { return hex.DecodeString(s) }

// HexDecodeStr 解码并返回字符串。
func HexDecodeStr(s string) (string, error) {
	b, err := HexDecode(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// URLEncode URL 编码。
func URLEncode(s string) string { return url.QueryEscape(s) }

// URLDecode URL 解码。
func URLDecode(s string) (string, error) { return url.QueryUnescape(s) }

package validator

import "testing"

func TestValidators(t *testing.T) {
	if !IsEmail("a@b.com") || IsEmail("abc") {
		t.Fatalf("IsEmail failed")
	}
	if !IsMobile("13812345678") || IsMobile("12812345678") {
		t.Fatalf("IsMobile failed")
	}
	if !IsURL("https://hutool.cn") || !IsURL("ftp://x") || IsURL("/relative/path") || IsURL(" https://hutool.cn") {
		t.Fatalf("IsURL failed")
	}
	if !IsIPv4("127.0.0.1") || IsIPv4("256.0.0.1") {
		t.Fatalf("IsIPv4 failed")
	}
	if !IsChinese("你好") || IsChinese("hello") {
		t.Fatalf("IsChinese failed")
	}
	if !IsNumberStr("-3.14") || IsNumberStr("ab") {
		t.Fatalf("IsNumberStr failed")
	}
}

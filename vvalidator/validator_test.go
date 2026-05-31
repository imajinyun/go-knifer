package vvalidator

import "testing"

func TestValidatorFacade(t *testing.T) {
	if !IsEmail("a@b.com") || IsEmail("bad") {
		t.Fatal("IsEmail failed")
	}
	if !IsMobile("13812345678") || IsMobile("12812345678") {
		t.Fatal("IsMobile failed")
	}
	if !IsURL("https://example.com") || !IsURL("ftp://example.com") || IsURL("/relative/path") {
		t.Fatal("IsURL failed")
	}
	if !IsIPv4("127.0.0.1") || IsIPv4("256.0.0.1") {
		t.Fatal("IsIPv4 failed")
	}
	if !IsChinese("你好") || IsChinese("hello") {
		t.Fatal("IsChinese failed")
	}
	if !IsNumberStr("-3.14") || IsNumberStr("x") {
		t.Fatal("IsNumberStr failed")
	}
}

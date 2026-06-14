package http

import (
	"strings"
	"testing"
)

func TestDecodeParams(t *testing.T) {
	paramsStr := "uuuu=0&a=b&c=%3F%23%40!%24%25%5E%26%3Ddsssss555555"
	m := DecodeParams(paramsStr)
	if m["uuuu"][0] != "0" {
		t.Fatalf("uuuu: %v", m["uuuu"])
	}
	if m["a"][0] != "b" {
		t.Fatalf("a: %v", m["a"])
	}
	if m["c"][0] != "?#@!$%^&=dsssss555555" {
		t.Fatalf("c: %v", m["c"])
	}
}

func TestDecodeParamMap(t *testing.T) {
	m := DecodeParamMap("aa=123&f_token=NzBkMjQxNDM1MDVlMDliZTk1OTU3ZDI1OTI0NTBiOWQ=")
	if m["aa"] != "123" {
		t.Fatalf("aa: %q", m["aa"])
	}
	if m["f_token"] != "NzBkMjQxNDM1MDVlMDliZTk1OTU3ZDI1OTI0NTBiOWQ=" {
		t.Fatalf("f_token: %q", m["f_token"])
	}
}

func TestEncodeParams(t *testing.T) {
	got := EncodeParams("http://www.abc.dd?a=b&c=d")
	if !strings.Contains(got, "a=b") || !strings.Contains(got, "c=d") {
		t.Fatalf("encoded: %q", got)
	}
	if EncodeParams("https://www.example.com/") != "https://www.example.com/" {
		t.Fatal("URL without query should be returned unchanged")
	}
}

func TestToParams(t *testing.T) {
	m := map[string]any{"a": "1"}
	got := ToParams(m)
	if got != "a=1" {
		t.Fatalf("ToParams: %q", got)
	}
}

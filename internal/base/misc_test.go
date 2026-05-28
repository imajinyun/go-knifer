package base

import "testing"

// 对应 hutool-core CharUtilTest / BooleanUtilTest / HashUtilTest / ReUtilTest / ValidatorTest / EscapeUtilTest。

func TestCharUtil(t *testing.T) {
	if !IsBlankChar(' ') || !IsBlankChar('\u00A0') || IsBlankChar('a') {
		t.Fatalf("IsBlankChar failed")
	}
	if !IsLetter('A') || IsLetter('1') {
		t.Fatalf("IsLetter failed")
	}
	if !IsDigit('1') || IsDigit('a') {
		t.Fatalf("IsDigit failed")
	}
	if !IsAscii('A') || IsAscii('中') {
		t.Fatalf("IsAscii failed")
	}
	if !IsLetterOrDigit('a') || !IsLetterOrDigit('1') || IsLetterOrDigit('?') {
		t.Fatalf("IsLetterOrDigit failed")
	}
}

func TestBooleanUtil(t *testing.T) {
	if !BoolNegate(false) || BoolNegate(true) {
		t.Fatalf("Negate failed")
	}
	if BoolToInt(true) != 1 || BoolToInt(false) != 0 {
		t.Fatalf("BoolToInt failed")
	}
	if !BoolAnd(true, true) || BoolAnd(true, false) {
		t.Fatalf("BoolAnd failed")
	}
	if !BoolOr(false, true) || BoolOr(false, false) {
		t.Fatalf("BoolOr failed")
	}
}

func TestHashFunctions(t *testing.T) {
	if MD5Hex("abc") != "900150983cd24fb0d6963f7d28e17f72" {
		t.Fatalf("MD5 failed")
	}
	if SHA1Hex("abc") != "a9993e364706816aba3e25717850c26c9cd0d89d" {
		t.Fatalf("SHA1 failed")
	}
	if SHA256Hex("abc") != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Fatalf("SHA256 failed")
	}
	if FnvHash("abc") == 0 {
		t.Fatalf("FnvHash zero")
	}
	if AdditiveHash("abc", 31) < 0 {
		t.Fatalf("AdditiveHash failed")
	}
}

func TestValidators(t *testing.T) {
	if !IsEmail("a@b.com") || IsEmail("abc") {
		t.Fatalf("IsEmail failed")
	}
	if !IsMobile("13812345678") || IsMobile("12812345678") {
		t.Fatalf("IsMobile failed")
	}
	if !IsURL("https://hutool.cn") || IsURL("ftp://x") {
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

func TestRegex(t *testing.T) {
	if !ReMatch(`^\d+$`, "123") || ReMatch(`^\d+$`, "12a") {
		t.Fatalf("ReMatch failed")
	}
	if ReFind(`\d+`, "ab123cd") != "123" {
		t.Fatalf("ReFind failed")
	}
	all := ReFindAll(`\d+`, "a1b22c333")
	if len(all) != 3 || all[2] != "333" {
		t.Fatalf("ReFindAll failed: %v", all)
	}
	if ReReplace(`\d`, "a1b2", "*") != "a*b*" {
		t.Fatalf("ReReplace failed")
	}
}

func TestObjUtilDefault(t *testing.T) {
	if DefaultIfEmpty("", "x") != "x" || DefaultIfEmpty("a", "x") != "a" {
		t.Fatalf("DefaultIfEmpty failed")
	}
	if DefaultIfBlank("  ", "x") != "x" || DefaultIfBlank("a", "x") != "a" {
		t.Fatalf("DefaultIfBlank failed")
	}
	v := 5
	if DefaultIfNil(&v, 10) != 5 {
		t.Fatalf("DefaultIfNil non-nil failed")
	}
	var p *int
	if DefaultIfNil(p, 10) != 10 {
		t.Fatalf("DefaultIfNil nil failed")
	}
}

func TestEscapeHTML(t *testing.T) {
	got := EscapeHTML(`<a href="x">&'</a>`)
	want := "&lt;a href=&quot;x&quot;&gt;&amp;&#39;&lt;/a&gt;"
	if got != want {
		t.Fatalf("EscapeHTML: %q", got)
	}
	if UnescapeHTML(want) != `<a href="x">&'</a>` {
		t.Fatalf("UnescapeHTML failed")
	}
	if UnescapeHTML("&nbsp;") != "\u00A0" {
		t.Fatalf("UnescapeHTML nbsp failed")
	}
}

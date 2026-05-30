package vchar

import "testing"

func TestCharFacade(t *testing.T) {
	if !IsBlankChar(' ') || !IsBlankChar('\u00A0') || IsBlankChar('a') {
		t.Fatal("IsBlankChar failed")
	}
	if !IsLetter('A') || IsLetter('1') {
		t.Fatal("IsLetter failed")
	}
	if !IsDigit('1') || IsDigit('a') {
		t.Fatal("IsDigit failed")
	}
	if !IsAscii('A') || IsAscii('中') {
		t.Fatal("IsAscii failed")
	}
	if !IsLetterOrDigit('a') || !IsLetterOrDigit('1') || IsLetterOrDigit('?') {
		t.Fatal("IsLetterOrDigit failed")
	}
}

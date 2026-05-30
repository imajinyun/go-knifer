package char

import "testing"

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

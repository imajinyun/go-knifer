package num

import "testing"

func TestIsNumber(t *testing.T) {
	if !IsNumber("123") || !IsNumber("-3.14") || IsNumber("abc") {
		t.Fatalf("IsNumber failed")
	}
	if !IsInteger("-12") || IsInteger("12.3") {
		t.Fatalf("IsInteger failed")
	}
	if !IsDigits("12345") || IsDigits("-12") {
		t.Fatalf("IsDigits failed")
	}
}

func TestNumberPredicatesComprehensive(t *testing.T) {
	validNumbers := []string{"  +12.30 ", "-0x1f", "6.02e23", "1F", "2d", "3L"}
	for _, s := range validNumbers {
		if !IsNumber(s) {
			t.Fatalf("IsNumber(%q) should be true", s)
		}
	}
	invalidNumbers := []string{"", "  ", "0x", "0xz", "1e-", "abc", "1.2.3"}
	for _, s := range invalidNumbers {
		if IsNumber(s) {
			t.Fatalf("IsNumber(%q) should be false", s)
		}
	}
	if IsInteger("") || IsInteger("12.0") || !IsInteger("-12") {
		t.Fatal("IsInteger edge cases failed")
	}
	if IsLong("9223372036854775808") || !IsLong("-9223372036854775808") {
		t.Fatal("IsLong boundary cases failed")
	}
	if IsDouble("1") || IsDouble("") || !IsDouble("1.0") || !IsDouble("-0.25") {
		t.Fatal("IsDouble edge cases failed")
	}
	if IsDigits("") || IsDigits("12a") || !IsDigits("00123") {
		t.Fatal("IsDigits edge cases failed")
	}
	primeCases := map[int]bool{-1: false, 0: false, 1: false, 2: true, 3: true, 4: false, 25: false, 7919: true}
	for n, want := range primeCases {
		if got := IsPrimes(n); got != want {
			t.Fatalf("IsPrimes(%d) = %v, want %v", n, got, want)
		}
	}
}

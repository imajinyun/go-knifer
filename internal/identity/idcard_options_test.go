package identity

import "testing"

func TestIDCardMatcherOptions(t *testing.T) {
	if IsValidIDCard15WithOptions("130503670401001", WithDigitsMatcher(func(string) bool { return false })) {
		t.Fatal("custom digits matcher should reject 15-digit ID card")
	}
	if IsValidIDCard18WithOptions("11010519491231002X", WithDigitsMatcher(func(string) bool { return false })) {
		t.Fatal("custom digits matcher should reject 18-digit ID card")
	}
	if _, ok := ParseRegionCardWithOptions("1571234(5)", WithMacauCardMatcher(func(string) bool { return false })); ok {
		t.Fatal("custom Macau matcher should reject region card")
	}
	if IsValidTWIDCardWithOptions("A123456789", WithTWCardMatcher(func(string) bool { return false })) {
		t.Fatal("custom Taiwan matcher should reject card")
	}
	if IsValidHKIDCardWithOptions("A123456(3)", WithHKCardMatcher(func(string) bool { return false })) {
		t.Fatal("custom Hong Kong matcher should reject card")
	}
	if CheckCode18WithOptions("11010519491231002", WithDigitsMatcher(func(string) bool { return false })) != ' ' {
		t.Fatal("custom digits matcher should reject check code input")
	}
}

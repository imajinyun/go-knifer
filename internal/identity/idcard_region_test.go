package identity

import "testing"

func TestRegionCards(t *testing.T) {
	if !IsValidTWIDCard("A123456789") {
		t.Fatal("expected valid Taiwan card")
	}
	if !IsValidHKIDCard("A123456(3)") {
		t.Fatal("expected valid Hong Kong card")
	}
	info, ok := ParseRegionCard("A123456789")
	if !ok || info.Region != "台湾" || info.Gender != "M" || !info.Valid {
		t.Fatalf("ParseRegionCard Taiwan = %+v, %v", info, ok)
	}
	info, ok = ParseRegionCard("1571234(5)")
	if !ok || info.Region != "澳门" || !info.Valid {
		t.Fatalf("ParseRegionCard Macau = %+v, %v", info, ok)
	}
}

func TestRegionCardsRejectMalformedValues(t *testing.T) {
	if info, ok := ParseRegionCard(""); ok || info != (RegionCardInfo{}) {
		t.Fatalf("ParseRegionCard(empty) = %+v, %v; want zero false", info, ok)
	}
	if info, ok := ParseRegionCard("ABC"); ok || info != (RegionCardInfo{}) {
		t.Fatalf("ParseRegionCard(short) = %+v, %v; want zero false", info, ok)
	}
	if info, ok := ParseRegionCardWithOptions("A323456789", WithTWCardMatcher(func(string) bool { return true })); !ok || info.Valid || info.Gender != "N" {
		t.Fatalf("ParseRegionCard(TW unknown gender) = %+v, %v; want invalid N", info, ok)
	}
	if IsValidTWIDCard("Z12345678A") {
		t.Fatal("Taiwan card should reject non-digit check code")
	}
	if IsValidTWIDCard("A12345678A") {
		t.Fatal("Taiwan card should reject non-digit body/check code")
	}
	if IsValidHKIDCard("A12345A(3)") {
		t.Fatal("Hong Kong card should reject non-digit body")
	}
	if IsValidHKIDCardWithOptions("A123456(Z)", WithHKCardMatcher(func(string) bool { return true })) {
		t.Fatal("Hong Kong card should reject unsupported check code")
	}
}

func TestRegionCardsHandleFullWidthParentheses(t *testing.T) {
	info, ok := ParseRegionCard("1571234（5）")
	if !ok || info.Region != "澳门" || !info.Valid {
		t.Fatalf("ParseRegionCard(full-width Macau) = %+v, %v", info, ok)
	}
}

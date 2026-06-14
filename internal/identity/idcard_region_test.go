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

package identity

import "testing"

func TestParseIDCardAndHide(t *testing.T) {
	info, ok := ParseIDCard("11010519491231002X")
	if !ok || info.ProvinceCode != "11" || info.CityCode != "1101" || info.DistrictCode != "110105" {
		t.Fatalf("ParseIDCard() = %+v, %v", info, ok)
	}
	if got := Hide("11010519491231002X", 6, 14); got != "110105********002X" {
		t.Fatalf("Hide() = %q", got)
	}
}

func TestParseIDCardRejectsInvalidID(t *testing.T) {
	if info, ok := ParseIDCard("11010519490231002X"); ok || info != (IDCardInfo{}) {
		t.Fatalf("ParseIDCard(invalid) = %+v, %v; want zero false", info, ok)
	}
}

func TestHideClampsRuneIndexes(t *testing.T) {
	if got := Hide("身份证ABC", -2, 3); got != "***ABC" {
		t.Fatalf("Hide(clamp start) = %q", got)
	}
	if got := Hide("身份证ABC", 3, 99); got != "身份证***" {
		t.Fatalf("Hide(clamp end) = %q", got)
	}
	if got := Hide("身份证ABC", 4, 2); got != "身份证ABC" {
		t.Fatalf("Hide(no-op) = %q", got)
	}
}

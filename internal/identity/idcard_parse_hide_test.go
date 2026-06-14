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

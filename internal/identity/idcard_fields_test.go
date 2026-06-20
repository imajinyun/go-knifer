package identity

import (
	"errors"
	"testing"
	"time"
)

func TestIDCardFields(t *testing.T) {
	const id = "11010519491231002X"
	birth, ok := BirthString(id)
	if !ok || birth != "19491231" {
		t.Fatalf("BirthString() = %q, %v", birth, ok)
	}
	year, ok := Year(id)
	if !ok || year != 1949 {
		t.Fatalf("Year() = %d, %v", year, ok)
	}
	month, ok := Month(id)
	if !ok || month != 12 {
		t.Fatalf("Month() = %d, %v", month, ok)
	}
	day, ok := Day(id)
	if !ok || day != 31 {
		t.Fatalf("Day() = %d, %v", day, ok)
	}
	age, ok := AgeAt(id, time.Date(2024, 12, 30, 0, 0, 0, 0, time.Local))
	if !ok || age != 74 {
		t.Fatalf("AgeAt(before birthday) = %d, %v", age, ok)
	}
	age, ok = AgeAt(id, time.Date(2024, 12, 31, 0, 0, 0, 0, time.Local))
	if !ok || age != 75 {
		t.Fatalf("AgeAt(on birthday) = %d, %v", age, ok)
	}
	age, ok = AgeWithOptions(id, WithAgeTime(time.Date(2024, 12, 31, 0, 0, 0, 0, time.Local)))
	if !ok || age != 75 {
		t.Fatalf("AgeWithOptions(WithAgeTime) = %d, %v", age, ok)
	}
	age, ok = AgeWithOptions(id, WithAgeClock(func() time.Time {
		return time.Date(2024, 12, 30, 0, 0, 0, 0, time.Local)
	}))
	if !ok || age != 74 {
		t.Fatalf("AgeWithOptions(WithAgeClock) = %d, %v", age, ok)
	}
	loc := time.FixedZone("birthday", 8*3600)
	birthDate, ok := BirthDateWithOptions(id, WithBirthLocation(loc))
	if !ok || birthDate.Location() != loc || birthDate.Format("2006-01-02") != "1949-12-31" {
		t.Fatalf("BirthDateWithOptions() = %v, %v", birthDate, ok)
	}
	if !IsValidBirthdayWithOptions("19491231", WithBirthLocation(loc)) || IsValidBirthdayWithOptions("19490231", WithBirthLocation(loc)) {
		t.Fatal("IsValidBirthdayWithOptions failed")
	}
	if IsValidBirthdayWithOptions("19491231", WithBirthDigitsMatcher(func(string) bool { return false })) {
		t.Fatal("custom birthday digits matcher should reject birthday")
	}
	if _, ok := BirthStringWithOptions(id, WithBirthParser(func(string, string, *time.Location) (time.Time, error) {
		return time.Time{}, errors.New("boom")
	})); ok {
		t.Fatal("custom birthday parser error should reject birth string")
	}
	parsedWithCustomParser := false
	birthDate, ok = BirthDateWithOptions(id, WithBirthParser(func(layout, value string, location *time.Location) (time.Time, error) {
		parsedWithCustomParser = true
		return time.ParseInLocation(layout, value, location)
	}))
	if !ok || !parsedWithCustomParser || birthDate.Format("2006-01-02") != "1949-12-31" {
		t.Fatalf("BirthDateWithOptions custom parser = %v, %v, called=%v", birthDate, ok, parsedWithCustomParser)
	}
	gender, ok := GenderOf(id)
	if !ok || gender != GenderFemale {
		t.Fatalf("GenderOf() = %d, %v", gender, ok)
	}
	province, ok := Province(id)
	if !ok || province != "北京" {
		t.Fatalf("Province() = %q, %v", province, ok)
	}
	district, ok := DistrictCode(id)
	if !ok || district != "110105" {
		t.Fatalf("DistrictCode() = %q, %v", district, ok)
	}
}

func TestInvalidIDCardFieldsReturnFalse(t *testing.T) {
	if birth, ok := BirthString("123"); ok || birth != "" {
		t.Fatalf("BirthString(short) = %q, %v; want empty false", birth, ok)
	}
	if birth, ok := BirthString("11010519490231002X"); ok || birth != "19490231" {
		t.Fatalf("BirthString(invalid date) = %q, %v; want birth false", birth, ok)
	}
	if _, ok := BirthDate("11010519490231002X"); ok {
		t.Fatal("BirthDate should reject invalid birthday")
	}
	if age, ok := AgeAt("123", time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)); ok || age != 0 {
		t.Fatalf("AgeAt(invalid) = %d, %v; want zero false", age, ok)
	}
	if year, ok := Year("123"); ok || year != 0 {
		t.Fatalf("Year(invalid) = %d, %v; want zero false", year, ok)
	}
	if month, ok := Month("123"); ok || month != 0 {
		t.Fatalf("Month(invalid) = %d, %v; want zero false", month, ok)
	}
	if day, ok := Day("123"); ok || day != 0 {
		t.Fatalf("Day(invalid) = %d, %v; want zero false", day, ok)
	}
}

func TestGenderOfCoversFifteenDigitAndInvalidInputs(t *testing.T) {
	gender, ok := GenderOf("130503670401001")
	if !ok || gender != GenderMale {
		t.Fatalf("GenderOf(15-digit) = %d, %v; want male true", gender, ok)
	}
	gender, ok = GenderOf("123")
	if ok || gender != GenderUnknown {
		t.Fatalf("GenderOf(short) = %d, %v; want unknown false", gender, ok)
	}
	gender, ok = GenderOf("1101051949123100XX")
	if ok || gender != GenderUnknown {
		t.Fatalf("GenderOf(non-digit sequence) = %d, %v; want unknown false", gender, ok)
	}
}

func TestRegionCodesRejectInvalidLengths(t *testing.T) {
	if code, ok := ProvinceCode("123"); ok || code != "" {
		t.Fatalf("ProvinceCode(short) = %q, %v; want empty false", code, ok)
	}
	if name, ok := Province("99010519491231002X"); ok || name != "" {
		t.Fatalf("Province(unknown) = %q, %v; want empty false", name, ok)
	}
	if code, ok := CityCode("123"); ok || code != "" {
		t.Fatalf("CityCode(short) = %q, %v; want empty false", code, ok)
	}
	if code, ok := DistrictCode("123"); ok || code != "" {
		t.Fatalf("DistrictCode(short) = %q, %v; want empty false", code, ok)
	}
}

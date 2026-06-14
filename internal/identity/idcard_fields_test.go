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

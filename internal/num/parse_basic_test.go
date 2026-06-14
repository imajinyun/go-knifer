package num

import "testing"

func TestParseBasicHelpers(t *testing.T) {
	if ParseInt("0x10") != 16 || ParseInt("123.56") != 123 || ParseLong("123.56") != 123 {
		t.Fatal("parse integer helpers failed")
	}
	if ParseFloat(".125") != 0.125 || ParseDouble("1,234.5") != 1234.5 || ParseIntDefault("bad", 7) != 7 {
		t.Fatal("parse float/default helpers failed")
	}
}

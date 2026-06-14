package num

import "testing"

func TestParseEdgeCases(t *testing.T) {
	if ParseInt("") != 0 || ParseInt(".5") != 0 || ParseInt("1e3") != 0 || ParseInt("1,234.9") != 1234 {
		t.Fatal("ParseInt edge cases failed")
	}
	if ParseLong("") != 0 || ParseLong(".5") != 0 || ParseLong("0x7f") != 127 || ParseLong("1,234.9") != 1234 {
		t.Fatal("ParseLong edge cases failed")
	}
	if ParseDouble("") != 0 || ParseDouble("1,234.5") != 1234.5 || ParseFloat("2.5") != 2.5 {
		t.Fatal("ParseFloat/ParseDouble edge cases failed")
	}
	if got, err := ParseNumber("+1,234.5"); err != nil || got != 1234.5 {
		t.Fatalf("ParseNumber plus/comma failed: %v %v", got, err)
	}
	if got, err := ParseNumber("0x10"); err != nil || got != 16 {
		t.Fatalf("ParseNumber hex failed: %v %v", got, err)
	}
	if _, err := ParseNumber("bad"); err == nil {
		t.Fatal("ParseNumber should reject invalid input")
	}
}

func TestParseDefaultEdgeCases(t *testing.T) {
	if ParseIntDefault("", 7) != 7 || ParseIntDefault("bad", 7) != 7 || ParseIntDefault("1,234", 7) != 1234 {
		t.Fatal("ParseIntDefault cases failed")
	}
	if ParseLongDefault("", 8) != 8 || ParseLongDefault("bad", 8) != 8 || ParseLongDefault("1,234", 8) != 1234 {
		t.Fatal("ParseLongDefault cases failed")
	}
	if ParseFloatDefault("", 1.5) != 1.5 || ParseFloatDefault("bad", 1.5) != 1.5 || ParseFloatDefault("1,234.5", 1.5) != 1234.5 {
		t.Fatal("ParseFloatDefault cases failed")
	}
	if ParseDoubleDefault("", 2.5) != 2.5 || ParseDoubleDefault("bad", 2.5) != 2.5 || ParseDoubleDefault("1,234.5", 2.5) != 1234.5 {
		t.Fatal("ParseDoubleDefault cases failed")
	}
}

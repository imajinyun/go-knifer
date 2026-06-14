package identity

import "testing"

func TestIsValidIDCard18(t *testing.T) {
	tests := []struct {
		id   string
		want bool
	}{
		{"11010519491231002X", true},
		{"11010519491231002x", true},
		{"81000019980902013X", true},
		{"820000200009100032", true},
		{"83000019810715006X", true},
		{"11010519490231002X", false},
		{"99010519491231002X", false},
		{"110105194912310021", false},
	}
	for _, tt := range tests {
		if got := IsValidIDCard18(tt.id); got != tt.want {
			t.Fatalf("IsValidIDCard18(%q) = %v, want %v", tt.id, got, tt.want)
		}
	}
	if IsValidIDCard18WithIgnoreCase("11010519491231002x", false) {
		t.Fatal("IsValidIDCard18WithIgnoreCase should reject lowercase x when ignoreCase=false")
	}
}

func TestIsValidIDCard15(t *testing.T) {
	if !IsValidIDCard15("130503670401001") {
		t.Fatal("expected valid 15-digit ID card")
	}
	if IsValidIDCard15("130503990230001") {
		t.Fatal("expected invalid birthday to be rejected")
	}
}

package conv

import "testing"

func TestToBool(t *testing.T) {
	cases := map[string]bool{
		"true": true, "yes": true, "y": true, "ok": true, "1": true, "on": true,
		"false": false, "no": false, "n": false, "0": false, "off": false,
	}
	for s, want := range cases {
		if ToBool(s) != want {
			t.Fatalf("ToBool(%q)", s)
		}
	}
	if ToBool(1) != true || ToBool(0) != false {
		t.Fatalf("ToBool int")
	}
	if ToBoolDefault("xx", true) != true {
		t.Fatalf("ToBool default")
	}
}

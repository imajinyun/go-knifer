package regex

import (
	"reflect"
	"testing"
)

func TestBasicRegexCompatibility(t *testing.T) {
	if !ReMatch(`^\d+$`, "123") || ReMatch(`^\d+$`, "12a") || ReMatch(`(`, "x") {
		t.Fatalf("ReMatch failed")
	}
	if ReFind(`\d+`, "ab123cd") != "123" || ReFind(`(`, "x") != "" {
		t.Fatalf("ReFind failed")
	}
	all := ReFindAll(`\d+`, "a1b22c333")
	if !reflect.DeepEqual(all, []string{"1", "22", "333"}) {
		t.Fatalf("ReFindAll failed: %v", all)
	}
	if ReReplace(`\d`, "a1b2", "*") != "a*b*" || ReReplace(`(`, "x", "*") != "x" {
		t.Fatalf("ReReplace failed")
	}
}

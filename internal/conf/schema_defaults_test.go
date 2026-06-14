package conf

import "testing"

func TestApplyDefaults(t *testing.T) {
	c := New()
	c.Set("name", "override")
	withDefaults := c.ApplyDefaults(Schema{Fields: []FieldRule{{Key: "region", Default: "cn"}}})
	if got := withDefaults.Get("region"); got != "cn" {
		t.Fatalf("ApplyDefaults region = %q", got)
	}
}

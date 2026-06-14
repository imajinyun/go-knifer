package bean

import "testing"

func TestCopyPropertiesStructToStructWithAliasAndWeakConversion(t *testing.T) {
	src := sourceProfile{
		embeddedProfile: embeddedProfile{Trace: "t-1"},
		Name:            "alice",
		Age:             "42",
		Admin:           "yes",
		Skip:            "ignored",
	}
	var dst targetProfile
	if err := CopyProperties(src, &dst, WithIgnoreEmpty(true)); err != nil {
		t.Fatalf("CopyProperties() error = %v", err)
	}
	if dst.Name != "alice" || dst.Age != 42 || !dst.Admin || dst.Trace != "t-1" || dst.Empty != "" {
		t.Fatalf("dst = %+v", dst)
	}
}

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

func TestToStructAndCopyPropertiesPointerBoundaries(t *testing.T) {
	type nested struct {
		Value string `bean:"value"`
	}
	type source struct {
		Name   string  `bean:"name"`
		Nested *nested `bean:"nested"`
		Drop   *nested `bean:"drop"`
	}
	type targetNested struct {
		Value string `bean:"value"`
	}
	type target struct {
		Name   *string       `bean:"name"`
		Nested *targetNested `bean:"nested"`
		Drop   *targetNested `bean:"drop"`
	}

	var dst target
	err := ToStruct(source{Name: "alice", Nested: &nested{Value: "ok"}}, &dst)
	if err != nil {
		t.Fatalf("ToStruct() error = %v", err)
	}
	if dst.Name == nil || *dst.Name != "alice" {
		t.Fatalf("Name = %#v", dst.Name)
	}
	if dst.Nested == nil || dst.Nested.Value != "ok" {
		t.Fatalf("Nested = %#v", dst.Nested)
	}
	if dst.Drop != nil {
		t.Fatalf("Drop = %#v, want nil", dst.Drop)
	}
}

func TestCopyPropertiesEmbeddedPointerAndCaseSensitivity(t *testing.T) {
	type embedded struct {
		Trace string `bean:"trace"`
	}
	type source struct {
		*embedded
		Display string `bean:"DISPLAY"`
	}
	type target struct {
		Trace string `bean:"trace"`
		Label string `bean:"display"`
	}

	var insensitive target
	err := CopyProperties(source{embedded: &embedded{Trace: "t-1"}, Display: "alice"}, &insensitive)
	if err != nil {
		t.Fatalf("CopyProperties() case-insensitive error = %v", err)
	}
	if insensitive.Trace != "t-1" || insensitive.Label != "alice" {
		t.Fatalf("case-insensitive dst = %+v", insensitive)
	}

	var sensitive target
	err = CopyProperties(source{embedded: &embedded{Trace: "t-2"}, Display: "bob"}, &sensitive, WithCaseInsensitive(false))
	if err != nil {
		t.Fatalf("CopyProperties() case-sensitive error = %v", err)
	}
	if sensitive.Trace != "t-2" || sensitive.Label != "" {
		t.Fatalf("case-sensitive dst = %+v", sensitive)
	}
}

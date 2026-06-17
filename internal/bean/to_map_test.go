package bean

import "testing"

func TestToMapUsesPrimaryTagAndOmit(t *testing.T) {
	got, err := ToMap(sourceProfile{Name: "alice", Age: "18", Skip: "hidden"})
	if err != nil {
		t.Fatalf("ToMap() error = %v", err)
	}
	if got["name"] != "alice" || got["age"] != "18" {
		t.Fatalf("map = %#v", got)
	}
	if _, ok := got["Skip"]; ok {
		t.Fatalf("omit field leaked: %#v", got)
	}
}

func TestFillMapWithTagOptionsAndIgnoreEmpty(t *testing.T) {
	type source struct {
		Name  string `bean:"name"`
		Alias string `bean:"alias,omitempty,aliases=legacy;display"`
		Empty []int  `bean:"empty"`
		Zero  int    `bean:"zero"`
	}

	dst := map[string]any{"existing": true}
	err := FillMap(source{Name: "alice", Alias: "bob", Empty: []int{}, Zero: 0}, dst,
		WithTagNames("bean"),
		WithIgnoreEmpty(true),
	)
	if err != nil {
		t.Fatalf("FillMap() error = %v", err)
	}
	if dst["name"] != "alice" || dst["alias"] != "bob" || dst["existing"] != true {
		t.Fatalf("dst = %#v", dst)
	}
	if _, ok := dst["empty"]; ok {
		t.Fatalf("empty field leaked: %#v", dst)
	}
	if _, ok := dst["zero"]; !ok {
		t.Fatalf("zero numeric field skipped by IgnoreEmpty: %#v", dst)
	}
}

func TestFillMapRejectsUnsupportedSources(t *testing.T) {
	tests := []struct {
		name string
		src  any
	}{
		{name: "nil pointer", src: (*sourceProfile)(nil)},
		{name: "unsupported", src: 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FillMap(tt.src, map[string]any{})
			assertBeanInvalidInput(t, err)
		})
	}
}

package sets

import "testing"

func TestGenericSetOperations(t *testing.T) {
	s := New("a", "b", "b")
	s.Add("c")
	s.Remove("a", "missing")

	if !s.Equal(New("b", "c")) {
		t.Fatalf("generic string set = %v, want b/c", s.Members())
	}
	if got := s.Sub(New("c")); !got.Equal(New("b")) {
		t.Fatalf("generic Sub() = %v, want b", got.Members())
	}
	if got := s.Union(New("d")); !got.Equal(New("b", "c", "d")) {
		t.Fatalf("generic Union() = %v, want b/c/d", got.Members())
	}
	if got := s.Intersect(New("c", "d")); !got.Equal(New("c")) {
		t.Fatalf("generic Intersect() = %v, want c", got.Members())
	}
}

func TestGenericSetWithStructValues(t *testing.T) {
	type key struct {
		ID   int
		Name string
	}

	a := key{ID: 1, Name: "a"}
	b := key{ID: 2, Name: "b"}
	c := key{ID: 3, Name: "c"}

	s := New(a, b)
	if !s.Contains(a) {
		t.Fatal("generic struct set should contain inserted key")
	}
	if got := s.Union(New(c)); !got.Equal(New(a, b, c)) {
		t.Fatalf("struct Union() = %v, want all keys", got.Members())
	}
}

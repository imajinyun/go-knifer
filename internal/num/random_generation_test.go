package num

import (
	"reflect"
	"testing"
)

func TestRandomGeneration(t *testing.T) {
	randoms := GenRandomNumber(1, 10, 5)
	if len(randoms) != 5 {
		t.Fatalf("GenRandomNumber length: %v", randoms)
	}
	seen := map[int]struct{}{}
	for _, v := range randoms {
		if v < 1 || v >= 10 {
			t.Fatalf("GenRandomNumber value out of range: %v", randoms)
		}
		if _, ok := seen[v]; ok {
			t.Fatalf("GenRandomNumber duplicated value: %v", randoms)
		}
		seen[v] = struct{}{}
	}
	bySet := GenBySet(1, 10, 5)
	if len(bySet) != 5 {
		t.Fatalf("GenBySet length: %v", bySet)
	}
}

func TestRandomGenerationEdges(t *testing.T) {
	if got := GenRandomNumber(10, 1, 3); len(got) != 0 {
		t.Fatalf("GenRandomNumber reversed bounds should be empty because default seed is empty: %v", got)
	}
	if got := GenRandomNumber(1, 3, 5); len(got) != 0 {
		t.Fatalf("GenRandomNumber oversize should be empty: %v", got)
	}
	if got := GenRandomNumberWithSeed(1, 10, 2, []int{7}); len(got) != 0 {
		t.Fatalf("GenRandomNumberWithSeed short seed should be empty: %v", got)
	}
	seed := []int{1, 2, 3, 4}
	got := GenRandomNumberWithSeed(1, 5, 2, seed)
	if len(got) != 2 || !reflect.DeepEqual(seed, []int{1, 2, 3, 4}) {
		t.Fatalf("GenRandomNumberWithSeed should not mutate seed: got=%v seed=%v", got, seed)
	}
	if got := GenBySet(5, 1, 0); len(got) != 0 {
		t.Fatalf("GenBySet zero size should be empty: %v", got)
	}
	if got := GenBySet(1, 2, 3); len(got) != 0 {
		t.Fatalf("GenBySet oversize should be empty: %v", got)
	}
}

func TestRandomGenerationWithOptions(t *testing.T) {
	seed := []int{10, 20, 30, 40}
	got := GenRandomNumberWithSeedWithOptions(0, 4, 3, seed, WithRandomReader(&sequenceReader{}))
	if !reflect.DeepEqual(got, []int{10, 20, 40}) {
		t.Fatalf("GenRandomNumberWithSeedWithOptions deterministic = %v", got)
	}
	if !reflect.DeepEqual(seed, []int{10, 20, 30, 40}) {
		t.Fatalf("GenRandomNumberWithSeedWithOptions should not mutate seed: %v", seed)
	}

	got = GenRandomNumberWithOptions(0, 5, 3, WithRandomReader(&sequenceReader{}))
	if !reflect.DeepEqual(got, []int{0, 1, 2}) {
		t.Fatalf("GenRandomNumberWithOptions deterministic = %v", got)
	}

	got = GenBySetWithOptions(0, 5, 3, WithRandomReader(&sequenceReader{}))
	if len(got) != 3 {
		t.Fatalf("GenBySetWithOptions length = %v", got)
	}
	seen := map[int]bool{}
	for _, v := range got {
		seen[v] = true
	}
	for _, want := range []int{0, 1, 2} {
		if !seen[want] {
			t.Fatalf("GenBySetWithOptions missing %d in %v", want, got)
		}
	}

	if got := GenRandomNumberWithOptions(0, 5, 2, WithRandomReader(errReader{})); !reflect.DeepEqual(got, []int{0, 4}) {
		t.Fatalf("random failure should preserve fallback index behavior: %v", got)
	}
}

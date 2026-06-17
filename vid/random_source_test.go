package vid

import (
	"errors"
	mathrand "math/rand"
	"testing"
	"time"
)

func TestIDFacadeFallbackRandomSourceProvider(t *testing.T) {
	ResetDefaultFallbackRandomSource()
	t.Cleanup(ResetDefaultFallbackRandomSource)

	ConfigureDefaultFallbackRandomSourceProvider(func() *mathrand.Rand {
		return mathrand.New(mathrand.NewSource(13))
	})
	first := SimpleUUIDWithOptions(WithRandomReader(errReader{}))
	ConfigureDefaultFallbackRandomSourceProvider(func() *mathrand.Rand {
		return mathrand.New(mathrand.NewSource(13))
	})
	if got := SimpleUUIDWithOptions(WithRandomReader(errReader{})); got != first {
		t.Fatalf("SimpleUUIDWithOptions after provider reset = %s, want %s", got, first)
	}

	SetFallbackRandomSeed(14)
	seeded := SimpleUUIDWithOptions(WithRandomReader(errReader{}))
	SetFallbackRandomSeed(14)
	if got := SimpleUUIDWithOptions(WithRandomReader(errReader{})); got != seeded {
		t.Fatalf("SimpleUUIDWithOptions after seed reset = %s, want %s", got, seeded)
	}
}

func TestIDFacadeRandomFallbackBoundaries(t *testing.T) {
	source := mathrand.New(mathrand.NewSource(21))
	first := RandomUUIDWithOptions(WithRandomReader(errReader{}), WithFallbackRandomSource(source))
	source = mathrand.New(mathrand.NewSource(21))
	if got := RandomUUIDWithOptions(WithRandomReader(errReader{}), WithFallbackRandomSource(source)); got != first {
		t.Fatalf("RandomUUIDWithOptions fallback = %s, want %s", got, first)
	}

	objectSource := mathrand.New(mathrand.NewSource(22))
	objectID := ObjectIdWithOptions(
		WithObjectIDRandomReader(errReader{}),
		WithObjectIDFallbackRandomSource(objectSource),
		WithObjectIDTimeFunc(func() time.Time { return time.Unix(1_700_000_000, 0) }),
		WithObjectIDCounter(func() uint32 { return 1 }),
	)
	if len(objectID) != 24 {
		t.Fatalf("ObjectIdWithOptions len = %d", len(objectID))
	}

	nanoSource := mathrand.New(mathrand.NewSource(23))
	if got := NanoIdNWithOptions(0, WithNanoIDRandomReader(errReader{}), WithNanoIDFallbackRandomSource(nanoSource)); got != "" {
		t.Fatalf("NanoIdNWithOptions(0) = %q, want empty", got)
	}
	nanoSource = mathrand.New(mathrand.NewSource(23))
	nanoID := NanoIdNWithOptions(6,
		WithNanoIDRandomReader(errReader{}),
		WithNanoIDFallbackRandomSource(nanoSource),
		WithNanoIDAlphabet("ab"),
	)
	if len(nanoID) != 6 {
		t.Fatalf("NanoIdNWithOptions len = %d", len(nanoID))
	}
	for _, r := range nanoID {
		if r != 'a' && r != 'b' {
			t.Fatalf("NanoIdNWithOptions alphabet = %q", nanoID)
		}
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errors.New("boom") }

package id

import (
	"bytes"
	"testing"
)

func TestNanoId(t *testing.T) {
	id := NanoId()
	if len(id) != 21 {
		t.Fatalf("NanoId default len: %s", id)
	}
	id = NanoIdN(10)
	if len(id) != 10 {
		t.Fatalf("NanoIdN len: %s", id)
	}
}

func TestNanoIDOptions(t *testing.T) {
	nid := NanoIdWithOptions(
		WithNanoIDLength(5),
		WithNanoIDAlphabet("ab"),
		WithNanoIDRandomReader(bytes.NewReader([]byte{0, 1, 0, 1, 1})),
	)
	if nid != "ababb" {
		t.Fatalf("NanoIdWithOptions = %q", nid)
	}
}

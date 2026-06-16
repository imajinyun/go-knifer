package mail

import (
	"errors"
	"strings"
	"testing"
)

func TestParseAddressListRejectsCRLF(t *testing.T) {
	_, err := ParseAddressList("to@example.com\nBcc: attacker@example.com")
	if !errors.Is(err, ErrInvalidAddress) {
		t.Fatalf("ParseAddressList() error = %v, want %v", err, ErrInvalidAddress)
	}
}

func TestAddressAndHeaderHelpers(t *testing.T) {
	addr, err := NewAddress("Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("NewAddress() error = %v", err)
	}
	if addr.String() != `"Alice" <alice@example.com>` {
		t.Fatalf("Address.String() = %q", addr.String())
	}
	if empty := (*Address)(nil).String(); empty != "" {
		t.Fatalf("nil Address.String() = %q, want empty", empty)
	}
	if _, err := NewAddress("bad\nname", "alice@example.com"); !errors.Is(err, ErrInvalidAddress) {
		t.Fatalf("NewAddress(CRLF) error = %v, want %v", err, ErrInvalidAddress)
	}
	if _, err := NewAddress("Alice", "Alice <alice@example.com>"); !errors.Is(err, ErrInvalidAddress) {
		t.Fatalf("NewAddress(display email) error = %v, want %v", err, ErrInvalidAddress)
	}

	var header Header
	if err := header.Add("X-Test", "one"); err != nil {
		t.Fatalf("Header.Add() error = %v", err)
	}
	if err := header.Set("X-Test", "two", "three"); err != nil {
		t.Fatalf("Header.Set(existing) error = %v", err)
	}
	if err := header.Set("X-Other", "value"); err != nil {
		t.Fatalf("Header.Set(new) error = %v", err)
	}
	values := header.Values("x-test")
	if strings.Join(values, ",") != "two,three" {
		t.Fatalf("Header.Values() = %v", values)
	}
	values[0] = "mutated"
	if got := header.Values("X-Test")[0]; got != "two" {
		t.Fatalf("Header.Values() returned mutable slice, got %q", got)
	}
	if got := header.Values("missing"); got != nil {
		t.Fatalf("Header.Values(missing) = %v, want nil", got)
	}
	if err := header.Add("Bad:Name", "value"); !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("Header.Add(bad name) error = %v, want %v", err, ErrInvalidHeader)
	}
	if err := header.Add("X-Bad", "value\r\nInjected: true"); !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("Header.Add(bad value) error = %v, want %v", err, ErrInvalidHeader)
	}
}

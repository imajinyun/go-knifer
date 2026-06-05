package vskt_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vskt"
)

func TestFacadeSocketConfig(t *testing.T) {
	cfg := vskt.NewSocketConfig()
	if cfg == nil {
		t.Fatal("expected non-nil socket config")
	}
}

func TestFacadeSocketIsConnected(t *testing.T) {
	// nil conn should not be connected
	if vskt.SocketIsConnected(nil) {
		t.Fatal("expected nil conn to be disconnected")
	}
}

func TestFacadeSocketError(t *testing.T) {
	err := vskt.NewSocketErrorMsg("test error")
	if err == nil {
		t.Fatal("expected non-nil socket error")
	}
	if err.Error() != "test error" {
		t.Fatalf("expected 'test error', got %q", err.Error())
	}
}

func TestFacadeOperations(t *testing.T) {
	// verify operation constants are accessible
	_ = vskt.OpRead
	_ = vskt.OpWrite
	_ = vskt.OpConnect
	_ = vskt.OpAccept
}

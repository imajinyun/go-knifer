package resty

import (
	"testing"
	"time"
)

func TestSnapshotGlobalConfigAndExplicitRequestConfig(t *testing.T) {
	oldTimeout := GetGlobalTimeout()
	oldMaxRedirects := GetGlobalMaxRedirects()
	oldMaxResponse := GetGlobalMaxResponseBytes()
	oldFollow := GetGlobalFollowRedirects()
	oldUA := GetGlobalUserAgent()
	defer SetGlobalTimeout(oldTimeout)
	defer SetGlobalMaxRedirects(oldMaxRedirects)
	defer SetGlobalMaxResponseBytes(oldMaxResponse)
	defer SetGlobalFollowRedirects(oldFollow)
	defer SetGlobalUserAgent(oldUA)
	defer RemoveGlobalHeader("X-Snapshot")

	SetGlobalTimeout(123 * time.Millisecond)
	SetGlobalMaxRedirects(3)
	SetGlobalMaxResponseBytes(321)
	SetGlobalFollowRedirects(false)
	SetGlobalUserAgent("snapshot-agent")
	SetGlobalHeader("X-Snapshot", "one")

	cfg := SnapshotGlobalConfig()
	SetGlobalHeader("X-Snapshot", "mutated")
	cfg.Headers["X-Snapshot"][0] = "cfg"

	req := NewRequestWithConfig(MethodGet, "http://example.com", cfg)
	if req.timeout != 123*time.Millisecond || req.maxRedirects != 3 || req.maxResponse != 321 || req.followRedir == nil || *req.followRedir || req.userAgent != "snapshot-agent" {
		t.Fatalf("request config not applied: timeout=%v max=%d maxResponse=%d follow=%v ua=%q", req.timeout, req.maxRedirects, req.maxResponse, req.followRedir, req.userAgent)
	}
	if got := req.headers["X-Snapshot"]; len(got) != 1 || got[0] != "cfg" {
		t.Fatalf("explicit config headers = %v, want [cfg]", got)
	}
	if got := CloneGlobalHeaders()["X-Snapshot"]; len(got) != 1 || got[0] != "mutated" {
		t.Fatalf("snapshot should be detached from globals, global header = %v", got)
	}
}

func TestClientUsesCapturedConfig(t *testing.T) {
	oldUA := GetGlobalUserAgent()
	oldFollow := GetGlobalFollowRedirects()
	defer SetGlobalUserAgent(oldUA)
	defer SetGlobalFollowRedirects(oldFollow)

	SetGlobalUserAgent("client-agent")
	SetGlobalFollowRedirects(false)
	client := NewClient()
	SetGlobalUserAgent("mutated-agent")
	SetGlobalFollowRedirects(true)

	req := client.Get("https://example.com")
	if req.userAgent != "client-agent" {
		t.Fatalf("client request userAgent = %q, want captured client-agent", req.userAgent)
	}
	if req.followRedir == nil || *req.followRedir {
		t.Fatalf("client request followRedirects = %v, want captured false", req.followRedir)
	}

	isolated := NewIsolatedClient().Get("https://example.com")
	if isolated.userAgent != "" || isolated.followRedir == nil || !*isolated.followRedir {
		t.Fatalf("isolated client defaults ua=%q follow=%v", isolated.userAgent, isolated.followRedir)
	}
}

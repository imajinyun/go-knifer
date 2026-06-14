package resty

import "testing"

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

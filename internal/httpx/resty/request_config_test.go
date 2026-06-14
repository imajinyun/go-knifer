package resty

import (
	"testing"
	"time"
)

func TestNewRequestWithOptionsAppliesRequestOptions(t *testing.T) {
	getReq := Get("http://example.com", WithFollowRedirects(false), WithHeader("X-Create", "get"), WithUserAgent("create-get-agent"))
	if getReq.followRedir == nil || *getReq.followRedir {
		t.Fatalf("followRedir: %v", getReq.followRedir)
	}
	if got := getReq.headers["X-Create"]; len(got) != 1 || got[0] != "get" {
		t.Fatalf("Get header = %q, want get", got)
	}
	if got := getReq.userAgent; got != "create-get-agent" {
		t.Fatalf("Get userAgent = %q", got)
	}

	postReq := Post("http://example.com", WithHeader("X-Create", "post"))
	if postReq.method != MethodPost {
		t.Fatalf("Post method = %v, want POST", postReq.method)
	}
	if got := postReq.headers["X-Create"]; len(got) != 1 || got[0] != "post" {
		t.Fatalf("Post header = %q, want post", got)
	}
}

func TestNewIsolatedRequestDoesNotReadGlobals(t *testing.T) {
	oldTimeout := GetGlobalTimeout()
	oldMaxRedirects := GetGlobalMaxRedirects()
	oldFollow := GetGlobalFollowRedirects()
	oldUA := GetGlobalUserAgent()
	defer SetGlobalTimeout(oldTimeout)
	defer SetGlobalMaxRedirects(oldMaxRedirects)
	defer SetGlobalFollowRedirects(oldFollow)
	defer SetGlobalUserAgent(oldUA)
	defer RemoveGlobalHeader("X-Isolated")

	SetGlobalTimeout(time.Second)
	SetGlobalMaxRedirects(1)
	SetGlobalFollowRedirects(false)
	SetGlobalUserAgent("global-agent")
	SetGlobalHeader("X-Isolated", "global")

	req := NewIsolatedRequest(MethodGet, "http://example.com")
	if req.timeout != defaultGlobalTimeout || req.maxRedirects != 10 || req.maxResponse != defaultGlobalMaxResponseBytes || req.followRedir == nil || !*req.followRedir || req.userAgent != "" {
		t.Fatalf("isolated request leaked globals: timeout=%v max=%d maxResponse=%d follow=%v ua=%q", req.timeout, req.maxRedirects, req.maxResponse, req.followRedir, req.userAgent)
	}
	if got := req.headers["X-Isolated"]; len(got) != 0 {
		t.Fatalf("isolated request should not include global header: %v", got)
	}
}

func TestWithGlobalConfigOptionOverridesConstructionDefaults(t *testing.T) {
	cfg := GlobalConfig{
		Timeout:          250 * time.Millisecond,
		MaxRedirects:     2,
		MaxResponseBytes: 456,
		FollowRedirects:  false,
		DefaultUserAgent: "option-agent",
		Headers:          HeaderValues{"X-Config": []string{"yes"}},
	}
	req := NewIsolatedRequest(MethodGet, "http://example.com", WithGlobalConfig(cfg), WithHeader("X-Req", "ok"))
	if req.timeout != 250*time.Millisecond || req.maxRedirects != 2 || req.maxResponse != 456 || req.followRedir == nil || *req.followRedir || req.userAgent != "option-agent" {
		t.Fatalf("WithGlobalConfig not applied: timeout=%v max=%d maxResponse=%d follow=%v ua=%q", req.timeout, req.maxRedirects, req.maxResponse, req.followRedir, req.userAgent)
	}
	if got := req.headers["X-Config"]; len(got) != 1 || got[0] != "yes" {
		t.Fatalf("config header = %v, want [yes]", got)
	}
	if got := req.headers["X-Req"]; len(got) != 1 || got[0] != "ok" {
		t.Fatalf("request header after config = %v, want [ok]", got)
	}
}

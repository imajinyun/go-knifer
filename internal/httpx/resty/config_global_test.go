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

func TestDefaultGlobalTimeoutIsBounded(t *testing.T) {
	previous := SnapshotGlobalConfig()
	defer ConfigureGlobalConfig(previous)

	ResetGlobalConfig()
	if got := GetGlobalTimeout(); got != defaultGlobalTimeout || got <= 0 {
		t.Fatalf("default timeout = %v, want positive %v", got, defaultGlobalTimeout)
	}
}

func TestResetGlobalConfigRestoresDefaults(t *testing.T) {
	previous := SnapshotGlobalConfig()
	defer ConfigureGlobalConfig(previous)

	SetGlobalTimeout(time.Second)
	SetGlobalMaxRedirects(2)
	SetGlobalMaxResponseBytes(3)
	SetGlobalFollowRedirects(false)
	SetGlobalUserAgent("mutated-agent")
	SetGlobalHeader("X-Reset", "mutated")
	CloseCookie()

	ResetGlobalConfig()
	cfg := SnapshotGlobalConfig()
	if cfg.Timeout != defaultGlobalTimeout || cfg.MaxRedirects != 10 || cfg.MaxResponseBytes != defaultGlobalMaxResponseBytes || !cfg.FollowRedirects || cfg.DefaultUserAgent != "" || cfg.CookieDisabled {
		t.Fatalf("reset scalar config = %#v", cfg)
	}
	if got := cfg.Headers["X-Reset"]; len(got) != 0 {
		t.Fatalf("reset retained X-Reset header: %v", got)
	}
	if got := cfg.Headers[string(HeaderUserAgent)]; len(got) == 0 || got[0] == "" {
		t.Fatalf("reset default User-Agent header = %v", got)
	}
}

func TestWithScopedGlobalConfigRestoresPreviousDefaults(t *testing.T) {
	previous := SnapshotGlobalConfig()
	defer ConfigureGlobalConfig(previous)

	ConfigureGlobalConfig(GlobalConfig{
		Timeout:          time.Second,
		MaxRedirects:     4,
		MaxResponseBytes: 64,
		FollowRedirects:  true,
		DefaultUserAgent: "outer-agent",
		Headers:          HeaderValues{"X-Scope": []string{"outer"}},
	})

	WithScopedGlobalConfig(GlobalConfig{
		Timeout:          2 * time.Second,
		MaxRedirects:     1,
		MaxResponseBytes: 32,
		FollowRedirects:  false,
		DefaultUserAgent: "inner-agent",
		Headers:          HeaderValues{"X-Scope": []string{"inner"}},
		CookieDisabled:   true,
	}, func() {
		cfg := SnapshotGlobalConfig()
		if cfg.Timeout != 2*time.Second || cfg.MaxRedirects != 1 || cfg.MaxResponseBytes != 32 || cfg.FollowRedirects || cfg.DefaultUserAgent != "inner-agent" || cfg.Headers["X-Scope"][0] != "inner" || !cfg.CookieDisabled {
			t.Fatalf("scoped inner config = %#v", cfg)
		}
	})

	cfg := SnapshotGlobalConfig()
	if cfg.Timeout != time.Second || cfg.MaxRedirects != 4 || cfg.MaxResponseBytes != 64 || !cfg.FollowRedirects || cfg.DefaultUserAgent != "outer-agent" || cfg.Headers["X-Scope"][0] != "outer" || cfg.CookieDisabled {
		t.Fatalf("scoped restored config = %#v", cfg)
	}
}

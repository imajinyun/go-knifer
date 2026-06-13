package resty

import (
	grestry "resty.dev/v3"
	"testing"
	"time"
)

func TestGlobalHeadersArePlainValues(t *testing.T) {
	SetGlobalHeader("X-Resty-Plain", "one")
	AddGlobalHeader("X-Resty-Plain", "two")
	defer RemoveGlobalHeader("X-Resty-Plain")

	headers := CloneGlobalHeaders()
	if got := headers["X-Resty-Plain"]; len(got) != 2 || got[0] != "one" || got[1] != "two" {
		t.Fatalf("CloneGlobalHeaders()[X-Resty-Plain] = %v, want [one two]", got)
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

func TestRestyClientFactoryProviderLifecycle(t *testing.T) {
	ResetDefaultRestyClientProvider()
	t.Cleanup(ResetDefaultRestyClientProvider)

	defaultCalled := 0
	ConfigureDefaultRestyClientProvider(func() *grestry.Client {
		defaultCalled++
		return grestry.New()
	})
	client := NewIsolatedRequest(MethodGet, "http://example.com").buildClient()
	if client == nil || defaultCalled != 1 {
		t.Fatalf("default provider client=%v called=%d", client, defaultCalled)
	}

	perCallCalled := 0
	client = NewIsolatedRequest(MethodGet, "http://example.com", WithRestyClientFactory(func() *grestry.Client {
		perCallCalled++
		return grestry.New()
	})).buildClient()
	if client == nil || perCallCalled != 1 || defaultCalled != 1 {
		t.Fatalf("per-call factory client=%v perCall=%d default=%d", client, perCallCalled, defaultCalled)
	}

	client = NewIsolatedRequest(MethodGet, "http://example.com", WithRestyClientFactory(func() *grestry.Client { return nil })).buildClient()
	if client == nil || defaultCalled != 2 {
		t.Fatalf("nil per-call factory client=%v default=%d", client, defaultCalled)
	}

	ResetDefaultRestyClientProvider()
	client = NewIsolatedRequest(MethodGet, "http://example.com").buildClient()
	if client == nil {
		t.Fatal("reset default provider should create a client")
	}
}

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

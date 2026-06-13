package resty

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetWithQueryAndHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "go" {
			t.Fatalf("query q = %q, want go", r.URL.Query().Get("q"))
		}
		if r.Header.Get("X-Test") != "yes" {
			t.Fatalf("X-Test = %q, want yes", r.Header.Get("X-Test"))
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	resp := Get(srv.URL).Query("q", "go").Header("X-Test", "yes").Execute()
	if resp.Err() != nil {
		t.Fatalf("Execute() error = %v", resp.Err())
	}
	if !resp.IsOK() || resp.Body() != "ok" {
		t.Fatalf("status/body = %d/%q, want 2xx/ok", resp.Status(), resp.Body())
	}
}

func TestPostForm(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm() error = %v", err)
		}
		if got := r.Form.Get("name"); got != "resty" {
			t.Fatalf("form name = %q, want resty", got)
		}
		_, _ = w.Write([]byte("posted"))
	}))
	defer srv.Close()

	got, err := PostFormE(srv.URL, map[string]any{"name": "resty"})
	if err != nil {
		t.Fatalf("PostFormE() error = %v", err)
	}
	if got != "posted" {
		t.Fatalf("PostFormE() = %q, want posted", got)
	}
}

func TestPostJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), string(ContentTypeJSON)) {
			t.Fatalf("Content-Type = %q, want application/json", r.Header.Get("Content-Type"))
		}
		_, _ = w.Write([]byte("json"))
	}))
	defer srv.Close()

	got, err := PostJSONE(srv.URL, `{"ok":true}`)
	if err != nil {
		t.Fatalf("PostJSONE() error = %v", err)
	}
	if got != "json" {
		t.Fatalf("PostJSONE() = %q, want json", got)
	}
}

func TestRequestJSONMarshalProvider(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	called := false
	resp := Post(srv.URL, WithJSONMarshalFunc(func(any) ([]byte, error) {
		called = true
		return []byte(`{"provided":true}`), nil
	})).BodyJSONValue(map[string]any{"ignored": true}).Execute()
	if resp.Err() != nil {
		t.Fatalf("Execute() error = %v", resp.Err())
	}
	if !called || resp.Body() != `{"provided":true}` {
		t.Fatalf("marshal provider called=%v body=%q", called, resp.Body())
	}
}

func TestRequestJSONUnmarshalProvider(t *testing.T) {
	type result struct {
		Name string `json:"name"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"ignored"}`))
	}))
	defer srv.Close()

	called := false
	out := &result{}
	resp := Get(srv.URL, WithJSONUnmarshalFunc(func(_ []byte, dst any) error {
		called = true
		return json.Unmarshal([]byte(`{"name":"provided"}`), dst)
	})).Result(out).Execute()
	if resp.Err() != nil {
		t.Fatalf("Execute() error = %v", resp.Err())
	}
	if !called || out.Name != "provided" || resp.Result() == nil {
		t.Fatalf("unmarshal provider called=%v result=%+v raw=%v", called, out, resp.Result())
	}
}

func TestRequestJSONDecodeReadOptions(t *testing.T) {
	type result struct {
		Name string `json:"name"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"abcdef"}`))
	}))
	defer srv.Close()

	tooLarge := &result{}
	resp := Get(srv.URL,
		WithMaxDecodeBytes(3),
		WithJSONUnmarshalFunc(json.Unmarshal),
	).Result(tooLarge).Execute()
	if resp.Err() == nil {
		t.Fatal("Execute() with max decode bytes error = nil")
	}

	readAllCalled := false
	out := &result{}
	resp = Get(srv.URL,
		WithJSONDecodeReadAllFunc(func(io.Reader) ([]byte, error) {
			readAllCalled = true
			return []byte(`{"name":"provided"}`), nil
		}),
		WithJSONUnmarshalFunc(json.Unmarshal),
	).Result(out).Execute()
	if resp.Err() != nil || !readAllCalled || out.Name != "provided" {
		t.Fatalf("custom decode readAll called=%v out=%+v err=%v", readAllCalled, out, resp.Err())
	}
}

func TestTimeoutReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		_, _ = w.Write([]byte("late"))
	}))
	defer srv.Close()

	resp := Get(srv.URL).Timeout(time.Millisecond).Execute()
	if resp.Err() == nil {
		t.Fatal("Execute() error is nil, want timeout error")
	}
}

func TestRequestOptionsOverrideGlobalDefaults(t *testing.T) {
	oldUA := GetGlobalUserAgent()
	oldFollow := GetGlobalFollowRedirects()
	defer SetGlobalUserAgent(oldUA)
	defer SetGlobalFollowRedirects(oldFollow)

	SetGlobalUserAgent("global-resty-agent")
	SetGlobalFollowRedirects(false)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/start" {
			http.Redirect(w, r, "/end", http.StatusFound)
			return
		}
		_, _ = w.Write([]byte(r.Header.Get("X-Req") + ":" + r.Header.Get("User-Agent")))
	}))
	defer srv.Close()

	resp := Get(srv.URL+"/start",
		WithHeader("X-Req", "per-call"),
		WithUserAgent("request-resty-agent"),
		WithFollowRedirects(true),
	).Execute()
	if resp.Err() != nil {
		t.Fatalf("Execute() error = %v", resp.Err())
	}
	if got := resp.Body(); got != "per-call:request-resty-agent" {
		t.Fatalf("Body() = %q, want per-call options to override globals", got)
	}
}

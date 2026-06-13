package conf

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	knifer "github.com/imajinyun/go-knifer"
)

func TestLoadProfileAndParseByExt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.toml")
	if err := os.WriteFile(path, []byte("name='base'\n[profile.test]\nname='test'"), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := LoadProfile(path, "test")
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Get("name"); got != "test" {
		t.Fatalf("LoadProfile name = %q", got)
	}

	yamlConf, err := ParseByExt("app.yaml", []byte("app:\n  name: demo"))
	if err != nil {
		t.Fatal(err)
	}
	if got := yamlConf.GetByGroup("app", "name"); got != "demo" {
		t.Fatalf("ParseByExt yaml app.name = %q", got)
	}

	profileYAML, err := ParseYAMLFull(`
app:
  name: base
server:
  port: 8080
profile:
  dev:
    app:
      name: dev
    server:
      port: 9090
`)
	if err != nil {
		t.Fatal(err)
	}
	dev := profileYAML.ApplyProfile("dev")
	if got := dev.GetByGroup("app", "name"); got != "dev" {
		t.Fatalf("YAML profile app.name = %q", got)
	}
	if got := dev.GetByGroup("server", "port"); got != "9090" {
		t.Fatalf("YAML profile server.port = %q", got)
	}

	custom, err := ParseByExtWithOptions("app.custom", []byte("ignored"), WithParserForExt("custom", func([]byte) (*Conf, error) {
		c := New()
		c.Set("name", "custom-parser")
		return c, nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	if got := custom.Get("name"); got != "custom-parser" {
		t.Fatalf("custom parser name = %q", got)
	}
}

func TestLoadWithOptionsPassesParseOptions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.yaml")
	if err := os.WriteFile(path, []byte("ignored"), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := LoadWithOptions(path, LoadOptions{ParseOptions: []ParseOption{WithYAMLUnmarshalFunc(func([]byte, any) error {
		return errors.New("custom yaml error")
	})}})
	if err == nil {
		t.Fatalf("LoadWithOptions = %#v, nil error", c)
	}
	assertConfCode(t, err, knifer.ErrCodeInvalidInput)
}

func TestLoadWithOptionsIncludesMergeDecryptAndSchema(t *testing.T) {
	dir := t.TempDir()
	common := filepath.Join(dir, "common.setting")
	main := filepath.Join(dir, "main.setting")
	secret := base64.StdEncoding.EncodeToString([]byte("s3cr3t"))
	if err := os.WriteFile(common, []byte("name=common\n[server]\nhost=127.0.0.1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(main, []byte("include=common.setting\nname=main\nsecret=ENC(base64:"+secret+")\n[server]\nport=8080"), 0o644); err != nil {
		t.Fatal(err)
	}

	c, err := LoadWithOptions(main, LoadOptions{AllowInclude: true})
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Get("name"); got != "main" {
		t.Fatalf("merged name = %q", got)
	}
	if got := c.Get("secret"); got != "s3cr3t" {
		t.Fatalf("decrypted secret = %q", got)
	}
	if got := c.GetByGroup("server", "host"); got != "127.0.0.1" {
		t.Fatalf("included server.host = %q", got)
	}
	if _, ok := c.Lookup("", "include"); ok {
		t.Fatal("include key should be removed after loading")
	}
	if err := c.ValidateSchema(Schema{Fields: []FieldRule{
		{Key: "name", Required: true, Type: TypeString},
		{Group: "server", Key: "port", Required: true, Type: TypeInt},
		{Group: "server", Key: "host", Required: true},
	}}); err != nil {
		t.Fatalf("ValidateSchema() error = %v", err)
	}
	if err := c.ValidateSchema(Schema{Fields: []FieldRule{{Group: "server", Key: "debug", Required: true}}}); err == nil {
		t.Fatal("ValidateSchema() missing required error = nil")
	}
}

func TestLoadWithOptionsIncludeRejectsPathTraversal(t *testing.T) {
	dir := t.TempDir()
	outside := filepath.Join(dir, "outside.setting")
	confDir := filepath.Join(dir, "conf")
	if err := os.Mkdir(confDir, 0o755); err != nil {
		t.Fatal(err)
	}
	main := filepath.Join(confDir, "main.setting")
	if err := os.WriteFile(outside, []byte("secret=outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(main, []byte("include=../outside.setting\nname=main"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadWithOptions(main, LoadOptions{AllowInclude: true})
	assertConfCode(t, err, knifer.ErrCodeInvalidInput)
}

func TestLoadWithOptionsIncludeRejectsAbsolutePathByDefault(t *testing.T) {
	dir := t.TempDir()
	common := filepath.Join(dir, "common.setting")
	main := filepath.Join(dir, "main.setting")
	if err := os.WriteFile(common, []byte("name=common"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(main, []byte("include="+common+"\nname=main"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadWithOptions(main, LoadOptions{AllowInclude: true})
	assertConfCode(t, err, knifer.ErrCodeInvalidInput)
}

func TestLoadWithOptionsIncludeRootAllowsConfiguredDirectory(t *testing.T) {
	dir := t.TempDir()
	root := filepath.Join(dir, "root")
	serviceDir := filepath.Join(root, "service")
	commonDir := filepath.Join(root, "common")
	if err := os.MkdirAll(serviceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(commonDir, 0o755); err != nil {
		t.Fatal(err)
	}
	common := filepath.Join(commonDir, "base.setting")
	main := filepath.Join(serviceDir, "main.setting")
	if err := os.WriteFile(common, []byte("name=common"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(main, []byte("include=../common/base.setting\nmode=service"), 0o644); err != nil {
		t.Fatal(err)
	}

	c, err := LoadWithOptions(main, LoadOptions{AllowInclude: true, IncludeRoot: root})
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Get("name"); got != "common" {
		t.Fatalf("included name = %q", got)
	}
	if got := c.Get("mode"); got != "service" {
		t.Fatalf("main mode = %q", got)
	}
}

func TestLoadFilesAndApplyDefaults(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "base.setting")
	override := filepath.Join(dir, "override.toml")
	if err := os.WriteFile(base, []byte("name=base\nmode=dev"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(override, []byte("name='override'\n[server]\nport=9090"), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := LoadFiles(base, override)
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Get("name"); got != "override" {
		t.Fatalf("LoadFiles merged name = %q", got)
	}
	withDefaults := c.ApplyDefaults(Schema{Fields: []FieldRule{{Key: "region", Default: "cn"}}})
	if got := withDefaults.Get("region"); got != "cn" {
		t.Fatalf("ApplyDefaults region = %q", got)
	}
}

func TestLoadRemoteWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Config-Token") != "secret" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		_, _ = w.Write([]byte("app:\n  name: remote"))
	}))
	defer server.Close()
	calledFactory := false
	c, err := LoadRemoteWithOptions(server.URL+"/app.yaml", LoadOptions{
		Timeout: time.Second,
		Headers: http.Header{"X-Config-Token": []string{"secret"}},
		RequestFactory: func(ctx context.Context, rawURL string) (*http.Request, error) {
			calledFactory = true
			return http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !calledFactory {
		t.Fatal("request factory was not called")
	}
	if got := c.GetByGroup("app", "name"); got != "remote" {
		t.Fatalf("remote app.name = %q", got)
	}
	if _, err := LoadRemoteWithOptions(server.URL+"/app.yaml", LoadOptions{MaxBytes: 3, Headers: http.Header{"X-Config-Token": []string{"secret"}}}); err == nil {
		t.Fatal("LoadRemoteWithOptions max bytes error = nil")
	}
}

func TestLoadRemoteSafeRejectsPrivateHostsAndUnsafeRedirects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("app:\n  name: remote"))
	}))
	defer server.Close()

	if _, err := LoadRemoteSafe(server.URL + "/app.yaml"); err == nil {
		t.Fatal("LoadRemoteSafe should reject private hosts by default")
	}
	if _, err := LoadRemoteSafe("http://224.0.0.1/app.yaml"); err == nil {
		t.Fatal("LoadRemoteSafe should reject multicast hosts by default")
	}
	if _, err := LoadRemoteSafe("http://0.0.0.0/app.yaml"); err == nil {
		t.Fatal("LoadRemoteSafe should reject unspecified hosts by default")
	}
	remoteURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := LoadRemoteSafeWithOptions(server.URL+"/app.yaml", LoadOptions{RemoteAllowedHosts: []string{remoteURL.Hostname()}}); err == nil {
		t.Fatal("LoadRemoteSafeWithOptions should reject allowlisted private hosts")
	}

	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://127.0.0.1/private.yaml", http.StatusFound)
	}))
	defer redirect.Close()
	redirectURL, err := url.Parse(redirect.URL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := LoadRemoteSafeWithOptions(redirect.URL+"/app.yaml", LoadOptions{RemoteAllowedHosts: []string{redirectURL.Hostname()}}); err == nil {
		t.Fatal("LoadRemoteSafeWithOptions should reject unsafe redirect target")
	}
}

func TestLoadRemoteSafeAllowedHostsDoesNotBypassPrivateRejection(t *testing.T) {
	if _, err := LoadRemoteSafeWithOptions("http://127.0.0.1/app.yaml", LoadOptions{RemoteAllowedHosts: []string{"127.0.0.1"}}); err == nil {
		t.Fatal("LoadRemoteSafeWithOptions should reject allowlisted loopback host")
	}

	lookupCount := 0
	_, err := LoadRemoteSafeWithOptions("http://config.example/app.yaml", LoadOptions{
		RemoteAllowedHosts: []string{"config.example"},
		LookupIP: func(context.Context, string) ([]net.IP, error) {
			lookupCount++
			return []net.IP{net.ParseIP("10.0.0.1")}, nil
		},
	})
	if err == nil {
		t.Fatal("LoadRemoteSafeWithOptions should reject allowlisted host resolving to private address")
	}
	if lookupCount == 0 {
		t.Fatal("LoadRemoteSafeWithOptions did not resolve allowlisted host for private-address validation")
	}
}

func TestLoadRemoteSafeAllowsAllowedPublicHost(t *testing.T) {
	client := &http.Client{Transport: confRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("app:\n  name: remote")),
			Request:    r,
		}, nil
	})}
	c, err := LoadRemoteSafeWithOptions("http://config.example/app.yaml", LoadOptions{
		RemoteClient:       client,
		RemoteAllowedHosts: []string{"config.example"},
		LookupIP: func(context.Context, string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("93.184.216.34")}, nil
		},
	})
	if err != nil {
		t.Fatalf("LoadRemoteSafeWithOptions allowed public host: %v", err)
	}
	if got := c.GetByGroup("app", "name"); got != "remote" {
		t.Fatalf("remote app.name = %q", got)
	}
}

func TestLoadRemoteSafeRevalidatesHostAtRoundTrip(t *testing.T) {
	lookups := [][]net.IP{{net.ParseIP("93.184.216.34")}, {net.ParseIP("127.0.0.1")}}
	lookupCount := 0
	client := &http.Client{Transport: confRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("unsafe request reached base transport")
		return nil, nil
	})}
	_, err := LoadRemoteSafeWithOptions("http://example.com/app.yaml", LoadOptions{
		RemoteClient: client,
		LookupIP: func(context.Context, string) ([]net.IP, error) {
			if lookupCount >= len(lookups) {
				return lookups[len(lookups)-1], nil
			}
			ips := lookups[lookupCount]
			lookupCount++
			return ips, nil
		},
	})
	if err == nil {
		t.Fatal("LoadRemoteSafeWithOptions should reject a host that resolves private during RoundTrip")
	}
	if lookupCount != 2 {
		t.Fatalf("lookup count = %d, want 2", lookupCount)
	}
}

func TestLoadWithOptionsReadFileProvider(t *testing.T) {
	c, err := LoadWithOptions("virtual.setting", LoadOptions{
		MaxBytes: 16,
		ReadFile: func(path string, maxBytes int64) ([]byte, error) {
			if path != "virtual.setting" || maxBytes != 16 {
				t.Fatalf("read path=%q maxBytes=%d", path, maxBytes)
			}
			return []byte("name=fake"), nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Get("name"); got != "fake" {
		t.Fatalf("loaded name = %q", got)
	}
}

func TestLoadWithOptionsReadFileProviderUsesDefaultMaxBytes(t *testing.T) {
	_, err := LoadWithOptions("virtual.setting", LoadOptions{
		ReadFile: func(path string, maxBytes int64) ([]byte, error) {
			if maxBytes != DefaultMaxBytes {
				t.Fatalf("default maxBytes=%d, want %d", maxBytes, DefaultMaxBytes)
			}
			return []byte("name=fake"), nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadWithOptionsAllowsExplicitUnlimitedMaxBytes(t *testing.T) {
	_, err := LoadWithOptions("virtual.setting", LoadOptions{
		MaxBytes: -1,
		ReadFile: func(path string, maxBytes int64) ([]byte, error) {
			if maxBytes != -1 {
				t.Fatalf("maxBytes=%d, want -1", maxBytes)
			}
			return []byte("name=fake"), nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadWithOptionsReadFileProviderEnforcesMaxBytes(t *testing.T) {
	_, err := LoadWithOptions("virtual.setting", LoadOptions{
		MaxBytes: 4,
		ReadFile: func(path string, maxBytes int64) ([]byte, error) {
			return []byte("name=fake"), nil
		},
	})
	assertConfCode(t, err, knifer.ErrCodeInvalidInput)
}

type confRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f confRoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

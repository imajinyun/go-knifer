package vconf_test

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	knifer "github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vconf"
)

func TestAdvancedLoadAndSchemaFacade(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "base.setting")
	main := filepath.Join(dir, "main.setting")
	secret := base64.StdEncoding.EncodeToString([]byte("token"))
	if err := os.WriteFile(base, []byte("name=base\n[server]\nhost=localhost"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(main, []byte("import=base.setting\nname=main\nsecret=ENC(base64:"+secret+")\n[server]\nport=8080"), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := vconf.LoadWithOptions(main, vconf.LoadOptions{AllowInclude: true})
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Get("name"); got != "main" {
		t.Fatalf("LoadWithOptions name = %q", got)
	}
	if got := c.Get("secret"); got != "token" {
		t.Fatalf("LoadWithOptions secret = %q", got)
	}
	if err := c.ValidateSchema(vconf.Schema{Fields: []vconf.FieldRule{
		{Key: "name", Required: true},
		{Group: "server", Key: "port", Required: true, Type: vconf.TypeInt},
	}}); err != nil {
		t.Fatalf("ValidateSchema() error = %v", err)
	}
	merged, err := vconf.LoadFiles(base, main)
	if err != nil {
		t.Fatal(err)
	}
	if got := merged.Get("name"); got != "main" {
		t.Fatalf("LoadFiles name = %q", got)
	}
	type cfg struct {
		Name string `conf:"name,required"`
	}
	schema, err := vconf.SchemaFromStruct(cfg{})
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Fields) != 1 {
		t.Fatalf("SchemaFromStruct fields = %d", len(schema.Fields))
	}
}

func TestLoadWithOptionsIncludeRootFacade(t *testing.T) {
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

	_, err := vconf.LoadWithOptions(main, vconf.LoadOptions{AllowInclude: true})
	if err == nil {
		t.Fatal("LoadWithOptions path traversal error = nil")
	}
	if !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("errors.Is(err, ErrCodeInvalidInput) = false: %v", err)
	}

	c, err := vconf.LoadWithOptions(main, vconf.LoadOptions{AllowInclude: true, IncludeRoot: root})
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

func TestLoadFilesAndRemoteWithOptionsFacade(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "base.setting")
	main := filepath.Join(dir, "main.setting")
	if err := os.WriteFile(base, []byte("name=base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(main, []byte("name=main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	merged, err := vconf.LoadFilesWithOptions(vconf.LoadOptions{MaxBytes: 64}, base, main)
	if err != nil {
		t.Fatalf("LoadFilesWithOptions() error = %v", err)
	}
	if got := merged.Get("name"); got != "main" {
		t.Fatalf("LoadFilesWithOptions() name = %q, want main", got)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Config-Token"); got != "secret" {
			t.Fatalf("remote header X-Config-Token = %q, want secret", got)
		}
		_, _ = w.Write([]byte("remote: true\n"))
	}))
	defer server.Close()
	remote, err := vconf.LoadRemoteWithOptions(server.URL+"/app.yaml", vconf.LoadOptions{
		Headers:  http.Header{"X-Config-Token": []string{"secret"}},
		Timeout:  time.Second,
		MaxBytes: 64,
	})
	if err != nil {
		t.Fatalf("LoadRemoteWithOptions() error = %v", err)
	}
	if got := remote.Get("remote"); got != "true" {
		t.Fatalf("LoadRemoteWithOptions() remote = %q, want true", got)
	}
}

func TestFacadeRemoteSafeAndParseWrappers(t *testing.T) {
	trustedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("trusted=ok\n"))
	}))
	defer trustedServer.Close()
	trusted, err := vconf.LoadRemote(trustedServer.URL + "/app.setting")
	if err != nil {
		t.Fatalf("LoadRemote() error = %v", err)
	}
	if got := trusted.Get("trusted"); got != "ok" {
		t.Fatalf("LoadRemote trusted = %q", got)
	}

	client := &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if got := req.Header.Get("X-Remote-Token"); got != "token" {
			t.Fatalf("remote header = %q, want token", got)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/plain"}},
			Body:       io.NopCloser(strings.NewReader("remote=ok\n")),
			Request:    req,
		}, nil
	})}
	lookupPublic := func(context.Context, string) ([]net.IP, error) {
		return []net.IP{net.ParseIP("8.8.8.8")}, nil
	}
	opts := vconf.LoadOptions{
		RemoteClient:       client,
		Headers:            http.Header{"X-Remote-Token": []string{"token"}},
		RemoteAllowedHosts: []string{"config.example"},
		LookupIP:           lookupPublic,
		Timeout:            time.Second,
		MaxBytes:           64,
	}

	remote, err := vconf.LoadRemoteWithOptions("http://config.example/app.setting", opts)
	if err != nil {
		t.Fatalf("LoadRemoteWithOptions() error = %v", err)
	}
	if got := remote.Get("remote"); got != "ok" {
		t.Fatalf("LoadRemoteWithOptions remote = %q", got)
	}
	safe, err := vconf.LoadRemoteSafeWithOptions("http://config.example/app.setting", opts)
	if err != nil {
		t.Fatalf("LoadRemoteSafeWithOptions() error = %v", err)
	}
	if got := safe.Get("remote"); got != "ok" {
		t.Fatalf("LoadRemoteSafeWithOptions remote = %q", got)
	}
	if _, err := vconf.LoadRemoteSafe("http://127.0.0.1/app.setting"); err == nil {
		t.Fatal("LoadRemoteSafe private host error = nil")
	}

	parsed, err := vconf.ParseByExt("app.setting", []byte("name=parse"))
	if err != nil || parsed.Get("name") != "parse" {
		t.Fatalf("ParseByExt = %#v, %v", parsed, err)
	}
	yaml, err := vconf.ParseYAMLFull("server:\n  port: 8080\n")
	if err != nil || yaml.GetByGroup("server", "port") != "8080" {
		t.Fatalf("ParseYAMLFull = %#v, %v", yaml, err)
	}
	decoded, err := vconf.Base64Decrypt(base64.StdEncoding.EncodeToString([]byte("secret")))
	if err != nil || decoded != "secret" {
		t.Fatalf("Base64Decrypt = %q, %v", decoded, err)
	}
}

package vresty_test

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	knifer "github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vresty"
	grestry "resty.dev/v3"
)

func TestFacadeGetString(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("facade"))
	}))
	defer srv.Close()

	got, err := vresty.GetStringE(srv.URL)
	if err != nil {
		t.Fatalf("GetStringE() error = %v", err)
	}
	if got != "facade" {
		t.Fatalf("GetStringE() = %q, want facade", got)
	}
}

func TestFacadeBuildBasicAuth(t *testing.T) {
	if got := vresty.BuildBasicAuth("u", "p"); got != "Basic dTpw" {
		t.Fatalf("BuildBasicAuth() = %q, want Basic dTpw", got)
	}
}

func TestFacadeCloneGlobalHeaders(t *testing.T) {
	previous := vresty.SnapshotGlobalConfig()
	defer vresty.ConfigureGlobalConfig(previous)

	vresty.SetGlobalHeader("X-Facade", "one")
	vresty.AddGlobalHeader("X-Facade", "two")

	headers := vresty.CloneGlobalHeaders()
	if got := headers["X-Facade"]; len(got) != 2 || got[0] != "one" || got[1] != "two" {
		t.Fatalf("CloneGlobalHeaders()[X-Facade] = %v, want [one two]", got)
	}
}

func TestFacadeRequestFollowRedirectOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Header.Get("X-Opt") + ":" + r.Header.Get("User-Agent")))
	}))
	defer srv.Close()

	resp := vresty.Get(srv.URL,
		vresty.WithHeader("X-Opt", "yes"),
		vresty.WithUserAgent("vresty-test/1.0"),
	).Execute()
	if resp.Err() != nil {
		t.Fatal(resp.Err())
	}
	if got := resp.Body(); got != "yes:vresty-test/1.0" {
		t.Fatalf("Body() = %q, want option headers", got)
	}
}

func TestFacadeRestyClientFactoryProvider(t *testing.T) {
	vresty.ResetDefaultRestyClientProvider()
	t.Cleanup(vresty.ResetDefaultRestyClientProvider)

	called := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Header.Get("X-Factory")))
	}))
	defer server.Close()

	resp := vresty.Get(server.URL,
		vresty.WithRestyClientFactory(func() *grestry.Client {
			called++
			return grestry.New().SetHeader("X-Factory", "per-call")
		}),
	).Execute()
	if resp.Err() != nil {
		t.Fatal(resp.Err())
	}
	if called != 1 || resp.Body() != "per-call" {
		t.Fatalf("factory called=%d body=%q", called, resp.Body())
	}

	vresty.ConfigureDefaultRestyClientProvider(func() *grestry.Client {
		called++
		return grestry.New().SetHeader("X-Factory", "default")
	})
	resp = vresty.Get(server.URL).Execute()
	if resp.Err() != nil {
		t.Fatal(resp.Err())
	}
	if called != 2 || resp.Body() != "default" {
		t.Fatalf("default provider called=%d body=%q", called, resp.Body())
	}
}

func TestFacadeRequestOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/final", http.StatusFound)
			return
		}
		_, _ = w.Write([]byte(r.Method + ":" + r.Header.Get("X-Create")))
	}))
	defer srv.Close()

	getResp := vresty.Get(srv.URL+"/redirect", vresty.WithFollowRedirects(false), vresty.WithHeader("X-Create", "get")).Execute()
	if getResp.Err() != nil {
		t.Fatal(getResp.Err())
	}
	if got := getResp.Status(); got != http.StatusFound {
		t.Fatalf("Get status = %d, want 302", got)
	}

	postResp := vresty.Post(srv.URL, vresty.WithHeader("X-Create", "post")).Execute()
	if postResp.Err() != nil {
		t.Fatal(postResp.Err())
	}
	if got := postResp.Body(); got != "POST:post" {
		t.Fatalf("Post body = %q, want POST:post", got)
	}
}

func TestFacadeAdditionalMethodsGlobalAndContentHelpers(t *testing.T) {
	var lastMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastMethod = r.Method
		w.Header().Set("X-Method", r.Method)
		if r.Method != http.MethodHead {
			_, _ = w.Write([]byte(r.Method))
		}
	}))
	defer srv.Close()

	tests := []struct {
		name   string
		method string
		req    *vresty.Request
	}{
		{name: "put", method: http.MethodPut, req: vresty.Put(srv.URL)},
		{name: "delete", method: http.MethodDelete, req: vresty.Delete(srv.URL)},
		{name: "patch", method: http.MethodPatch, req: vresty.Patch(srv.URL)},
		{name: "head", method: http.MethodHead, req: vresty.Head(srv.URL)},
		{name: "options", method: http.MethodOptions, req: vresty.Options(srv.URL)},
		{name: "new request", method: http.MethodTrace, req: vresty.NewRequest(vresty.MethodTrace, srv.URL)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := tt.req.Execute()
			if resp.Err() != nil {
				t.Fatalf("Execute: %v", resp.Err())
			}
			if lastMethod != tt.method {
				t.Fatalf("server method = %q, want %q", lastMethod, tt.method)
			}
			if got := resp.Header("X-Method"); got != tt.method {
				t.Fatalf("response method header = %q, want %q", got, tt.method)
			}
		})
	}

	previous := vresty.SnapshotGlobalConfig()
	defer vresty.ConfigureGlobalConfig(previous)
	vresty.SetGlobalMaxRedirects(4)
	vresty.SetGlobalMaxResponseBytes(123)
	vresty.SetGlobalFollowRedirects(false)
	vresty.SetGlobalUserAgent("vresty-extra/1.0")
	vresty.CloseCookie()
	vresty.SetGlobalHeader("X-Extra", "one")
	vresty.RemoveGlobalHeader("X-Extra")
	cfg := vresty.SnapshotGlobalConfig()
	if vresty.GetGlobalMaxRedirects() != 4 || vresty.GetGlobalMaxResponseBytes() != 123 || vresty.GetGlobalFollowRedirects() || vresty.GetGlobalUserAgent() != "vresty-extra/1.0" || !cfg.CookieDisabled || len(cfg.Headers["X-Extra"]) != 0 {
		t.Fatalf("global config = %#v", cfg)
	}

	if got := vresty.BuildContentType("application/json", "utf-8"); got != "application/json;charset=utf-8" {
		t.Fatalf("BuildContentType = %q", got)
	}
	if got := vresty.GuessContentType("<root/>"); got != vresty.ContentTypeXML {
		t.Fatalf("GuessContentType = %q", got)
	}
	if !vresty.IsDefaultContentType("") || !vresty.IsFormURLEncoded("application/x-www-form-urlencoded; charset=utf-8") {
		t.Fatal("content type predicates returned unexpected result")
	}
	if got := vresty.URLWithForm("https://example.com/path?x=1", map[string]any{"q": "go"}); got != "https://example.com/path?x=1&q=go" {
		t.Fatalf("URLWithForm = %q", got)
	}
	if got := vresty.GetCharsetFromContentTypeWithOptions("text/plain; enc=gbk", vresty.WithCharsetRegexp(regexp.MustCompile(`enc=([a-z0-9-]+)`))); got != "gbk" {
		t.Fatalf("GetCharsetFromContentTypeWithOptions = %q", got)
	}
	if got := vresty.GetCharsetFromHTMLWithOptions(`<meta data-charset="big5">`, vresty.WithMetaCharsetRegexp(regexp.MustCompile(`data-charset="([^"]+)"`))); got != "big5" {
		t.Fatalf("GetCharsetFromHTMLWithOptions = %q", got)
	}
	if got := vresty.GetMimeType("payload.zip"); got != "application/zip" {
		t.Fatalf("GetMimeType = %q", got)
	}
}

func TestFacadeSafeShortcutAndDownloadHelpers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, _ := io.ReadAll(r.Body)
			_, _ = w.Write([]byte(r.Method + ":" + string(body)))
			return
		}
		_, _ = w.Write([]byte("download"))
	}))
	defer srv.Close()

	allowLocal := vresty.WithURLPolicy(vresty.URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false})
	if _, err := vresty.GetStringSafeE(srv.URL); err == nil {
		t.Fatal("GetStringSafeE(localhost default policy) error = nil")
	}
	if got, err := vresty.GetStringSafeE(srv.URL, allowLocal); err != nil || got != "download" {
		t.Fatalf("GetStringSafeE allowed = %q, %v", got, err)
	}
	if got, err := vresty.PostStringSafeE(srv.URL, "body", allowLocal); err != nil || got != "POST:body" {
		t.Fatalf("PostStringSafeE allowed = %q, %v", got, err)
	}
	if got, err := vresty.DownloadBytesSafeE(srv.URL, allowLocal); err != nil || string(got) != "download" {
		t.Fatalf("DownloadBytesSafeE allowed = %q, %v", got, err)
	}
	var buf bytes.Buffer
	if n, err := vresty.Download(srv.URL, &buf); err != nil || n != int64(len("download")) || buf.String() != "download" {
		t.Fatalf("Download n=%d body=%q err=%v", n, buf.String(), err)
	}
}

func TestFacadeUtilityWrappers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, _ := io.ReadAll(r.Body)
			_, _ = w.Write([]byte("post:" + string(body) + ":" + r.Header.Get("X-Util")))
			return
		}
		_, _ = w.Write([]byte(r.URL.Query().Get("q") + ":" + r.Header.Get("X-Util")))
	}))
	defer srv.Close()

	got, err := vresty.GetWithParamsEWithOptions(srv.URL, map[string]any{"q": "go"}, vresty.WithHeader("X-Util", "get"))
	if err != nil {
		t.Fatalf("GetWithParamsEWithOptions() error = %v", err)
	}
	if got != "go:get" {
		t.Fatalf("GetWithParamsEWithOptions() = %q, want go:get", got)
	}
	got, err = vresty.PostStringEWithOptions(srv.URL, "body", vresty.WithHeader("X-Util", "post"))
	if err != nil {
		t.Fatalf("PostStringEWithOptions() error = %v", err)
	}
	if got != "post:body:post" {
		t.Fatalf("PostStringEWithOptions() = %q, want post:body:post", got)
	}
	if !vresty.IsHTTP("http://example.com") || !vresty.IsHTTPS("https://example.com") {
		t.Fatal("IsHTTP/IsHTTPS wrappers returned false")
	}
	if got := vresty.ToParams(map[string]any{"q": "go"}); got != "q=go" {
		t.Fatalf("ToParams() = %q, want q=go", got)
	}
}

func TestFacadeRequestGlobalConfigAPIs(t *testing.T) {
	previous := vresty.SnapshotGlobalConfig()
	defer vresty.ConfigureGlobalConfig(previous)

	vresty.SetGlobalTimeout(321 * time.Millisecond)
	vresty.SetGlobalHeader("X-Facade-Config", "global")

	cfg := vresty.SnapshotGlobalConfig()
	cfg.Headers["X-Facade-Config"][0] = "snapshot"
	cfg.DefaultUserAgent = "facade-config-agent"
	cfg.Headers["User-Agent"] = []string{"facade-config-agent"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Header.Get("X-Facade-Config") + ":" + r.Header.Get("User-Agent")))
	}))
	defer srv.Close()

	resp := vresty.NewRequestWithConfig(vresty.MethodGet, srv.URL, cfg).Execute()
	if resp.Err() != nil {
		t.Fatal(resp.Err())
	}
	if got := resp.Body(); got != "snapshot:facade-config-agent" {
		t.Fatalf("NewRequestWithConfig body = %q", got)
	}

	resp = vresty.NewIsolatedRequest(vresty.MethodGet, srv.URL, vresty.WithGlobalConfig(cfg)).Execute()
	if resp.Err() != nil {
		t.Fatal(resp.Err())
	}
	if got := resp.Body(); got != "snapshot:facade-config-agent" {
		t.Fatalf("NewIsolatedRequest WithGlobalConfig body = %q", got)
	}
}

func TestFacadeScopedGlobalConfig(t *testing.T) {
	previous := vresty.SnapshotGlobalConfig()
	defer vresty.ConfigureGlobalConfig(previous)

	vresty.ResetGlobalConfig()
	vresty.WithScopedGlobalConfig(vresty.GlobalConfig{
		Timeout:          3 * time.Second,
		MaxRedirects:     1,
		MaxResponseBytes: 32,
		FollowRedirects:  false,
		DefaultUserAgent: "facade-scope-agent",
		Headers:          vresty.HeaderValues{"X-Facade-Scope": []string{"inner"}},
		CookieDisabled:   true,
	}, func() {
		cfg := vresty.SnapshotGlobalConfig()
		if cfg.Timeout != 3*time.Second || cfg.MaxRedirects != 1 || cfg.MaxResponseBytes != 32 || cfg.FollowRedirects || cfg.DefaultUserAgent != "facade-scope-agent" || cfg.Headers["X-Facade-Scope"][0] != "inner" || !cfg.CookieDisabled {
			t.Fatalf("facade scoped config = %#v", cfg)
		}
	})

	cfg := vresty.SnapshotGlobalConfig()
	if cfg.Timeout != 30*time.Second || cfg.MaxRedirects != 10 || cfg.MaxResponseBytes != 64<<20 || !cfg.FollowRedirects || len(cfg.Headers["X-Facade-Scope"]) != 0 || cfg.CookieDisabled {
		t.Fatalf("facade config not restored after scoped helper: %#v", cfg)
	}
}

func TestFacadeSaveProviderOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("vresty-save"))
	}))
	defer server.Close()

	var mkdirPath string
	var mkdirPerm fs.FileMode
	var openPath string
	var openFlag int
	var openPerm fs.FileMode
	var written bytes.Buffer
	n, err := vresty.Get(server.URL).Execute().SaveAs("/virtual/out.txt",
		vresty.WithSaveMkdirAll(func(path string, perm fs.FileMode) error {
			mkdirPath, mkdirPerm = path, perm
			return nil
		}),
		vresty.WithSaveOpenFile(func(path string, flag int, perm fs.FileMode) (io.WriteCloser, error) {
			openPath, openFlag, openPerm = path, flag, perm
			return nopWriteCloser{Writer: &written}, nil
		}),
		vresty.WithSaveDirPerm(0o700), vresty.WithSaveFilePerm(0o600),
	)
	if err != nil || n != int64(len("vresty-save")) {
		t.Fatalf("SaveAs provider n=%d err=%v", n, err)
	}
	if mkdirPath != "/virtual" || mkdirPerm != 0o700 || openPath != "/virtual/out.txt" || openPerm != 0o600 || openFlag&os.O_CREATE == 0 || written.String() != "vresty-save" {
		t.Fatalf("providers mkdir=%q/%v open=%q flag=%#x perm=%v content=%q", mkdirPath, mkdirPerm, openPath, openFlag, openPerm, written.String())
	}
}

func TestFacadeClientAndSafeRequestWrappers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Method", r.Method)
		if r.Method != http.MethodHead {
			_, _ = w.Write([]byte(r.Method + ":" + r.Header.Get("X-Client-Default")))
		}
	}))
	defer server.Close()

	client := vresty.NewClient(vresty.WithClientRequestOptions(vresty.WithHeader("X-Client-Default", "shared")))
	if got := client.Get(server.URL).Execute().Body(); got != "GET:shared" {
		t.Fatalf("client.Get body = %q", got)
	}
	if got := client.Post(server.URL).Execute().Body(); got != "POST:shared" {
		t.Fatalf("client.Post body = %q", got)
	}
	if got := client.NewRequest(vresty.MethodPut, server.URL).Execute().Body(); got != "PUT:shared" {
		t.Fatalf("client.NewRequest body = %q", got)
	}

	cfg := vresty.SnapshotGlobalConfig()
	cfg.Headers["X-Client-Default"] = []string{"configured"}
	if got := vresty.NewClientWithConfig(cfg).Get(server.URL).Execute().Body(); got != "GET:configured" {
		t.Fatalf("NewClientWithConfig body = %q", got)
	}
	if got := vresty.NewIsolatedClient(vresty.WithClientGlobalConfig(cfg)).Get(server.URL).Execute().Body(); got != "GET:configured" {
		t.Fatalf("NewIsolatedClient body = %q", got)
	}

	if resp := vresty.GetSafe(server.URL).Execute(); resp.Err() == nil {
		t.Fatal("GetSafe(localhost default policy) error = nil")
	}
	allowLocal := vresty.WithURLPolicy(vresty.URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false})
	tests := []struct {
		name   string
		method string
		req    *vresty.Request
	}{
		{name: "post safe", method: http.MethodPost, req: vresty.PostSafe(server.URL, allowLocal)},
		{name: "put safe", method: http.MethodPut, req: vresty.PutSafe(server.URL, allowLocal)},
		{name: "delete safe", method: http.MethodDelete, req: vresty.DeleteSafe(server.URL, allowLocal)},
		{name: "patch safe", method: http.MethodPatch, req: vresty.PatchSafe(server.URL, allowLocal)},
		{name: "head safe", method: http.MethodHead, req: vresty.HeadSafe(server.URL, allowLocal)},
		{name: "options safe", method: http.MethodOptions, req: vresty.OptionsSafe(server.URL, allowLocal)},
		{name: "new safe", method: http.MethodTrace, req: vresty.NewSafeRequest(vresty.MethodTrace, server.URL, allowLocal)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := tt.req.Execute()
			if resp.Err() != nil {
				t.Fatalf("Execute: %v", resp.Err())
			}
			if got := resp.Header("X-Method"); got != tt.method {
				t.Fatalf("method header = %q, want %q", got, tt.method)
			}
		})
	}
}

func TestFacadeRequestOptionWrappers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/final", http.StatusFound)
			return
		}
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		_, _ = w.Write([]byte(r.Method + ":" + string(body) + ":" + r.Header.Get("X-A") + ":" + r.Header.Get("X-B")))
	}))
	defer server.Close()

	resp := vresty.Post(server.URL,
		vresty.WithTimeout(time.Second),
		vresty.WithHeaders(map[string]string{"X-A": "one", "X-B": "two"}),
		vresty.WithContentType(string(vresty.ContentTypeTextPlain)),
		vresty.WithCharset("utf-8"),
		vresty.WithCookieDisabled(true),
		vresty.WithMaxResponseBytes(1024),
		vresty.WithMaxDecodeBytes(1024),
		vresty.WithJSONMarshalFunc(func(v any) ([]byte, error) { return []byte(`"custom"`), nil }),
		vresty.WithJSONUnmarshalFunc(func([]byte, any) error { return nil }),
		vresty.WithJSONDecodeReadAllFunc(io.ReadAll),
	).Body([]byte("payload")).Execute()
	if resp.Err() != nil {
		t.Fatalf("Post Execute: %v", resp.Err())
	}
	if got := resp.Body(); got != "POST:payload:one:two" {
		t.Fatalf("option body = %q", got)
	}

	redirect := vresty.Get(server.URL+"/redirect", vresty.WithMaxRedirects(0), vresty.WithFollowRedirects(false)).Execute()
	if redirect.Err() != nil {
		t.Fatalf("redirect Execute: %v", redirect.Err())
	}
	if got := redirect.Status(); got != http.StatusFound {
		t.Fatalf("redirect status = %d, want 302", got)
	}

	restyClient := grestry.New().SetHeader("X-A", "resty")
	withClient := vresty.Get(server.URL, vresty.WithRestyClient(restyClient)).Execute()
	if withClient.Err() != nil {
		t.Fatalf("WithRestyClient Execute: %v", withClient.Err())
	}
	if !strings.Contains(withClient.Body(), ":resty:") {
		t.Fatalf("WithRestyClient body = %q", withClient.Body())
	}
}

func TestFacadeDownloadAndFileWrappers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("download-text"))
	}))
	defer server.Close()

	allowLocal := vresty.WithURLPolicy(vresty.URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false})
	var buf bytes.Buffer
	if n, err := vresty.DownloadSafe(server.URL, &buf, allowLocal); err != nil || n != int64(len("download-text")) || buf.String() != "download-text" {
		t.Fatalf("DownloadSafe n=%d body=%q err=%v", n, buf.String(), err)
	}
	if b, err := vresty.DownloadBytesE(server.URL); err != nil || string(b) != "download-text" {
		t.Fatalf("DownloadBytesE = %q, %v", b, err)
	}
	if b, err := vresty.DownloadBytesEWithOptions(server.URL, vresty.WithMaxResponseBytes(64)); err != nil || string(b) != "download-text" {
		t.Fatalf("DownloadBytesEWithOptions = %q, %v", b, err)
	}
	if got, err := vresty.DownloadStringE(server.URL, ""); err != nil || got != "download-text" {
		t.Fatalf("DownloadStringE = %q, %v", got, err)
	}
	if got, err := vresty.DownloadStringEWithOptions(server.URL, "", vresty.WithMaxResponseBytes(64)); err != nil || got != "download-text" {
		t.Fatalf("DownloadStringEWithOptions = %q, %v", got, err)
	}
	if got, err := vresty.DownloadStringSafeE(server.URL, "", allowLocal); err != nil || got != "download-text" {
		t.Fatalf("DownloadStringSafeE = %q, %v", got, err)
	}

	dir := t.TempDir()
	file := filepath.Join(dir, "plain.txt")
	if n, err := vresty.DownloadFile(server.URL, file); err != nil || n != int64(len("download-text")) {
		t.Fatalf("DownloadFile n=%d err=%v", n, err)
	}
	if data, err := os.ReadFile(file); err != nil || string(data) != "download-text" {
		t.Fatalf("DownloadFile content = %q, %v", data, err)
	}
	fileWithOpts := filepath.Join(dir, "with-options.txt")
	if n, err := vresty.DownloadFileWithOptions(server.URL, fileWithOpts, []vresty.RequestOption{vresty.WithMaxResponseBytes(64)}, vresty.WithSaveOverwrite(true)); err != nil || n != int64(len("download-text")) {
		t.Fatalf("DownloadFileWithOptions n=%d err=%v", n, err)
	}
	safeFile := filepath.Join(dir, "safe.txt")
	if n, err := vresty.DownloadFileSafe(server.URL, safeFile, vresty.WithSaveOverwrite(true)); err == nil || n != 0 {
		t.Fatalf("DownloadFileSafe default policy n=%d err=%v, want private host rejection", n, err)
	}
	if n, err := vresty.DownloadFileSafeWithOptions(server.URL, safeFile, []vresty.RequestOption{allowLocal}, vresty.WithSaveOverwrite(true)); err != nil || n != int64(len("download-text")) {
		t.Fatalf("DownloadFileSafeWithOptions n=%d err=%v", n, err)
	}
}

func TestFacadeErrorsCharsetAndSaveOptions(t *testing.T) {
	cause := errors.New("network closed")
	err := vresty.NewHTTPError("request failed", cause)
	if !errors.Is(err, cause) || !errors.Is(err, knifer.ErrCodeInternal) {
		t.Fatalf("NewHTTPError does not unwrap cause or code: %v", err)
	}
	if got := vresty.HTTPErrorf("status %d", http.StatusBadGateway).Error(); got != "status 502" {
		t.Fatalf("HTTPErrorf = %q", got)
	}
	if got := vresty.GetCharsetFromContentType("text/plain; charset=gb18030"); got != "gb18030" {
		t.Fatalf("GetCharsetFromContentType = %q", got)
	}
	if got := vresty.GetCharsetFromHTML(`<meta charset="utf-8">`); got != "utf-8" {
		t.Fatalf("GetCharsetFromHTML = %q", got)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("stat"))
	}))
	defer server.Close()

	info := fakeFileInfo{isDir: true}
	statCalled := false
	if _, err := vresty.Get(server.URL).Execute().SaveAs("/virtual-dir",
		vresty.WithSaveStat(func(path string) (os.FileInfo, error) {
			statCalled = path == "/virtual-dir"
			return info, nil
		}),
		vresty.WithSaveDefaultFilename("fallback.txt"),
		vresty.WithSaveCreateParents(false),
		vresty.WithSaveOpenFile(func(string, int, fs.FileMode) (io.WriteCloser, error) {
			return nopWriteCloser{Writer: io.Discard}, nil
		}),
	); err != nil {
		t.Fatalf("SaveAs with stat provider: %v", err)
	}
	if !statCalled {
		t.Fatal("WithSaveStat provider was not called")
	}
}

type fakeFileInfo struct{ isDir bool }

func (f fakeFileInfo) Name() string       { return "fake" }
func (f fakeFileInfo) Size() int64        { return 0 }
func (f fakeFileInfo) Mode() fs.FileMode  { return fs.ModeDir }
func (f fakeFileInfo) ModTime() time.Time { return time.Time{} }
func (f fakeFileInfo) IsDir() bool        { return f.isDir }
func (f fakeFileInfo) Sys() any           { return nil }

type nopWriteCloser struct{ io.Writer }

func (w nopWriteCloser) Close() error { return nil }

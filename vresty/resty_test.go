package vresty_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imajinyun/go-knifer/vresty"
)

func TestFacadeGetString(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("facade"))
	}))
	defer srv.Close()

	if got := vresty.GetString(srv.URL); got != "facade" {
		t.Fatalf("GetString() = %q, want facade", got)
	}
}

func TestFacadeBuildBasicAuth(t *testing.T) {
	if got := vresty.BuildBasicAuth("u", "p"); got != "Basic dTpw" {
		t.Fatalf("BuildBasicAuth() = %q, want Basic dTpw", got)
	}
}

func TestFacadeCloneGlobalHeaders(t *testing.T) {
	vresty.SetGlobalHeader("X-Facade", "one")
	vresty.AddGlobalHeader("X-Facade", "two")
	defer vresty.RemoveGlobalHeader("X-Facade")

	headers := vresty.CloneGlobalHeaders()
	if got := headers["X-Facade"]; len(got) != 2 || got[0] != "one" || got[1] != "two" {
		t.Fatalf("CloneGlobalHeaders()[X-Facade] = %v, want [one two]", got)
	}
}

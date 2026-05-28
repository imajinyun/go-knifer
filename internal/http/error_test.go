package http

import (
	"errors"
	"testing"
)

// 对应 hutool-http HttpException 相关用例

func TestHTTPErrorMessage(t *testing.T) {
	e := NewHTTPError("read failed", errors.New("conn closed"))
	if e.Error() != "read failed: conn closed" {
		t.Fatalf("error: %q", e.Error())
	}
	if e.Unwrap() == nil {
		t.Fatal("unwrap nil")
	}
}

func TestHTTPErrorf(t *testing.T) {
	e := HTTPErrorf("status %d", 500)
	if e.Error() != "status 500" {
		t.Fatalf("error: %q", e.Error())
	}
	if e.Unwrap() != nil {
		t.Fatal("unwrap should be nil")
	}
}

func TestStatusHelpers(t *testing.T) {
	if !IsRedirected(301) {
		t.Fatal("301")
	}
	if !IsRedirected(302) {
		t.Fatal("302")
	}
	if IsRedirected(200) {
		t.Fatal("200 should not")
	}
	if IsRedirected(404) {
		t.Fatal("404 should not")
	}
}

func TestMethodString(t *testing.T) {
	if MethodGet.String() != "GET" {
		t.Fatalf("get: %q", MethodGet.String())
	}
	if MethodPost.String() != "POST" {
		t.Fatalf("post: %q", MethodPost.String())
	}
}

func TestHeaderString(t *testing.T) {
	if HeaderContentType.String() != "Content-Type" {
		t.Fatalf("ct: %q", HeaderContentType.String())
	}
	if HeaderUserAgent.String() != "User-Agent" {
		t.Fatalf("ua: %q", HeaderUserAgent.String())
	}
}

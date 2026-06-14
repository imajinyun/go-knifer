package http

import "testing"

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

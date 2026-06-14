package http

import "testing"

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

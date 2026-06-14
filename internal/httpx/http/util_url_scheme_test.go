package http

import "testing"

func TestIsHTTP(t *testing.T) {
	if !IsHTTP("Http://aaa.bbb") {
		t.Fatal("Http://")
	}
	if !IsHTTP("HTTP://aaa.bbb") {
		t.Fatal("HTTP://")
	}
	if IsHTTP("FTP://aaa.bbb") {
		t.Fatal("FTP://")
	}
}

func TestIsHTTPS(t *testing.T) {
	if !IsHTTPS("Https://aaa.bbb") {
		t.Fatal("Https://")
	}
	if !IsHTTPS("HTTPS://aaa.bbb") {
		t.Fatal("HTTPS://")
	}
	if !IsHTTPS("https://aaa.bbb") {
		t.Fatal("https://")
	}
	if IsHTTPS("ftp://aaa.bbb") {
		t.Fatal("ftp://")
	}
}

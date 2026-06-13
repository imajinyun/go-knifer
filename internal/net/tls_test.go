package net

import (
	"crypto/tls"
	"io"
	"strings"
	"testing"
)

func TestTLSReaderWithOptionsUsesReadAll(t *testing.T) {
	b := NewTLSConfigBuilder()
	called := false
	err := b.AddRootCAReaderWithOptions(strings.NewReader("ignored"), WithTLSReadAll(func(r io.Reader) ([]byte, error) {
		called = true
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		if string(data) != "ignored" {
			t.Fatalf("reader data = %q", data)
		}
		return []byte("not a certificate"), nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("custom TLS readAll provider was not called")
	}
}

func TestTLSHelpers(t *testing.T) {
	cfg := NewTLSConfigBuilder().SetMinVersion(tls.VersionTLS12).SetServerName("example.com").Build()
	if cfg.MinVersion != tls.VersionTLS12 || cfg.ServerName != "example.com" {
		t.Fatalf("TLS builder failed: %#v", cfg)
	}
	if TLSVersion(TLSv13) != tls.VersionTLS13 {
		t.Fatal("TLSVersion failed")
	}
}

func TestTLSRootCAProviderOptions(t *testing.T) {
	const certPEM = `-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIRAPWQSq0Qr7yZD5twH61BxFIwCgYIKoZIzj0EAwIwEjEQ
MA4GA1UEChMHZ28tdGVzdDAeFw0yNjA2MDYwMDAwMDBaFw0yNzA2MDYwMDAwMDBa
MBIxEDAOBgNVBAoTB2dvLXRlc3QwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASm
1YPqMC7UTw4R7ovbHYgk4+LALoU6hr61VnsBiKCdsMCMScpLob8ldIl+6o4f/ntM
5kmXvEFd9Mp6FfaHkgnbo0IwQDAOBgNVHQ8BAf8EBAMCAqQwDwYDVR0TAQH/BAUw
AwEB/zAdBgNVHQ4EFgQUX90U1OkOXbGUzD2JNoWlqQtk3/0wCgYIKoZIzj0EAwID
SQAwRgIhANw7UzN0vtxOfygWqANg00uGOo7y98q1/Ac3N1wQxVBkAiEA7QjQRHtH
LA6wKo8yoCnW36b+nvxlhHvzrIxwWCgwCWM=
-----END CERTIFICATE-----`
	readPath := ""
	b := NewTLSConfigBuilder()
	if err := b.AddRootCAFileWithOptions("ca.pem", WithTLSReadFile(func(path string) ([]byte, error) {
		readPath = path
		return []byte(certPEM), nil
	})); err != nil {
		t.Fatalf("AddRootCAFileWithOptions: %v", err)
	}
	if readPath != "ca.pem" || b.Build().RootCAs == nil {
		t.Fatalf("TLS read provider not applied path=%q cfg=%#v", readPath, b.Build())
	}

	b = NewTLSConfigBuilder()
	if err := b.AddRootCAReader(strings.NewReader(certPEM)); err != nil {
		t.Fatalf("AddRootCAReader: %v", err)
	}
	if b.Build().RootCAs == nil {
		t.Fatal("AddRootCAReader should initialize RootCAs")
	}
}

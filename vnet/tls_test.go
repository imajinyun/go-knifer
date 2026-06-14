package vnet_test

import (
	"strings"
	"testing"

	"github.com/imajinyun/go-knifer/vnet"
)

func TestVNetTLSFileOptionsFacade(t *testing.T) {
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
	b := vnet.NewTLSConfigBuilder()
	if err := b.AddRootCAFileWithOptions("ca.pem", vnet.WithTLSReadFile(func(path string) ([]byte, error) {
		readPath = path
		return []byte(certPEM), nil
	})); err != nil {
		t.Fatalf("AddRootCAFileWithOptions: %v", err)
	}
	if readPath != "ca.pem" || b.Build().RootCAs == nil {
		t.Fatalf("TLS read provider not applied path=%q cfg=%#v", readPath, b.Build())
	}

	b = vnet.NewTLSConfigBuilder()
	if err := b.AddRootCAReader(strings.NewReader(certPEM)); err != nil || b.Build().RootCAs == nil {
		t.Fatalf("AddRootCAReader rootCAs=%#v err=%v", b.Build().RootCAs, err)
	}
}

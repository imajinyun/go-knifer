package vnet

import (
	"crypto/tls"
	"crypto/x509"
	"io"

	netimpl "github.com/imajinyun/go-knifer/internal/net"
)

func NewTLSConfigBuilder() *TLSConfigBuilder { return netimpl.NewTLSConfigBuilder() }

func WithTLSReadFile(readFile func(string) ([]byte, error)) TLSFileOption {
	return netimpl.WithTLSReadFile(readFile)
}

func WithTLSReadAll(readAll func(io.Reader) ([]byte, error)) TLSFileOption {
	return netimpl.WithTLSReadAll(readAll)
}

func AddRootCAFileWithOptions(b *TLSConfigBuilder, path string, opts ...TLSFileOption) error {
	return b.AddRootCAFileWithOptions(path, opts...)
}

func AddRootCAReader(b *TLSConfigBuilder, r io.Reader) error { return b.AddRootCAReader(r) }

func AddRootCAReaderWithOptions(b *TLSConfigBuilder, r io.Reader, opts ...TLSFileOption) error {
	return b.AddRootCAReaderWithOptions(r, opts...)
}

func AddRootCABytes(b *TLSConfigBuilder, pem []byte) error { return b.AddRootCABytes(pem) }

func CreateTLSConfig(insecureSkipVerify bool) *tls.Config {
	return netimpl.CreateTLSConfig(insecureSkipVerify)
}

func InsecureTLSConfig() *tls.Config { return netimpl.InsecureTLSConfig() }

func TLSVersion(protocol string) uint16 { return netimpl.TLSVersion(protocol) }

func NewCertPool() *x509.CertPool { return x509.NewCertPool() }

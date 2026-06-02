package net

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

const (
	// SSL is a legacy SSL protocol label.
	SSL = "SSL"
	// SSLv2 is a legacy SSLv2 protocol label.
	SSLv2 = "SSLv2"
	// SSLv3 is a legacy SSLv3 protocol label.
	SSLv3 = "SSLv3"
	// TLS is the TLS protocol label.
	TLS = "TLS"
	// TLSv1 is TLS 1.0.
	TLSv1 = "TLSv1"
	// TLSv11 is TLS 1.1.
	TLSv11 = "TLSv1.1"
	// TLSv12 is TLS 1.2.
	TLSv12 = "TLSv1.2"
	// TLSv13 is TLS 1.3.
	TLSv13 = "TLSv1.3"
)

// TLSConfigBuilder builds tls.Config values.
type TLSConfigBuilder struct {
	config tls.Config
}

// NewTLSConfigBuilder creates a TLS config builder.
func NewTLSConfigBuilder() *TLSConfigBuilder { return &TLSConfigBuilder{} }

// SetMinVersion sets the minimum TLS version.
func (b *TLSConfigBuilder) SetMinVersion(version uint16) *TLSConfigBuilder {
	b.config.MinVersion = version
	return b
}

// SetMaxVersion sets the maximum TLS version.
func (b *TLSConfigBuilder) SetMaxVersion(version uint16) *TLSConfigBuilder {
	b.config.MaxVersion = version
	return b
}

// SetInsecureSkipVerify controls certificate verification.
func (b *TLSConfigBuilder) SetInsecureSkipVerify(skip bool) *TLSConfigBuilder {
	b.config.InsecureSkipVerify = skip //nolint:gosec // Caller explicitly chooses trust behavior for this helper.
	return b
}

// SetServerName sets the TLS server name.
func (b *TLSConfigBuilder) SetServerName(name string) *TLSConfigBuilder {
	b.config.ServerName = name
	return b
}

// SetRootCAs sets root CAs.
func (b *TLSConfigBuilder) SetRootCAs(pool *x509.CertPool) *TLSConfigBuilder {
	b.config.RootCAs = pool
	return b
}

// AddRootCAFile appends PEM certificates from path to RootCAs.
func (b *TLSConfigBuilder) AddRootCAFile(path string) error {
	pem, err := os.ReadFile(path) // #nosec G304 -- caller controls certificate path.
	if err != nil {
		return err
	}
	pool := b.config.RootCAs
	if pool == nil {
		pool = x509.NewCertPool()
	}
	pool.AppendCertsFromPEM(pem)
	b.config.RootCAs = pool
	return nil
}

// SetCertificates sets client certificates.
func (b *TLSConfigBuilder) SetCertificates(certs []tls.Certificate) *TLSConfigBuilder {
	b.config.Certificates = certs
	return b
}

// Build returns a cloned tls.Config.
func (b *TLSConfigBuilder) Build() *tls.Config { return b.config.Clone() }

// CreateTLSConfig creates a TLS config with optional insecure verification.
func CreateTLSConfig(insecureSkipVerify bool) *tls.Config {
	return (&tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: insecureSkipVerify}).Clone() //nolint:gosec // Caller explicitly chooses trust behavior.
}

// InsecureTLSConfig creates a TLS config that skips certificate verification.
func InsecureTLSConfig() *tls.Config { return CreateTLSConfig(true) }

// TLSVersion maps a protocol label to crypto/tls version constants.
func TLSVersion(protocol string) uint16 {
	switch protocol {
	case TLSv1:
		return tls.VersionTLS10
	case TLSv11:
		return tls.VersionTLS11
	case TLSv12:
		return tls.VersionTLS12
	case TLSv13:
		return tls.VersionTLS13
	default:
		return tls.VersionTLS12
	}
}

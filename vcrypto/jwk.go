package vcrypto

import (
	"crypto/rsa"

	cryptoimpl "github.com/imajinyun/knifer-go/internal/crypto"
)

// JWK represents a JSON Web Key. RSA keys are supported in this package.
type JWK = cryptoimpl.JWK

// JWKS represents a JSON Web Key Set.
type JWKS = cryptoimpl.JWKS

// RSAPublicKeyToJWK exports an RSA public key as a JWK.
func RSAPublicKeyToJWK(pub *rsa.PublicKey, kid string) (JWK, error) {
	return cryptoimpl.RSAPublicKeyToJWK(pub, kid)
}

// RSAPrivateKeyToJWK exports an RSA private key as a private JWK.
func RSAPrivateKeyToJWK(priv *rsa.PrivateKey, kid string) (JWK, error) {
	return cryptoimpl.RSAPrivateKeyToJWK(priv, kid)
}

// JWKToRSAPublicKey parses an RSA public key from jwk.
func JWKToRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	return cryptoimpl.JWKToRSAPublicKey(jwk)
}

// JWKToRSAPrivateKey parses an RSA private key from jwk.
func JWKToRSAPrivateKey(jwk JWK) (*rsa.PrivateKey, error) {
	return cryptoimpl.JWKToRSAPrivateKey(jwk)
}

// MarshalJWK marshals jwk to compact JSON.
func MarshalJWK(jwk JWK) ([]byte, error) { return cryptoimpl.MarshalJWK(jwk) }

// ParseJWK parses a JWK JSON document.
func ParseJWK(data []byte) (JWK, error) { return cryptoimpl.ParseJWK(data) }

// MarshalJWKS marshals keys to a JWKS JSON document.
func MarshalJWKS(keys []JWK) ([]byte, error) { return cryptoimpl.MarshalJWKS(keys) }

// ParseJWKS parses a JWKS JSON document.
func ParseJWKS(data []byte) (JWKS, error) { return cryptoimpl.ParseJWKS(data) }

// SelectJWKByKeyID returns the key with kid from set.
func SelectJWKByKeyID(set JWKS, kid string) (JWK, error) {
	return cryptoimpl.SelectJWKByKeyID(set, kid)
}

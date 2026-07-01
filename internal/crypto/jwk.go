package crypto

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"strings"

	knifer "github.com/imajinyun/knifer-go"
)

// JWK represents a JSON Web Key. RSA keys are supported in this package.
type JWK struct {
	KeyType       string   `json:"kty"`
	KeyID         string   `json:"kid,omitempty"`
	Use           string   `json:"use,omitempty"`
	KeyOperations []string `json:"key_ops,omitempty"`
	Algorithm     string   `json:"alg,omitempty"`
	Modulus       string   `json:"n,omitempty"`
	Exponent      string   `json:"e,omitempty"`
	Private       string   `json:"d,omitempty"`
	PrimeP        string   `json:"p,omitempty"`
	PrimeQ        string   `json:"q,omitempty"`
	ExponentP     string   `json:"dp,omitempty"`
	ExponentQ     string   `json:"dq,omitempty"`
	Coefficient   string   `json:"qi,omitempty"`
}

// JWKS represents a JSON Web Key Set.
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// RSAPublicKeyToJWK exports an RSA public key as a JWK.
func RSAPublicKeyToJWK(pub *rsa.PublicKey, kid string) (JWK, error) {
	if pub == nil || pub.N == nil || pub.E <= 0 {
		return JWK{}, knifer.WrapError(knifer.ErrCodeInvalidInput, "rsa public key is invalid", ErrInvalidJWK)
	}
	return JWK{
		KeyType:  "RSA",
		KeyID:    strings.TrimSpace(kid),
		Modulus:  base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
		Exponent: base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
	}, nil
}

// RSAPrivateKeyToJWK exports an RSA private key as a private JWK.
func RSAPrivateKeyToJWK(priv *rsa.PrivateKey, kid string) (JWK, error) {
	if priv == nil {
		return JWK{}, knifer.WrapError(knifer.ErrCodeInvalidInput, "rsa private key is invalid", ErrInvalidJWK)
	}
	if err := priv.Validate(); err != nil {
		return JWK{}, knifer.WrapError(knifer.ErrCodeInvalidInput, "rsa private key validation failed", ErrInvalidJWK)
	}
	jwk, err := RSAPublicKeyToJWK(&priv.PublicKey, kid)
	if err != nil {
		return JWK{}, err
	}
	jwk.Private = base64.RawURLEncoding.EncodeToString(priv.D.Bytes())
	if len(priv.Primes) >= 2 {
		jwk.PrimeP = base64.RawURLEncoding.EncodeToString(priv.Primes[0].Bytes())
		jwk.PrimeQ = base64.RawURLEncoding.EncodeToString(priv.Primes[1].Bytes())
	}
	if priv.Precomputed.Dp != nil && priv.Precomputed.Dq != nil && priv.Precomputed.Qinv != nil {
		jwk.ExponentP = base64.RawURLEncoding.EncodeToString(priv.Precomputed.Dp.Bytes())
		jwk.ExponentQ = base64.RawURLEncoding.EncodeToString(priv.Precomputed.Dq.Bytes())
		jwk.Coefficient = base64.RawURLEncoding.EncodeToString(priv.Precomputed.Qinv.Bytes())
	}
	return jwk, nil
}

// JWKToRSAPublicKey parses an RSA public key from jwk.
func JWKToRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	if jwk.KeyType != "RSA" {
		return nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "jwk key type must be RSA", ErrInvalidJWK)
	}
	n, err := jwkBigInt(jwk.Modulus, "rsa modulus")
	if err != nil {
		return nil, err
	}
	e, err := jwkBigInt(jwk.Exponent, "rsa exponent")
	if err != nil {
		return nil, err
	}
	if !e.IsInt64() || e.Sign() <= 0 {
		return nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "rsa exponent is invalid", ErrInvalidJWK)
	}
	return &rsa.PublicKey{N: n, E: int(e.Int64())}, nil
}

// JWKToRSAPrivateKey parses an RSA private key from jwk.
func JWKToRSAPrivateKey(jwk JWK) (*rsa.PrivateKey, error) {
	pub, err := JWKToRSAPublicKey(jwk)
	if err != nil {
		return nil, err
	}
	d, err := jwkBigInt(jwk.Private, "rsa private exponent")
	if err != nil {
		return nil, err
	}
	p, err := jwkBigInt(jwk.PrimeP, "rsa prime p")
	if err != nil {
		return nil, err
	}
	q, err := jwkBigInt(jwk.PrimeQ, "rsa prime q")
	if err != nil {
		return nil, err
	}
	priv := &rsa.PrivateKey{
		PublicKey: *pub,
		D:         d,
		Primes:    []*big.Int{p, q},
	}
	if err := priv.Validate(); err != nil {
		return nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "rsa private jwk validation failed", ErrInvalidJWK)
	}
	priv.Precompute()
	return priv, nil
}

// MarshalJWK marshals jwk to compact JSON.
func MarshalJWK(jwk JWK) ([]byte, error) {
	if jwk.KeyType == "" {
		return nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "jwk key type is required", ErrInvalidJWK)
	}
	return json.Marshal(jwk)
}

// ParseJWK parses a JWK JSON document.
func ParseJWK(data []byte) (JWK, error) {
	var jwk JWK
	if err := json.Unmarshal(data, &jwk); err != nil {
		return JWK{}, knifer.WrapError(knifer.ErrCodeInvalidInput, "jwk json is malformed", ErrInvalidJWK)
	}
	if jwk.KeyType == "" {
		return JWK{}, knifer.WrapError(knifer.ErrCodeInvalidInput, "jwk key type is required", ErrInvalidJWK)
	}
	return jwk, nil
}

// MarshalJWKS marshals keys to a JWKS JSON document.
func MarshalJWKS(keys []JWK) ([]byte, error) {
	if len(keys) == 0 {
		return nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "jwks must contain at least one key", ErrInvalidJWK)
	}
	for _, key := range keys {
		if key.KeyType == "" {
			return nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "jwks contains key without key type", ErrInvalidJWK)
		}
	}
	return json.Marshal(JWKS{Keys: append([]JWK(nil), keys...)})
}

// ParseJWKS parses a JWKS JSON document.
func ParseJWKS(data []byte) (JWKS, error) {
	var set JWKS
	if err := json.Unmarshal(data, &set); err != nil {
		return JWKS{}, knifer.WrapError(knifer.ErrCodeInvalidInput, "jwks json is malformed", ErrInvalidJWK)
	}
	if len(set.Keys) == 0 {
		return JWKS{}, knifer.WrapError(knifer.ErrCodeInvalidInput, "jwks must contain at least one key", ErrInvalidJWK)
	}
	return set, nil
}

// SelectJWKByKeyID returns the key with kid from set.
func SelectJWKByKeyID(set JWKS, kid string) (JWK, error) {
	kid = strings.TrimSpace(kid)
	if kid == "" {
		return JWK{}, knifer.WrapError(knifer.ErrCodeInvalidInput, "jwk kid must be non-empty", ErrInvalidJWK)
	}
	for _, key := range set.Keys {
		if key.KeyID == kid {
			return key, nil
		}
	}
	return JWK{}, knifer.WrapError(knifer.ErrCodeNotFound, "jwk kid not found", ErrInvalidJWK)
}

func jwkBigInt(value, name string) (*big.Int, error) {
	if strings.TrimSpace(value) == "" {
		return nil, knifer.WrapError(knifer.ErrCodeInvalidInput, name+" is required", ErrInvalidJWK)
	}
	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil || len(decoded) == 0 {
		return nil, knifer.WrapError(knifer.ErrCodeInvalidInput, name+" is malformed", ErrInvalidJWK)
	}
	return new(big.Int).SetBytes(decoded), nil
}

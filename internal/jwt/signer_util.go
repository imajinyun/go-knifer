package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
)

// Convenience factory functions matching the utility toolkit-jwt JWTSignerUtil.
// They are provided as package-level functions in Go style.

// HS256 creates an HS256 signer.
func HS256(key []byte) JWTSigner { return MustHMACSigner(AlgHS256, key) }

// HS384 creates an HS384 signer.
func HS384(key []byte) JWTSigner { return MustHMACSigner(AlgHS384, key) }

// HS512 creates an HS512 signer.
func HS512(key []byte) JWTSigner { return MustHMACSigner(AlgHS512, key) }

// PS256 creates an RSA-PSS signer.
func PS256(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	return PS256WithOptions(priv, pub)
}

// PS256WithOptions creates a configurable RSA-PSS signer.
func PS256WithOptions(priv *rsa.PrivateKey, pub *rsa.PublicKey, opts ...SignerOption) JWTSigner {
	return mustRSAPSSWithOptions(AlgPS256, priv, pub, opts...)
}

// PS384 same as above.
func PS384(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	return PS384WithOptions(priv, pub)
}

// PS384WithOptions creates a configurable RSA-PSS signer.
func PS384WithOptions(priv *rsa.PrivateKey, pub *rsa.PublicKey, opts ...SignerOption) JWTSigner {
	return mustRSAPSSWithOptions(AlgPS384, priv, pub, opts...)
}

// PS512 same as above.
func PS512(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	return PS512WithOptions(priv, pub)
}

// PS512WithOptions creates a configurable RSA-PSS signer.
func PS512WithOptions(priv *rsa.PrivateKey, pub *rsa.PublicKey, opts ...SignerOption) JWTSigner {
	return mustRSAPSSWithOptions(AlgPS512, priv, pub, opts...)
}

// ES256 creates an ECDSA(P-256) signer.
func ES256(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) JWTSigner {
	return ES256WithOptions(priv, pub)
}

// ES256WithOptions creates a configurable ECDSA(P-256) signer.
func ES256WithOptions(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey, opts ...SignerOption) JWTSigner {
	return mustECDSAWithOptions(AlgES256, priv, pub, opts...)
}

// ES384 creates an ECDSA(P-384) signer.
func ES384(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) JWTSigner {
	return ES384WithOptions(priv, pub)
}

// ES384WithOptions creates a configurable ECDSA(P-384) signer.
func ES384WithOptions(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey, opts ...SignerOption) JWTSigner {
	return mustECDSAWithOptions(AlgES384, priv, pub, opts...)
}

// ES512 creates an ECDSA(P-521) signer.
func ES512(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) JWTSigner {
	return ES512WithOptions(priv, pub)
}

// ES512WithOptions creates a configurable ECDSA(P-521) signer.
func ES512WithOptions(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey, opts ...SignerOption) JWTSigner {
	return mustECDSAWithOptions(AlgES512, priv, pub, opts...)
}

func mustRSAPSSWithOptions(alg string, priv *rsa.PrivateKey, pub *rsa.PublicKey, opts ...SignerOption) JWTSigner {
	s, err := NewRSAPSSSignerWithOptions(alg, priv, pub, opts...)
	if err != nil {
		panic(err)
	}
	return s
}

func mustECDSAWithOptions(alg string, priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey, opts ...SignerOption) JWTSigner {
	s, err := NewECDSASignerWithOptions(alg, priv, pub, opts...)
	if err != nil {
		panic(err)
	}
	return s
}

package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
)

// 对应 the utility toolkit-jwt JWTSignerUtil 的便捷工厂函数。
// 在 Go 风格中以包级函数提供。

// HS256 创建 HS256 签名器。
func HS256(key []byte) JWTSigner { return MustHMACSigner(AlgHS256, key) }

// HS384 创建 HS384 签名器。
func HS384(key []byte) JWTSigner { return MustHMACSigner(AlgHS384, key) }

// HS512 创建 HS512 签名器。
func HS512(key []byte) JWTSigner { return MustHMACSigner(AlgHS512, key) }

// RS256 创建 RS256 签名器；priv 用于签名、pub 用于验签，传 nil 表示该方向不可用。
func RS256(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner { return mustRSA(AlgRS256, priv, pub) }

// RS384 同上。
func RS384(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner { return mustRSA(AlgRS384, priv, pub) }

// RS512 同上。
func RS512(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner { return mustRSA(AlgRS512, priv, pub) }

// PS256 创建 RSA-PSS 签名器。
func PS256(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	return mustRSAPSS(AlgPS256, priv, pub)
}

// PS384 同上。
func PS384(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	return mustRSAPSS(AlgPS384, priv, pub)
}

// PS512 同上。
func PS512(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	return mustRSAPSS(AlgPS512, priv, pub)
}

// ES256 创建 ECDSA(P-256) 签名器。
func ES256(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) JWTSigner {
	return mustECDSA(AlgES256, priv, pub)
}

// ES384 创建 ECDSA(P-384) 签名器。
func ES384(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) JWTSigner {
	return mustECDSA(AlgES384, priv, pub)
}

// ES512 创建 ECDSA(P-521) 签名器。
func ES512(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) JWTSigner {
	return mustECDSA(AlgES512, priv, pub)
}

// None 返回无签名签名器（the utility toolkit JWTSignerUtil.none）。
func None() JWTSigner { return NoneSigner() }

func mustRSA(alg string, priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	s, err := NewRSASigner(alg, priv, pub)
	if err != nil {
		panic(err)
	}
	return s
}

func mustRSAPSS(alg string, priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	s, err := NewRSAPSSSigner(alg, priv, pub)
	if err != nil {
		panic(err)
	}
	return s
}

func mustECDSA(alg string, priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) JWTSigner {
	s, err := NewECDSASigner(alg, priv, pub)
	if err != nil {
		panic(err)
	}
	return s
}

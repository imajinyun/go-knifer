package vjwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"time"

	jwtimpl "github.com/imajinyun/go-knifer/internal/jwt"
)

// NewJWTError delegates to the internal jwt implementation.
func NewJWTError(msg string) *JWTError {
	return jwtimpl.NewJWTError(msg)
}

// JWTErrorf delegates to the internal jwt implementation.
func JWTErrorf(format string, args ...any) *JWTError {
	return jwtimpl.JWTErrorf(format, args...)
}

// New delegates to the internal jwt implementation.
func New() *JWT {
	return jwtimpl.New()
}

// NoneSigner delegates to the internal jwt implementation.
func NoneSigner() JWTSigner {
	return jwtimpl.NoneSigner()
}

// IsNoneAlg delegates to the internal jwt implementation.
func IsNoneAlg(alg string) bool {
	return jwtimpl.IsNoneAlg(alg)
}

// NewHMACSigner delegates to the internal jwt implementation.
func NewHMACSigner(algorithm string, key []byte) (JWTSigner, error) {
	return jwtimpl.NewHMACSigner(algorithm, key)
}

// MustHMACSigner delegates to the internal jwt implementation.
func MustHMACSigner(algorithm string, key []byte) JWTSigner {
	return jwtimpl.MustHMACSigner(algorithm, key)
}

// CreateSigner delegates to the internal jwt implementation.
func CreateSigner(algorithmID string, key []byte) (JWTSigner, error) {
	return jwtimpl.CreateSigner(algorithmID, key)
}

// AlgorithmName delegates to the internal jwt implementation.
func AlgorithmName(idOrAlgorithm string) string {
	return jwtimpl.AlgorithmName(idOrAlgorithm)
}

// NewRSASigner delegates to the internal jwt implementation.
func NewRSASigner(algorithm string, priv *rsa.PrivateKey, pub *rsa.PublicKey) (JWTSigner, error) {
	return jwtimpl.NewRSASigner(algorithm, priv, pub)
}

// NewRSAPSSSigner delegates to the internal jwt implementation.
func NewRSAPSSSigner(algorithm string, priv *rsa.PrivateKey, pub *rsa.PublicKey) (JWTSigner, error) {
	return jwtimpl.NewRSAPSSSigner(algorithm, priv, pub)
}

// NewECDSASigner delegates to the internal jwt implementation.
func NewECDSASigner(algorithm string, priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) (JWTSigner, error) {
	return jwtimpl.NewECDSASigner(algorithm, priv, pub)
}

// HS256 delegates to the internal jwt implementation.
func HS256(key []byte) JWTSigner {
	return jwtimpl.HS256(key)
}

// HS384 delegates to the internal jwt implementation.
func HS384(key []byte) JWTSigner {
	return jwtimpl.HS384(key)
}

// HS512 delegates to the internal jwt implementation.
func HS512(key []byte) JWTSigner {
	return jwtimpl.HS512(key)
}

// RS256 delegates to the internal jwt implementation.
func RS256(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	return jwtimpl.RS256(priv, pub)
}

// RS384 delegates to the internal jwt implementation.
func RS384(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	return jwtimpl.RS384(priv, pub)
}

// RS512 delegates to the internal jwt implementation.
func RS512(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	return jwtimpl.RS512(priv, pub)
}

// PS256 delegates to the internal jwt implementation.
func PS256(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	return jwtimpl.PS256(priv, pub)
}

// PS384 delegates to the internal jwt implementation.
func PS384(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	return jwtimpl.PS384(priv, pub)
}

// PS512 delegates to the internal jwt implementation.
func PS512(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	return jwtimpl.PS512(priv, pub)
}

// ES256 delegates to the internal jwt implementation.
func ES256(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) JWTSigner {
	return jwtimpl.ES256(priv, pub)
}

// ES384 delegates to the internal jwt implementation.
func ES384(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) JWTSigner {
	return jwtimpl.ES384(priv, pub)
}

// ES512 delegates to the internal jwt implementation.
func ES512(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) JWTSigner {
	return jwtimpl.ES512(priv, pub)
}

// None delegates to the internal jwt implementation.
func None() JWTSigner {
	return jwtimpl.None()
}

// CreateToken delegates to the internal jwt implementation.
func CreateToken(payload map[string]any, key []byte) (string, error) {
	return jwtimpl.CreateToken(payload, key)
}

// CreateTokenWithHeaders delegates to the internal jwt implementation.
func CreateTokenWithHeaders(headers, payload map[string]any, key []byte) (string, error) {
	return jwtimpl.CreateTokenWithHeaders(headers, payload, key)
}

// CreateTokenWithSigner delegates to the internal jwt implementation.
func CreateTokenWithSigner(payload map[string]any, signer JWTSigner) (string, error) {
	return jwtimpl.CreateTokenWithSigner(payload, signer)
}

// CreateTokenWithHeadersAndSigner delegates to the internal jwt implementation.
func CreateTokenWithHeadersAndSigner(headers, payload map[string]any, signer JWTSigner) (string, error) {
	return jwtimpl.CreateTokenWithHeadersAndSigner(headers, payload, signer)
}

// ParseToken delegates to the internal jwt implementation.
func ParseToken(token string) (*JWT, error) {
	return jwtimpl.ParseToken(token)
}

// Verify delegates to the internal jwt implementation.
func Verify(token string, key []byte) bool {
	return jwtimpl.Verify(token, key)
}

// VerifyWithSigner delegates to the internal jwt implementation.
func VerifyWithSigner(token string, signer JWTSigner) bool {
	return jwtimpl.VerifyWithSigner(token, signer)
}

// OfValidator delegates to the internal jwt implementation.
func OfValidator(token string) *JWTValidator {
	return jwtimpl.OfValidator(token)
}

// OfValidatorJWT delegates to the internal jwt implementation.
func OfValidatorJWT(j *JWT) *JWTValidator {
	return jwtimpl.OfValidatorJWT(j)
}

// ValidateAlgorithm delegates to the internal jwt implementation.
func ValidateAlgorithm(token string, signer JWTSigner) error {
	return jwtimpl.ValidateAlgorithm(token, signer)
}

// ValidateDate delegates to the internal jwt implementation.
func ValidateDate(j *JWT, now time.Time, leeway int64) error {
	return jwtimpl.ValidateDate(j, now, leeway)
}

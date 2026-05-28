package vjwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"time"

	jwtimpl "github.com/imajinyun/go-knifer/internal/jwt"
)

// JWT is a JSON Web Token object.
type JWT = jwtimpl.JWT

// JWTSigner signs and verifies JWT tokens.
type JWTSigner = jwtimpl.JWTSigner

// JWTValidator validates JWT claims.
type JWTValidator = jwtimpl.JWTValidator

// JWTError is the JWT module error type.
type JWTError = jwtimpl.JWTError

const (
	// JWTAlgNone is the none algorithm identifier.
	JWTAlgNone = jwtimpl.AlgNone
	// JWTAlgHS256 is the HS256 algorithm identifier.
	JWTAlgHS256 = jwtimpl.AlgHS256
	// JWTAlgRS256 is the RS256 algorithm identifier.
	JWTAlgRS256 = jwtimpl.AlgRS256
	// JWTAlgES256 is the ES256 algorithm identifier.
	JWTAlgES256 = jwtimpl.AlgES256
	// JWTHeaderAlgorithm is the alg header key.
	JWTHeaderAlgorithm = jwtimpl.HeaderAlgorithm
	// JWTPayloadIssuer is the iss payload key.
	JWTPayloadIssuer = jwtimpl.PayloadIssuer
	// JWTPayloadSubject is the sub payload key.
	JWTPayloadSubject = jwtimpl.PayloadSubject
	// JWTPayloadExpiresAt is the exp payload key.
	JWTPayloadExpiresAt = jwtimpl.PayloadExpiresAt
)

// NewJWT creates a new JWT object.
func NewJWT() *JWT { return jwtimpl.New() }

// ParseJWT parses a token string.
func ParseJWT(token string) (*JWT, error) { return jwtimpl.ParseToken(token) }

// JWTOf parses a token string.
func JWTOf(token string) (*JWT, error) { return jwtimpl.Of(token) }

// CreateJWTToken creates a signed token using HMAC key.
func CreateJWTToken(payload map[string]any, key []byte) (string, error) {
	return jwtimpl.CreateToken(payload, key)
}

// CreateJWTTokenWithSigner creates a signed token using signer.
func CreateJWTTokenWithSigner(payload map[string]any, signer JWTSigner) (string, error) {
	return jwtimpl.CreateTokenWithSigner(payload, signer)
}

// VerifyJWT verifies a token using HMAC key.
func VerifyJWT(token string, key []byte) bool { return jwtimpl.Verify(token, key) }

// VerifyJWTWithSigner verifies a token using signer.
func VerifyJWTWithSigner(token string, signer JWTSigner) bool {
	return jwtimpl.VerifyWithSigner(token, signer)
}

// ValidateJWTDate validates time based JWT claims.
func ValidateJWTDate(j *JWT, now time.Time, leeway int64) error {
	return jwtimpl.ValidateDate(j, now, leeway)
}

// JWTSignerHMAC creates an HMAC signer.
func JWTSignerHMAC(algorithm string, key []byte) (JWTSigner, error) {
	return jwtimpl.NewHMACSigner(algorithm, key)
}

// JWTSignerHS256 creates an HS256 signer.
func JWTSignerHS256(key []byte) JWTSigner { return jwtimpl.HS256(key) }

// JWTSignerRSA creates an RSA signer.
func JWTSignerRSA(algorithm string, priv *rsa.PrivateKey, pub *rsa.PublicKey) (JWTSigner, error) {
	return jwtimpl.NewRSASigner(algorithm, priv, pub)
}

// JWTSignerRS256 creates an RS256 signer.
func JWTSignerRS256(priv *rsa.PrivateKey, pub *rsa.PublicKey) JWTSigner {
	return jwtimpl.RS256(priv, pub)
}

// JWTSignerECDSA creates an ECDSA signer.
func JWTSignerECDSA(algorithm string, priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) (JWTSigner, error) {
	return jwtimpl.NewECDSASigner(algorithm, priv, pub)
}

// JWTSignerES256 creates an ES256 signer.
func JWTSignerES256(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) JWTSigner {
	return jwtimpl.ES256(priv, pub)
}

// JWTSignerNone creates a none signer.
func JWTSignerNone() JWTSigner { return jwtimpl.None() }

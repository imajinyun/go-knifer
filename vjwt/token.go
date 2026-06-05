package vjwt

import jwtimpl "github.com/imajinyun/go-knifer/internal/jwt"

// NewJWT creates a new JWT object.
func NewJWT() *JWT { return jwtimpl.New() }

// ParseJWT parses a token string.
func ParseJWT(token string) (*JWT, error) { return jwtimpl.ParseToken(token) }

// JWTOf parses a token string.
func JWTOf(token string) (*JWT, error) { return jwtimpl.Of(token) }

// New delegates to the internal jwt implementation.
func New() *JWT {
	return jwtimpl.New()
}

// CreateJWTToken creates a signed token using HMAC key.
func CreateJWTToken(payload map[string]any, key []byte) (string, error) {
	return jwtimpl.CreateToken(payload, key)
}

// CreateJWTTokenWithSigner creates a signed token using signer.
func CreateJWTTokenWithSigner(payload map[string]any, signer JWTSigner) (string, error) {
	return jwtimpl.CreateTokenWithSigner(payload, signer)
}

// CreateToken delegates to the internal jwt implementation.
func CreateToken(payload map[string]any, key []byte) (string, error) {
	return jwtimpl.CreateToken(payload, key)
}

// CreateTokenWithHeaders delegates to the internal jwt implementation.
func CreateTokenWithHeaders(headers, payload map[string]any, key []byte) (string, error) {
	return jwtimpl.CreateTokenWithHeaders(headers, payload, key)
}

// CreateTokenWithAlgorithm creates a token with an explicit HMAC algorithm.
func CreateTokenWithAlgorithm(payload map[string]any, key []byte, algorithm string) (string, error) {
	return jwtimpl.CreateTokenWithAlgorithm(payload, key, algorithm)
}

// CreateTokenWithHeadersAndAlgorithm creates a token with headers and an explicit HMAC algorithm.
func CreateTokenWithHeadersAndAlgorithm(headers, payload map[string]any, key []byte, algorithm string) (string, error) {
	return jwtimpl.CreateTokenWithHeadersAndAlgorithm(headers, payload, key, algorithm)
}

// CreateTokenWithSigner delegates to the internal jwt implementation.
func CreateTokenWithSigner(payload map[string]any, signer JWTSigner) (string, error) {
	return jwtimpl.CreateTokenWithSigner(payload, signer)
}

// CreateTokenWithHeadersAndSigner delegates to the internal jwt implementation.
func CreateTokenWithHeadersAndSigner(headers, payload map[string]any, signer JWTSigner) (string, error) {
	return jwtimpl.CreateTokenWithHeadersAndSigner(headers, payload, signer)
}

// WithTokenHeaders sets JWT header fields for CreateTokenWithOptions.
func WithTokenHeaders(headers map[string]any) TokenOption { return jwtimpl.WithTokenHeaders(headers) }

// WithTokenPayload sets JWT payload fields for CreateTokenWithOptions.
func WithTokenPayload(payload map[string]any) TokenOption { return jwtimpl.WithTokenPayload(payload) }

// WithTokenKey sets the HMAC key used by CreateTokenWithOptions.
func WithTokenKey(key []byte) TokenOption { return jwtimpl.WithTokenKey(key) }

// WithTokenAlgorithm sets the HMAC algorithm used by CreateTokenWithOptions.
func WithTokenAlgorithm(algorithm string) TokenOption { return jwtimpl.WithTokenAlgorithm(algorithm) }

// WithTokenSigner sets the signer used by CreateTokenWithOptions and takes precedence over key/algorithm options.
func WithTokenSigner(signer JWTSigner) TokenOption { return jwtimpl.WithTokenSigner(signer) }

// CreateTokenWithOptions creates a token from functional options and avoids adding more overload variants.
func CreateTokenWithOptions(opts ...TokenOption) (string, error) {
	return jwtimpl.CreateTokenWithOptions(opts...)
}

// ParseToken delegates to the internal jwt implementation.
func ParseToken(token string) (*JWT, error) {
	return jwtimpl.ParseToken(token)
}

// OfValidator delegates to the internal jwt implementation.
func OfValidator(token string) *JWTValidator {
	return jwtimpl.OfValidator(token)
}

// OfValidatorJWT delegates to the internal jwt implementation.
func OfValidatorJWT(j *JWT) *JWTValidator {
	return jwtimpl.OfValidatorJWT(j)
}

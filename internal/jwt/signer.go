package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"hash"
	"strings"
)

// JWTSigner is the JWT signer interface matching the utility toolkit-jwt JWTSigner.
type JWTSigner interface {
	// Algorithm returns the algorithm ID, such as HS256.
	Algorithm() string
	// Sign signs headerBase64.payloadBase64 and returns an unpadded base64url string.
	Sign(headerB64, payloadB64 string) string
	// Verify verifies whether the signature matches.
	Verify(headerB64, payloadB64, signB64 string) bool
}

// Algorithm ID constants.
const (
	AlgNone  = "none"
	AlgHS256 = "HS256"
	AlgHS384 = "HS384"
	AlgHS512 = "HS512"
)

// b64URLEncode encodes base64url without padding as used by standard JWT.
func b64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// b64URLDecode decodes base64url and accepts padded input.
func b64URLDecode(s string) ([]byte, error) {
	// Prefer unpadded decoding and try padded decoding on failure.
	if b, err := base64.RawURLEncoding.DecodeString(s); err == nil {
		return b, nil
	}
	return base64.URLEncoding.DecodeString(s)
}

func isNoneAlg(alg string) bool { return strings.EqualFold(strings.TrimSpace(alg), AlgNone) }

// hmacSigner HMAC signer family.
type hmacSigner struct {
	alg    string
	key    []byte
	hashFn func() hash.Hash
}

// NewHMACSigner creates an HMAC signer. algorithm only supports HS256, HS384, and HS512.
func NewHMACSigner(algorithm string, key []byte) (JWTSigner, error) {
	algorithm = strings.ToUpper(strings.TrimSpace(algorithm))
	switch algorithm {
	case AlgHS256:
		return &hmacSigner{alg: AlgHS256, key: append([]byte{}, key...), hashFn: sha256.New}, nil
	case AlgHS384:
		return &hmacSigner{alg: AlgHS384, key: append([]byte{}, key...), hashFn: sha512.New384}, nil
	case AlgHS512:
		return &hmacSigner{alg: AlgHS512, key: append([]byte{}, key...), hashFn: sha512.New}, nil
	}
	// Accept SHA384 as a direct hash alias for compatibility.
	return nil, unsupportedJWTErrorf("unsupported HMAC algorithm: %s", algorithm)
}

// MustHMACSigner creates an HMAC signer and panics on failure.
func MustHMACSigner(algorithm string, key []byte) JWTSigner {
	s, err := NewHMACSigner(algorithm, key)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *hmacSigner) Algorithm() string { return s.alg }

func (s *hmacSigner) Sign(headerB64, payloadB64 string) string {
	mac := hmac.New(s.hashFn, s.key)
	mac.Write([]byte(headerB64))
	mac.Write([]byte{'.'})
	mac.Write([]byte(payloadB64))
	return b64URLEncode(mac.Sum(nil))
}

func (s *hmacSigner) Verify(headerB64, payloadB64, signB64 string) bool {
	expected := s.Sign(headerB64, payloadB64)
	return subtle.ConstantTimeCompare([]byte(expected), []byte(signB64)) == 1
}

// CreateSigner selects a signer from the algorithm ID and HMAC key; only HS* is supported.
// The none algorithm is always rejected.
//
// Use NewRSAPSSSigner or NewECDSASigner for asymmetric algorithms,
// or use convenience factories such as PS256 and ES256 provided by JWTSignerUtil.
func CreateSigner(algorithmID string, key []byte) (JWTSigner, error) {
	if isNoneAlg(algorithmID) {
		return nil, unsupportedJWTErrorf("jwt alg=none is not supported")
	}
	return NewHMACSigner(algorithmID, key)
}

// AlgorithmName returns the standard algorithm name for a JWT algorithm ID, matching the utility toolkit AlgorithmUtil.getAlgorithm.
// returns unknown IDs unchanged.
func AlgorithmName(idOrAlgorithm string) string {
	id := strings.ToUpper(strings.TrimSpace(idOrAlgorithm))
	switch id {
	case AlgHS256:
		return "HmacSHA256"
	case AlgHS384:
		return "HmacSHA384"
	case AlgHS512:
		return "HmacSHA512"
	case AlgPS256:
		return "SHA256withRSA_PSS"
	case AlgPS384:
		return "SHA384withRSA_PSS"
	case AlgPS512:
		return "SHA512withRSA_PSS"
	case AlgES256:
		return "SHA256withECDSA"
	case AlgES384:
		return "SHA384withECDSA"
	case AlgES512:
		return "SHA512withECDSA"
	}
	return idOrAlgorithm
}

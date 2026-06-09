package jwt

// Registered standard Header field names matching the utility toolkit-jwt JWTHeader.
const (
	HeaderAlgorithm   = "alg" // signing algorithm, such as HS256
	HeaderType        = "typ" // type, usually JWT
	HeaderContentType = "cty" // content type
	HeaderKeyID       = "kid" // key ID
)

// Registered standard Payload field names matching the utility toolkit-jwt RegisteredPayload.
const (
	PayloadIssuer    = "iss" // issuer
	PayloadSubject   = "sub" // subject
	PayloadAudience  = "aud" // audience
	PayloadExpiresAt = "exp" // expiration time
	PayloadNotBefore = "nbf" // not-before time
	PayloadIssuedAt  = "iat" // issued-at time
	PayloadJWTID     = "jti" // unique JWT ID
)

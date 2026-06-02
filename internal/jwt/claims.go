package jwt

// 注册的标准 Header 字段名（对应 the utility toolkit-jwt JWTHeader）。
const (
	HeaderAlgorithm   = "alg" // 加密算法（如 HS256）
	HeaderType        = "typ" // 类型（一般为 JWT）
	HeaderContentType = "cty" // 内容类型
	HeaderKeyID       = "kid" // 密钥编号
)

// 注册的标准 Payload 字段名（对应 the utility toolkit-jwt RegisteredPayload）。
const (
	PayloadIssuer    = "iss" // 签发者
	PayloadSubject   = "sub" // 面向用户
	PayloadAudience  = "aud" // 接收方
	PayloadExpiresAt = "exp" // 过期时间
	PayloadNotBefore = "nbf" // 生效时间
	PayloadIssuedAt  = "iat" // 签发时间
	PayloadJWTID     = "jti" // 唯一身份标识
)

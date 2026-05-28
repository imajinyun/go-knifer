package jwt

// 包级便捷函数（对应 hutool-jwt JWTUtil）。

// CreateToken 用 HS256 创建 token。
func CreateToken(payload map[string]any, key []byte) (string, error) {
	return CreateTokenWithHeaders(nil, payload, key)
}

// CreateTokenWithHeaders 创建带 header 的 token（HS256）。
func CreateTokenWithHeaders(headers, payload map[string]any, key []byte) (string, error) {
	j := New().AddHeaders(headers).AddPayloads(payload).SetKey(key)
	return j.Sign()
}

// CreateTokenWithSigner 使用自定义签名器创建 token。
func CreateTokenWithSigner(payload map[string]any, signer JWTSigner) (string, error) {
	return CreateTokenWithHeadersAndSigner(nil, payload, signer)
}

// CreateTokenWithHeadersAndSigner 使用自定义签名器与 header 创建 token。
func CreateTokenWithHeadersAndSigner(headers, payload map[string]any, signer JWTSigner) (string, error) {
	j := New().AddHeaders(headers).AddPayloads(payload).SetSigner(signer)
	return j.Sign()
}

// ParseToken 解析 token。
func ParseToken(token string) (*JWT, error) { return Of(token) }

// Verify 使用 HS256 密钥校验 token。
func Verify(token string, key []byte) bool {
	j, err := Of(token)
	if err != nil {
		return false
	}
	return j.SetKey(key).Verify()
}

// VerifyWithSigner 使用自定义 signer 校验 token。
func VerifyWithSigner(token string, signer JWTSigner) bool {
	j, err := Of(token)
	if err != nil {
		return false
	}
	return j.VerifyWith(signer)
}

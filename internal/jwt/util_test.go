package jwt

import "testing"

// Simplified utility toolkit-jwt JWTUtilTest.

func TestUtil_CreateAndVerify(t *testing.T) {
	key := []byte("1234567890")
	payload := map[string]any{
		"sub":  "1234567890",
		"name": "looly",
	}
	tok, err := CreateToken(payload, key)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !Verify(tok, key) {
		t.Fatalf("verify failed")
	}
	if Verify(tok, []byte("wrong")) {
		t.Fatalf("verify should fail with wrong key")
	}
}

func TestUtil_CreateWithSigner(t *testing.T) {
	signer := MustHMACSigner(AlgHS512, []byte("secret"))
	tok, err := CreateTokenWithSigner(map[string]any{"a": 1}, signer)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !VerifyWithSigner(tok, signer) {
		t.Fatalf("verify failed")
	}
}

func TestUtil_CreateAndVerifyStrictWithAlgorithm(t *testing.T) {
	key := []byte("secret")
	tok, err := CreateTokenWithAlgorithm(map[string]any{"a": 1}, key, AlgHS512)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	parsed, err := ParseToken(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if parsed.Algorithm() != AlgHS512 {
		t.Fatalf("alg = %q, want %q", parsed.Algorithm(), AlgHS512)
	}
	if !VerifyStrict(tok, key) {
		t.Fatal("VerifyStrict failed")
	}
	if !Verify(tok, key) {
		t.Fatal("Verify should use header algorithm without fallback")
	}
	if _, err := CreateTokenWithAlgorithm(map[string]any{"a": 1}, key, "bad"); err == nil {
		t.Fatal("CreateTokenWithAlgorithm bad alg error = nil")
	}
}

func TestUtil_VerifyRejectsUnsupportedHeaderAlgorithm(t *testing.T) {
	key := []byte("secret")
	tok, err := New().SetHeader(HeaderAlgorithm, "BAD").SetPayload("a", 1).SetSigner(HS256(key)).Sign()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if Verify(tok, key) {
		t.Fatal("Verify should reject unsupported header alg instead of falling back")
	}
}

func TestUtil_ParseToken(t *testing.T) {
	tok := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9" +
		".eyJsb2dpblR5cGUiOiJsb2dpbiIsImxvZ2luSWQiOiJhZG1pbiIsImRldmljZSI6ImRlZmF1bHQtZGV2aWNlIiwiZWZmIjoxNjc4Mjg1NzEzOTM1LCJyblN0ciI6IkVuMTczWFhvWUNaaVZUWFNGOTNsN1pabGtOalNTd0pmIn0" +
		".wRe2soTaWYPhwcjxdzesDi1BgEm9D61K-mMT3fPc4YM"
	j, err := ParseToken(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	p := j.Payloads()
	if p["loginType"] != "login" {
		t.Fatalf("loginType: %v", p["loginType"])
	}
	// JSON numbers parse as float64 by default.
	if v, ok := p["eff"].(float64); !ok || int64(v) != 1678285713935 {
		t.Fatalf("eff: %v (%T)", p["eff"], p["eff"])
	}
}

func TestUtil_CreateTokenWithHeaders(t *testing.T) {
	headers := map[string]any{HeaderKeyID: "kid-1"}
	payload := map[string]any{"a": 1}
	tok, err := CreateTokenWithHeaders(headers, payload, []byte("k"))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	j, err := ParseToken(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if j.Header(HeaderKeyID) != "kid-1" {
		t.Fatalf("kid: %v", j.Header(HeaderKeyID))
	}
}

func TestUtil_ParseInvalid(t *testing.T) {
	if _, err := ParseToken(""); err == nil {
		t.Fatalf("expected error for blank token")
	}
	if _, err := ParseToken("not.a.jwt.too.many"); err == nil {
		t.Fatalf("expected error for malformed token")
	}
	if Verify("bad", []byte("k")) {
		t.Fatalf("expected verify=false for bad token")
	}
}

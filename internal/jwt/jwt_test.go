package jwt

import (
	"strings"
	"testing"
	"time"
)

// 对应 the utility toolkit-jwt JWTTest。

func TestCreateHS256(t *testing.T) {
	key := []byte("1234567890")
	j := New().
		SetPayload("sub", "1234567890").
		SetPayload("name", "looly").
		SetPayload("admin", true).
		SetExpiresAt(time.Unix(1640966400, 0)).
		SetKey(key)

	tok, err := j.Sign()
	if err != nil {
		t.Fatalf("sign err: %v", err)
	}
	parts := strings.Split(tok, ".")
	if len(parts) != 3 {
		t.Fatalf("token parts: %d", len(parts))
	}
	// 解析回来后能验证通过即可
	parsed, err := Of(tok)
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	if !parsed.SetKey(key).Verify() {
		t.Fatalf("verify failed")
	}
	if parsed.Payload("name") != "looly" {
		t.Fatalf("payload name: %v", parsed.Payload("name"))
	}
	if parsed.Algorithm() != AlgHS256 {
		t.Fatalf("alg: %s", parsed.Algorithm())
	}
	if parsed.Type() != "JWT" {
		t.Fatalf("typ: %s", parsed.Type())
	}
}

func TestParseAndVerifyKnownToken(t *testing.T) {
	// 来自 the utility toolkit 的固定测试 token
	rightToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9." +
		"eyJzdWIiOiIxMjM0NTY3ODkwIiwiYWRtaW4iOnRydWUsIm5hbWUiOiJsb29seSJ9." +
		"U2aQkC2THYV9L0fTN-yBBI7gmo5xhmvMhATtu8v0zEA"

	j, err := Of(rightToken)
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	if !j.SetKey([]byte("1234567890")).Verify() {
		t.Fatalf("verify failed")
	}
	if j.Header(HeaderType) != "JWT" {
		t.Fatalf("type: %v", j.Header(HeaderType))
	}
	if j.Header(HeaderAlgorithm) != "HS256" {
		t.Fatalf("alg: %v", j.Header(HeaderAlgorithm))
	}
	if j.Header(HeaderContentType) != nil {
		t.Fatalf("cty should be nil")
	}
	if j.Payload("sub") != "1234567890" {
		t.Fatalf("sub: %v", j.Payload("sub"))
	}
	if j.Payload("name") != "looly" {
		t.Fatalf("name: %v", j.Payload("name"))
	}
	if j.Payload("admin") != true {
		t.Fatalf("admin: %v", j.Payload("admin"))
	}
}

func TestCreateNone(t *testing.T) {
	j := New().
		SetPayload("sub", "1234567890").
		SetPayload("name", "looly").
		SetPayload("admin", true).
		SetSigner(NoneSigner())

	tok, err := j.Sign()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	parts := strings.Split(tok, ".")
	if len(parts) != 3 || parts[2] != "" {
		t.Fatalf("none signature should be empty: %q", tok)
	}
	parsed, err := Of(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if !parsed.SetSigner(NoneSigner()).Verify() {
		t.Fatalf("verify failed for none")
	}
}

func TestNeedSigner(t *testing.T) {
	j := New().SetPayload("sub", "x")
	if _, err := j.Sign(); err == nil {
		t.Fatalf("expected error when no signer set")
	}
}

func TestVerifyMismatchKey(t *testing.T) {
	tok, err := New().SetPayload("a", 1).SetKey([]byte("right")).Sign()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	j, _ := Of(tok)
	if j.SetKey([]byte("wrong")).Verify() {
		t.Fatalf("should fail with wrong key")
	}
}

func TestAlgMismatch(t *testing.T) {
	// alg=none 时使用非 None signer 应失败
	tok, _ := New().SetSigner(NoneSigner()).SetPayload("a", 1).Sign()
	j, _ := Of(tok)
	hs, _ := NewHMACSigner(AlgHS256, []byte("k"))
	if j.VerifyWith(hs) {
		t.Fatalf("none alg with HS256 signer should fail")
	}
	// alg=HS256 时使用 None signer 应失败
	tok2, _ := New().SetKey([]byte("k")).SetPayload("a", 1).Sign()
	j2, _ := Of(tok2)
	if j2.VerifyWith(NoneSigner()) {
		t.Fatalf("HS256 token with None signer should fail")
	}
}

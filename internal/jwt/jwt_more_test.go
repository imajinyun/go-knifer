package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"strings"
	"testing"
	"time"
)

func TestJWTSetIssuer(t *testing.T) {
	j := New()
	result := j.SetIssuer("my-issuer")
	if result != j {
		t.Error("SetIssuer should return the JWT for chaining")
	}
	if j.Payload(PayloadIssuer) != "my-issuer" {
		t.Errorf("SetIssuer payload = %v, want %q", j.Payload(PayloadIssuer), "my-issuer")
	}
}

func TestJWTSetSubject(t *testing.T) {
	j := New()
	j.SetSubject("my-subject")
	if j.Payload(PayloadSubject) != "my-subject" {
		t.Errorf("SetSubject payload = %v, want %q", j.Payload(PayloadSubject), "my-subject")
	}
}

func TestJWTSetAudience(t *testing.T) {
	t.Run("single audience", func(t *testing.T) {
		j := New()
		j.SetAudience("aud1")
		if j.Payload(PayloadAudience) != "aud1" {
			t.Errorf("SetAudience single = %v, want %q", j.Payload(PayloadAudience), "aud1")
		}
	})
	t.Run("multiple audiences", func(t *testing.T) {
		j := New()
		j.SetAudience("aud1", "aud2")
		aud, ok := j.Payload(PayloadAudience).([]string)
		if !ok {
			t.Fatalf("SetAudience multiple should produce []string, got %T", j.Payload(PayloadAudience))
		}
		if len(aud) != 2 || aud[0] != "aud1" || aud[1] != "aud2" {
			t.Errorf("SetAudience multiple = %v, want [aud1 aud2]", aud)
		}
	})
}

func TestJWTSetJWTID(t *testing.T) {
	j := New()
	j.SetJWTID("my-jti")
	if j.Payload(PayloadJWTID) != "my-jti" {
		t.Errorf("SetJWTID payload = %v, want %q", j.Payload(PayloadJWTID), "my-jti")
	}
}

func TestJWTHeaders(t *testing.T) {
	j := New()
	j.SetHeader("kid", "key-1")
	j.SetHeader("typ", "JWT")
	hdrs := j.Headers()
	if len(hdrs) != 2 {
		t.Fatalf("Headers() length = %d, want 2", len(hdrs))
	}
	if hdrs["kid"] != "key-1" || hdrs["typ"] != "JWT" {
		t.Errorf("Headers() = %v, want {kid:key-1 typ:JWT}", hdrs)
	}
	// Verify returned map is a clone.
	delete(hdrs, "kid")
	if j.Header("kid") != "key-1" {
		t.Error("Headers() should return a copy, not the original map")
	}
}

func TestJWTParse(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		key := []byte("test-key-123456")
		tok, err := CreateToken(map[string]any{"sub": "u1"}, key)
		if err != nil {
			t.Fatal(err)
		}
		j := New()
		if err := j.Parse(tok); err != nil {
			t.Fatalf("Parse: %v", err)
		}
		if j.Payload("sub") != "u1" {
			t.Errorf("Parse payload sub = %v, want u1", j.Payload("sub"))
		}
	})
	t.Run("empty token", func(t *testing.T) {
		j := New()
		if err := j.Parse(""); err == nil {
			t.Fatal("Parse empty should return error")
		}
	})
	t.Run("invalid parts", func(t *testing.T) {
		j := New()
		if err := j.Parse("a.b"); err == nil {
			t.Fatal("Parse invalid should return error")
		}
	})
}

func TestJWTSignWith(t *testing.T) {
	key := []byte("test-key-123456")
	signer := MustHMACSigner(AlgHS256, key)
	j := New().SetPayload("sub", "u1")
	tok, err := j.SignWith(signer)
	if err != nil {
		t.Fatalf("SignWith: %v", err)
	}
	if tok == "" {
		t.Fatal("SignWith returned empty token")
	}
	if !strings.Contains(tok, ".") {
		t.Fatal("SignWith returned invalid token format")
	}
	// Verify the token.
	parsed, err := Of(tok)
	if err != nil {
		t.Fatal(err)
	}
	if !parsed.VerifyWith(signer) {
		t.Error("SignWith token should verify successfully")
	}
}

func TestJWTMustSign(t *testing.T) {
	t.Run("signs successfully", func(t *testing.T) {
		key := []byte("test-key-123456")
		j := New().SetPayload("sub", "u1").SetKey(key)
		tok := j.MustSign()
		if tok == "" {
			t.Fatal("MustSign returned empty token")
		}
	})
	t.Run("panics without signer", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("MustSign should panic without signer")
			}
		}()
		New().MustSign()
	})
}

func TestJWTValidateAt(t *testing.T) {
	key := []byte("test-key-123456")
	j := New().SetPayload("sub", "u1").SetKey(key)
	tok, err := j.Sign()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := Of(tok)
	if err != nil {
		t.Fatal(err)
	}
	// Set the key on the parsed JWT so that Verify succeeds.
	parsed.SetKey(key)
	now := time.Now()
	if !parsed.ValidateAt(now, 60) {
		t.Error("ValidateAt should pass for valid token")
	}
}

func TestParseTokenWithOptions(t *testing.T) {
	key := []byte("test-key-123456")
	tok, err := CreateToken(map[string]any{"sub": "u1"}, key)
	if err != nil {
		t.Fatal(err)
	}
	j, err := ParseTokenWithOptions(tok)
	if err != nil {
		t.Fatalf("ParseTokenWithOptions: %v", err)
	}
	if j == nil {
		t.Fatal("ParseTokenWithOptions returned nil JWT")
	}
	if j.Payload("sub") != "u1" {
		t.Errorf("ParseTokenWithOptions sub = %v, want u1", j.Payload("sub"))
	}

	// Test with custom JSON option.
	customTok, err := CreateToken(map[string]any{"sub": "u2"}, key)
	if err != nil {
		t.Fatal(err)
	}
	j2, err := ParseTokenWithOptions(customTok, WithJSONMarshalFunc(nil)) // opt is nil-guarded
	if err != nil {
		t.Fatalf("ParseTokenWithOptions with nil marshal: %v", err)
	}
	if j2.Payload("sub") != "u2" {
		t.Errorf("ParseTokenWithOptions opts sub = %v, want u2", j2.Payload("sub"))
	}
}

func TestOfValidatorJWT(t *testing.T) {
	j := New()
	v := OfValidatorJWT(j)
	if v == nil {
		t.Fatal("OfValidatorJWT returned nil")
	}
	if v.JWT() != j {
		t.Error("OfValidatorJWT.JWT() should return the same JWT")
	}
	if v.Err() != nil {
		t.Errorf("OfValidatorJWT.Err() should be nil, got %v", v.Err())
	}
}

func TestSignerUtilPS384(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	signer := PS384(priv, nil)
	if signer == nil {
		t.Fatal("PS384 returned nil")
	}
	if signer.Algorithm() != AlgPS384 {
		t.Errorf("PS384 algorithm = %q, want %q", signer.Algorithm(), AlgPS384)
	}
}

func TestSignerUtilPS384WithOptions(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	signer := PS384WithOptions(priv, nil)
	if signer == nil || signer.Algorithm() != AlgPS384 {
		t.Fatal("PS384WithOptions failed")
	}
}

func TestSignerUtilPS512(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	signer := PS512(priv, nil)
	if signer == nil || signer.Algorithm() != AlgPS512 {
		t.Fatal("PS512 failed")
	}
}

func TestSignerUtilPS512WithOptions(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	signer := PS512WithOptions(priv, nil)
	if signer == nil || signer.Algorithm() != AlgPS512 {
		t.Fatal("PS512WithOptions failed")
	}
}

func TestSignerUtilES384(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer := ES384(priv, nil)
	if signer == nil || signer.Algorithm() != AlgES384 {
		t.Fatal("ES384 failed")
	}
}

func TestSignerUtilES384WithOptions(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer := ES384WithOptions(priv, nil)
	if signer == nil || signer.Algorithm() != AlgES384 {
		t.Fatal("ES384WithOptions failed")
	}
}

func TestSignerUtilES512(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer := ES512(priv, nil)
	if signer == nil || signer.Algorithm() != AlgES512 {
		t.Fatal("ES512 failed")
	}
}

func TestSignerUtilES512WithOptions(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer := ES512WithOptions(priv, nil)
	if signer == nil || signer.Algorithm() != AlgES512 {
		t.Fatal("ES512WithOptions failed")
	}
}

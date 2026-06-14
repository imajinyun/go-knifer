package vjwt_test

import (
	"encoding/json"
	"testing"

	"github.com/imajinyun/go-knifer/vjwt"
)

func TestSignVerifyWithJSONProviders(t *testing.T) {
	marshalCalled := false
	unmarshalCalled := false
	marshal := func(v any) ([]byte, error) {
		marshalCalled = true
		return json.Marshal(v)
	}
	unmarshal := func(data []byte, v any) error {
		unmarshalCalled = true
		return json.Unmarshal(data, v)
	}

	jwt := vjwt.New().
		SetPayload(vjwt.JWTPayloadSubject, "alice").
		SetKey([]byte("secret"))
	token, err := jwt.SignOptsWithOptions(true, vjwt.WithJSONMarshalFunc(marshal))
	if err != nil {
		t.Fatalf("SignOptsWithOptions: %v", err)
	}
	parsed, err := vjwt.ParseTokenWithOptions(token, vjwt.WithJSONUnmarshalFunc(unmarshal))
	if err != nil {
		t.Fatalf("ParseTokenWithOptions: %v", err)
	}
	if !parsed.SetKey([]byte("secret")).Verify() {
		t.Fatal("parsed token should verify")
	}
	if !marshalCalled || !unmarshalCalled {
		t.Fatalf("JSON providers called marshal=%v unmarshal=%v", marshalCalled, unmarshalCalled)
	}
}

func TestCreateTokenWithOptionsStrictKey(t *testing.T) {
	if token, err := vjwt.CreateTokenWithOptions(
		vjwt.WithTokenPayload(map[string]any{vjwt.JWTPayloadSubject: "alice"}),
		vjwt.WithTokenKey([]byte("weak")),
		vjwt.WithTokenStrictKey(),
	); err == nil || token != "" {
		t.Fatalf("strict weak key token=%q err=%v, want error", token, err)
	}
}

func TestFacadeTokenConstructorsAndValidators(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	payload := map[string]any{vjwt.JWTPayloadSubject: "alice"}

	token, err := vjwt.CreateJWTToken(payload, key)
	if err != nil {
		t.Fatalf("CreateJWTToken: %v", err)
	}
	if !vjwt.VerifyJWT(token, key) {
		t.Fatal("VerifyJWT(CreateJWTToken) = false")
	}
	parsed := vjwt.NewJWT()
	if err := parsed.Parse(token); err != nil {
		t.Fatalf("NewJWT.Parse: %v", err)
	}
	if parsed.Payload(vjwt.JWTPayloadSubject) != "alice" {
		t.Fatalf("parsed subject = %#v", parsed.Payload(vjwt.JWTPayloadSubject))
	}
	if _, err := vjwt.JWTOf(token); err != nil {
		t.Fatalf("JWTOf: %v", err)
	}
	if _, err := vjwt.JWTOfWithOptions(token); err != nil {
		t.Fatalf("JWTOfWithOptions: %v", err)
	}
	if _, err := vjwt.ParseJWT(token); err != nil {
		t.Fatalf("ParseJWT: %v", err)
	}

	signer := vjwt.HS256(key)
	token, err = vjwt.CreateJWTTokenWithSigner(payload, signer)
	if err != nil || !vjwt.VerifyJWTWithSigner(token, signer) {
		t.Fatalf("CreateJWTTokenWithSigner token=%q err=%v", token, err)
	}
	token, err = vjwt.CreateToken(payload, key)
	if err != nil || !vjwt.Verify(token, key) {
		t.Fatalf("CreateToken token=%q err=%v", token, err)
	}
	token, err = vjwt.CreateTokenWithHeaders(map[string]any{vjwt.JWTHeaderKeyID: "kid-1"}, payload, key)
	if err != nil {
		t.Fatalf("CreateTokenWithHeaders: %v", err)
	}
	headered, err := vjwt.ParseToken(token)
	if err != nil || headered.Header(vjwt.JWTHeaderKeyID) != "kid-1" {
		t.Fatalf("CreateTokenWithHeaders parsed=%#v err=%v", headered, err)
	}
	token, err = vjwt.CreateTokenWithAlgorithm(payload, key, vjwt.JWTAlgHS384)
	if err != nil || !vjwt.VerifyWithSigner(token, vjwt.HS384(key)) {
		t.Fatalf("CreateTokenWithAlgorithm token=%q err=%v", token, err)
	}
	token, err = vjwt.CreateTokenWithHeadersAndAlgorithm(map[string]any{vjwt.JWTHeaderKeyID: "kid-2"}, payload, key, vjwt.JWTAlgHS512)
	if err != nil || !vjwt.VerifyWithSigner(token, vjwt.HS512(key)) {
		t.Fatalf("CreateTokenWithHeadersAndAlgorithm token=%q err=%v", token, err)
	}
	token, err = vjwt.CreateTokenWithSigner(payload, signer)
	if err != nil || !vjwt.VerifyWithSigner(token, signer) {
		t.Fatalf("CreateTokenWithSigner token=%q err=%v", token, err)
	}
	token, err = vjwt.CreateTokenWithHeadersAndSigner(map[string]any{vjwt.JWTHeaderKeyID: "kid-3"}, payload, signer)
	if err != nil || !vjwt.VerifyWithSigner(token, signer) {
		t.Fatalf("CreateTokenWithHeadersAndSigner token=%q err=%v", token, err)
	}
	token, err = vjwt.CreateTokenWithOptions(
		vjwt.WithTokenHeaders(map[string]any{vjwt.JWTHeaderKeyID: "kid-4"}),
		vjwt.WithTokenPayload(payload),
		vjwt.WithTokenSigner(signer),
		vjwt.WithTokenJSONOptions(),
	)
	if err != nil || vjwt.OfValidator(token).ValidateAlgorithm(signer).Err() != nil {
		t.Fatalf("CreateTokenWithOptions token=%q err=%v", token, err)
	}
	if vjwt.OfValidatorJWT(parsed).JWT() != parsed {
		t.Fatal("OfValidatorJWT did not retain JWT pointer")
	}
}

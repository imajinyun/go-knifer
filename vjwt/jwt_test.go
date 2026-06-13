package vjwt_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io"
	"testing"
	"time"

	knifer "github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vjwt"
)

func TestVerifyRejectsNoneToken(t *testing.T) {
	token := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiJwdWJsaWMifQ."

	if vjwt.Verify(token, []byte("ignored")) {
		t.Fatal("Verify should reject alg=none tokens")
	}
	if vjwt.VerifyStrict(token, []byte("ignored")) {
		t.Fatal("VerifyStrict should reject alg=none tokens")
	}
	if _, err := vjwt.CreateSigner("none", []byte("ignored")); !errors.Is(err, knifer.ErrCodeUnsupported) {
		t.Fatalf("CreateSigner(none) error = %v, want unsupported", err)
	}
}

func TestStrictHMACSignerRejectsWeakKey(t *testing.T) {
	if _, err := vjwt.NewHMACSignerStrict(vjwt.JWTAlgHS256, []byte("weak")); err == nil {
		t.Fatal("NewHMACSignerStrict should reject weak key")
	}
	if _, err := vjwt.CreateSignerStrict(vjwt.JWTAlgHS256, []byte("weak")); err == nil {
		t.Fatal("CreateSignerStrict should reject weak key")
	}
	if minBytes, err := vjwt.MinHMACKeyBytes(vjwt.JWTAlgHS256); err != nil || minBytes != vjwt.MinHMACKeyBytesHS256 {
		t.Fatalf("MinHMACKeyBytes = %d, %v", minBytes, err)
	}
}

func TestVerifyWithSignerRejectsAlgorithmMismatch(t *testing.T) {
	key := []byte("secret")
	token, err := vjwt.CreateTokenWithOptions(
		vjwt.WithTokenAlgorithm(vjwt.JWTAlgHS256),
		vjwt.WithTokenKey(key),
		vjwt.WithTokenPayload(map[string]any{vjwt.JWTPayloadSubject: "alice"}),
	)
	if err != nil {
		t.Fatalf("CreateTokenWithOptions: %v", err)
	}

	wrongAlgSigner, err := vjwt.NewHMACSigner(vjwt.JWTAlgHS512, key)
	if err != nil {
		t.Fatalf("NewHMACSigner: %v", err)
	}
	if vjwt.VerifyWithSigner(token, wrongAlgSigner) {
		t.Fatal("VerifyWithSigner should reject signer/token algorithm mismatch")
	}
	if err := vjwt.ValidateAlgorithm(token, wrongAlgSigner); err == nil {
		t.Fatal("ValidateAlgorithm should reject signer/token algorithm mismatch")
	}
}

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

func TestFacadeSignerFactoriesAndAlgorithms(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	for _, tt := range []struct {
		name string
		fn   func([]byte) vjwt.JWTSigner
		alg  string
	}{
		{name: "JWTSignerHS256", fn: vjwt.JWTSignerHS256, alg: vjwt.JWTAlgHS256},
		{name: "HS256", fn: vjwt.HS256, alg: vjwt.JWTAlgHS256},
		{name: "HS384", fn: vjwt.HS384, alg: vjwt.JWTAlgHS384},
		{name: "HS512", fn: vjwt.HS512, alg: vjwt.JWTAlgHS512},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(key).Algorithm(); got != tt.alg {
				t.Fatalf("Algorithm = %q, want %q", got, tt.alg)
			}
		})
	}
	signer, err := vjwt.JWTSignerHMAC(vjwt.JWTAlgHS384, key)
	if err != nil || signer.Algorithm() != vjwt.JWTAlgHS384 {
		t.Fatalf("JWTSignerHMAC alg=%q err=%v", signer.Algorithm(), err)
	}
	if got := vjwt.MustHMACSigner(vjwt.JWTAlgHS512, key).Algorithm(); got != vjwt.JWTAlgHS512 {
		t.Fatalf("MustHMACSigner alg = %q", got)
	}
	if got := vjwt.AlgorithmName(vjwt.JWTAlgPS256); got != "SHA256withRSA_PSS" {
		t.Fatalf("AlgorithmName(PS256) = %q", got)
	}
	if _, err := vjwt.NewRSAPSSSigner(vjwt.JWTAlgPS256, nil, nil); err == nil {
		t.Fatal("NewRSAPSSSigner(nil keys) error = nil")
	}
	if _, err := vjwt.NewECDSASigner(vjwt.JWTAlgES256, nil, nil); err == nil {
		t.Fatal("NewECDSASigner(nil keys) error = nil")
	}

	rsaKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatalf("rsa.GenerateKey: %v", err)
	}
	reader := zeroReader{}
	psSigner, err := vjwt.NewRSAPSSSignerWithOptions(vjwt.JWTAlgPS256, rsaKey, &rsaKey.PublicKey, vjwt.WithSignerRandomReader(reader), vjwt.WithRSAPSSOptions(nil))
	if err != nil || psSigner.Algorithm() != vjwt.JWTAlgPS256 {
		t.Fatalf("NewRSAPSSSignerWithOptions alg=%q err=%v", psSigner.Algorithm(), err)
	}
	if got := vjwt.PS256WithOptions(rsaKey, &rsaKey.PublicKey, vjwt.WithSignerRandomReader(reader)).Algorithm(); got != vjwt.JWTAlgPS256 {
		t.Fatalf("PS256WithOptions alg = %q", got)
	}
	if got := vjwt.PS384(rsaKey, &rsaKey.PublicKey).Algorithm(); got != vjwt.JWTAlgPS384 {
		t.Fatalf("PS384 alg = %q", got)
	}
	if got := vjwt.PS512WithOptions(rsaKey, &rsaKey.PublicKey).Algorithm(); got != vjwt.JWTAlgPS512 {
		t.Fatalf("PS512WithOptions alg = %q", got)
	}

	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("ecdsa.GenerateKey: %v", err)
	}
	ecSigner, err := vjwt.JWTSignerECDSA(vjwt.JWTAlgES256, ecdsaKey, &ecdsaKey.PublicKey)
	if err != nil || ecSigner.Algorithm() != vjwt.JWTAlgES256 {
		t.Fatalf("JWTSignerECDSA alg=%q err=%v", ecSigner.Algorithm(), err)
	}
	if got := vjwt.JWTSignerES256(ecdsaKey, &ecdsaKey.PublicKey).Algorithm(); got != vjwt.JWTAlgES256 {
		t.Fatalf("JWTSignerES256 alg = %q", got)
	}
	if got := vjwt.ES256WithOptions(ecdsaKey, &ecdsaKey.PublicKey, vjwt.WithSignerRandomReader(reader)).Algorithm(); got != vjwt.JWTAlgES256 {
		t.Fatalf("ES256WithOptions alg = %q", got)
	}
	ecdsa384Key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("ecdsa.GenerateKey P384: %v", err)
	}
	if got := vjwt.ES384(ecdsa384Key, &ecdsa384Key.PublicKey).Algorithm(); got != vjwt.JWTAlgES384 {
		t.Fatalf("ES384 alg = %q", got)
	}
	ecdsa521Key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatalf("ecdsa.GenerateKey P521: %v", err)
	}
	if got := vjwt.ES512WithOptions(ecdsa521Key, &ecdsa521Key.PublicKey).Algorithm(); got != vjwt.JWTAlgES512 {
		t.Fatalf("ES512WithOptions alg = %q", got)
	}
}

func TestFacadeDateValidationOptions(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	j := vjwt.New().
		SetPayload(vjwt.JWTPayloadNotBefore, now.Add(-time.Minute).Unix()).
		SetPayload(vjwt.JWTPayloadExpiresAt, now.Add(time.Minute).Unix()).
		SetPayload(vjwt.JWTPayloadIssuedAt, now.Add(-time.Second).Unix()).
		SetKey([]byte("0123456789abcdef0123456789abcdef"))
	token, err := j.Sign()
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	parsed, err := vjwt.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	parsed.SetKey([]byte("0123456789abcdef0123456789abcdef"))
	if err := vjwt.ValidateJWTDate(parsed, now, 0); err != nil {
		t.Fatalf("ValidateJWTDate: %v", err)
	}
	if err := vjwt.ValidateDate(parsed, now, 0); err != nil {
		t.Fatalf("ValidateDate: %v", err)
	}
	if !parsed.ValidateWithOptions(vjwt.WithValidateTime(now), vjwt.WithValidateClock(func() time.Time { return now }), vjwt.WithValidateLeeway(0)) {
		t.Fatal("ValidateWithOptions = false")
	}

	expired := vjwt.New().SetPayload(vjwt.JWTPayloadExpiresAt, now.Add(-2*time.Second).Unix()).SetKey([]byte("0123456789abcdef0123456789abcdef"))
	if err := vjwt.ValidateDate(expired, now, 1); err == nil {
		t.Fatal("ValidateDate(expired) error = nil")
	}
}

type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 1
	}
	return len(p), nil
}

var _ io.Reader = zeroReader{}

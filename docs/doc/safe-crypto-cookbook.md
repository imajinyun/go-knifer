# Safe Crypto Cookbook

Use this cookbook when an application needs hashing, authentication, encryption, signing, token creation, or secure random material across a trust boundary. Start with `vrand` for secret bytes, use `vcrypto` for digest, HMAC, AES-GCM, RSA, SM2, SM3, and SM4 workflows, and use `vjwt` when the artifact is a signed JWT rather than arbitrary bytes.

## Crypto Decision Matrix

| Task | First package | API family | Guardrail |
| --- | --- | --- | --- |
| Generate keys, tokens, salts, or nonces | `vrand` | `SecureBytes`, `SecureBytesWithOptions` | Encode secure bytes after generation; do not use pseudo-random string helpers for secrets. |
| Hash non-secret data | `vcrypto` | `SHA256Hex`, `SHA512Hex`, `SM3Hex` | Hashes are fingerprints, not authentication codes or encryption. |
| Authenticate a message with a shared secret | `vcrypto` | `HMACSHA256Hex`, `HMACBytes`, `HMACEqual`, `HMACSM3Hex` | Compare MACs with constant-time helpers and keep keys outside source code. |
| Encrypt arbitrary bytes | `vcrypto` | `AESSealGCM`, `AESOpenGCM`, `SM4SealGCM`, `SM4DecryptGCM` | Prefer authenticated encryption and never reuse a nonce with the same key. |
| Sign or verify structured application tokens | `vjwt` | `CreateTokenWithOptions`, `CreateTokenWithSigner`, `VerifyWithSigner`, `ValidateWithOptions` | Parse, verify, then validate time and application claims before authorizing. |
| Sign or verify raw payloads | `vcrypto` | `RSASignPSS`, `RSAVerifyPSS`, `SM2Sign`, `SM2Verify` | Match the algorithm to the interoperability contract and keep private keys secret. |

## Secret Boundary Checklist

- Treat keys, tokens, salts, nonces, password-derived keys, private keys, and bearer credentials as secrets.
- Use `vrand.SecureBytes` or `vcrypto.RandomBytes` for secret material; pseudo-random helpers are for simulations, examples, and non-secret values.
- Keep nonce policy explicit. AES-GCM and SM4-GCM require a fresh nonce for each key/message pair.
- Prefer HMAC for message authentication with a shared secret; do not use plain hashes as tamper checks.
- Prefer AES-GCM or SM4-GCM over unauthenticated encryption modes for new designs.
- Keep SM2, SM3, and SM4 usage tied to national-crypto interoperability or policy requirements.
- Verify JWT signatures before trusting claims, and validate `exp`, `nbf`, `iat`, issuer, audience, subject, tenant, scope, and key id according to the application policy.
- Do not store credentials, private data, or secrets in readable JWT payloads unless a separate encryption layer protects them.
- Do not log raw secrets, private keys, nonces, salts, derived credentials, signatures, or bearer tokens.
- Treat errors such as `ErrInvalidCipherText` and invalid JWT signatures as security-relevant failures.

## Recipes

### Generate Secret Material

Use secure bytes for tokens, keys, salts, and nonces. Encode bytes only after the entropy budget has been selected.

```go
package main

import (
	"encoding/base64"
	"fmt"

	"github.com/imajinyun/knifer-go/vrand"
)

func main() {
	tokenBytes, err := vrand.SecureBytes(32)
	if err != nil {
		panic(err)
	}

	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	fmt.Println(len(token) > 0)
}
```

### Authenticate Messages

Use HMAC when both sides share a secret and the receiver needs to detect tampering.

```go
package main

import (
	"fmt"

	"github.com/imajinyun/knifer-go/vcrypto"
)

func main() {
	key := []byte("shared secret")
	message := []byte("important payload")

	mac := vcrypto.HMACSHA256Hex(key, message)
	expected := vcrypto.HMACSHA256Hex(key, message)

	fmt.Println(vcrypto.HMACEqual([]byte(mac), []byte(expected)))
}
```

### Encrypt Data

Use AES-GCM for common authenticated encryption workflows. Keep associated data stable across encryption and decryption.

```go
package main

import (
	"fmt"

	"github.com/imajinyun/knifer-go/vcrypto"
)

func main() {
	key, err := vcrypto.GenAESKey(32)
	if err != nil {
		panic(err)
	}

	aad := []byte("tenant:demo")
	nonce, cipherText, err := vcrypto.AESSealGCM([]byte("secret data"), key, aad)
	if err != nil {
		panic(err)
	}

	plain, err := vcrypto.AESOpenGCM(cipherText, key, nonce, aad)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(plain))
}
```

### Sign And Verify Tokens

Use JWT helpers when the artifact is a signed token with claims. Verification is not the end of validation: time and application claims still need policy checks.

```go
package main

import (
	"fmt"
	"time"

	"github.com/imajinyun/knifer-go/vjwt"
)

func main() {
	now := time.Unix(1700000000, 0)
	key := []byte("shared secret")

	token, err := vjwt.CreateTokenWithOptions(
		vjwt.WithTokenPayload(map[string]any{
			vjwt.JWTPayloadIssuer:    "issuer",
			vjwt.JWTPayloadSubject:   "alice",
			vjwt.JWTPayloadIssuedAt:  now.Unix(),
			vjwt.JWTPayloadExpiresAt: now.Add(time.Hour).Unix(),
		}),
		vjwt.WithTokenKey(key),
	)
	if err != nil {
		panic(err)
	}

	parsed, err := vjwt.ParseToken(token)
	if err != nil {
		panic(err)
	}

	fmt.Println(vjwt.Verify(token, key))
	fmt.Println(parsed.ValidateWithOptions(vjwt.WithValidateTime(now)) == nil)
}
```

## Validation

Run these checks after changing Safe Crypto cookbook guidance, examples, or governance metadata:

```bash
go test ./vcrypto ./vrand ./vjwt ./internal/crypto ./internal/rand ./internal/jwt
make docs-check
make ai-context-check
make governance-maturity-check
make agent-security-check
```

The cookbook is governed by `safe_crypto_cookbook_governance` in `ai-context.json`; `make governance-maturity-check` verifies the document path, covered packages, required scenarios, scorecard comparison and cookbook status, and Sprint 24 roadmap state.

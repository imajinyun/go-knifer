# vjwt Quickstart

`vjwt` provides JWT creation, parsing, signature verification, date-claim validation, and multiple signers including HMAC, RSA-PSS, and ECDSA.

## Which helper should I use?

Choose helpers by the trust boundary: token creation, signature verification, claim validation, or key/algorithm selection.

| Need | Use | Notes |
| --- | --- | --- |
| Create a simple HMAC token | `CreateToken` | Good for trusted internal tokens when the shared secret is managed outside source code. |
| Create a token with explicit headers, payload, algorithm, and key | `CreateTokenWithOptions` | Prefer this when `kid`, issuer, audience, or algorithm policy must be visible at the call site. |
| Use asymmetric signing | `CreateTokenWithSigner`, `PS256`, ECDSA/RSA signer helpers | Prefer asymmetric signers when verifiers should not hold signing keys. |
| Parse token structure without trusting it | `ParseToken`, `JWTOf` | Parsing exposes headers and claims; it does not by itself make the token trustworthy. |
| Verify a token signature | `Verify`, `VerifyWithSigner` | Verify before authorizing a request or trusting claims. |
| Validate time-based claims | `ValidateWithOptions`, `WithValidateTime`, `WithValidateLeeway` | Use deterministic validation time in tests and small leeway for clock skew. |

## JWT safety checklist

- Verify signatures before trusting any header or payload claim for authorization decisions.
- Keep accepted algorithms explicit. Do not let untrusted token headers silently choose an unexpected signing method.
- Validate `exp`, `nbf`, and `iat` where applicable, and keep clock-skew leeway small and documented.
- Validate application claims such as `iss`, `aud`, `sub`, tenant, scope, and key id against your own policy.
- Treat JWT payloads as readable metadata, not encrypted secrets. Do not store credentials or sensitive personal data in plain JWT claims.
- Rotate and scope signing keys outside source code. Prefer asymmetric signers when many services only need verification.

## FAQ

### Is a JWT encrypted?

No. A signed JWT protects integrity, not confidentiality. Anyone who can read the token can decode its header and payload unless you use a separate encryption scheme.

### How should I choose HMAC vs RSA/ECDSA signers?

Use HMAC when a small trusted set of services can safely share one secret. Use asymmetric signers when signing and verification responsibilities should be separated, such as many verifiers and one issuer.

### Is parsing enough before using claims?

No. Parse, verify the signature, then validate time and application claims before authorizing access.

## Create and verify tokens with HMAC

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vjwt"
)

func main() {
	key := []byte("secret")
	token, err := vjwt.CreateToken(map[string]any{vjwt.JWTPayloadSubject: "alice"}, key)
	if err != nil {
		panic(err)
	}

	parsed, err := vjwt.ParseToken(token)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsed.Payload(vjwt.JWTPayloadSubject))
	fmt.Println(vjwt.Verify(token, key))
}
```

## Set headers, payloads, and algorithms with options

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vjwt"
)

func main() {
	token, err := vjwt.CreateTokenWithOptions(
		vjwt.WithTokenHeaders(map[string]any{vjwt.JWTHeaderKeyID: "key-1"}),
		vjwt.WithTokenPayload(map[string]any{vjwt.JWTPayloadIssuer: "issuer", vjwt.JWTPayloadSubject: "alice"}),
		vjwt.WithTokenAlgorithm(vjwt.JWTAlgHS384),
		vjwt.WithTokenKey([]byte("secret")),
	)
	if err != nil {
		panic(err)
	}

	parsed, err := vjwt.JWTOf(token)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsed.Header(vjwt.JWTHeaderKeyID), parsed.Algorithm())
}
```

## Build JWTs fluently and validate date claims

```go
package main

import (
	"fmt"
	"time"

	"github.com/imajinyun/go-knifer/vjwt"
)

func main() {
	now := time.Now()
	j := vjwt.New().
		SetKey([]byte("secret")).
		SetIssuer("go-knifer").
		SetSubject("alice").
		SetIssuedAt(now).
		SetNotBefore(now.Add(-time.Minute)).
		SetExpiresAt(now.Add(time.Hour))

	token, err := j.Sign()
	if err != nil {
		panic(err)
	}

	parsed, err := vjwt.ParseToken(token)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsed.ValidateWithOptions(vjwt.WithValidateTime(now), vjwt.WithValidateLeeway(30)))
}
```

## Use an RSA-PSS signer

```go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/imajinyun/go-knifer/vjwt"
)

func main() {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	signer := vjwt.PS256(priv, &priv.PublicKey)

	token, err := vjwt.CreateTokenWithSigner(map[string]any{vjwt.JWTPayloadSubject: "alice"}, signer)
	if err != nil {
		panic(err)
	}

	fmt.Println(vjwt.VerifyWithSigner(token, signer))
}
```

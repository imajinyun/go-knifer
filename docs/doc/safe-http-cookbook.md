# Safe HTTP Cookbook

Use this cookbook when an application needs to fetch, inspect, or download HTTP(S) resources across a trust boundary. Start with `vurl` for URL shape and safe probing, use `vhttp` for dependency-light request and download helpers, and use `vresty` when the codebase already standardizes on Resty-style request chains.

## Decision Matrix

| Task | First package | API family | Guardrail |
| --- | --- | --- | --- |
| Normalize or build a URL before any request | `vurl` | `NormalizeUsingOptions`, `NewHTTPURLBuilder`, `AppendQuery` | Normalization is not safety; apply safe open/probe or request policy before network access. |
| Probe an untrusted HTTP(S) URL | `vurl` | `OpenSafeWithOptions`, `ContentLengthSafeWithOptions` | Set allowed schemes, allowed hosts, timeout, byte limit, and deterministic resolver in tests. |
| Fetch a trusted response with standard-library semantics | `vhttp` | `GetStringE`, `Get`, `NewRequest` | Prefer explicit-error helpers when callers must branch on failure. |
| Fetch an untrusted response with standard-library semantics | `vhttp` | `GetStringSafeE`, `GetSafe`, `DownloadBytesSafeE`, `DownloadFileSafe` | Keep URL policy and response-size limits at the call site. |
| Fetch through a Resty-based application boundary | `vresty` | `GetStringE`, `Get`, `NewClient`, `WithRestyClientFactory` | Use isolated clients or injected Resty clients in tests. |
| Fetch untrusted URLs through Resty | `vresty` | `GetStringSafeE`, `GetSafe`, `DownloadBytesSafeE`, `DownloadFileSafe` | Pair safe constructors with `WithAllowedHosts`, `WithURLPolicy`, and response limits. |

## Trust Boundary Checklist

- Treat URLs from users, configuration files, webhooks, service discovery, queues, and partner payloads as untrusted.
- Allow only `http` and `https` unless a narrower caller contract exists.
- Prefer allow-listed hosts for known upstreams; host allow-lists are easier to audit than broad DNS rules.
- Keep private, loopback, link-local, multicast, and unspecified address rejection enabled for internet-facing inputs.
- Re-check redirect targets; safe helpers must not validate only the first URL.
- Bound response reads with `WithMaxBytes`, `WithMaxResponseBytes`, or download helpers that expose size policy.
- Use timeouts for every remote probe or request.
- Inject lookup or client providers in tests so examples and contracts do not depend on external DNS or internet access.
- Treat file destinations as a second boundary when downloading; review overwrite, parent creation, permissions, and safe filenames.

## Recipes

### Validate and Probe Before Requesting

Use `vurl` when the application needs to decide whether a remote resource is reachable or acceptable before choosing an HTTP client facade.

```go
package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/imajinyun/knifer-go/vurl"
)

func main() {
	size, err := vurl.ContentLengthSafeWithOptions(
		"https://api.example.com/report.csv",
		vurl.WithAllowedSchemes("https"),
		vurl.WithAllowedHosts("api.example.com"),
		vurl.WithLookupIP(func(context.Context, string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("203.0.113.10")}, nil
		}),
		vurl.WithTimeout(3*time.Second),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(size >= 0)
}
```

### Fetch With Standard-Library Style Helpers

Use `vhttp` when the project wants a small dependency surface and `net/http`-style behavior.

```go
package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/imajinyun/knifer-go/vhttp"
)

func main() {
	body, err := vhttp.GetStringSafeE(
		"https://api.example.com/users",
		vhttp.WithAllowedHosts("api.example.com"),
		vhttp.WithLookupIP(func(context.Context, string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("203.0.113.10")}, nil
		}),
		vhttp.WithTimeout(3*time.Second),
		vhttp.WithMaxResponseBytes(1<<20),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(body))
}
```

### Fetch Through Resty-Style Clients

Use `vresty` when Resty request chains are already the application convention.

```go
package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/imajinyun/knifer-go/vresty"
)

func main() {
	resp := vresty.GetSafe(
		"https://api.example.com/users",
		vresty.WithAllowedHosts("api.example.com"),
		vresty.WithLookupIP(func(context.Context, string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("203.0.113.10")}, nil
		}),
		vresty.WithTimeout(3*time.Second),
		vresty.WithMaxResponseBytes(1<<20),
	).Header("Accept", "application/json").Execute()
	if resp.Err() != nil {
		panic(resp.Err())
	}
	fmt.Println(resp.Status())
}
```

### Download To Files

For user-controlled URLs, validate the URL and keep destination behavior explicit. Use `vhttp.DownloadFileSafe` for the standard-library path or `vresty.DownloadFileSafeWithOptions` for Resty-backed flows.

```go
package main

import (
	"fmt"

	"github.com/imajinyun/knifer-go/vhttp"
)

func main() {
	n, err := vhttp.DownloadFileSafe(
		"https://files.example.com/export.csv",
		"./downloads/export.csv",
		vhttp.WithAllowedHosts("files.example.com"),
		vhttp.WithMaxResponseBytes(10<<20),
		vhttp.WithSaveOverwrite(false),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
}
```

## Validation

Run these checks after changing Safe HTTP cookbook guidance, examples, or governance metadata:

```bash
go test ./vurl ./vhttp ./vresty ./internal/url ./internal/httpx/...
make docs-check
make ai-context-check
make governance-maturity-check
make agent-security-check
```

The cookbook is governed by `safe_http_cookbook_governance` in `ai-context.json`; `make governance-maturity-check` verifies the document path, covered packages, required scenarios, and roadmap scorecard status.

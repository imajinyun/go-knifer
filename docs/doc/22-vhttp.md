# vhttp Quickstart

`vhttp` provides chained HTTP requests, shortcut GET/POST/download helpers, response saving, global/client configuration, safe URL validation, and simple HTTP server wrappers.

## Recommended HTTP entry points

| Scenario | Recommended package | Recommended API family | Why |
| --- | --- | --- | --- |
| Standard-library style request | `vhttp` | `Get`, `Post`, request builders | Keeps dependencies light and behavior explicit. |
| Resty-based fluent client | `vresty` | client/request helpers | Uses Resty ergonomics while keeping go-knifer safety docs. |
| Untrusted URL | `vhttp`/`vresty` plus safe APIs | `Safe`/`E` variants | Applies validation and explicit errors before network access. |
| File download | `vhttp`/`vresty` safe download helpers | `DownloadFileSafe` family | Keeps path and transfer risks visible. |

## Benchmarks and trade-offs

Use the HTTP benchmark suite to measure the convenience and safety overhead on your machine:

```bash
go test -bench=. -benchmem -run=^$ ./internal/httpx/... ./vhttp ./vresty
```

The suite uses `httptest.Server` and temporary files only. It covers simple GET requests, JSON response decode, bounded body reads, safe URL validation, and safe file downloads. Treat the output as a local baseline rather than a universal performance claim.

`vhttp` does not replace `net/http`; it provides repeatable convenience helpers and safe entry points for common request, response, and download flows.

`vresty` does not replace Resty; it keeps Resty ergonomics while documenting go-knifer's safety boundaries and generated examples.

Safe APIs may add validation overhead. Use the benchmark commands in this document to measure the trade-off on your workload.

## FAQ

### Why not use only `net/http`?

Use `net/http` directly when you need full control. Use `vhttp` when the common request, bounded read, or safe download path matches your use case and you want less boilerplate.

### How do I choose `vhttp` vs `vresty`?

Choose `vhttp` for lightweight standard-library style helpers. Choose `vresty` when your codebase already uses Resty or needs Resty's fluent request/client model.

### Are safe APIs free?

No. Safe APIs perform validation before work that can touch untrusted network or filesystem boundaries. Measure with the documented benchmark commands.

## Send chained requests

```go
package main

import (
	"fmt"
	"time"

	"github.com/imajinyun/go-knifer/vhttp"
)

func main() {
	resp := vhttp.Get("https://example.com",
		vhttp.WithTimeout(5*time.Second),
		vhttp.WithHeader("Accept", "text/html"),
	).
		Query("page", 1).
		Execute()
	defer resp.Close()

	if err := resp.Err(); err != nil {
		panic(err)
	}
	fmt.Println(resp.Status(), resp.ContentType())
}
```

## Submit data with shortcut helpers

```go
package main

import (
	"fmt"
	"time"

	"github.com/imajinyun/go-knifer/vhttp"
)

func main() {
	body, err := vhttp.PostJSONE(
		"https://example.com/api",
		`{"name":"alice"}`,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(body)

	body, err = vhttp.GetWithTimeoutE("https://example.com", 3*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(body))
}
```

## Download response content

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/imajinyun/go-knifer/vhttp"
)

func main() {
	data, err := vhttp.DownloadBytesE("https://example.com/file.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(len(data))

	var buf bytes.Buffer
	n, err := vhttp.DownloadWithOptions("https://example.com/file.txt", &buf, vhttp.WithMaxResponseBytes(1<<20))
	if err != nil {
		panic(err)
	}
	fmt.Println(n, buf.Len())
}
```

## Create a simple HTTP service

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/imajinyun/go-knifer/vhttp"
)

func main() {
	server := vhttp.NewSimpleServerWithOptions(8080, vhttp.WithReadHeaderTimeout(5*time.Second))
	server.AddAction("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "ok")
	})
	server.SetRoot("./public")

	// errCh := server.StartAsync()
	// _ = server.Stop(5 * time.Second)
	_ = server
}
```

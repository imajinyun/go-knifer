# vconf Quickstart

`vconf` reads, parses, and manages grouped configuration, with support for setting/properties, YAML, TOML, profile overrides, environment expansion, schema validation, and file watching.

## Which helper should I use?

Choose the helper by source format, layering needs, and how strictly configuration must be validated before use.

| Need | Use | Notes |
| --- | --- | --- |
| Parse inline setting/properties text | `Parse`, `ParseBytes` | Good for tests, generated config, and simple key/value files. |
| Parse by file extension | `ParseByExt`, `ParseByExtWithOptions` | Keeps format dispatch explicit when callers accept multiple config formats. |
| Parse YAML or TOML | `ParseYAML`, `ParseYAMLFull`, `ParseTOML` | Use the full parser variants when nested structures matter. |
| Load one or more files | `Load`, `LoadFiles`, `LoadWithOptions` | Use ordered file lists to make precedence reviewable. |
| Apply environment or profile overrides | `GetExpandedWithOptions`, `ApplyProfile`, `LoadProfile` | Inject env lookup in tests for deterministic expansion. |
| Bind configuration to structs | `Bind`, `BindGroup`, bind options | Prefer binding before application startup code uses configuration values. |
| Validate configuration | `SchemaFromStruct`, validation helpers | Validate required fields, ranges, and type expectations before starting long-running work. |
| Watch a file for changes | `Watch`, `WatchWithOptions` | Ensure callbacks are idempotent and safe to run more than once. |
| Load remote configuration | `LoadRemoteSafe`, `LoadRemoteSafeWithOptions` | Prefer safe remote helpers for URLs from config, users, or service discovery. |

## Configuration safety checklist

- Keep configuration precedence explicit: defaults, files, profiles, environment expansion, and remote sources should be reviewable in order.
- Use safe remote loading for any URL that is not a compile-time constant owned by the application.
- Inject environment lookup in tests instead of depending on the host process environment.
- Validate required fields and ranges before starting services, opening network listeners, or launching background workers.
- Treat decrypted or expanded values as secrets when they contain credentials; do not log raw configuration maps.
- Make file watchers resilient: callbacks should tolerate partial writes, repeated events, and invalid intermediate config.

## FAQ

### Does vconf replace Viper?

No. `vconf` is a lightweight grouped configuration helper for common parsing, binding, profile, and validation workflows. Use Viper when you need broad ecosystem features such as Cobra integration, remote providers, or complex precedence stacks.

### When should I use ParseYAMLFull instead of ParseYAML?

Use `ParseYAMLFull` when nested YAML structures and standard YAML behavior matter. Use `ParseYAML` for the smaller supported subset when simple grouped config is enough.

### How should I test environment expansion?

Use `WithEnvLookup` to inject deterministic values. This keeps tests independent of the developer machine, CI environment, and secret-bearing process variables.

## Parse TOML and read grouped values

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vconf"
)

func main() {
	c, err := vconf.ParseTOML(`
name = "demo"
[server]
port = 8080
debug = true
`)
	if err != nil {
		panic(err)
	}

	fmt.Println(c.Get("name"))
	fmt.Println(c.GetIntByGroup("server", "port", 0))
	fmt.Println(c.GetBoolByGroup("server", "debug", false))
}
```

## Expand environment variables

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vconf"
)

func main() {
	c, err := vconf.Parse("base=http://${ENV:HOST}\n")
	if err != nil {
		panic(err)
	}

	value := c.GetExpandedWithOptions("base", vconf.WithEnvLookup(func(name string) string {
		if name == "HOST" {
			return "localhost:8080"
		}
		return ""
	}))
	fmt.Println(value)
}
```

## Bind to a struct

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vconf"
)

type Server struct {
	Port  int      `conf:"port"`
	Debug bool     `conf:"debug"`
	Tags  []string `conf:"tags"`
}

func main() {
	c, err := vconf.ParseTOML(`
[server]
port = 8080
debug = true
tags = ["api", "admin"]
`)
	if err != nil {
		panic(err)
	}

	var server Server
	if err := c.BindGroup("server", &server); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", server)
}
```

## Apply profile overrides

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vconf"
)

func main() {
	c, err := vconf.ParseTOML(`
[server]
port = 8080
[profile.prod.server]
port = 9090
`)
	if err != nil {
		panic(err)
	}

	prod := c.ApplyProfile("prod")
	fmt.Println(prod.GetByGroup("server", "port"))
}
```

# vform Quickstart

`vform` validates common form fields such as email, mobile numbers, URLs, IPs, ID cards, Chinese text, and numeric strings, and also supports injecting matchers for selected rules.

## Validate email, mobile, and URL values

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vform"
)

func main() {
	fmt.Println(vform.IsEmail("alice@example.com"))
	fmt.Println(vform.IsMobile("13800138000"))
	fmt.Println(vform.IsURL("https://example.com/path"))
}
```

## Validate IP, Chinese text, and numeric strings

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vform"
)

func main() {
	fmt.Println(vform.IsIPv4("127.0.0.1"))
	fmt.Println(vform.IsIPv6("::1"))
	fmt.Println(vform.IsChinese("\u4e2d\u6587"))
	fmt.Println(vform.IsNumberStr("-12.34"))
}
```

## Validate ID card numbers

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vform"
)

func main() {
	fmt.Println(vform.IsIDCard("11010519491231002X"))
	fmt.Println(vform.IsIDCardWithOptions(
		"custom-id",
		vform.WithIDCardMatcher(func(s string) bool { return s == "custom-id" }),
	))
}
```

## Inject custom matchers

```go
package main

import (
	"fmt"
	"strings"

	"github.com/imajinyun/go-knifer/vform"
)

func main() {
	fmt.Println(vform.IsEmailWithOptions(
		"user@internal",
		vform.WithEmailMatcher(func(s string) bool { return strings.HasSuffix(s, "@internal") }),
	))

	fmt.Println(vform.IsNumberStrWithOptions(
		"N/A",
		vform.WithNumberMatcher(func(s string) bool { return s == "N/A" }),
	))
}
```

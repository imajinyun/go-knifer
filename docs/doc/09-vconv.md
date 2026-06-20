# vconv Quickstart

`vconv` provides loose type conversion helpers that convert common inputs to string, int, int64, float64, bool, and []byte. Each scalar family has zero-value helpers, default-value helpers, and explicit-error `E` helpers for code that must distinguish invalid input from a valid zero value.

## Conversion contract

- `ToXxx` helpers return the destination type zero value when conversion fails.
- `ToXxxDefault` helpers return the caller-provided fallback when conversion fails.
- `ToIntE`, `ToInt64E`, `ToFloat64E`, and `ToBoolE` return `ErrInvalidConversion` and match `knifer.ErrCodeInvalidInput` when conversion fails.
- String-to-int conversion trims spaces, tries integer parsing first, then accepts float strings by truncating toward zero, so `"42.9"` becomes `42`.
- Bool conversion accepts `true`, `yes`, `y`, `ok`, `1`, `on`, `false`, `no`, `n`, `0`, and `off` case-insensitively after trimming spaces. Non-string numerics convert to `true` when nonzero.
- `ToBytes` returns `nil` for `nil`, returns an existing `[]byte` as-is, converts strings directly, and stringifies other values.

## Convert to numbers

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vconv"
)

func main() {
	fmt.Println(vconv.ToInt("42"))
	fmt.Println(vconv.ToIntDefault("bad", 7))
	fmt.Println(vconv.ToFloat64("3.14"))
}
```

## Return explicit conversion errors

```go
package main

import (
	"errors"
	"fmt"

	knifer "github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vconv"
)

func main() {
	value, err := vconv.ToInt64E("42.9")
	fmt.Println(value, err)

	_, err = vconv.ToBoolE("maybe")
	fmt.Println(errors.Is(err, vconv.ErrInvalidConversion))
	fmt.Println(errors.Is(err, knifer.ErrCodeInvalidInput))
}
```

## Convert to bool

```go
package main

import (
	"fmt"
	"strings"

	"github.com/imajinyun/go-knifer/vconv"
)

func main() {
	fmt.Println(vconv.ToBool("true"))
	fmt.Println(vconv.ToBoolWithOptions("YES", vconv.WithBoolParser(func(s string) (bool, error) {
		return strings.EqualFold(s, "yes"), nil
	})))
}
```

## Convert to strings

```go
package main

import (
	"fmt"
	"strconv"

	"github.com/imajinyun/go-knifer/vconv"
)

func main() {
	fmt.Println(vconv.ToString(123))
	fmt.Println(vconv.ToStringDefault(nil, "fallback"))
	fmt.Println(vconv.ToStringWithOptions(true, vconv.WithFormatBoolFunc(strconv.FormatBool)))
}
```

## Convert to byte slices

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vconv"
)

func main() {
	b := vconv.ToBytes("hello")
	fmt.Println(string(b))
}
```

# vbean Quickstart

`vbean` maps, copies, and loosely converts fields between structs and maps, with options for tags, case matching, and skipping empty values.

## Convert a struct to a map

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vbean"
)

type UserDTO struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

func main() {
	m, err := vbean.ToMap(UserDTO{Name: "alice", Age: "18"})
	if err != nil {
		panic(err)
	}
	fmt.Println(m["name"], m["age"])
}
```

## Fill a struct from a map

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vbean"
)

type User struct {
	Name string `json:"full_name"`
	Age  int    `json:"age"`
}

func main() {
	var user User
	err := vbean.ToStruct(map[string]any{"FULL_NAME": "drew", "age": "21"}, &user,
		vbean.WithCaseInsensitive(true),
		vbean.WithWeaklyTyped(true),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s:%d\n", user.Name, user.Age)
}
```

## Use custom tags and skip zero values

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vbean"
)

type Row struct {
	Name string `db:"user_name"`
	Age  int    `db:"age"`
}

func main() {
	m, err := vbean.ToMap(Row{Name: "casey", Age: 0},
		vbean.WithTagNames("db"),
		vbean.WithIgnoreZero(true),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(m) // age is skipped
}
```

## Copy fields while preserving existing non-empty values

```go
package main

import (
	"fmt"

	"github.com/imajinyun/go-knifer/vbean"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	dst := User{Name: "existing", Age: 30}
	err := vbean.Copy(map[string]any{"name": "", "age": "22"}, &dst,
		vbean.WithIgnoreEmpty(true),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", dst)
}
```

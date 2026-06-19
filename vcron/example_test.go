package vcron_test

import (
	"fmt"
	"time"

	"github.com/imajinyun/go-knifer/vcron"
)

func ExampleNewPattern() {
	p, err := vcron.NewPattern("* * * * *")
	fmt.Println(p != nil)
	fmt.Println(err)
	// Output:
	// true
	// <nil>
}

func ExampleMustNewPattern() {
	p := vcron.MustNewPattern("0 9 * * *")
	t := time.Date(2026, 6, 15, 9, 0, 30, 0, time.UTC)

	fmt.Println(p.Raw())
	fmt.Println(p.Match(t, false))
	// Output:
	// 0 9 * * *
	// true
}

func ExampleNewConfigWithOptions() {
	loc := time.FixedZone("docs", 8*60*60)
	cfg := vcron.NewConfigWithOptions(
		vcron.WithConfigLocation(loc),
		vcron.WithConfigMatchSecond(true),
	)

	fmt.Println(cfg.Location.String())
	fmt.Println(cfg.MatchSecond)
	// Output:
	// docs
	// true
}

// Package knifer is the root package of the go-knifer utility toolkit.
//
// This module is split into multiple public subpackages by domain.
// Import only the packages you need:
//
//	import "github.com/imajinyun/go-knifer/vstr"
//	import "github.com/imajinyun/go-knifer/vslice"
//	import "github.com/imajinyun/go-knifer/vcache"
//	import "github.com/imajinyun/go-knifer/vcrypto"
//	import "github.com/imajinyun/go-knifer/vzip"
//	import "github.com/imajinyun/go-knifer/vhttp"
//	import "github.com/imajinyun/go-knifer/vconf"
//
// Subpackages are independent from each other, and the root package does not
// expose business APIs.
//
// The project follows an internal implementation plus public facade layout:
// concrete implementations live in internal/* packages, while application code
// should import the public v* packages. This keeps domain boundaries explicit
// and allows internal implementations to evolve without exposing every helper as
// public API.
package knifer

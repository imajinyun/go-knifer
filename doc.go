// Package knifer is the root package of the go-knifer utility toolkit.
//
// This module is split into multiple public subpackages by domain.
// Import only the packages you need:
//
// import "github.com/imajinyun/go-knifer/vbean"
// import "github.com/imajinyun/go-knifer/vblf"
// import "github.com/imajinyun/go-knifer/vbool"
// import "github.com/imajinyun/go-knifer/vcache"
// import "github.com/imajinyun/go-knifer/vcaptcha"
// import "github.com/imajinyun/go-knifer/vcodec"
// import "github.com/imajinyun/go-knifer/vconf"
// import "github.com/imajinyun/go-knifer/vconv"
// import "github.com/imajinyun/go-knifer/vcrypto"
// import "github.com/imajinyun/go-knifer/vdate"
// import "github.com/imajinyun/go-knifer/vdb"
// import "github.com/imajinyun/go-knifer/vdes"
// import "github.com/imajinyun/go-knifer/verr"
// import "github.com/imajinyun/go-knifer/vfile"
// import "github.com/imajinyun/go-knifer/vhash"
// import "github.com/imajinyun/go-knifer/vhttp"
// import "github.com/imajinyun/go-knifer/vid"
// import "github.com/imajinyun/go-knifer/vident"
// import "github.com/imajinyun/go-knifer/vjob"
// import "github.com/imajinyun/go-knifer/vjson"
// import "github.com/imajinyun/go-knifer/vjwt"
// import "github.com/imajinyun/go-knifer/vlog"
// import "github.com/imajinyun/go-knifer/vmap"
// import "github.com/imajinyun/go-knifer/vnet"
// import "github.com/imajinyun/go-knifer/vnum"
// import "github.com/imajinyun/go-knifer/vstr"
// import "github.com/imajinyun/go-knifer/vobj"
// import "github.com/imajinyun/go-knifer/vpoi"
// import "github.com/imajinyun/go-knifer/vrand"
// import "github.com/imajinyun/go-knifer/vref"
// import "github.com/imajinyun/go-knifer/vregex"
// import "github.com/imajinyun/go-knifer/vresty"
// import "github.com/imajinyun/go-knifer/vsem"
// import "github.com/imajinyun/go-knifer/vset"
// import "github.com/imajinyun/go-knifer/vstk"
// import "github.com/imajinyun/go-knifer/vslice"
// import "github.com/imajinyun/go-knifer/vstr"
// import "github.com/imajinyun/go-knifer/vsys"
// import "github.com/imajinyun/go-knifer/vtpl"
// import "github.com/imajinyun/go-knifer/vurl"
// import "github.com/imajinyun/go-knifer/vver"
// import "github.com/imajinyun/go-knifer/vxml"
// import "github.com/imajinyun/go-knifer/vzip"
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

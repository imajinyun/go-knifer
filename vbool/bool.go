package vbool

import boolimpl "github.com/imajinyun/go-knifer/internal/boolutil"

func Negate(b bool) bool  { return boolimpl.BoolNegate(b) }
func ToInt(b bool) int    { return boolimpl.BoolToInt(b) }
func And(bs ...bool) bool { return boolimpl.BoolAnd(bs...) }
func Or(bs ...bool) bool  { return boolimpl.BoolOr(bs...) }

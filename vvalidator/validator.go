package vvalidator

import validatorimpl "github.com/imajinyun/go-knifer/internal/validator"

func IsEmail(s string) bool     { return validatorimpl.IsEmail(s) }
func IsMobile(s string) bool    { return validatorimpl.IsMobile(s) }
func IsURL(s string) bool       { return validatorimpl.IsURL(s) }
func IsIPv4(s string) bool      { return validatorimpl.IsIPv4(s) }
func IsChinese(s string) bool   { return validatorimpl.IsChinese(s) }
func IsNumberStr(s string) bool { return validatorimpl.IsNumberStr(s) }

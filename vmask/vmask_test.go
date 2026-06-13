package vmask_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vmask"
)

func TestFacadeBuiltInRules(t *testing.T) {
	if got := vmask.Masked("18049531999", vmask.MobilePhoneType); got != "180****1999" {
		t.Fatalf("mobile: %q", got)
	}
	if got := vmask.ChineseName("段正淳"); got != "段**" {
		t.Fatalf("name: %q", got)
	}
	if got := vmask.Email("duandazhi-jack@gmail.com.cn"); got != "d*************@gmail.com.cn" {
		t.Fatalf("email: %q", got)
	}
	if got := vmask.BankCard("11011111222233333256"); got != "1101 **** **** **** 3256" {
		t.Fatalf("bank: %q", got)
	}
	if got := vmask.Passport("PJ1234567"); got != "PJ*****67" {
		t.Fatalf("passport: %q", got)
	}
	if vmask.MaskedPtr("x", vmask.ClearToNullType) != nil {
		t.Fatal("ClearToNullType should return nil")
	}
}

func TestFacadeMaskDispatchAndClearHelpers(t *testing.T) {
	tests := []struct {
		name string
		in   string
		typ  vmask.Type
		want string
	}{
		{name: "user id", in: "12345", typ: vmask.UserID, want: "0"},
		{name: "id card", in: "11010519491231002X", typ: vmask.IDCard, want: "1***************2X"},
		{name: "fixed phone", in: "01012345678", typ: vmask.FixedPhoneType, want: "0101*****78"},
		{name: "address", in: "北京市朝阳区望京街道", typ: vmask.AddressType, want: "北京********"},
		{name: "password", in: "secret", typ: vmask.PasswordType, want: "******"},
		{name: "car license", in: "京A12345", typ: vmask.CarLicenseType, want: "京A1***5"},
		{name: "ipv4", in: "192.168.1.10", typ: vmask.IPv4Type, want: "192.*.*.*"},
		{name: "ipv6", in: "2001:db8::1", typ: vmask.IPv6Type, want: "2001:*:*:*:*:*:*:*:*"},
		{name: "credit code", in: "91350211M000100Y43", typ: vmask.CreditCodeType, want: "9135**********0Y43"},
		{name: "first mask", in: "hello", typ: vmask.FirstMaskType, want: "h****"},
		{name: "clear empty", in: "hello", typ: vmask.ClearToEmptyType, want: ""},
		{name: "unknown type", in: "hello", typ: vmask.Type(999), want: "hello"},
		{name: "blank input", in: "  ", typ: vmask.PasswordType, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vmask.Masked(tt.in, tt.typ); got != tt.want {
				t.Fatalf("Masked(%q, %v) = %q, want %q", tt.in, tt.typ, got, tt.want)
			}
		})
	}

	if got := vmask.Clear(); got != "" {
		t.Fatalf("Clear = %q", got)
	}
	if got := vmask.ClearToNil(); got != nil {
		t.Fatalf("ClearToNil = %v", got)
	}
	if got := vmask.UserIDValue(); got != 0 {
		t.Fatalf("UserIDValue = %d", got)
	}
	ptr := vmask.MaskedPtr("hello", vmask.FirstMaskType)
	if ptr == nil || *ptr != "h****" {
		t.Fatalf("MaskedPtr FirstMask = %v", ptr)
	}
}

func TestFacadeDirectMaskHelpers(t *testing.T) {
	if got := vmask.FirstMask("中文ab"); got != "中***" {
		t.Fatalf("FirstMask = %q", got)
	}
	if got := vmask.IDCardNum("1234567890", 2, 3); got != "12*****890" {
		t.Fatalf("IDCardNum = %q", got)
	}
	if got := vmask.IDCardNum("123", 2, 2); got != "" {
		t.Fatalf("IDCardNum invalid = %q", got)
	}
	if got := vmask.FixedPhone("01012345678"); got != "0101*****78" {
		t.Fatalf("FixedPhone = %q", got)
	}
	if got := vmask.MobilePhone("18049531999"); got != "180****1999" {
		t.Fatalf("MobilePhone = %q", got)
	}
	if got := vmask.Address("北京市朝阳区望京街道", 4); got != "北京市朝阳区****" {
		t.Fatalf("Address = %q", got)
	}
	if got := vmask.Email("a@example.com"); got != "a@example.com" {
		t.Fatalf("Email short local = %q", got)
	}
	if got := vmask.Password("密码ab"); got != "****" {
		t.Fatalf("Password = %q", got)
	}
	if got := vmask.CarLicense("粤B123456"); got != "粤B1****6" {
		t.Fatalf("CarLicense new energy = %q", got)
	}
	if got := vmask.CarLicense("too-long-plate"); got != "too-long-plate" {
		t.Fatalf("CarLicense invalid len = %q", got)
	}
	if got := vmask.BankCard("1234 5678 9012 3456"); got != "1234 **** **** 3456" {
		t.Fatalf("BankCard spaced = %q", got)
	}
	if got := vmask.BankCard("12345678"); got != "12345678" {
		t.Fatalf("BankCard short = %q", got)
	}
	if got := vmask.IPv4("localhost"); got != "localhost.*.*.*" {
		t.Fatalf("IPv4 no dot = %q", got)
	}
	if got := vmask.IPv6("localhost"); got != "localhost:*:*:*:*:*:*:*:*" {
		t.Fatalf("IPv6 no colon = %q", got)
	}
	if got := vmask.Passport("AB"); got != "**" {
		t.Fatalf("Passport short = %q", got)
	}
	if got := vmask.CreditCode("1234"); got != "****" {
		t.Fatalf("CreditCode short = %q", got)
	}
	if got := vmask.Hide("abcdef", -1, 99); got != "******" {
		t.Fatalf("Hide clamped = %q", got)
	}
	if got := vmask.Hide("abcdef", 4, 2); got != "abcdef" {
		t.Fatalf("Hide reversed = %q", got)
	}
}

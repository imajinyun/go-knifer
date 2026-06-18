package mask

import "testing"

func TestBuiltInRules(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"user", Masked("100", UserID), "0"},
		{"name", ChineseName("段正淳"), "段**"},
		{"id", IDCardNum("51343620000320711X", 1, 2), "5***************1X"},
		{"fixed", FixedPhone("09157518479"), "0915*****79"},
		{"mobile", MobilePhone("18049531999"), "180****1999"},
		{"address", Address("北京市海淀区马连洼街道289号", 8), "北京市海淀区马********"},
		{"email", Email("duandazhi-jack@gmail.com.cn"), "d*************@gmail.com.cn"},
		{"password", Password("1234567890"), "**********"},
		{"car7", CarLicense("苏D40000"), "苏D4***0"},
		{"car8", CarLicense("陕A12345D"), "陕A1****D"},
		{"bank", BankCard("11011111222233333256"), "1101 **** **** **** 3256"},
		{"ipv4", IPv4("192.168.1.1"), "192.*.*.*"},
		{"ipv6", IPv6("2001:0db8:86a3:08d3:1319:8a2e:0370:7344"), "2001:*:*:*:*:*:*:*:*"},
		{"passport", Passport("PJ1234567"), "PJ*****67"},
		{"credit", CreditCode("91110108MA01ABCDE7"), "9111**********CDE7"},
		{"first", FirstMask("123456789"), "1********"},
		{"clear", Clear(), ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Fatalf("got %q want %q", tc.got, tc.want)
			}
		})
	}
}

func TestMaskedPtrAndBoundary(t *testing.T) {
	if got := Masked("18049531999", MobilePhoneType); got != "180****1999" {
		t.Fatalf("Masked: %q", got)
	}
	if MaskedPtr("x", ClearToNullType) != nil {
		t.Fatal("ClearToNullType should return nil pointer")
	}
	if got := IDCardNum("123", 2, 2); got != "" {
		t.Fatalf("invalid id mask: %q", got)
	}
	if got := BankCard("1234 5678"); got != "12345678" {
		t.Fatalf("short bank card: %q", got)
	}
}

func TestMaskedDispatch(t *testing.T) {
	cases := []struct {
		name string
		in   string
		typ  Type
		want string
	}{
		{"blank", " \t\n", MobilePhoneType, ""},
		{"user id", "100", UserID, "0"},
		{"chinese name", "段正淳", ChineseNameType, "段**"},
		{"id card", "51343620000320711X", IDCard, "5***************1X"},
		{"fixed phone", "09157518479", FixedPhoneType, "0915*****79"},
		{"mobile phone", "18049531999", MobilePhoneType, "180****1999"},
		{"address", "北京市海淀区马连洼街道289号", AddressType, "北京市海淀区马********"},
		{"email", "duandazhi-jack@gmail.com.cn", EmailType, "d*************@gmail.com.cn"},
		{"password", "1234567890", PasswordType, "**********"},
		{"car license", "陕A12345D", CarLicenseType, "陕A1****D"},
		{"bank card", "11011111222233333256", BankCardType, "1101 **** **** **** 3256"},
		{"ipv4", "192.168.1.1", IPv4Type, "192.*.*.*"},
		{"ipv6", "2001:0db8:86a3:08d3:1319:8a2e:0370:7344", IPv6Type, "2001:*:*:*:*:*:*:*:*"},
		{"passport", "PJ1234567", PassportType, "PJ*****67"},
		{"credit code", "91110108MA01ABCDE7", CreditCodeType, "9111**********CDE7"},
		{"first mask", "123456789", FirstMaskType, "1********"},
		{"clear empty", "secret", ClearToEmptyType, ""},
		{"clear null", "secret", ClearToNullType, ""},
		{"unknown", "keep", Type(99), "keep"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Masked(tc.in, tc.typ); got != tc.want {
				t.Fatalf("Masked(%q, %d) = %q, want %q", tc.in, tc.typ, got, tc.want)
			}
		})
	}
}

func TestPointerAndClearHelpers(t *testing.T) {
	got := MaskedPtr("18049531999", MobilePhoneType)
	if got == nil {
		t.Fatal("MaskedPtr returned nil for non-clear type")
	}
	if *got != "180****1999" {
		t.Fatalf("MaskedPtr value = %q, want %q", *got, "180****1999")
	}
	if ClearToNil() != nil {
		t.Fatal("ClearToNil should return nil")
	}
	if got := UserIDValue(); got != 0 {
		t.Fatalf("UserIDValue() = %d, want 0", got)
	}
}

func TestMaskBoundaries(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"chinese name blank", ChineseName(" \t"), ""},
		{"id card negative front", IDCardNum("123456", -1, 1), ""},
		{"id card negative end", IDCardNum("123456", 1, -1), ""},
		{"fixed phone short", FixedPhone("12345"), "12345"},
		{"mobile phone short", MobilePhone("123456"), "123456"},
		{"address shorter than sensitive size", Address("短址", 8), "**"},
		{"email blank", Email(" \t"), ""},
		{"email one character local part", Email("a@example.com"), "a@example.com"},
		{"email without at", Email("plain"), "plain"},
		{"password blank", Password(" \t"), ""},
		{"car license unsupported length", CarLicense("ABC"), "ABC"},
		{"bank card blank", BankCard(" \t"), " \t"},
		{"bank card non multiple group", BankCard("1234 5678 9"), "1234 **** 9"},
		{"ipv4 no delimiter", IPv4("localhost"), "localhost.*.*.*"},
		{"ipv6 no delimiter", IPv6("node"), "node:*:*:*:*:*:*:*:*"},
		{"passport blank", Passport(" \t"), " \t"},
		{"passport short", Passport("AB"), "**"},
		{"credit code blank", CreditCode(" \t"), " \t"},
		{"credit code short", CreditCode("ABCD"), "****"},
		{"first mask blank", FirstMask(" \t"), ""},
		{"hide negative start", Hide("敏感数据", -2, 2), "**数据"},
		{"hide oversized end", Hide("abcd", 1, 10), "a***"},
		{"hide start after end", Hide("abcd", 3, 1), "abcd"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Fatalf("got %q want %q", tc.got, tc.want)
			}
		})
	}
}

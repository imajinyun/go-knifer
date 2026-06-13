package bloomfilter

import "testing"

func TestHashAlgorithms(t *testing.T) {
	s := "test-bloomFilter"
	checks := map[string]int32{
		"rs":   RsHash(s),
		"js":   JsHash(s),
		"pjw":  PjwHash(s),
		"elf":  ElfHash(s),
		"bkdr": BkdrHash(s),
		"sdbm": SdbmHash(s),
		"djb":  DjbHash(s),
		"ap":   ApHash(s),
		"fnv":  FnvHashString(s),
	}
	for name, v := range checks {
		// Only verify stability: the same string should produce the same result.
		if v != checkAgain(name, s) {
			t.Fatalf("%s is unstable", name)
		}
	}
	if JavaDefaultHash("a") != 97 {
		t.Fatal("javaDefault hash 'a' should be 97")
	}
	if TianlHash("") != 0 {
		t.Fatal("tianl empty should be 0")
	}
}

// checkAgain runs the same algorithm again for stability tests.
func checkAgain(name, s string) int32 {
	switch name {
	case "rs":
		return RsHash(s)
	case "js":
		return JsHash(s)
	case "pjw":
		return PjwHash(s)
	case "elf":
		return ElfHash(s)
	case "bkdr":
		return BkdrHash(s)
	case "sdbm":
		return SdbmHash(s)
	case "djb":
		return DjbHash(s)
	case "ap":
		return ApHash(s)
	case "fnv":
		return FnvHashString(s)
	}
	return 0
}

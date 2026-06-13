package system

import "testing"

func TestResetInfoCacheClearsSingletons(t *testing.T) {
	ResetInfoCache()
	firstHost := GetHostInfo()
	firstOS := GetOsInfo()
	firstUser := GetUserInfo()
	firstGo := GetGoInfo()
	firstRuntime := GetRuntimeInfo()
	ResetInfoCache()
	if got := GetHostInfo(); got == nil || got == firstHost {
		t.Fatalf("GetHostInfo after reset = %p, first %p", got, firstHost)
	}
	if got := GetOsInfo(); got == nil || got == firstOS {
		t.Fatalf("GetOsInfo after reset = %p, first %p", got, firstOS)
	}
	if got := GetUserInfo(); got == nil || got == firstUser {
		t.Fatalf("GetUserInfo after reset = %p, first %p", got, firstUser)
	}
	if got := GetGoInfo(); got == nil || got == firstGo {
		t.Fatalf("GetGoInfo after reset = %p, first %p", got, firstGo)
	}
	if got := GetRuntimeInfo(); got == nil || got == firstRuntime {
		t.Fatalf("GetRuntimeInfo after reset = %p, first %p", got, firstRuntime)
	}
}

package id

import "testing"

func TestNormalizeSnowflakeIDUsesProvidedMax(t *testing.T) {
	if got := normalizeSnowflakeID(14, 5); got != 2 {
		t.Fatalf("normalizeSnowflakeID should use provided max: got %d", got)
	}
	if got := normalizeSnowflakeID(-14, 5); got != 2 {
		t.Fatalf("normalizeSnowflakeID should normalize negative values with provided max: got %d", got)
	}
	if got := normalizeSnowflakeID(14, 0); got != 0 {
		t.Fatalf("normalizeSnowflakeID should return 0 when max is not positive: got %d", got)
	}
}

func TestWorkerAndDatacenterID(t *testing.T) {
	if dc := GetDataCenterID(31); dc < 0 || dc > 31 {
		t.Fatalf("datacenter id out of range: %d", dc)
	}
	if worker := GetWorkerID(1, 31); worker < 0 || worker > 31 {
		t.Fatalf("worker id out of range: %d", worker)
	}
}

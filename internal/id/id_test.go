package id

import (
	"bytes"
	"encoding/hex"
	"errors"
	mathrand "math/rand"
	"net"
	"strings"
	"testing"
	"time"
)

func TestSimpleUUID(t *testing.T) {
	u1 := SimpleUUID()
	u2 := SimpleUUID()
	if len(u1) != 32 || len(u2) != 32 {
		t.Fatalf("UUID length wrong")
	}
	if u1 == u2 {
		t.Fatalf("UUID collision")
	}
	// Version 4 marker: the 13th character is '4'.
	if u1[12] != '4' {
		t.Fatalf("UUID version: %s", u1)
	}
}

func TestRandomUUIDAndFastSimpleUUID(t *testing.T) {
	u := RandomUUID()
	if len(u) != 36 || strings.Count(u, "-") != 4 {
		t.Fatalf("RandomUUID format: %s", u)
	}
	s := FastSimpleUUID()
	if len(s) != 32 || strings.Contains(s, "-") || s[12] != '4' {
		t.Fatalf("FastSimpleUUID format: %s", s)
	}
}

func TestFastUUID(t *testing.T) {
	u := FastUUID()
	if len(u) != 36 || strings.Count(u, "-") != 4 {
		t.Fatalf("FastUUID format: %s", u)
	}
}

func TestObjectId(t *testing.T) {
	o := ObjectId()
	if len(o) != 24 {
		t.Fatalf("ObjectId length: %s", o)
	}
}

func TestNanoId(t *testing.T) {
	id := NanoId()
	if len(id) != 21 {
		t.Fatalf("NanoId default len: %s", id)
	}
	id = NanoIdN(10)
	if len(id) != 10 {
		t.Fatalf("NanoIdN len: %s", id)
	}
}

func TestSnowflake(t *testing.T) {
	sf := CreateSnowflake(1, 2)
	if sf.WorkerID() != 1 || sf.DatacenterID() != 2 {
		t.Fatalf("snowflake ids: worker=%d datacenter=%d", sf.WorkerID(), sf.DatacenterID())
	}
	id1 := sf.NextID()
	id2 := sf.NextID()
	if id1 <= 0 || id2 <= id1 {
		t.Fatalf("snowflake should be positive and increasing: %d %d", id1, id2)
	}
	if sf.NextIDStr() == "" {
		t.Fatal("snowflake string id should not be empty")
	}
	first := GetSnowflakeWithWorkerDataCenter(1, 2)
	second := GetSnowflakeWithWorkerDataCenter(1, 2)
	if first != second {
		t.Fatal("same worker/datacenter pair should return singleton")
	}
	if GetSnowflakeWithWorker(3) == nil || GetSnowflake() == nil {
		t.Fatal("snowflake singleton helpers should not return nil")
	}
	if GetSnowflakeNextID() <= 0 || GetSnowflakeNextIDStr() == "" {
		t.Fatal("default snowflake next id helpers failed")
	}
}

func TestSnowflakeOptions(t *testing.T) {
	now := int64(1288834974657)
	sf := CreateSnowflakeWithOptions(
		WithSnowflakeWorkerID(3),
		WithSnowflakeDatacenterID(4),
		WithSnowflakeTimeFunc(func() int64 { return now }),
	)
	if sf.WorkerID() != 3 || sf.DatacenterID() != 4 {
		t.Fatalf("snowflake option ids: worker=%d datacenter=%d", sf.WorkerID(), sf.DatacenterID())
	}
	id1 := sf.NextID()
	id2 := sf.NextID()
	if id1 <= 0 || id2 <= id1 {
		t.Fatalf("snowflake option IDs should be positive and increasing: %d %d", id1, id2)
	}
}

func TestDefaultSnowflakeOptions(t *testing.T) {
	t.Cleanup(func() { ConfigureDefaultSnowflake() })
	now := int64(1288834974657)
	sf := ConfigureDefaultSnowflake(
		WithSnowflakeWorkerID(5),
		WithSnowflakeDatacenterID(6),
		WithSnowflakeTimeFunc(func() int64 { return now }),
	)
	if sf.WorkerID() != 5 || sf.DatacenterID() != 6 {
		t.Fatalf("default snowflake option ids: worker=%d datacenter=%d", sf.WorkerID(), sf.DatacenterID())
	}
	if GetSnowflake() != sf || GetSnowflakeWithOptions(WithSnowflakeWorkerID(7)) != sf {
		t.Fatal("default snowflake singleton should keep configured instance")
	}
	first := sf.NextID()
	second := GetSnowflakeNextID()
	if first <= 0 || second <= first {
		t.Fatalf("configured default snowflake should generate increasing ids: %d %d", first, second)
	}
	if got := GetSnowflakeNextIDStr(); got == "" {
		t.Fatal("configured default snowflake string id should not be empty")
	}
}

func TestSnowflakeRuntimeOptionsBypassSingletonCache(t *testing.T) {
	t.Cleanup(func() { ConfigureDefaultSnowflake() })
	configured := ConfigureDefaultSnowflake(WithSnowflakeWorkerID(1), WithSnowflakeDatacenterID(1))

	now := int64(1288834974657)
	one := GetSnowflakeWithOptions(
		WithSnowflakeWorkerID(2),
		WithSnowflakeDatacenterID(3),
		WithSnowflakeTimeFunc(func() int64 { return now }),
	)
	two := GetSnowflakeWithOptions(
		WithSnowflakeWorkerID(2),
		WithSnowflakeDatacenterID(3),
		WithSnowflakeTimeFunc(func() int64 { return now }),
	)
	if one == configured || two == configured || one == two {
		t.Fatalf("runtime options should bypass default singleton/cache: configured=%p one=%p two=%p", configured, one, two)
	}
	if one.WorkerID() != 2 || one.DatacenterID() != 3 {
		t.Fatalf("runtime options ids = worker %d datacenter %d", one.WorkerID(), one.DatacenterID())
	}
}

func TestSnowflakeCacheOptionRetainsSingletonBehavior(t *testing.T) {
	now := int64(1288834974657)
	one := GetSnowflakeWithWorkerDataCenterWithOptions(9, 10,
		WithSnowflakeTimeFunc(func() int64 { return now }),
		WithSnowflakeCache(true),
	)
	two := GetSnowflakeWithWorkerDataCenterWithOptions(9, 10,
		WithSnowflakeTimeFunc(func() int64 { return now }),
		WithSnowflakeCache(true),
	)
	if one != two {
		t.Fatal("explicit cache option should retain singleton behavior")
	}

	isolated := GetSnowflakeWithWorkerDataCenterWithOptions(9, 10, WithSnowflakeCache(false))
	if isolated == one {
		t.Fatal("WithSnowflakeCache(false) should bypass cached worker/datacenter generator")
	}
}

func TestNewIsolatedSnowflake(t *testing.T) {
	t.Cleanup(func() { ConfigureDefaultSnowflake() })
	configured := ConfigureDefaultSnowflake(WithSnowflakeWorkerID(1), WithSnowflakeDatacenterID(1))
	isolated := NewIsolatedSnowflake(WithSnowflakeWorkerID(4), WithSnowflakeDatacenterID(5))
	if isolated == configured || isolated.WorkerID() != 4 || isolated.DatacenterID() != 5 {
		t.Fatalf("isolated snowflake = %p worker %d datacenter %d", isolated, isolated.WorkerID(), isolated.DatacenterID())
	}
}

func TestDefaultSnowflakeDerivedProviderOptions(t *testing.T) {
	t.Cleanup(func() { ConfigureDefaultSnowflake() })
	interfacesCalled := false
	pidCalled := false
	sf := ConfigureDefaultSnowflake(
		WithSnowflakeInterfacesFunc(func() ([]net.Interface, error) {
			interfacesCalled = true
			return []net.Interface{{HardwareAddr: net.HardwareAddr{0, 0, 0, 0, 0x40, 0}}}, nil
		}),
		WithSnowflakePIDFunc(func() int {
			pidCalled = true
			return 12345
		}),
		WithSnowflakeTimeFunc(func() int64 { return 1288834974657 }),
	)
	if !interfacesCalled || !pidCalled {
		t.Fatalf("derived providers called interfaces=%v pid=%v", interfacesCalled, pidCalled)
	}
	if sf.DatacenterID() != 1 {
		t.Fatalf("derived datacenter id = %d, want 1", sf.DatacenterID())
	}
	if want := getWorkerID(1, 31, func() int { return 12345 }); sf.WorkerID() != want {
		t.Fatalf("derived worker id = %d, want %d", sf.WorkerID(), want)
	}
	if sf.WorkerID() < 0 || sf.WorkerID() > 31 {
		t.Fatalf("derived worker id out of range: %d", sf.WorkerID())
	}
	if GetSnowflake() != sf {
		t.Fatal("configured default snowflake should be installed")
	}
}

func TestSnowflakeDerivedProviderOptionsDoNotOverrideExplicitIDs(t *testing.T) {
	sf := CreateSnowflakeWithOptions(
		WithSnowflakeWorkerID(7),
		WithSnowflakeDatacenterID(8),
		WithSnowflakeInterfacesFunc(func() ([]net.Interface, error) {
			t.Fatal("interfaces provider should not be used when datacenter id is explicit")
			return nil, nil
		}),
		WithSnowflakePIDFunc(func() int {
			t.Fatal("pid provider should not be used when worker id is explicit")
			return 0
		}),
	)
	if sf.WorkerID() != 7 || sf.DatacenterID() != 8 {
		t.Fatalf("explicit ids = worker %d datacenter %d", sf.WorkerID(), sf.DatacenterID())
	}
}

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

func TestIDOptions(t *testing.T) {
	reader := bytes.NewReader(bytes.Repeat([]byte{0x11}, 32))
	u := SimpleUUIDWithOptions(WithRandomReader(reader))
	if len(u) != 32 || u[12] != '4' || u[16] != '9' {
		t.Fatalf("SimpleUUIDWithOptions format: %s", u)
	}

	obj := ObjectIdWithOptions(
		WithObjectIDTimeFunc(func() time.Time { return time.Unix(1, 0) }),
		WithObjectIDRandomReader(bytes.NewReader([]byte{1, 2, 3, 4, 5})),
		WithObjectIDCounter(func() uint32 { return 0xabcdef }),
	)
	if obj != "000000010102030405abcdef" {
		t.Fatalf("ObjectIdWithOptions = %s", obj)
	}
	if _, err := hex.DecodeString(obj); err != nil {
		t.Fatalf("ObjectIdWithOptions is not hex: %v", err)
	}

	nid := NanoIdWithOptions(
		WithNanoIDLength(5),
		WithNanoIDAlphabet("ab"),
		WithNanoIDRandomReader(bytes.NewReader([]byte{0, 1, 0, 1, 1})),
	)
	if nid != "ababb" {
		t.Fatalf("NanoIdWithOptions = %q", nid)
	}
}

func TestDefaultFallbackRandomSourceProviderCanBeConfiguredAndReset(t *testing.T) {
	ResetDefaultFallbackRandomSource()
	t.Cleanup(ResetDefaultFallbackRandomSource)

	ConfigureDefaultFallbackRandomSourceProvider(func() *mathrand.Rand {
		return mathrand.New(mathrand.NewSource(11))
	})
	first := SimpleUUIDWithOptions(WithRandomReader(errReader{}))
	second := SimpleUUIDWithOptions(WithRandomReader(errReader{}))
	ConfigureDefaultFallbackRandomSourceProvider(func() *mathrand.Rand {
		return mathrand.New(mathrand.NewSource(11))
	})
	if got := SimpleUUIDWithOptions(WithRandomReader(errReader{})); got != first {
		t.Fatalf("SimpleUUIDWithOptions after provider reset = %s, want %s", got, first)
	}
	if got := SimpleUUIDWithOptions(WithRandomReader(errReader{})); got != second {
		t.Fatalf("second SimpleUUIDWithOptions after provider reset = %s, want %s", got, second)
	}

	SetFallbackRandomSeed(12)
	seeded := SimpleUUIDWithOptions(WithRandomReader(errReader{}))
	SetFallbackRandomSeed(12)
	if got := SimpleUUIDWithOptions(WithRandomReader(errReader{})); got != seeded {
		t.Fatalf("SimpleUUIDWithOptions after seed reset = %s, want %s", got, seeded)
	}

	ConfigureDefaultFallbackRandomSourceProvider(func() *mathrand.Rand { return nil })
	if got := SimpleUUIDWithOptions(WithRandomReader(errReader{})); len(got) != 32 || got[12] != '4' {
		t.Fatalf("nil provider fallback UUID = %s", got)
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errors.New("boom") }

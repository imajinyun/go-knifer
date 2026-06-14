package id

import (
	"net"
	"testing"
)

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

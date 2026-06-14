package id

import "testing"

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

package vid

import (
	"io"
	"time"

	idimpl "github.com/imajinyun/go-knifer/internal/id"
)

type (
	Snowflake      = idimpl.Snowflake
	RandomOption   = idimpl.RandomOption
	ObjectIDOption = idimpl.ObjectIDOption
	NanoIDOption   = idimpl.NanoIDOption
)

func RandomUUID() string     { return idimpl.RandomUUID() }
func SimpleUUID() string     { return idimpl.SimpleUUID() }
func FastUUID() string       { return idimpl.FastUUID() }
func FastSimpleUUID() string { return idimpl.FastSimpleUUID() }
func UUID() string           { return idimpl.SimpleUUID() }
func ObjectId() string       { return idimpl.ObjectId() }

func WithRandomReader(reader io.Reader) RandomOption { return idimpl.WithRandomReader(reader) }

func RandomUUIDWithOptions(opts ...RandomOption) string { return idimpl.RandomUUIDWithOptions(opts...) }

func SimpleUUIDWithOptions(opts ...RandomOption) string { return idimpl.SimpleUUIDWithOptions(opts...) }

func WithObjectIDRandomReader(reader io.Reader) ObjectIDOption {
	return idimpl.WithObjectIDRandomReader(reader)
}

func WithObjectIDTimeFunc(now func() time.Time) ObjectIDOption {
	return idimpl.WithObjectIDTimeFunc(now)
}

func WithObjectIDCounter(counter func() uint32) ObjectIDOption {
	return idimpl.WithObjectIDCounter(counter)
}

func ObjectIdWithOptions(opts ...ObjectIDOption) string { return idimpl.ObjectIdWithOptions(opts...) }

func CreateSnowflake(workerID, datacenterID int64) *Snowflake {
	return idimpl.CreateSnowflake(workerID, datacenterID)
}

func GetSnowflake() *Snowflake { return idimpl.GetSnowflake() }

func GetSnowflakeWithWorker(workerID int64) *Snowflake {
	return idimpl.GetSnowflakeWithWorker(workerID)
}

func GetSnowflakeWithWorkerDataCenter(workerID, datacenterID int64) *Snowflake {
	return idimpl.GetSnowflakeWithWorkerDataCenter(workerID, datacenterID)
}

func GetDataCenterID(maxDatacenterID int64) int64 { return idimpl.GetDataCenterID(maxDatacenterID) }
func GetWorkerID(datacenterID, maxWorkerID int64) int64 {
	return idimpl.GetWorkerID(datacenterID, maxWorkerID)
}

func NanoId() string       { return idimpl.NanoId() }
func NanoIdN(n int) string { return idimpl.NanoIdN(n) }

func WithNanoIDRandomReader(reader io.Reader) NanoIDOption {
	return idimpl.WithNanoIDRandomReader(reader)
}

func WithNanoIDAlphabet(alphabet string) NanoIDOption { return idimpl.WithNanoIDAlphabet(alphabet) }

func WithNanoIDLength(length int) NanoIDOption { return idimpl.WithNanoIDLength(length) }

func NanoIdWithOptions(opts ...NanoIDOption) string { return idimpl.NanoIdWithOptions(opts...) }

func GetSnowflakeNextID() int64     { return idimpl.GetSnowflakeNextID() }
func GetSnowflakeNextIDStr() string { return idimpl.GetSnowflakeNextIDStr() }

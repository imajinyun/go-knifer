package sets

// Int is a set of int values.
type Int = Set[int]

// Int32 is a set of int32 values.
type Int32 = Set[int32]

// Int64 is a set of int64 values.
type Int64 = Set[int64]

// Uint is a set of uint values.
type Uint = Set[uint]

// Uint32 is a set of uint32 values.
type Uint32 = Set[uint32]

// Uint64 is a set of uint64 values.
type Uint64 = Set[uint64]

// String is a set of string values.
type String = Set[string]

// NewInt creates an int set.
func NewInt(items ...int) Int { return New(items...) }

// NewInt32 creates an int32 set.
func NewInt32(items ...int32) Int32 { return New(items...) }

// NewInt64 creates an int64 set.
func NewInt64(items ...int64) Int64 { return New(items...) }

// NewUint creates a uint set.
func NewUint(items ...uint) Uint { return New(items...) }

// NewUint32 creates a uint32 set.
func NewUint32(items ...uint32) Uint32 { return New(items...) }

// NewUint64 creates a uint64 set.
func NewUint64(items ...uint64) Uint64 { return New(items...) }

// NewString creates a string set.
func NewString(items ...string) String { return New(items...) }

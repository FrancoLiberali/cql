package condition

import (
	"database/sql"
	"sync"

	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/model"
)

// Per-type sync.Pools of scan destinations. Generated ScanValues acquires
// wrappers via AcquireX; the runtime returns them via ReleaseX after
// AssignValues has copied their contents into the destination model.
//
// This mirrors gorm's schema.normalPool pattern (gorm/schema/pool.go) — the
// goal is the same: amortize per-cell allocation cost across all rows of
// the result set. Pools are pointer-safe because sql.Null* are values and
// AssignValues copies them out by value before Release is called.
var (
	nullBoolPool    = sync.Pool{New: func() any { return new(sql.NullBool) }}
	nullStringPool  = sync.Pool{New: func() any { return new(sql.NullString) }}
	nullInt16Pool   = sync.Pool{New: func() any { return new(sql.NullInt16) }}
	nullInt32Pool   = sync.Pool{New: func() any { return new(sql.NullInt32) }}
	nullInt64Pool   = sync.Pool{New: func() any { return new(sql.NullInt64) }}
	nullBytePool    = sync.Pool{New: func() any { return new(sql.NullByte) }}
	nullFloat64Pool = sync.Pool{New: func() any { return new(sql.NullFloat64) }}
	nullTimePool    = sync.Pool{New: func() any { return new(sql.NullTime) }}
	uuidPool        = sync.Pool{New: func() any { return new(model.UUID) }}
	deletedAtPool   = sync.Pool{New: func() any { return new(gorm.DeletedAt) }}
	nullSinkPool    = sync.Pool{New: func() any { return new(NullSink) }}
)

// AcquireNullBool returns a zeroed *sql.NullBool from the pool. Generated
// ScanValues calls this in place of `new(sql.NullBool)`.
func AcquireNullBool() *sql.NullBool {
	v, _ := nullBoolPool.Get().(*sql.NullBool)
	*v = sql.NullBool{}

	return v
}

// ReleaseNullBool returns the wrapper to its pool. Generated ReleaseValues
// calls this after AssignValues has copied the cell's value into the model.
// Safe because sql.Null* are pure values — AssignValues does `dest.X =
// v.Bool`, leaving no pointer to v inside dest.
func ReleaseNullBool(v *sql.NullBool) { nullBoolPool.Put(v) }

func AcquireNullString() *sql.NullString {
	v, _ := nullStringPool.Get().(*sql.NullString)
	*v = sql.NullString{}

	return v
}

func ReleaseNullString(v *sql.NullString) { nullStringPool.Put(v) }

func AcquireNullInt16() *sql.NullInt16 {
	v, _ := nullInt16Pool.Get().(*sql.NullInt16)
	*v = sql.NullInt16{}

	return v
}

func ReleaseNullInt16(v *sql.NullInt16) { nullInt16Pool.Put(v) }

func AcquireNullInt32() *sql.NullInt32 {
	v, _ := nullInt32Pool.Get().(*sql.NullInt32)
	*v = sql.NullInt32{}

	return v
}

func ReleaseNullInt32(v *sql.NullInt32) { nullInt32Pool.Put(v) }

func AcquireNullInt64() *sql.NullInt64 {
	v, _ := nullInt64Pool.Get().(*sql.NullInt64)
	*v = sql.NullInt64{}

	return v
}

func ReleaseNullInt64(v *sql.NullInt64) { nullInt64Pool.Put(v) }

func AcquireNullByte() *sql.NullByte {
	v, _ := nullBytePool.Get().(*sql.NullByte)
	*v = sql.NullByte{}

	return v
}

func ReleaseNullByte(v *sql.NullByte) { nullBytePool.Put(v) }

func AcquireNullFloat64() *sql.NullFloat64 {
	v, _ := nullFloat64Pool.Get().(*sql.NullFloat64)
	*v = sql.NullFloat64{}

	return v
}

func ReleaseNullFloat64(v *sql.NullFloat64) { nullFloat64Pool.Put(v) }

func AcquireNullTime() *sql.NullTime {
	v, _ := nullTimePool.Get().(*sql.NullTime)
	*v = sql.NullTime{}

	return v
}

func ReleaseNullTime(v *sql.NullTime) { nullTimePool.Put(v) }

func AcquireUUID() *model.UUID {
	v, _ := uuidPool.Get().(*model.UUID)
	*v = model.NilUUID

	return v
}

// ReleaseUUID returns the wrapper. Safe because AssignValues copies the
// UUID by value (`dest.ID = *v`) — a [16]byte array — before Release runs.
func ReleaseUUID(v *model.UUID) { uuidPool.Put(v) }

// AcquireDeletedAt returns a zeroed *gorm.DeletedAt. NOTE: gorm.DeletedAt
// contains a Time and a Valid bool — pure value, safe to pool.
func AcquireDeletedAt() *gorm.DeletedAt {
	v, _ := deletedAtPool.Get().(*gorm.DeletedAt)
	*v = gorm.DeletedAt{}

	return v
}

func ReleaseDeletedAt(v *gorm.DeletedAt) { deletedAtPool.Put(v) }

// AcquireNullSink returns a NullSink (stateless — no reset needed).
func AcquireNullSink() *NullSink {
	v, _ := nullSinkPool.Get().(*NullSink)
	return v
}

// ReleaseNullSink returns the sink to its pool.
func ReleaseNullSink(v *NullSink) { nullSinkPool.Put(v) }

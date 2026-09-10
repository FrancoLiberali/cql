package condition

import (
	"database/sql"
	"sync"

	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/model"
)

// scanPool is a sync.Pool of *T scan destinations. Generated ScanValues
// acquires a wrapper via one of the AcquireX helpers below; the runtime returns
// it via the matching ReleaseX after AssignValues has copied its contents into
// the destination model.
//
// This mirrors gorm's schema.normalPool pattern (gorm/schema/pool.go): amortize
// per-cell allocation cost across all rows of the result set. Pooling is
// pointer-safe because the pooled types are pure values and AssignValues copies
// them out by value (e.g. dest.X = v.Bool, dest.ID = *uuid) before Release runs.
type scanPool[T any] struct {
	p sync.Pool
}

func newScanPool[T any]() *scanPool[T] {
	return &scanPool[T]{p: sync.Pool{New: func() any { return new(T) }}}
}

// acquire returns a zeroed *T from the pool (the zero value of every pooled
// type is its "empty" state, including model.NilUUID for model.UUID).
func (pl *scanPool[T]) acquire() *T {
	var zero T

	v, _ := pl.p.Get().(*T)
	*v = zero

	return v
}

// release returns v to the pool. It takes any (rather than *T) so the generated
// ReleaseValues can hand it the raw values[i] cell without a per-column type
// assertion — the assertion lives here, once, instead of in every scanner.
func (pl *scanPool[T]) release(v any) {
	if t, ok := v.(*T); ok {
		pl.p.Put(t)
	}
}

// One pool per scanned cell type. Per-type static pools are required: a single
// generic Acquire[T]() can't dispatch to the right pool without a per-call
// reflect lookup, which would land on the per-row hot path.
var (
	nullBoolPool    = newScanPool[sql.NullBool]()
	nullStringPool  = newScanPool[sql.NullString]()
	nullInt16Pool   = newScanPool[sql.NullInt16]()
	nullInt32Pool   = newScanPool[sql.NullInt32]()
	nullInt64Pool   = newScanPool[sql.NullInt64]()
	nullBytePool    = newScanPool[sql.NullByte]()
	nullFloat64Pool = newScanPool[sql.NullFloat64]()
	nullTimePool    = newScanPool[sql.NullTime]()
	uuidPool        = newScanPool[model.UUID]()
	deletedAtPool   = newScanPool[gorm.DeletedAt]()
	nullSinkPool    = newScanPool[NullSink]()
)

// The generated scanners call these by name (Acquire<Type> / Release<Type>) in
// place of new(<Type>).

func AcquireNullBool() *sql.NullBool       { return nullBoolPool.acquire() }
func ReleaseNullBool(v any)                { nullBoolPool.release(v) }
func AcquireNullString() *sql.NullString   { return nullStringPool.acquire() }
func ReleaseNullString(v any)              { nullStringPool.release(v) }
func AcquireNullInt16() *sql.NullInt16     { return nullInt16Pool.acquire() }
func ReleaseNullInt16(v any)               { nullInt16Pool.release(v) }
func AcquireNullInt32() *sql.NullInt32     { return nullInt32Pool.acquire() }
func ReleaseNullInt32(v any)               { nullInt32Pool.release(v) }
func AcquireNullInt64() *sql.NullInt64     { return nullInt64Pool.acquire() }
func ReleaseNullInt64(v any)               { nullInt64Pool.release(v) }
func AcquireNullByte() *sql.NullByte       { return nullBytePool.acquire() }
func ReleaseNullByte(v any)                { nullBytePool.release(v) }
func AcquireNullFloat64() *sql.NullFloat64 { return nullFloat64Pool.acquire() }
func ReleaseNullFloat64(v any)             { nullFloat64Pool.release(v) }
func AcquireNullTime() *sql.NullTime       { return nullTimePool.acquire() }
func ReleaseNullTime(v any)                { nullTimePool.release(v) }
func AcquireUUID() *model.UUID             { return uuidPool.acquire() }
func ReleaseUUID(v any)                    { uuidPool.release(v) }
func AcquireDeletedAt() *gorm.DeletedAt    { return deletedAtPool.acquire() }
func ReleaseDeletedAt(v any)               { deletedAtPool.release(v) }
func AcquireNullSink() *NullSink           { return nullSinkPool.acquire() }
func ReleaseNullSink(v any)                { nullSinkPool.release(v) }

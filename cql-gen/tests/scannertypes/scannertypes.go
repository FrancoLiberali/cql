package scannertypes

import (
	"database/sql"
	"time"

	"github.com/FrancoLiberali/cql/model"
)

// ScannerTypes is a scanner kitchen-sink: it exercises the fast-scan generator
// paths that the narrower fixtures don't reach — a []byte column, a time.Time
// value and pointer, and database/sql nullable wrappers used directly as
// fields — alongside a couple of plain basic value/pointer columns.
type ScannerTypes struct {
	model.UUIDModel

	Int    int
	Float  float64
	String string

	PtrInt *int

	Time    time.Time
	PtrTime *time.Time

	Blob []byte

	NullString sql.NullString
	NullInt64  sql.NullInt64
	NullTime   sql.NullTime
}

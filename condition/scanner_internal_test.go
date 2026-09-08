package condition

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/model"
)

// TestValueHasData exercises every branch of valueHasData — the per-cell
// non-NULL check allValuesNull uses to detect LEFT-JOIN-no-match rows. Through
// the query path it is only reached via joined preloads, whose child models
// don't cover every scanned cell type, so it is unit-tested directly here.
func TestValueHasData(t *testing.T) {
	t.Parallel()

	now := time.Now()
	nilUUID := model.NilUUID
	nonNilUUID := model.UUID{1}

	tests := []struct {
		name string
		v    any
		want bool
	}{
		{"NullBool NULL", &sql.NullBool{}, false},
		{"NullBool set", &sql.NullBool{Bool: true, Valid: true}, true},
		{"NullString NULL", &sql.NullString{}, false},
		{"NullString set", &sql.NullString{String: "x", Valid: true}, true},
		{"NullInt64 NULL", &sql.NullInt64{}, false},
		{"NullInt64 set", &sql.NullInt64{Int64: 1, Valid: true}, true},
		{"NullInt32 set", &sql.NullInt32{Int32: 1, Valid: true}, true},
		{"NullInt16 set", &sql.NullInt16{Int16: 1, Valid: true}, true},
		{"NullByte set", &sql.NullByte{Byte: 1, Valid: true}, true},
		{"NullFloat64 set", &sql.NullFloat64{Float64: 1, Valid: true}, true},
		{"NullTime NULL", &sql.NullTime{}, false},
		{"NullTime set", &sql.NullTime{Time: now, Valid: true}, true},
		{"DeletedAt NULL", &gorm.DeletedAt{}, false},
		{"DeletedAt set", &gorm.DeletedAt{Time: now, Valid: true}, true},
		{"NullSink is always NULL", &NullSink{}, false},
		{"NullableScanner is conservative (has data)", &NullableScanner{Inner: &NullSink{}}, true},
		{"UUID nil sentinel", &nilUUID, false},
		{"UUID non-nil", &nonNilUUID, true},
		{"unknown type assumed to have data", &struct{ X int }{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, valueHasData(tt.v))
		})
	}
}

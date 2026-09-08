package test

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/models"
)

type ScannerTypesIntTestSuite struct {
	testSuite
}

func NewScannerTypesIntTestSuite(
	db *cql.DB,
) *ScannerTypesIntTestSuite {
	return &ScannerTypesIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

// TestFastScanRoundTripsEveryScalarTypeWithValues inserts a row with every
// supported scalar field populated — including each integer size's boundary
// value and non-nil pointers — then reads it back through the fast scanner and
// asserts every value round-trips. Distinct values per field also catch any
// column/field mis-wiring in the generated scanner.
func (ts *ScannerTypesIntTestSuite) TestFastScanRoundTripsEveryScalarTypeWithValues() {
	// second precision + no sub-second so every dialect's datetime column
	// preserves it; compared by instant below.
	valTime := time.Date(2021, 3, 14, 15, 9, 26, 0, time.UTC)

	i, i8, i16, i32, i64 := -7, int8(-8), int16(-16), int32(-32), int64(-64)
	u, u8, u16, u32, u64 := uint(7), uint8(8), uint16(16), uint32(32), uint64(64)
	f32, f64 := float32(1.5), 6.25
	bp := true
	sp := "ptr"
	tp := valTime

	in := &models.AllTypes{
		ValInt:     -42,
		ValInt8:    math.MaxInt8,
		ValInt16:   math.MaxInt16,
		ValInt32:   math.MaxInt32,
		ValInt64:   math.MaxInt64,
		ValUint:    42,
		ValUint8:   math.MaxUint8,
		ValUint16:  math.MaxUint16,
		ValUint32:  math.MaxUint32,
		ValUint64:  math.MaxInt64, // largest uint64 that still fits signed storage
		ValFloat32: 3.5,
		ValFloat64: 2.718281828459045,
		ValBool:    true,
		ValString:  "kitchen sink",
		ValTime:    valTime,

		PtrInt:     &i,
		PtrInt8:    &i8,
		PtrInt16:   &i16,
		PtrInt32:   &i32,
		PtrInt64:   &i64,
		PtrUint:    &u,
		PtrUint8:   &u8,
		PtrUint16:  &u16,
		PtrUint32:  &u32,
		PtrUint64:  &u64,
		PtrFloat32: &f32,
		PtrFloat64: &f64,
		PtrBool:    &bp,
		PtrString:  &sp,
		PtrTime:    &tp,

		NullInt16: sql.NullInt16{Int16: 1616, Valid: true},
		NullInt32: sql.NullInt32{Int32: 323232, Valid: true},
		NullByte:  sql.NullByte{Byte: 200, Valid: true},
	}

	create(&ts.testSuite, in)

	// no conditions: the scanner is resolved by type from the registry, so
	// this still takes the fast-scan path.
	got, err := cql.Query[models.AllTypes](
		context.Background(),
		ts.db,
	).First()
	ts.Require().NoError(err)

	// value fields
	ts.Equal(in.ValInt, got.ValInt)
	ts.Equal(in.ValInt8, got.ValInt8)
	ts.Equal(in.ValInt16, got.ValInt16)
	ts.Equal(in.ValInt32, got.ValInt32)
	ts.Equal(in.ValInt64, got.ValInt64)
	ts.Equal(in.ValUint, got.ValUint)
	ts.Equal(in.ValUint8, got.ValUint8)
	ts.Equal(in.ValUint16, got.ValUint16)
	ts.Equal(in.ValUint32, got.ValUint32)
	ts.Equal(in.ValUint64, got.ValUint64)
	ts.InDelta(in.ValFloat32, got.ValFloat32, 1e-6)
	ts.InDelta(in.ValFloat64, got.ValFloat64, 1e-9)
	ts.Equal(in.ValBool, got.ValBool)
	ts.Equal(in.ValString, got.ValString)
	ts.WithinDuration(in.ValTime, got.ValTime, time.Second)

	// pointer fields (non-nil, correct value)
	ts.Require().NotNil(got.PtrInt)
	ts.Equal(*in.PtrInt, *got.PtrInt)
	ts.Require().NotNil(got.PtrInt8)
	ts.Equal(*in.PtrInt8, *got.PtrInt8)
	ts.Require().NotNil(got.PtrInt16)
	ts.Equal(*in.PtrInt16, *got.PtrInt16)
	ts.Require().NotNil(got.PtrInt32)
	ts.Equal(*in.PtrInt32, *got.PtrInt32)
	ts.Require().NotNil(got.PtrInt64)
	ts.Equal(*in.PtrInt64, *got.PtrInt64)
	ts.Require().NotNil(got.PtrUint)
	ts.Equal(*in.PtrUint, *got.PtrUint)
	ts.Require().NotNil(got.PtrUint8)
	ts.Equal(*in.PtrUint8, *got.PtrUint8)
	ts.Require().NotNil(got.PtrUint16)
	ts.Equal(*in.PtrUint16, *got.PtrUint16)
	ts.Require().NotNil(got.PtrUint32)
	ts.Equal(*in.PtrUint32, *got.PtrUint32)
	ts.Require().NotNil(got.PtrUint64)
	ts.Equal(*in.PtrUint64, *got.PtrUint64)
	ts.Require().NotNil(got.PtrFloat32)
	ts.InDelta(*in.PtrFloat32, *got.PtrFloat32, 1e-6)
	ts.Require().NotNil(got.PtrFloat64)
	ts.InDelta(*in.PtrFloat64, *got.PtrFloat64, 1e-9)
	ts.Require().NotNil(got.PtrBool)
	ts.Equal(*in.PtrBool, *got.PtrBool)
	ts.Require().NotNil(got.PtrString)
	ts.Equal(*in.PtrString, *got.PtrString)
	ts.Require().NotNil(got.PtrTime)
	ts.WithinDuration(*in.PtrTime, *got.PtrTime, time.Second)

	// sql.Null* wrapper fields (scanned via the NullInt16/32/Byte pools)
	ts.Equal(in.NullInt16, got.NullInt16)
	ts.Equal(in.NullInt32, got.NullInt32)
	ts.Equal(in.NullByte, got.NullByte)
}

// TestFastScanRoundTripsNilPointers inserts a row whose every pointer field is
// nil and confirms the fast scanner materializes them back as nil (the NULL
// path) rather than a zero-valued pointer, while the value fields keep their
// zero values.
func (ts *ScannerTypesIntTestSuite) TestFastScanRoundTripsNilPointers() {
	in := &models.AllTypes{
		ValString: "no pointers",
		// a non-pointer time.Time value must be valid: mysql strict mode
		// rejects the zero time ('0000-00-00') on insert.
		ValTime: time.Date(2021, 3, 14, 15, 9, 26, 0, time.UTC),
	}

	create(&ts.testSuite, in)

	got, err := cql.Query[models.AllTypes](
		context.Background(),
		ts.db,
	).First()
	ts.Require().NoError(err)

	ts.Equal("no pointers", got.ValString)
	ts.Equal(0, got.ValInt)
	ts.Equal(uint64(0), got.ValUint64)
	ts.InDelta(0.0, got.ValFloat64, 1e-9)
	ts.False(got.ValBool)

	ts.Nil(got.PtrInt)
	ts.Nil(got.PtrInt8)
	ts.Nil(got.PtrInt16)
	ts.Nil(got.PtrInt32)
	ts.Nil(got.PtrInt64)
	ts.Nil(got.PtrUint)
	ts.Nil(got.PtrUint8)
	ts.Nil(got.PtrUint16)
	ts.Nil(got.PtrUint32)
	ts.Nil(got.PtrUint64)
	ts.Nil(got.PtrFloat32)
	ts.Nil(got.PtrFloat64)
	ts.Nil(got.PtrBool)
	ts.Nil(got.PtrString)
	ts.Nil(got.PtrTime)

	// zero sql.Null* wrappers round-trip as NULL (Valid == false)
	ts.False(got.NullInt16.Valid)
	ts.False(got.NullInt32.Valid)
	ts.False(got.NullByte.Valid)
}

// TestUnsupportedColumnFallsBackToGorm is the safety net for the fast scanner:
// models.WithUnsupportedColumn has a column type the scanner can't classify
// (models.Color, a named scalar), so no scanner is generated and the query
// must fall back to gorm. The unsupported column must round-trip — a partial
// fast scanner would silently drop it (leaving the zero value).
func (ts *ScannerTypesIntTestSuite) TestUnsupportedColumnFallsBackToGorm() {
	in := &models.WithUnsupportedColumn{
		Name:     "acme",
		Favorite: models.ColorBlue,
	}

	create(&ts.testSuite, in)

	got, err := cql.Query[models.WithUnsupportedColumn](
		context.Background(),
		ts.db,
	).First()
	ts.Require().NoError(err)

	ts.Equal("acme", got.Name)
	ts.Equal(models.ColorBlue, got.Favorite)
}

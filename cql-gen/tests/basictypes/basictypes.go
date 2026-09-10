package basictypes

import "github.com/FrancoLiberali/cql/model"

// BasicTypes covers every Go basic type that round-trips through SQL via gorm
// (and therefore through cql's fast scanner). Types with no SQL representation
// (uintptr, complex64, complex128) live in the unsupportedbasictypes fixture.
type BasicTypes struct {
	model.UUIDModel

	Bool    bool
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	UInt    uint
	UInt8   uint8
	UInt16  uint16
	UInt32  uint32
	UInt64  uint64
	Float32 float32
	Float64 float64
	String  string
	Byte    byte
}

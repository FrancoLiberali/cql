package condition

import (
	"time"

	"github.com/FrancoLiberali/cql/model"
)

func Int(value int) NumericValue[int] {
	return NumericValue[int]{Value: value}
}

func Int8(value int8) NumericValue[int8] {
	return NumericValue[int8]{Value: value}
}

func Int16(value int16) NumericValue[int16] {
	return NumericValue[int16]{Value: value}
}

func Int32(value int32) NumericValue[int32] {
	return NumericValue[int32]{Value: value}
}

func Int64(value int64) NumericValue[int64] {
	return NumericValue[int64]{Value: value}
}

func UInt(value uint) NumericValue[uint] {
	return NumericValue[uint]{Value: value}
}

func UInt8(value uint8) NumericValue[uint8] {
	return NumericValue[uint8]{Value: value}
}

func UInt16(value uint16) NumericValue[uint16] {
	return NumericValue[uint16]{Value: value}
}

func UInt32(value uint32) NumericValue[uint32] {
	return NumericValue[uint32]{Value: value}
}

func UInt64(value uint64) NumericValue[uint64] {
	return NumericValue[uint64]{Value: value}
}

func UIntPTR(value uintptr) NumericValue[uintptr] {
	return NumericValue[uintptr]{Value: value}
}

func Float32(value float32) NumericValue[float32] {
	return NumericValue[float32]{Value: value}
}

func Float64(value float64) NumericValue[float64] {
	return NumericValue[float64]{Value: value}
}

func Bool(value bool) BoolValue {
	return BoolValue{Value: value}
}

func String(value string) Value[string] {
	return Value[string]{Value: value}
}

func ByteArray(value []byte) Value[[]byte] {
	return Value[[]byte]{Value: value}
}

func Time(value time.Time) Value[time.Time] {
	return Value[time.Time]{Value: value}
}

func UUID(value model.UUID) Value[model.UUID] {
	return Value[model.UUID]{Value: value}
}

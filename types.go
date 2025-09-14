package cql

import (
	"time"

	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
)

func Int(value int) condition.NumericValue[int] {
	return condition.Int(value)
}

func Int8(value int8) condition.NumericValue[int8] {
	return condition.Int8(value)
}

func Int16(value int16) condition.NumericValue[int16] {
	return condition.Int16(value)
}

func Int32(value int32) condition.NumericValue[int32] {
	return condition.Int32(value)
}

func Int64(value int64) condition.NumericValue[int64] {
	return condition.Int64(value)
}

func UInt(value uint) condition.NumericValue[uint] {
	return condition.UInt(value)
}

func UInt8(value uint8) condition.NumericValue[uint8] {
	return condition.UInt8(value)
}

func UInt16(value uint16) condition.NumericValue[uint16] {
	return condition.UInt16(value)
}

func UInt32(value uint32) condition.NumericValue[uint32] {
	return condition.UInt32(value)
}

func UInt64(value uint64) condition.NumericValue[uint64] {
	return condition.UInt64(value)
}

func Float32(value float32) condition.NumericValue[float32] {
	return condition.Float32(value)
}

func Float64(value float64) condition.NumericValue[float64] {
	return condition.Float64(value)
}

func Bool(value bool) condition.BoolValue {
	return condition.Bool(value)
}

func String(value string) condition.Value[string] {
	return condition.String(value)
}

func ByteArray(value []byte) condition.Value[[]byte] {
	return condition.ByteArray(value)
}

func Time(value time.Time) condition.Value[time.Time] {
	return condition.Time(value)
}

func UUID(value model.UUID) condition.Value[model.UUID] {
	return condition.UUID(value)
}

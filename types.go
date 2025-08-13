package cql

import "github.com/FrancoLiberali/cql/condition"

func Int(value int) condition.NumericValue[int] {
	return condition.NumericValue[int]{Value: value}
}

func Int8(value int8) condition.NumericValue[int8] {
	return condition.NumericValue[int8]{Value: value}
}

func Int16(value int16) condition.NumericValue[int16] {
	return condition.NumericValue[int16]{Value: value}
}

func Int32(value int32) condition.NumericValue[int32] {
	return condition.NumericValue[int32]{Value: value}
}

func Int64(value int64) condition.NumericValue[int64] {
	return condition.NumericValue[int64]{Value: value}
}

func UInt(value uint) condition.NumericValue[uint] {
	return condition.NumericValue[uint]{Value: value}
}

func UInt8(value uint8) condition.NumericValue[uint8] {
	return condition.NumericValue[uint8]{Value: value}
}

func UInt16(value uint16) condition.NumericValue[uint16] {
	return condition.NumericValue[uint16]{Value: value}
}

func UInt32(value uint32) condition.NumericValue[uint32] {
	return condition.NumericValue[uint32]{Value: value}
}

func UInt64(value uint64) condition.NumericValue[uint64] {
	return condition.NumericValue[uint64]{Value: value}
}

func UIntPTR(value uintptr) condition.NumericValue[uintptr] {
	return condition.NumericValue[uintptr]{Value: value}
}

func Float32(value float32) condition.NumericValue[float32] {
	return condition.NumericValue[float32]{Value: value}
}

func Float64(value float64) condition.NumericValue[float64] {
	return condition.NumericValue[float64]{Value: value}
}

func Bool(value bool) condition.BoolValue {
	return condition.BoolValue{Value: value}
}

func String(value string) condition.Value[string] {
	return condition.Value[string]{Value: value}
}

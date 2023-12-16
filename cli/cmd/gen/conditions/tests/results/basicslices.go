// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	basicslices "github.com/ditrit/badaas-orm/cli/cmd/gen/conditions/tests/basicslices"
	orm "github.com/ditrit/badaas/orm"
	"time"
)

func BasicSlicesId(operator orm.Operator[orm.UUID]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, orm.UUID]{
		FieldIdentifier: orm.IDFieldID,
		Operator:        operator,
	}
}
func BasicSlicesCreatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, time.Time]{
		FieldIdentifier: orm.CreatedAtFieldID,
		Operator:        operator,
	}
}
func BasicSlicesUpdatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, time.Time]{
		FieldIdentifier: orm.UpdatedAtFieldID,
		Operator:        operator,
	}
}
func BasicSlicesDeletedAt(operator orm.Operator[time.Time]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, time.Time]{
		FieldIdentifier: orm.DeletedAtFieldID,
		Operator:        operator,
	}
}

var basicSlicesBoolFieldID = orm.FieldIdentifier{Field: "Bool"}

func BasicSlicesBool(operator orm.Operator[[]bool]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []bool]{
		FieldIdentifier: basicSlicesBoolFieldID,
		Operator:        operator,
	}
}

var basicSlicesIntFieldID = orm.FieldIdentifier{Field: "Int"}

func BasicSlicesInt(operator orm.Operator[[]int]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []int]{
		FieldIdentifier: basicSlicesIntFieldID,
		Operator:        operator,
	}
}

var basicSlicesInt8FieldID = orm.FieldIdentifier{Field: "Int8"}

func BasicSlicesInt8(operator orm.Operator[[]int8]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []int8]{
		FieldIdentifier: basicSlicesInt8FieldID,
		Operator:        operator,
	}
}

var basicSlicesInt16FieldID = orm.FieldIdentifier{Field: "Int16"}

func BasicSlicesInt16(operator orm.Operator[[]int16]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []int16]{
		FieldIdentifier: basicSlicesInt16FieldID,
		Operator:        operator,
	}
}

var basicSlicesInt32FieldID = orm.FieldIdentifier{Field: "Int32"}

func BasicSlicesInt32(operator orm.Operator[[]int32]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []int32]{
		FieldIdentifier: basicSlicesInt32FieldID,
		Operator:        operator,
	}
}

var basicSlicesInt64FieldID = orm.FieldIdentifier{Field: "Int64"}

func BasicSlicesInt64(operator orm.Operator[[]int64]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []int64]{
		FieldIdentifier: basicSlicesInt64FieldID,
		Operator:        operator,
	}
}

var basicSlicesUIntFieldID = orm.FieldIdentifier{Field: "UInt"}

func BasicSlicesUInt(operator orm.Operator[[]uint]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []uint]{
		FieldIdentifier: basicSlicesUIntFieldID,
		Operator:        operator,
	}
}

var basicSlicesUInt8FieldID = orm.FieldIdentifier{Field: "UInt8"}

func BasicSlicesUInt8(operator orm.Operator[[]uint8]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []uint8]{
		FieldIdentifier: basicSlicesUInt8FieldID,
		Operator:        operator,
	}
}

var basicSlicesUInt16FieldID = orm.FieldIdentifier{Field: "UInt16"}

func BasicSlicesUInt16(operator orm.Operator[[]uint16]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []uint16]{
		FieldIdentifier: basicSlicesUInt16FieldID,
		Operator:        operator,
	}
}

var basicSlicesUInt32FieldID = orm.FieldIdentifier{Field: "UInt32"}

func BasicSlicesUInt32(operator orm.Operator[[]uint32]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []uint32]{
		FieldIdentifier: basicSlicesUInt32FieldID,
		Operator:        operator,
	}
}

var basicSlicesUInt64FieldID = orm.FieldIdentifier{Field: "UInt64"}

func BasicSlicesUInt64(operator orm.Operator[[]uint64]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []uint64]{
		FieldIdentifier: basicSlicesUInt64FieldID,
		Operator:        operator,
	}
}

var basicSlicesUIntptrFieldID = orm.FieldIdentifier{Field: "UIntptr"}

func BasicSlicesUIntptr(operator orm.Operator[[]uintptr]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []uintptr]{
		FieldIdentifier: basicSlicesUIntptrFieldID,
		Operator:        operator,
	}
}

var basicSlicesFloat32FieldID = orm.FieldIdentifier{Field: "Float32"}

func BasicSlicesFloat32(operator orm.Operator[[]float32]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []float32]{
		FieldIdentifier: basicSlicesFloat32FieldID,
		Operator:        operator,
	}
}

var basicSlicesFloat64FieldID = orm.FieldIdentifier{Field: "Float64"}

func BasicSlicesFloat64(operator orm.Operator[[]float64]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []float64]{
		FieldIdentifier: basicSlicesFloat64FieldID,
		Operator:        operator,
	}
}

var basicSlicesComplex64FieldID = orm.FieldIdentifier{Field: "Complex64"}

func BasicSlicesComplex64(operator orm.Operator[[]complex64]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []complex64]{
		FieldIdentifier: basicSlicesComplex64FieldID,
		Operator:        operator,
	}
}

var basicSlicesComplex128FieldID = orm.FieldIdentifier{Field: "Complex128"}

func BasicSlicesComplex128(operator orm.Operator[[]complex128]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []complex128]{
		FieldIdentifier: basicSlicesComplex128FieldID,
		Operator:        operator,
	}
}

var basicSlicesStringFieldID = orm.FieldIdentifier{Field: "String"}

func BasicSlicesString(operator orm.Operator[[]string]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []string]{
		FieldIdentifier: basicSlicesStringFieldID,
		Operator:        operator,
	}
}

var basicSlicesByteFieldID = orm.FieldIdentifier{Field: "Byte"}

func BasicSlicesByte(operator orm.Operator[[]uint8]) orm.WhereCondition[basicslices.BasicSlices] {
	return orm.FieldCondition[basicslices.BasicSlices, []uint8]{
		FieldIdentifier: basicSlicesByteFieldID,
		Operator:        operator,
	}
}

var BasicSlicesPreloadAttributes = orm.NewPreloadCondition[basicslices.BasicSlices](basicSlicesBoolFieldID, basicSlicesIntFieldID, basicSlicesInt8FieldID, basicSlicesInt16FieldID, basicSlicesInt32FieldID, basicSlicesInt64FieldID, basicSlicesUIntFieldID, basicSlicesUInt8FieldID, basicSlicesUInt16FieldID, basicSlicesUInt32FieldID, basicSlicesUInt64FieldID, basicSlicesUIntptrFieldID, basicSlicesFloat32FieldID, basicSlicesFloat64FieldID, basicSlicesComplex64FieldID, basicSlicesComplex128FieldID, basicSlicesStringFieldID, basicSlicesByteFieldID)

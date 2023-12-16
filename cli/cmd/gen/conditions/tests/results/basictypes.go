// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	basictypes "github.com/ditrit/badaas-orm/cli/cmd/gen/conditions/tests/basictypes"
	orm "github.com/ditrit/badaas/orm"
	"time"
)

func BasicTypesId(operator orm.Operator[orm.UUID]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, orm.UUID]{
		FieldIdentifier: orm.IDFieldID,
		Operator:        operator,
	}
}
func BasicTypesCreatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, time.Time]{
		FieldIdentifier: orm.CreatedAtFieldID,
		Operator:        operator,
	}
}
func BasicTypesUpdatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, time.Time]{
		FieldIdentifier: orm.UpdatedAtFieldID,
		Operator:        operator,
	}
}
func BasicTypesDeletedAt(operator orm.Operator[time.Time]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, time.Time]{
		FieldIdentifier: orm.DeletedAtFieldID,
		Operator:        operator,
	}
}

var basicTypesBoolFieldID = orm.FieldIdentifier{Field: "Bool"}

func BasicTypesBool(operator orm.Operator[bool]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, bool]{
		FieldIdentifier: basicTypesBoolFieldID,
		Operator:        operator,
	}
}

var basicTypesIntFieldID = orm.FieldIdentifier{Field: "Int"}

func BasicTypesInt(operator orm.Operator[int]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, int]{
		FieldIdentifier: basicTypesIntFieldID,
		Operator:        operator,
	}
}

var basicTypesInt8FieldID = orm.FieldIdentifier{Field: "Int8"}

func BasicTypesInt8(operator orm.Operator[int8]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, int8]{
		FieldIdentifier: basicTypesInt8FieldID,
		Operator:        operator,
	}
}

var basicTypesInt16FieldID = orm.FieldIdentifier{Field: "Int16"}

func BasicTypesInt16(operator orm.Operator[int16]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, int16]{
		FieldIdentifier: basicTypesInt16FieldID,
		Operator:        operator,
	}
}

var basicTypesInt32FieldID = orm.FieldIdentifier{Field: "Int32"}

func BasicTypesInt32(operator orm.Operator[int32]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, int32]{
		FieldIdentifier: basicTypesInt32FieldID,
		Operator:        operator,
	}
}

var basicTypesInt64FieldID = orm.FieldIdentifier{Field: "Int64"}

func BasicTypesInt64(operator orm.Operator[int64]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, int64]{
		FieldIdentifier: basicTypesInt64FieldID,
		Operator:        operator,
	}
}

var basicTypesUIntFieldID = orm.FieldIdentifier{Field: "UInt"}

func BasicTypesUInt(operator orm.Operator[uint]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, uint]{
		FieldIdentifier: basicTypesUIntFieldID,
		Operator:        operator,
	}
}

var basicTypesUInt8FieldID = orm.FieldIdentifier{Field: "UInt8"}

func BasicTypesUInt8(operator orm.Operator[uint8]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, uint8]{
		FieldIdentifier: basicTypesUInt8FieldID,
		Operator:        operator,
	}
}

var basicTypesUInt16FieldID = orm.FieldIdentifier{Field: "UInt16"}

func BasicTypesUInt16(operator orm.Operator[uint16]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, uint16]{
		FieldIdentifier: basicTypesUInt16FieldID,
		Operator:        operator,
	}
}

var basicTypesUInt32FieldID = orm.FieldIdentifier{Field: "UInt32"}

func BasicTypesUInt32(operator orm.Operator[uint32]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, uint32]{
		FieldIdentifier: basicTypesUInt32FieldID,
		Operator:        operator,
	}
}

var basicTypesUInt64FieldID = orm.FieldIdentifier{Field: "UInt64"}

func BasicTypesUInt64(operator orm.Operator[uint64]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, uint64]{
		FieldIdentifier: basicTypesUInt64FieldID,
		Operator:        operator,
	}
}

var basicTypesUIntptrFieldID = orm.FieldIdentifier{Field: "UIntptr"}

func BasicTypesUIntptr(operator orm.Operator[uintptr]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, uintptr]{
		FieldIdentifier: basicTypesUIntptrFieldID,
		Operator:        operator,
	}
}

var basicTypesFloat32FieldID = orm.FieldIdentifier{Field: "Float32"}

func BasicTypesFloat32(operator orm.Operator[float32]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, float32]{
		FieldIdentifier: basicTypesFloat32FieldID,
		Operator:        operator,
	}
}

var basicTypesFloat64FieldID = orm.FieldIdentifier{Field: "Float64"}

func BasicTypesFloat64(operator orm.Operator[float64]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, float64]{
		FieldIdentifier: basicTypesFloat64FieldID,
		Operator:        operator,
	}
}

var basicTypesComplex64FieldID = orm.FieldIdentifier{Field: "Complex64"}

func BasicTypesComplex64(operator orm.Operator[complex64]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, complex64]{
		FieldIdentifier: basicTypesComplex64FieldID,
		Operator:        operator,
	}
}

var basicTypesComplex128FieldID = orm.FieldIdentifier{Field: "Complex128"}

func BasicTypesComplex128(operator orm.Operator[complex128]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, complex128]{
		FieldIdentifier: basicTypesComplex128FieldID,
		Operator:        operator,
	}
}

var basicTypesStringFieldID = orm.FieldIdentifier{Field: "String"}

func BasicTypesString(operator orm.Operator[string]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, string]{
		FieldIdentifier: basicTypesStringFieldID,
		Operator:        operator,
	}
}

var basicTypesByteFieldID = orm.FieldIdentifier{Field: "Byte"}

func BasicTypesByte(operator orm.Operator[uint8]) orm.WhereCondition[basictypes.BasicTypes] {
	return orm.FieldCondition[basictypes.BasicTypes, uint8]{
		FieldIdentifier: basicTypesByteFieldID,
		Operator:        operator,
	}
}

var BasicTypesPreloadAttributes = orm.NewPreloadCondition[basictypes.BasicTypes](basicTypesBoolFieldID, basicTypesIntFieldID, basicTypesInt8FieldID, basicTypesInt16FieldID, basicTypesInt32FieldID, basicTypesInt64FieldID, basicTypesUIntFieldID, basicTypesUInt8FieldID, basicTypesUInt16FieldID, basicTypesUInt32FieldID, basicTypesUInt64FieldID, basicTypesUIntptrFieldID, basicTypesFloat32FieldID, basicTypesFloat64FieldID, basicTypesComplex64FieldID, basicTypesComplex128FieldID, basicTypesStringFieldID, basicTypesByteFieldID)

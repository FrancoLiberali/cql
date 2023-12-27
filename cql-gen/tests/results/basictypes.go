// Code generated by cql-gen v0.0.5, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	basictypes "github.com/FrancoLiberali/cql/cql-gen/cmd/gen/conditions/tests/basictypes"
	model "github.com/FrancoLiberali/cql/model"
	"time"
)

type basicTypesConditions struct {
	ID         condition.Field[basictypes.BasicTypes, model.UUID]
	CreatedAt  condition.Field[basictypes.BasicTypes, time.Time]
	UpdatedAt  condition.Field[basictypes.BasicTypes, time.Time]
	DeletedAt  condition.Field[basictypes.BasicTypes, time.Time]
	Bool       condition.BoolField[basictypes.BasicTypes]
	Int        condition.UpdatableField[basictypes.BasicTypes, int]
	Int8       condition.UpdatableField[basictypes.BasicTypes, int8]
	Int16      condition.UpdatableField[basictypes.BasicTypes, int16]
	Int32      condition.UpdatableField[basictypes.BasicTypes, int32]
	Int64      condition.UpdatableField[basictypes.BasicTypes, int64]
	UInt       condition.UpdatableField[basictypes.BasicTypes, uint]
	UInt8      condition.UpdatableField[basictypes.BasicTypes, uint8]
	UInt16     condition.UpdatableField[basictypes.BasicTypes, uint16]
	UInt32     condition.UpdatableField[basictypes.BasicTypes, uint32]
	UInt64     condition.UpdatableField[basictypes.BasicTypes, uint64]
	UIntptr    condition.UpdatableField[basictypes.BasicTypes, uintptr]
	Float32    condition.UpdatableField[basictypes.BasicTypes, float32]
	Float64    condition.UpdatableField[basictypes.BasicTypes, float64]
	Complex64  condition.UpdatableField[basictypes.BasicTypes, complex64]
	Complex128 condition.UpdatableField[basictypes.BasicTypes, complex128]
	String     condition.StringField[basictypes.BasicTypes]
	Byte       condition.UpdatableField[basictypes.BasicTypes, uint8]
}

var BasicTypes = basicTypesConditions{
	Bool:       condition.BoolField[basictypes.BasicTypes]{UpdatableField: condition.UpdatableField[basictypes.BasicTypes, bool]{Field: condition.Field[basictypes.BasicTypes, bool]{Name: "Bool"}}},
	Byte:       condition.UpdatableField[basictypes.BasicTypes, uint8]{Field: condition.Field[basictypes.BasicTypes, uint8]{Name: "Byte"}},
	Complex128: condition.UpdatableField[basictypes.BasicTypes, complex128]{Field: condition.Field[basictypes.BasicTypes, complex128]{Name: "Complex128"}},
	Complex64:  condition.UpdatableField[basictypes.BasicTypes, complex64]{Field: condition.Field[basictypes.BasicTypes, complex64]{Name: "Complex64"}},
	CreatedAt:  condition.Field[basictypes.BasicTypes, time.Time]{Name: "CreatedAt"},
	DeletedAt:  condition.Field[basictypes.BasicTypes, time.Time]{Name: "DeletedAt"},
	Float32:    condition.UpdatableField[basictypes.BasicTypes, float32]{Field: condition.Field[basictypes.BasicTypes, float32]{Name: "Float32"}},
	Float64:    condition.UpdatableField[basictypes.BasicTypes, float64]{Field: condition.Field[basictypes.BasicTypes, float64]{Name: "Float64"}},
	ID:         condition.Field[basictypes.BasicTypes, model.UUID]{Name: "ID"},
	Int:        condition.UpdatableField[basictypes.BasicTypes, int]{Field: condition.Field[basictypes.BasicTypes, int]{Name: "Int"}},
	Int16:      condition.UpdatableField[basictypes.BasicTypes, int16]{Field: condition.Field[basictypes.BasicTypes, int16]{Name: "Int16"}},
	Int32:      condition.UpdatableField[basictypes.BasicTypes, int32]{Field: condition.Field[basictypes.BasicTypes, int32]{Name: "Int32"}},
	Int64:      condition.UpdatableField[basictypes.BasicTypes, int64]{Field: condition.Field[basictypes.BasicTypes, int64]{Name: "Int64"}},
	Int8:       condition.UpdatableField[basictypes.BasicTypes, int8]{Field: condition.Field[basictypes.BasicTypes, int8]{Name: "Int8"}},
	String:     condition.StringField[basictypes.BasicTypes]{UpdatableField: condition.UpdatableField[basictypes.BasicTypes, string]{Field: condition.Field[basictypes.BasicTypes, string]{Name: "String"}}},
	UInt:       condition.UpdatableField[basictypes.BasicTypes, uint]{Field: condition.Field[basictypes.BasicTypes, uint]{Name: "UInt"}},
	UInt16:     condition.UpdatableField[basictypes.BasicTypes, uint16]{Field: condition.Field[basictypes.BasicTypes, uint16]{Name: "UInt16"}},
	UInt32:     condition.UpdatableField[basictypes.BasicTypes, uint32]{Field: condition.Field[basictypes.BasicTypes, uint32]{Name: "UInt32"}},
	UInt64:     condition.UpdatableField[basictypes.BasicTypes, uint64]{Field: condition.Field[basictypes.BasicTypes, uint64]{Name: "UInt64"}},
	UInt8:      condition.UpdatableField[basictypes.BasicTypes, uint8]{Field: condition.Field[basictypes.BasicTypes, uint8]{Name: "UInt8"}},
	UIntptr:    condition.UpdatableField[basictypes.BasicTypes, uintptr]{Field: condition.Field[basictypes.BasicTypes, uintptr]{Name: "UIntptr"}},
	UpdatedAt:  condition.Field[basictypes.BasicTypes, time.Time]{Name: "UpdatedAt"},
}

// Preload allows preloading the BasicTypes when doing a query
func (basicTypesConditions basicTypesConditions) preload() condition.Condition[basictypes.BasicTypes] {
	return condition.NewPreloadCondition[basictypes.BasicTypes](basicTypesConditions.ID, basicTypesConditions.CreatedAt, basicTypesConditions.UpdatedAt, basicTypesConditions.DeletedAt, basicTypesConditions.Bool, basicTypesConditions.Int, basicTypesConditions.Int8, basicTypesConditions.Int16, basicTypesConditions.Int32, basicTypesConditions.Int64, basicTypesConditions.UInt, basicTypesConditions.UInt8, basicTypesConditions.UInt16, basicTypesConditions.UInt32, basicTypesConditions.UInt64, basicTypesConditions.UIntptr, basicTypesConditions.Float32, basicTypesConditions.Float64, basicTypesConditions.Complex64, basicTypesConditions.Complex128, basicTypesConditions.String, basicTypesConditions.Byte)
}

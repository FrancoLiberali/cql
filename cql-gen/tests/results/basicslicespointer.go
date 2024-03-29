// Code generated by cql-gen v0.1.0, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	basicslicespointer "github.com/FrancoLiberali/cql/cql-gen/cmd/gen/conditions/tests/basicslicespointer"
	model "github.com/FrancoLiberali/cql/model"
	"time"
)

type basicSlicesPointerConditions struct {
	ID         condition.Field[basicslicespointer.BasicSlicesPointer, model.UUID]
	CreatedAt  condition.Field[basicslicespointer.BasicSlicesPointer, time.Time]
	UpdatedAt  condition.Field[basicslicespointer.BasicSlicesPointer, time.Time]
	DeletedAt  condition.Field[basicslicespointer.BasicSlicesPointer, time.Time]
	Bool       condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []bool]
	Int        condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []int]
	Int8       condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []int8]
	Int16      condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []int16]
	Int32      condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []int32]
	Int64      condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []int64]
	UInt       condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []uint]
	UInt8      condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []uint8]
	UInt16     condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []uint16]
	UInt32     condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []uint32]
	UInt64     condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []uint64]
	UIntptr    condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []uintptr]
	Float32    condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []float32]
	Float64    condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []float64]
	Complex64  condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []complex64]
	Complex128 condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []complex128]
	String     condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []string]
	Byte       condition.UpdatableField[basicslicespointer.BasicSlicesPointer, []uint8]
}

var BasicSlicesPointer = basicSlicesPointerConditions{
	Bool:       condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []bool]("Bool", "", ""),
	Byte:       condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []uint8]("Byte", "", ""),
	Complex128: condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []complex128]("Complex128", "", ""),
	Complex64:  condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []complex64]("Complex64", "", ""),
	CreatedAt:  condition.NewField[basicslicespointer.BasicSlicesPointer, time.Time]("CreatedAt", "", ""),
	DeletedAt:  condition.NewField[basicslicespointer.BasicSlicesPointer, time.Time]("DeletedAt", "", ""),
	Float32:    condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []float32]("Float32", "", ""),
	Float64:    condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []float64]("Float64", "", ""),
	ID:         condition.NewField[basicslicespointer.BasicSlicesPointer, model.UUID]("ID", "", ""),
	Int:        condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []int]("Int", "", ""),
	Int16:      condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []int16]("Int16", "", ""),
	Int32:      condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []int32]("Int32", "", ""),
	Int64:      condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []int64]("Int64", "", ""),
	Int8:       condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []int8]("Int8", "", ""),
	String:     condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []string]("String", "", ""),
	UInt:       condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []uint]("UInt", "", ""),
	UInt16:     condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []uint16]("UInt16", "", ""),
	UInt32:     condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []uint32]("UInt32", "", ""),
	UInt64:     condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []uint64]("UInt64", "", ""),
	UInt8:      condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []uint8]("UInt8", "", ""),
	UIntptr:    condition.NewUpdatableField[basicslicespointer.BasicSlicesPointer, []uintptr]("UIntptr", "", ""),
	UpdatedAt:  condition.NewField[basicslicespointer.BasicSlicesPointer, time.Time]("UpdatedAt", "", ""),
}

// Preload allows preloading the BasicSlicesPointer when doing a query
func (basicSlicesPointerConditions basicSlicesPointerConditions) preload() condition.Condition[basicslicespointer.BasicSlicesPointer] {
	return condition.NewPreloadCondition[basicslicespointer.BasicSlicesPointer](basicSlicesPointerConditions.ID, basicSlicesPointerConditions.CreatedAt, basicSlicesPointerConditions.UpdatedAt, basicSlicesPointerConditions.DeletedAt, basicSlicesPointerConditions.Bool, basicSlicesPointerConditions.Int, basicSlicesPointerConditions.Int8, basicSlicesPointerConditions.Int16, basicSlicesPointerConditions.Int32, basicSlicesPointerConditions.Int64, basicSlicesPointerConditions.UInt, basicSlicesPointerConditions.UInt8, basicSlicesPointerConditions.UInt16, basicSlicesPointerConditions.UInt32, basicSlicesPointerConditions.UInt64, basicSlicesPointerConditions.UIntptr, basicSlicesPointerConditions.Float32, basicSlicesPointerConditions.Float64, basicSlicesPointerConditions.Complex64, basicSlicesPointerConditions.Complex128, basicSlicesPointerConditions.String, basicSlicesPointerConditions.Byte)
}

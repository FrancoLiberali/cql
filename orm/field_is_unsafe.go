package orm

import (
	"github.com/ditrit/badaas/orm/condition"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/query"
	"github.com/ditrit/badaas/orm/unsafe"
)

type UnsafeFieldIs[TObject model.Model, TAttribute any] struct {
	fieldID query.FieldIdentifier[TAttribute]
}

// EqualTo
func (is UnsafeFieldIs[TObject, TAttribute]) Eq(value any) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, unsafe.Eq[TAttribute](value))
}

// NotEqualTo
func (is UnsafeFieldIs[TObject, TAttribute]) NotEq(value any) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, unsafe.NotEq[TAttribute](value))
}

// LessThan
func (is UnsafeFieldIs[TObject, TAttribute]) Lt(value any) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, unsafe.Lt[TAttribute](value))
}

// LessThanOrEqualTo
func (is UnsafeFieldIs[TObject, TAttribute]) LtOrEq(value any) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, unsafe.LtOrEq[TAttribute](value))
}

// GreaterThan
func (is UnsafeFieldIs[TObject, TAttribute]) Gt(value any) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, unsafe.Gt[TAttribute](value))
}

// GreaterThanOrEqualTo
func (is UnsafeFieldIs[TObject, TAttribute]) GtOrEq(value any) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, unsafe.GtOrEq[TAttribute](value))
}

// Equivalent to field1 < value < field2
func (is UnsafeFieldIs[TObject, TAttribute]) Between(v1, v2 any) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, unsafe.Between[TAttribute](v1, v2))
}

// Equivalent to NOT (field1 < value < field2)
func (is UnsafeFieldIs[TObject, TAttribute]) NotBetween(v1, v2 any) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, unsafe.NotBetween[TAttribute](v1, v2))
}

// Not supported by: mysql
func (is UnsafeFieldIs[TObject, TAttribute]) Distinct(value any) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, unsafe.IsDistinct[TAttribute](value))
}

// Not supported by: mysql
func (is UnsafeFieldIs[TObject, TAttribute]) NotDistinct(value any) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, unsafe.IsNotDistinct[TAttribute](value))
}

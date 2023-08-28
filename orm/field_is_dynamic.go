package orm

import (
	"github.com/ditrit/badaas/orm/condition"
	"github.com/ditrit/badaas/orm/dynamic"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/query"
)

type DynamicFieldIs[TObject model.Model, TAttribute any] struct {
	fieldID query.FieldIdentifier[TAttribute]
}

// EqualTo
func (is DynamicFieldIs[TObject, TAttribute]) Eq(field query.FieldIdentifier[TAttribute]) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, dynamic.Eq(field))
}

// NotEqualTo
func (is DynamicFieldIs[TObject, TAttribute]) NotEq(field query.FieldIdentifier[TAttribute]) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, dynamic.NotEq(field))
}

// LessThan
func (is DynamicFieldIs[TObject, TAttribute]) Lt(field query.FieldIdentifier[TAttribute]) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, dynamic.Lt(field))
}

// LessThanOrEqualTo
func (is DynamicFieldIs[TObject, TAttribute]) LtOrEq(field query.FieldIdentifier[TAttribute]) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, dynamic.LtOrEq(field))
}

// GreaterThan
func (is DynamicFieldIs[TObject, TAttribute]) Gt(field query.FieldIdentifier[TAttribute]) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, dynamic.Gt(field))
}

// GreaterThanOrEqualTo
func (is DynamicFieldIs[TObject, TAttribute]) GtOrEq(field query.FieldIdentifier[TAttribute]) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, dynamic.GtOrEq(field))
}

// Equivalent to field1 < value < field2
func (is DynamicFieldIs[TObject, TAttribute]) Between(field1, field2 query.FieldIdentifier[TAttribute]) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, dynamic.Between(field1, field2))
}

// Equivalent to NOT (field1 < value < field2)
func (is DynamicFieldIs[TObject, TAttribute]) NotBetween(field1, field2 query.FieldIdentifier[TAttribute]) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, dynamic.NotBetween(field1, field2))
}

// Not supported by: mysql
func (is DynamicFieldIs[TObject, TAttribute]) Distinct(field query.FieldIdentifier[TAttribute]) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, dynamic.IsDistinct(field))
}

// Not supported by: mysql
func (is DynamicFieldIs[TObject, TAttribute]) NotDistinct(field query.FieldIdentifier[TAttribute]) condition.DynamicCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.fieldID, dynamic.IsNotDistinct(field))
}

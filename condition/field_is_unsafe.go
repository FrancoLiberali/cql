package condition

import (
	"github.com/FrancoLiberali/cql/model"
)

type UnsafeFieldIs[TObject model.Model, TAttribute any] struct {
	field Field[TObject, TAttribute]
}

// EqualTo
func (is UnsafeFieldIs[TObject, TAttribute]) Eq(value IValue) WhereCondition[TObject] {
	return NewFieldCondition[TObject, TAttribute](is.field, Eq[TAttribute](unsafeValue{Value: value}))
}

// NotEqualTo
func (is UnsafeFieldIs[TObject, TAttribute]) NotEq(value IValue) WhereCondition[TObject] {
	return NewFieldCondition[TObject, TAttribute](is.field, NotEq[TAttribute](unsafeValue{Value: value}))
}

// LessThan
func (is UnsafeFieldIs[TObject, TAttribute]) Lt(value IValue) WhereCondition[TObject] {
	return NewFieldCondition[TObject, TAttribute](is.field, Lt[TAttribute](unsafeValue{Value: value}))
}

// LessThanOrEqualTo
func (is UnsafeFieldIs[TObject, TAttribute]) LtOrEq(value IValue) WhereCondition[TObject] {
	return NewFieldCondition[TObject, TAttribute](is.field, LtOrEq[TAttribute](unsafeValue{Value: value}))
}

// GreaterThan
func (is UnsafeFieldIs[TObject, TAttribute]) Gt(value IValue) WhereCondition[TObject] {
	return NewFieldCondition[TObject, TAttribute](is.field, Gt[TAttribute](unsafeValue{Value: value}))
}

// GreaterThanOrEqualTo
func (is UnsafeFieldIs[TObject, TAttribute]) GtOrEq(value IValue) WhereCondition[TObject] {
	return NewFieldCondition[TObject, TAttribute](is.field, GtOrEq[TAttribute](unsafeValue{Value: value}))
}

// Equivalent to field1 < value < field2
func (is UnsafeFieldIs[TObject, TAttribute]) Between(v1, v2 IValue) WhereCondition[TObject] {
	return NewFieldCondition[TObject, TAttribute](is.field, Between[TAttribute](unsafeValue{Value: v1}, unsafeValue{Value: v2}))
}

// Equivalent to NOT (field1 < value < field2)
func (is UnsafeFieldIs[TObject, TAttribute]) NotBetween(v1, v2 IValue) WhereCondition[TObject] {
	return NewFieldCondition[TObject, TAttribute](is.field, NotBetween[TAttribute](unsafeValue{Value: v1}, unsafeValue{Value: v2}))
}

func (is UnsafeFieldIs[TObject, TAttribute]) Distinct(value IValue) WhereCondition[TObject] {
	return NewFieldCondition[TObject, TAttribute](is.field, IsDistinct[TAttribute](unsafeValue{Value: value}))
}

func (is UnsafeFieldIs[TObject, TAttribute]) NotDistinct(value IValue) WhereCondition[TObject] {
	return NewFieldCondition[TObject, TAttribute](is.field, IsNotDistinct[TAttribute](unsafeValue{Value: value}))
}

package condition

import (
	"github.com/FrancoLiberali/cql/model"
)

type UnsafeFieldIs[TObject model.Model, TAttribute any] struct {
	field Field[TObject, TAttribute]
}

// EqualTo
func (is UnsafeFieldIs[TObject, TAttribute]) Eq(value any) WhereCondition[TObject] {
	return NewFieldCondition(is.field, Eq[TAttribute](value))
}

// NotEqualTo
func (is UnsafeFieldIs[TObject, TAttribute]) NotEq(value any) WhereCondition[TObject] {
	return NewFieldCondition(is.field, NotEq[TAttribute](value))
}

// LessThan
func (is UnsafeFieldIs[TObject, TAttribute]) Lt(value any) WhereCondition[TObject] {
	return NewFieldCondition(is.field, Lt[TAttribute](value))
}

// LessThanOrEqualTo
func (is UnsafeFieldIs[TObject, TAttribute]) LtOrEq(value any) WhereCondition[TObject] {
	return NewFieldCondition(is.field, LtOrEq[TAttribute](value))
}

// GreaterThan
func (is UnsafeFieldIs[TObject, TAttribute]) Gt(value any) WhereCondition[TObject] {
	return NewFieldCondition(is.field, Gt[TAttribute](value))
}

// GreaterThanOrEqualTo
func (is UnsafeFieldIs[TObject, TAttribute]) GtOrEq(value any) WhereCondition[TObject] {
	return NewFieldCondition(is.field, GtOrEq[TAttribute](value))
}

// Equivalent to field1 < value < field2
func (is UnsafeFieldIs[TObject, TAttribute]) Between(v1, v2 any) WhereCondition[TObject] {
	return NewFieldCondition(is.field, Between[TAttribute](v1, v2))
}

// Equivalent to NOT (field1 < value < field2)
func (is UnsafeFieldIs[TObject, TAttribute]) NotBetween(v1, v2 any) WhereCondition[TObject] {
	return NewFieldCondition(is.field, NotBetween[TAttribute](v1, v2))
}

func (is UnsafeFieldIs[TObject, TAttribute]) Distinct(value any) WhereCondition[TObject] {
	return NewFieldCondition(is.field, IsDistinct[TAttribute](value))
}

func (is UnsafeFieldIs[TObject, TAttribute]) NotDistinct(value any) WhereCondition[TObject] {
	return NewFieldCondition(is.field, IsNotDistinct[TAttribute](value))
}

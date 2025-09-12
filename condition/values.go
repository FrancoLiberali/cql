package condition

import "golang.org/x/exp/constraints"

type IValue interface {
	ToSQL(query *CQLQuery) (string, []any, error)
}

type ValueOfType[T any] interface {
	IValue

	GetValue() T
}

type NumericOfType[T any] interface {
	IValue

	GetValue() float64

	GetNumericValue() T
}

type Numeric interface {
	constraints.Integer | constraints.Float
}

type NumericValue[T Numeric] struct {
	Value T
}

func (numericValue NumericValue[T]) GetNumericValue() T {
	return numericValue.Value
}

func (numericValue NumericValue[T]) GetValue() float64 {
	return float64(numericValue.Value)
}

func (numericValue NumericValue[T]) ToSQL(_ *CQLQuery) (string, []any, error) {
	return "", []any{numericValue.Value}, nil
}

type BoolValue struct {
	Value bool
}

func (boolValue BoolValue) GetValue() bool {
	return boolValue.Value
}

func (boolValue BoolValue) ToSQL(_ *CQLQuery) (string, []any, error) {
	return "", []any{boolValue.Value}, nil
}

type Value[T any] struct {
	Value T
}

func (value Value[T]) GetValue() T {
	return value.Value
}

func (value Value[T]) ToSQL(_ *CQLQuery) (string, []any, error) {
	return "", []any{value.Value}, nil
}

type unsafeValue struct {
	Value IValue
}

func (unsafeValue unsafeValue) ToSQL(query *CQLQuery) (string, []any, error) {
	return unsafeValue.Value.ToSQL(query)
}

package condition

import "golang.org/x/exp/constraints"

type Numeric interface {
	constraints.Integer | constraints.Float
}

type NumericValue[T Numeric] struct {
	Value T
}

func (numericValue NumericValue[T]) GetValue() float64 {
	return float64(numericValue.Value)
}

func (numericValue NumericValue[T]) ToSQL(_ *CQLQuery) (string, []any, error) {
	return "", []any{numericValue.GetValue()}, nil
}

type BoolValue struct {
	Value bool
}

func (boolValue BoolValue) GetValue() bool {
	return boolValue.Value
}

func (boolValue BoolValue) ToSQL(_ *CQLQuery) (string, []any, error) {
	return "", []any{boolValue.GetValue()}, nil
}

type Value[T any] struct {
	Value T
}

func (value Value[T]) GetValue() T {
	return value.Value
}

func (value Value[T]) ToSQL(_ *CQLQuery) (string, []any, error) {
	return "", []any{value.GetValue()}, nil
}

type unsafeValue struct {
	Value IValue
}

func (unsafeValue unsafeValue) ToSQL(query *CQLQuery) (string, []any, error) {
	return unsafeValue.Value.ToSQL(query)
}

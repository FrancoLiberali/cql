package condition

import "golang.org/x/exp/constraints"

type toSQLFunc func(query *GormQuery) (string, error)

type Numeric interface {
	constraints.Integer | constraints.Float
}

type NumericValue[T Numeric] struct {
	Value T
}

func (numericValue NumericValue[T]) getSQL() toSQLFunc {
	return nil
}

func (numericValue NumericValue[T]) getValue() float64 {
	return float64(numericValue.Value)
}

type BoolValue struct {
	Value bool
}

func (boolValue BoolValue) getValue() bool {
	return boolValue.Value
}

func (boolValue BoolValue) getSQL() toSQLFunc {
	return nil
}

type Value[T any] struct {
	Value T
}

func (value Value[T]) getValue() T {
	return value.Value
}

func (value Value[T]) getSQL() toSQLFunc {
	return nil
}

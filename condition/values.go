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

func (numericValue NumericValue[T]) GetValue() float64 {
	return float64(numericValue.Value)
}

func (numericValue NumericValue[T]) ToSQL(_ *GormQuery) (string, []any, error) {
	return "", []any{numericValue.Value}, nil
}

type BoolValue struct {
	Value bool
}

func (boolValue BoolValue) GetValue() bool {
	return boolValue.Value
}

func (boolValue BoolValue) getSQL() toSQLFunc {
	return nil
}

func (boolValue BoolValue) ToSQL(_ *GormQuery) (string, []any, error) {
	return "", []any{boolValue.Value}, nil
}

type Value[T any] struct {
	Value T
}

func (value Value[T]) GetValue() T {
	return value.Value
}

func (value Value[T]) getSQL() toSQLFunc {
	// TODO intentar eliminar esta funcion
	return nil
}

func (value Value[T]) ToSQL(_ *GormQuery) (string, []any, error) {
	return "", []any{value.Value}, nil
}

type unsafeValue struct {
	Value IValue
}

func (unsafeValue unsafeValue) ToSQL(query *GormQuery) (string, []any, error) {
	return unsafeValue.Value.ToSQL(query)
}

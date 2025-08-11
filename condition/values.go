package condition

import "golang.org/x/exp/constraints"

type Numeric interface {
	constraints.Integer | constraints.Float
}

type NumericValue[T Numeric] struct {
	Value T
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

package dynamic

import (
	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/sql"
)

// Comparison Operators
// ref: https://www.postgresql.org/docs/current/functions-comparison.html

// EqualTo
func Eq[T any](field orm.FieldIdentifier[T]) orm.DynamicOperator[T] {
	return NewValueOperator(sql.Eq, field)
}

// NotEqualTo
func NotEq[T any](field orm.FieldIdentifier[T]) orm.DynamicOperator[T] {
	return NewValueOperator(sql.NotEq, field)
}

// LessThan
func Lt[T any](field orm.FieldIdentifier[T]) orm.DynamicOperator[T] {
	return NewValueOperator(sql.Lt, field)
}

// LessThanOrEqualTo
func LtOrEq[T any](field orm.FieldIdentifier[T]) orm.DynamicOperator[T] {
	return NewValueOperator(sql.LtOrEq, field)
}

// GreaterThan
func Gt[T any](field orm.FieldIdentifier[T]) orm.DynamicOperator[T] {
	return NewValueOperator(sql.Gt, field)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](field orm.FieldIdentifier[T]) orm.DynamicOperator[T] {
	return NewValueOperator(sql.GtOrEq, field)
}

// Comparison Predicates
// ref: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to field1 < value < field2
func Between[T any](field1, field2 orm.FieldIdentifier[T]) orm.DynamicOperator[T] {
	return newBetweenOperator(sql.Between, field1, field2)
}

// Equivalent to NOT (field1 < value < field2)
func NotBetween[T any](field1, field2 orm.FieldIdentifier[T]) orm.DynamicOperator[T] {
	return newBetweenOperator(sql.NotBetween, field1, field2)
}

func newBetweenOperator[T any](sqlOperator sql.Operator, field1, field2 orm.FieldIdentifier[T]) orm.DynamicOperator[T] {
	operator := NewValueOperator(sqlOperator, field1)
	return operator.AddOperation(sql.And, field2)
}

// Boolean Comparison Predicates

func IsDistinct[T any](field orm.FieldIdentifier[T]) orm.DynamicOperator[T] {
	return NewValueOperator(sql.IsDistinct, field)
}

func IsNotDistinct[T any](field orm.FieldIdentifier[T]) orm.DynamicOperator[T] {
	return NewValueOperator(sql.IsNotDistinct, field)
}

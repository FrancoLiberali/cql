package unsafe

import (
	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/sql"
)

// Comparison Operators
// ref: https://www.postgresql.org/docs/current/functions-comparison.html

// EqualTo
func Eq[T any](value any) orm.DynamicOperator[T] {
	return orm.NewValueOperator[T](sql.Eq, value)
}

// NotEqualTo
func NotEq[T any](value any) orm.DynamicOperator[T] {
	return orm.NewValueOperator[T](sql.NotEq, value)
}

// LessThan
func Lt[T any](value any) orm.DynamicOperator[T] {
	return orm.NewValueOperator[T](sql.Lt, value)
}

// LessThanOrEqualTo
func LtOrEq[T any](value any) orm.DynamicOperator[T] {
	return orm.NewValueOperator[T](sql.LtOrEq, value)
}

// GreaterThan
func Gt[T any](value any) orm.DynamicOperator[T] {
	return orm.NewValueOperator[T](sql.Gt, value)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](value any) orm.DynamicOperator[T] {
	return orm.NewValueOperator[T](sql.GtOrEq, value)
}

// Comparison Predicates
// ref: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to v1 < value < v2
func Between[T any](v1, v2 any) orm.DynamicOperator[T] {
	return newBetweenOperator[T](sql.Between, v1, v2)
}

// Equivalent to NOT (v1 < value < v2)
func NotBetween[T any](v1, v2 any) orm.DynamicOperator[T] {
	return newBetweenOperator[T](sql.NotBetween, v1, v2)
}

func newBetweenOperator[T any](sqlOperator sql.Operator, v1, v2 any) orm.DynamicOperator[T] {
	operator := orm.NewValueOperator[T](sqlOperator, v1)
	return operator.AddOperation(sql.And, v2)
}

// Boolean Comparison Predicates

func IsDistinct[T any](value any) orm.DynamicOperator[T] {
	return orm.NewValueOperator[T](sql.IsDistinct, value)
}

func IsNotDistinct[T any](value any) orm.DynamicOperator[T] {
	return orm.NewValueOperator[T](sql.IsNotDistinct, value)
}

// Row and Array Comparisons

func ArrayIn[T any](values ...any) orm.DynamicOperator[T] {
	return orm.NewValueOperator[T](sql.ArrayIn, values)
}

func ArrayNotIn[T any](values ...any) orm.DynamicOperator[T] {
	return orm.NewValueOperator[T](sql.ArrayNotIn, values)
}

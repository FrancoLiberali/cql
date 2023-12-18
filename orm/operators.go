package orm

import (
	"github.com/ditrit/badaas/orm/operator"
	"github.com/ditrit/badaas/orm/sql"
)

// Comparison Operators
// ref: https://www.postgresql.org/docs/current/functions-comparison.html

// EqualTo
// IsNotDistinct must be used in cases where value can be NULL
func Eq[T any](value T) operator.Operator[T] {
	return operator.NewValueOperator[T](sql.Eq, value)
}

// NotEqualTo
// IsDistinct must be used in cases where value can be NULL
func NotEq[T any](value T) operator.Operator[T] {
	return operator.NewValueOperator[T](sql.NotEq, value)
}

// LessThan
func Lt[T any](value T) operator.Operator[T] {
	return operator.NewValueOperator[T](sql.Lt, value)
}

// LessThanOrEqualTo
func LtOrEq[T any](value T) operator.Operator[T] {
	return operator.NewValueOperator[T](sql.LtOrEq, value)
}

// GreaterThan
func Gt[T any](value T) operator.Operator[T] {
	return operator.NewValueOperator[T](sql.Gt, value)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](value T) operator.Operator[T] {
	return operator.NewValueOperator[T](sql.GtOrEq, value)
}

// Comparison Predicates
// refs: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to v1 < value < v2
func Between[T any](v1 T, v2 T) operator.Operator[T] {
	return newBetweenOperator(sql.Between, v1, v2)
}

// Equivalent to NOT (v1 < value < v2)
func NotBetween[T any](v1 T, v2 T) operator.Operator[T] {
	return newBetweenOperator(sql.NotBetween, v1, v2)
}

func newBetweenOperator[T any](sqlOperator sql.Operator, v1 T, v2 T) operator.Operator[T] {
	operator := operator.NewValueOperator[T](sqlOperator, v1)
	return operator.AddOperation(sql.And, v2)
}

func IsNull[T any]() operator.Operator[T] {
	return operator.NewPredicateOperator[T]("IS NULL")
}

func IsNotNull[T any]() operator.Operator[T] {
	return operator.NewPredicateOperator[T]("IS NOT NULL")
}

// Boolean Comparison Predicates

func IsTrue() operator.Operator[bool] {
	return operator.NewPredicateOperator[bool]("IS TRUE")
}

func IsNotTrue() operator.Operator[bool] {
	return operator.NewPredicateOperator[bool]("IS NOT TRUE")
}

func IsFalse() operator.Operator[bool] {
	return operator.NewPredicateOperator[bool]("IS FALSE")
}

func IsNotFalse() operator.Operator[bool] {
	return operator.NewPredicateOperator[bool]("IS NOT FALSE")
}

func IsUnknown() operator.Operator[bool] {
	return operator.NewPredicateOperator[bool]("IS UNKNOWN")
}

func IsNotUnknown() operator.Operator[bool] {
	return operator.NewPredicateOperator[bool]("IS NOT UNKNOWN")
}

func IsDistinct[T any](value T) operator.Operator[T] {
	return operator.NewValueOperator[T](sql.IsDistinct, value)
}

func IsNotDistinct[T any](value T) operator.Operator[T] {
	return operator.NewValueOperator[T](sql.IsNotDistinct, value)
}

// Row and Array Comparisons

func ArrayIn[T any](values ...T) operator.Operator[T] {
	return operator.NewValueOperator[T](sql.ArrayIn, values)
}

func ArrayNotIn[T any](values ...T) operator.Operator[T] {
	return operator.NewValueOperator[T](sql.ArrayNotIn, values)
}

// Pattern Matching

type LikeOperator struct {
	operator.ValueOperator[string]
}

func NewLikeOperator(sqlOperator sql.Operator, pattern string) LikeOperator {
	return LikeOperator{
		ValueOperator: *operator.NewValueOperator[string](sqlOperator, pattern),
	}
}

func (operator LikeOperator) Escape(escape rune) operator.ValueOperator[string] {
	return *operator.AddOperation(sql.Escape, string(escape))
}

// Patterns:
//   - An underscore (_) in pattern stands for (matches) any single character.
//   - A percent sign (%) matches any sequence of zero or more characters.
//
// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-LIKE
func Like(pattern string) LikeOperator {
	return NewLikeOperator(sql.Like, pattern)
}

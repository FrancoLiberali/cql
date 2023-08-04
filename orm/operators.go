package orm

import (
	"github.com/ditrit/badaas/orm/sql"
)

// Comparison Operators
// ref: https://www.postgresql.org/docs/current/functions-comparison.html

// EqualTo
// IsNotDistinct must be used in cases where value can be NULL
func Eq[T any](value T) Operator[T] {
	return NewValueOperator[T](sql.Eq, value)
}

// NotEqualTo
// IsDistinct must be used in cases where value can be NULL
func NotEq[T any](value T) Operator[T] {
	return NewValueOperator[T](sql.NotEq, value)
}

// LessThan
func Lt[T any](value T) Operator[T] {
	return NewValueOperator[T](sql.Lt, value)
}

// LessThanOrEqualTo
func LtOrEq[T any](value T) Operator[T] {
	return NewValueOperator[T](sql.LtOrEq, value)
}

// GreaterThan
func Gt[T any](value T) Operator[T] {
	return NewValueOperator[T](sql.Gt, value)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](value T) Operator[T] {
	return NewValueOperator[T](sql.GtOrEq, value)
}

// Comparison Predicates
// refs: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to v1 < value < v2
func Between[T any](v1 T, v2 T) Operator[T] {
	return newBetweenOperator(sql.Between, v1, v2)
}

// Equivalent to NOT (v1 < value < v2)
func NotBetween[T any](v1 T, v2 T) Operator[T] {
	return newBetweenOperator(sql.NotBetween, v1, v2)
}

func newBetweenOperator[T any](sqlOperator sql.Operator, v1 T, v2 T) Operator[T] {
	operator := NewValueOperator[T](sqlOperator, v1)
	return operator.AddOperation(sql.And, v2)
}

func IsNull[T any]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NULL")
}

func IsNotNull[T any]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NOT NULL")
}

// Boolean Comparison Predicates

func IsTrue() PredicateOperator[bool] {
	return NewPredicateOperator[bool]("IS TRUE")
}

func IsNotTrue() PredicateOperator[bool] {
	return NewPredicateOperator[bool]("IS NOT TRUE")
}

func IsFalse() PredicateOperator[bool] {
	return NewPredicateOperator[bool]("IS FALSE")
}

func IsNotFalse() PredicateOperator[bool] {
	return NewPredicateOperator[bool]("IS NOT FALSE")
}

func IsUnknown() PredicateOperator[bool] {
	return NewPredicateOperator[bool]("IS UNKNOWN")
}

func IsNotUnknown() PredicateOperator[bool] {
	return NewPredicateOperator[bool]("IS NOT UNKNOWN")
}

func IsDistinct[T any](value T) ValueOperator[T] {
	return NewValueOperator[T](sql.IsDistinct, value)
}

func IsNotDistinct[T any](value T) ValueOperator[T] {
	return NewValueOperator[T](sql.IsNotDistinct, value)
}

// Row and Array Comparisons

func ArrayIn[T any](values ...T) ValueOperator[T] {
	return NewValueOperator[T](sql.ArrayIn, values)
}

func ArrayNotIn[T any](values ...T) ValueOperator[T] {
	return NewValueOperator[T](sql.ArrayNotIn, values)
}

// Pattern Matching

type LikeOperator struct {
	ValueOperator[string]
}

func NewLikeOperator(sqlOperator sql.Operator, pattern string) LikeOperator {
	return LikeOperator{
		ValueOperator: NewValueOperator[string](sqlOperator, pattern),
	}
}

func (operator LikeOperator) Escape(escape rune) ValueOperator[string] {
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

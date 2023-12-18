package operator

import (
	"github.com/ditrit/badaas/orm/sql"
)

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

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
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/queries/predicates?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

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

func IsNull[T any]() Operator[T] {
	return NewPredicateOperator[T]("IS NULL")
}

func IsNotNull[T any]() Operator[T] {
	return NewPredicateOperator[T]("IS NOT NULL")
}

// Boolean Comparison Predicates

// Not supported by: sqlserver
func IsTrue() PredicateOperator[bool] {
	return NewPredicateOperator[bool]("IS TRUE")
}

// Not supported by: sqlserver
func IsNotTrue() Operator[bool] {
	return NewPredicateOperator[bool]("IS NOT TRUE")
}

// Not supported by: sqlserver
func IsFalse() Operator[bool] {
	return NewPredicateOperator[bool]("IS FALSE")
}

// Not supported by: sqlserver
func IsNotFalse() Operator[bool] {
	return NewPredicateOperator[bool]("IS NOT FALSE")
}

// Not supported by: sqlserver, sqlite
func IsUnknown() Operator[bool] {
	return NewPredicateOperator[bool]("IS UNKNOWN")
}

// Not supported by: sqlserver, sqlite
func IsNotUnknown() Operator[bool] {
	return NewPredicateOperator[bool]("IS NOT UNKNOWN")
}

// Not supported by: mysql
func IsDistinct[T any](value T) Operator[T] {
	return NewValueOperator[T](sql.IsDistinct, value)
}

// Not supported by: mysql
func IsNotDistinct[T any](value T) Operator[T] {
	return NewValueOperator[T](sql.IsNotDistinct, value)
}

// Row and Array Comparisons

func In[T any](values []T) Operator[T] {
	return NewValueOperator[T](sql.ArrayIn, values)
}

func NotIn[T any](values []T) Operator[T] {
	return NewValueOperator[T](sql.ArrayNotIn, values)
}

// Pattern Matching

type LikeOperator struct {
	ValueOperator[string]
}

func NewLikeOperator(sqlOperator sql.Operator, pattern string) LikeOperator {
	return LikeOperator{
		ValueOperator: *NewValueOperator[string](sqlOperator, pattern),
	}
}

func (operator LikeOperator) Escape(escape rune) ValueOperator[string] {
	return *operator.AddOperation(sql.Escape, string(escape))
}

// Pattern in all databases:
//   - An underscore (_) in pattern stands for (matches) any single character.
//   - A percent sign (%) matches any sequence of zero or more characters.
//
// Additionally in SQLServer:
//   - Square brackets ([ ]) matches any single character within the specified range ([a-f]) or set ([abcdef]).
//   - [^] matches any single character not within the specified range ([^a-f]) or set ([^abcdef]).
//
// WARNINGS:
//   - SQLite: LIKE is case-insensitive unless case_sensitive_like pragma (https://www.sqlite.org/pragma.html#pragma_case_sensitive_like) is true.
//   - SQLServer, MySQL: the case-sensitivity depends on the collation used in compared column.
//   - PostgreSQL: LIKE is always case-sensitive, if you want case-insensitive use the ILIKE operator (implemented in psql.ILike)
//
// refs:
//   - mysql: https://dev.mysql.com/doc/refman/8.0/en/string-comparison-functions.html#operator_like
//   - postgresql: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-LIKE
//   - sqlserver: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/like-transact-sql?view=sql-server-ver16
//   - sqlite: https://www.sqlite.org/lang_expr.html#like
func Like(pattern string) LikeOperator {
	return NewLikeOperator(sql.Like, pattern)
}

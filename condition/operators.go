package condition

import (
	"github.com/FrancoLiberali/cql/sql"
)

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
// IsNotDistinct must be used in cases where value can be NULL
func Eq[T any](value IValue) Operator[T] {
	return NewValueOperator[T](sql.Eq, value)
}

// NotEqualTo
// IsDistinct must be used in cases where value can be NULL
func NotEq[T any](value IValue) Operator[T] {
	return NewValueOperator[T](sql.NotEq, value)
}

// LessThan
func Lt[T any](value IValue) Operator[T] {
	return NewValueOperator[T](sql.Lt, value)
}

// LessThanOrEqualTo
func LtOrEq[T any](value IValue) Operator[T] {
	return NewValueOperator[T](sql.LtOrEq, value)
}

// GreaterThan
func Gt[T any](value IValue) Operator[T] {
	return NewValueOperator[T](sql.Gt, value)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](value IValue) Operator[T] {
	return NewValueOperator[T](sql.GtOrEq, value)
}

// Comparison Predicates
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/queries/predicates?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// Equivalent to v1 < value < v2
func Between[T any](v1, v2 IValue) Operator[T] {
	return newBetweenOperator[T](sql.Between, v1, v2)
}

// Equivalent to NOT (v1 < value < v2)
func NotBetween[T any](v1, v2 IValue) Operator[T] {
	return newBetweenOperator[T](sql.NotBetween, v1, v2)
}

func newBetweenOperator[T any](sqlOperator sql.Operator, v1, v2 IValue) Operator[T] {
	operator := NewValueOperator[T](sqlOperator, v1)
	return operator.AddOperation(sql.And, v2)
}

func IsNull[T any]() Operator[T] {
	return NewPredicateOperator[T]("IS NULL")
}

func IsNotNull[T any]() Operator[T] {
	return NewPredicateOperator[T]("IS NOT NULL")
}

func IsDistinct[T any](value IValue) Operator[T] {
	isNotDistinct := new(ValueOperator[T]).AddOperation(
		map[sql.Dialector]sql.Operator{
			sql.Postgres:  sql.IsDistinct,
			sql.SQLServer: sql.IsDistinct,
			sql.SQLite:    sql.IsDistinct,
			sql.MySQL:     sql.MySQLNullSafeEqual,
		},
		value,
	)
	isNotDistinct.Modifier = map[sql.Dialector]string{ //nolint:exhaustive // not present is expected for the other ones
		sql.MySQL: "NOT",
	}

	return isNotDistinct
}

func IsNotDistinct[T any](value IValue) Operator[T] {
	return new(ValueOperator[T]).AddOperation(
		map[sql.Dialector]sql.Operator{
			sql.Postgres:  sql.IsNotDistinct,
			sql.SQLServer: sql.IsNotDistinct,
			sql.SQLite:    sql.IsNotDistinct,
			sql.MySQL:     sql.MySQLNullSafeEqual,
		},
		value,
	)
}

// Row and Array Comparisons

type IValueList[T any] []ValueOfType[T]

func (values IValueList[T]) ToSQL(_ *CQLQuery) (string, []any, error) {
	valuesAny := make([]any, 0, len(values))
	for _, value := range values {
		valuesAny = append(valuesAny, value.GetValue())
	}

	return "", valuesAny, nil
}

func In[T any](values IValueList[T]) Operator[T] {
	return NewValueOperator[T](sql.ArrayIn, values)
}

func NotIn[T any](values IValueList[T]) Operator[T] {
	return NewValueOperator[T](sql.ArrayNotIn, values)
}

// Pattern Matching

type LikeOperator struct {
	ValueOperator[string]
}

func NewLikeOperator(sqlOperator sql.Operator, pattern string) LikeOperator {
	return LikeOperator{
		ValueOperator: *NewValueOperator[string](sqlOperator, String(pattern)),
	}
}

func (operator LikeOperator) Escape(escape rune) ValueOperator[string] {
	return *operator.AddOperation(sql.Escape, String(string(escape)))
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

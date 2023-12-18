package unsafe

import (
	"github.com/ditrit/badaas/orm/operator"
	"github.com/ditrit/badaas/orm/sql"
)

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
func Eq[T any](value any) operator.DynamicOperator[T] {
	return operator.NewValueOperator[T](sql.Eq, value)
}

// NotEqualTo
func NotEq[T any](value any) operator.DynamicOperator[T] {
	return operator.NewValueOperator[T](sql.NotEq, value)
}

// LessThan
func Lt[T any](value any) operator.DynamicOperator[T] {
	return operator.NewValueOperator[T](sql.Lt, value)
}

// LessThanOrEqualTo
func LtOrEq[T any](value any) operator.DynamicOperator[T] {
	return operator.NewValueOperator[T](sql.LtOrEq, value)
}

// GreaterThan
func Gt[T any](value any) operator.DynamicOperator[T] {
	return operator.NewValueOperator[T](sql.Gt, value)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](value any) operator.DynamicOperator[T] {
	return operator.NewValueOperator[T](sql.GtOrEq, value)
}

// Comparison Predicates
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/queries/predicates?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// Equivalent to v1 < value < v2
func Between[T any](v1, v2 any) operator.DynamicOperator[T] {
	return newBetweenOperator[T](sql.Between, v1, v2)
}

// Equivalent to NOT (v1 < value < v2)
func NotBetween[T any](v1, v2 any) operator.DynamicOperator[T] {
	return newBetweenOperator[T](sql.NotBetween, v1, v2)
}

func newBetweenOperator[T any](sqlOperator sql.Operator, v1, v2 any) operator.DynamicOperator[T] {
	operator := operator.NewValueOperator[T](sqlOperator, v1)
	return operator.AddOperation(sql.And, v2)
}

// Not supported by: mysql
func IsDistinct[T any](value any) operator.DynamicOperator[T] {
	return operator.NewValueOperator[T](sql.IsDistinct, value)
}

// Not supported by: mysql
func IsNotDistinct[T any](value any) operator.DynamicOperator[T] {
	return operator.NewValueOperator[T](sql.IsNotDistinct, value)
}

// Row and Array Comparisons

func ArrayIn[T any](values ...any) operator.DynamicOperator[T] {
	return operator.NewValueOperator[T](sql.ArrayIn, values)
}

func ArrayNotIn[T any](values ...any) operator.DynamicOperator[T] {
	return operator.NewValueOperator[T](sql.ArrayNotIn, values)
}

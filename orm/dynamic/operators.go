package dynamic

import (
	"github.com/ditrit/badaas/orm/operator"
	"github.com/ditrit/badaas/orm/query"
	"github.com/ditrit/badaas/orm/sql"
)

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
func Eq[T any](field query.FieldIdentifier[T]) operator.DynamicOperator[T] {
	return newValueOperator(sql.Eq, field)
}

// NotEqualTo
func NotEq[T any](field query.FieldIdentifier[T]) operator.DynamicOperator[T] {
	return newValueOperator(sql.NotEq, field)
}

// LessThan
func Lt[T any](field query.FieldIdentifier[T]) operator.DynamicOperator[T] {
	return newValueOperator(sql.Lt, field)
}

// LessThanOrEqualTo
func LtOrEq[T any](field query.FieldIdentifier[T]) operator.DynamicOperator[T] {
	return newValueOperator(sql.LtOrEq, field)
}

// GreaterThan
func Gt[T any](field query.FieldIdentifier[T]) operator.DynamicOperator[T] {
	return newValueOperator(sql.Gt, field)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](field query.FieldIdentifier[T]) operator.DynamicOperator[T] {
	return newValueOperator(sql.GtOrEq, field)
}

// Comparison Predicates
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/queries/predicates?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// Equivalent to field1 < value < field2
func Between[T any](field1, field2 query.FieldIdentifier[T]) operator.DynamicOperator[T] {
	return newBetweenOperator(sql.Between, field1, field2)
}

// Equivalent to NOT (field1 < value < field2)
func NotBetween[T any](field1, field2 query.FieldIdentifier[T]) operator.DynamicOperator[T] {
	return newBetweenOperator(sql.NotBetween, field1, field2)
}

func newBetweenOperator[T any](sqlOperator sql.Operator, field1, field2 query.FieldIdentifier[T]) operator.DynamicOperator[T] {
	operator := newValueOperator(sqlOperator, field1)
	return operator.AddOperation(sql.And, field2)
}

// Not supported by: mysql
func IsDistinct[T any](field query.FieldIdentifier[T]) operator.DynamicOperator[T] {
	return newValueOperator(sql.IsDistinct, field)
}

// Not supported by: mysql
func IsNotDistinct[T any](field query.FieldIdentifier[T]) operator.DynamicOperator[T] {
	return newValueOperator(sql.IsNotDistinct, field)
}

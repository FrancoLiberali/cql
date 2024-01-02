package condition

import (
	"github.com/FrancoLiberali/cql/model"
)

type DynamicFieldIs[TObject model.Model, TAttribute any] struct {
	field Field[TObject, TAttribute]
}

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
func (is DynamicFieldIs[TObject, TAttribute]) Eq(value ValueOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Eq[TAttribute](value))
}

// NotEqualTo
func (is DynamicFieldIs[TObject, TAttribute]) NotEq(value ValueOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, NotEq[TAttribute](value))
}

// LessThan
func (is DynamicFieldIs[TObject, TAttribute]) Lt(value ValueOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Lt[TAttribute](value))
}

// LessThanOrEqualTo
func (is DynamicFieldIs[TObject, TAttribute]) LtOrEq(value ValueOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, LtOrEq[TAttribute](value))
}

// GreaterThan
func (is DynamicFieldIs[TObject, TAttribute]) Gt(value ValueOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Gt[TAttribute](value))
}

// GreaterThanOrEqualTo
func (is DynamicFieldIs[TObject, TAttribute]) GtOrEq(value ValueOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, GtOrEq[TAttribute](value))
}

// Comparison Predicates
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/queries/predicates?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// Equivalent to field1 < value < field2
func (is DynamicFieldIs[TObject, TAttribute]) Between(value1, value2 ValueOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Between[TAttribute](value1, value2))
}

// Equivalent to NOT (field1 < value < field2)
func (is DynamicFieldIs[TObject, TAttribute]) NotBetween(value1, value2 ValueOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, NotBetween[TAttribute](value1, value2))
}

func (is DynamicFieldIs[TObject, TAttribute]) Distinct(value ValueOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, IsDistinct[TAttribute](value))
}

func (is DynamicFieldIs[TObject, TAttribute]) NotDistinct(value ValueOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, IsNotDistinct[TAttribute](value))
}

type NumericDynamicFieldIs[TObject model.Model, TAttribute any] struct {
	field Field[TObject, TAttribute]
}

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
func (is NumericDynamicFieldIs[TObject, TAttribute]) Eq(value ValueOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Eq[TAttribute](value))
}

// NotEqualTo
func (is NumericDynamicFieldIs[TObject, TAttribute]) NotEq(value ValueOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, NotEq[TAttribute](value))
}

// LessThan
func (is NumericDynamicFieldIs[TObject, TAttribute]) Lt(value ValueOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Lt[TAttribute](value))
}

// LessThanOrEqualTo
func (is NumericDynamicFieldIs[TObject, TAttribute]) LtOrEq(value ValueOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, LtOrEq[TAttribute](value))
}

// GreaterThan
func (is NumericDynamicFieldIs[TObject, TAttribute]) Gt(value ValueOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Gt[TAttribute](value))
}

// GreaterThanOrEqualTo
func (is NumericDynamicFieldIs[TObject, TAttribute]) GtOrEq(value ValueOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, GtOrEq[TAttribute](value))
}

// Comparison Predicates
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/queries/predicates?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// Equivalent to field1 < value < field2
func (is NumericDynamicFieldIs[TObject, TAttribute]) Between(value1, value2 ValueOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Between[TAttribute](value1, value2))
}

// Equivalent to NOT (field1 < value < field2)
func (is NumericDynamicFieldIs[TObject, TAttribute]) NotBetween(value1, value2 ValueOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, NotBetween[TAttribute](value1, value2))
}

func (is NumericDynamicFieldIs[TObject, TAttribute]) Distinct(value ValueOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, IsDistinct[TAttribute](value))
}

func (is NumericDynamicFieldIs[TObject, TAttribute]) NotDistinct(value ValueOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, IsNotDistinct[TAttribute](value))
}

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
func (is DynamicFieldIs[TObject, TAttribute]) Eq(field FieldOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Eq[TAttribute](field))
}

// NotEqualTo
func (is DynamicFieldIs[TObject, TAttribute]) NotEq(field FieldOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, NotEq[TAttribute](field))
}

// LessThan
func (is DynamicFieldIs[TObject, TAttribute]) Lt(field FieldOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Lt[TAttribute](field))
}

// LessThanOrEqualTo
func (is DynamicFieldIs[TObject, TAttribute]) LtOrEq(field FieldOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, LtOrEq[TAttribute](field))
}

// GreaterThan
func (is DynamicFieldIs[TObject, TAttribute]) Gt(field FieldOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Gt[TAttribute](field))
}

// GreaterThanOrEqualTo
func (is DynamicFieldIs[TObject, TAttribute]) GtOrEq(field FieldOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, GtOrEq[TAttribute](field))
}

// Comparison Predicates
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/queries/predicates?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// Equivalent to field1 < value < field2
func (is DynamicFieldIs[TObject, TAttribute]) Between(field1, field2 FieldOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Between[TAttribute](field1, field2))
}

// Equivalent to NOT (field1 < value < field2)
func (is DynamicFieldIs[TObject, TAttribute]) NotBetween(field1, field2 FieldOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, NotBetween[TAttribute](field1, field2))
}

func (is DynamicFieldIs[TObject, TAttribute]) Distinct(field FieldOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, IsDistinct[TAttribute](field))
}

func (is DynamicFieldIs[TObject, TAttribute]) NotDistinct(field FieldOfType[TAttribute]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, IsNotDistinct[TAttribute](field))
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
func (is NumericDynamicFieldIs[TObject, TAttribute]) Eq(field FieldOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Eq[TAttribute](field))
}

// NotEqualTo
func (is NumericDynamicFieldIs[TObject, TAttribute]) NotEq(field FieldOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, NotEq[TAttribute](field))
}

// LessThan
func (is NumericDynamicFieldIs[TObject, TAttribute]) Lt(field FieldOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Lt[TAttribute](field))
}

// LessThanOrEqualTo
func (is NumericDynamicFieldIs[TObject, TAttribute]) LtOrEq(field FieldOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, LtOrEq[TAttribute](field))
}

// GreaterThan
func (is NumericDynamicFieldIs[TObject, TAttribute]) Gt(field FieldOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Gt[TAttribute](field))
}

// GreaterThanOrEqualTo
func (is NumericDynamicFieldIs[TObject, TAttribute]) GtOrEq(field FieldOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, GtOrEq[TAttribute](field))
}

// Comparison Predicates
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/queries/predicates?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// Equivalent to field1 < value < field2
func (is NumericDynamicFieldIs[TObject, TAttribute]) Between(field1, field2 FieldOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, Between[TAttribute](field1, field2))
}

// Equivalent to NOT (field1 < value < field2)
func (is NumericDynamicFieldIs[TObject, TAttribute]) NotBetween(field1, field2 FieldOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, NotBetween[TAttribute](field1, field2))
}

func (is NumericDynamicFieldIs[TObject, TAttribute]) Distinct(field FieldOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, IsDistinct[TAttribute](field))
}

func (is NumericDynamicFieldIs[TObject, TAttribute]) NotDistinct(field FieldOfType[numeric]) DynamicCondition[TObject] {
	return NewFieldCondition(is.field, IsNotDistinct[TAttribute](field))
}

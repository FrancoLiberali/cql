package orm

import (
	"github.com/ditrit/badaas/orm/condition"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/operator"
	"github.com/ditrit/badaas/orm/query"
)

type FieldIs[TObject model.Model, TAttribute any] struct {
	FieldID query.FieldIdentifier[TAttribute]
}

type BoolFieldIs[TObject model.Model] struct {
	FieldIs[TObject, bool]
}

type StringFieldIs[TObject model.Model] struct {
	FieldIs[TObject, string]
}

// EqualTo
// NotDistinct must be used in cases where value can be NULL
func (is FieldIs[TObject, TAttribute]) Eq(value TAttribute) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.Eq(value))
}

// NotEqualTo
// Distinct must be used in cases where value can be NULL
func (is FieldIs[TObject, TAttribute]) NotEq(value TAttribute) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.NotEq(value))
}

// LessThan
func (is FieldIs[TObject, TAttribute]) Lt(value TAttribute) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.Lt(value))
}

// LessThanOrEqualTo
func (is FieldIs[TObject, TAttribute]) LtOrEq(value TAttribute) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.LtOrEq(value))
}

// GreaterThan
func (is FieldIs[TObject, TAttribute]) Gt(value TAttribute) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.Gt(value))
}

// GreaterThanOrEqualTo
func (is FieldIs[TObject, TAttribute]) GtOrEq(value TAttribute) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.GtOrEq(value))
}

// Equivalent to v1 < value < v2
func (is FieldIs[TObject, TAttribute]) Between(v1, v2 TAttribute) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.Between(v1, v2))
}

// Equivalent to NOT (v1 < value < v2)
func (is FieldIs[TObject, TAttribute]) NotBetween(v1, v2 TAttribute) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.NotBetween(v1, v2))
}

func (is FieldIs[TObject, TAttribute]) Null() condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.IsNull[TAttribute]())
}

func (is FieldIs[TObject, TAttribute]) NotNull() condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.IsNotNull[TAttribute]())
}

// Not supported by: sqlserver
func (is BoolFieldIs[TObject]) True() condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, bool](is.FieldID, operator.IsTrue())
}

// Not supported by: sqlserver
func (is BoolFieldIs[TObject]) NotTrue() condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, bool](is.FieldID, operator.IsNotTrue())
}

// Not supported by: sqlserver
func (is BoolFieldIs[TObject]) False() condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, bool](is.FieldID, operator.IsFalse())
}

// Not supported by: sqlserver
func (is BoolFieldIs[TObject]) NotFalse() condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, bool](is.FieldID, operator.IsNotFalse())
}

// Not supported by: sqlserver, sqlite
func (is BoolFieldIs[TObject]) Unknown() condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, bool](is.FieldID, operator.IsUnknown())
}

// Not supported by: sqlserver, sqlite
func (is BoolFieldIs[TObject]) NotUnknown() condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, bool](is.FieldID, operator.IsNotUnknown())
}

// Not supported by: mysql
func (is FieldIs[TObject, TAttribute]) Distinct(value TAttribute) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.IsDistinct(value))
}

// Not supported by: mysql
func (is FieldIs[TObject, TAttribute]) NotDistinct(value TAttribute) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.IsNotDistinct(value))
}

func (is FieldIs[TObject, TAttribute]) In(values ...TAttribute) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.In(values))
}

func (is FieldIs[TObject, TAttribute]) NotIn(values ...TAttribute) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, operator.NotIn(values))
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
func (is StringFieldIs[TObject]) Like(pattern string) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, string](is.FieldID, operator.Like(pattern))
}

// Custom can be used to use other Operators, like database specific operators
func (is FieldIs[TObject, TAttribute]) Custom(op operator.Operator[TAttribute]) condition.WhereCondition[TObject] {
	return condition.NewFieldCondition[TObject, TAttribute](is.FieldID, op)
}

// Dynamic transforms the FieldIs in a DynamicFieldIs to use dynamic operators
func (is FieldIs[TObject, TAttribute]) Dynamic() DynamicFieldIs[TObject, TAttribute] {
	return DynamicFieldIs[TObject, TAttribute]{
		fieldID: is.FieldID,
	}
}

// Unsafe transforms the FieldIs in an UnsafeFieldIs to use unsafe operators
func (is FieldIs[TObject, TAttribute]) Unsafe() UnsafeFieldIs[TObject, TAttribute] {
	return UnsafeFieldIs[TObject, TAttribute]{
		fieldID: is.FieldID,
	}
}

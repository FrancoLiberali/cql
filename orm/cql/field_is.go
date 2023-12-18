package cql

import (
	"github.com/ditrit/badaas/orm/model"
)

type FieldIs[TObject model.Model, TAttribute any] struct {
	Field Field[TObject, TAttribute]
}

type BoolFieldIs[TObject model.Model] struct {
	Field Field[TObject, bool]
}

type StringFieldIs[TObject model.Model] struct {
	FieldIs[TObject, string]
}

// EqualTo
// NotDistinct must be used in cases where value can be NULL
func (is FieldIs[TObject, TAttribute]) Eq(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.Field, Eq[TAttribute](value))
}

// NotEqualTo
// Distinct must be used in cases where value can be NULL
func (is FieldIs[TObject, TAttribute]) NotEq(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.Field, NotEq[TAttribute](value))
}

// LessThan
func (is FieldIs[TObject, TAttribute]) Lt(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.Field, Lt[TAttribute](value))
}

// LessThanOrEqualTo
func (is FieldIs[TObject, TAttribute]) LtOrEq(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.Field, LtOrEq[TAttribute](value))
}

// GreaterThan
func (is FieldIs[TObject, TAttribute]) Gt(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.Field, Gt[TAttribute](value))
}

// GreaterThanOrEqualTo
func (is FieldIs[TObject, TAttribute]) GtOrEq(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.Field, GtOrEq[TAttribute](value))
}

// Equivalent to v1 < value < v2
func (is FieldIs[TObject, TAttribute]) Between(v1, v2 TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.Field, Between[TAttribute](v1, v2))
}

// Equivalent to NOT (v1 < value < v2)
func (is FieldIs[TObject, TAttribute]) NotBetween(v1, v2 TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.Field, NotBetween[TAttribute](v1, v2))
}

func (is FieldIs[TObject, TAttribute]) Null() WhereCondition[TObject] {
	return NewFieldCondition(is.Field, IsNull[TAttribute]())
}

func (is FieldIs[TObject, TAttribute]) NotNull() WhereCondition[TObject] {
	return NewFieldCondition(is.Field, IsNotNull[TAttribute]())
}

func (is BoolFieldIs[TObject]) True() WhereCondition[TObject] {
	return NewFieldCondition[TObject, bool](is.Field, Eq[bool](true))
}

func (is BoolFieldIs[TObject]) NotTrue() WhereCondition[TObject] {
	return NewFieldCondition[TObject, bool](is.Field, IsDistinct[bool](true))
}

func (is BoolFieldIs[TObject]) False() WhereCondition[TObject] {
	return NewFieldCondition[TObject, bool](is.Field, Eq[bool](false))
}

func (is BoolFieldIs[TObject]) NotFalse() WhereCondition[TObject] {
	return NewFieldCondition[TObject, bool](is.Field, IsDistinct[bool](false))
}

func (is BoolFieldIs[TObject]) Unknown() WhereCondition[TObject] {
	return NewFieldCondition[TObject, bool](is.Field, IsNull[bool]())
}

func (is BoolFieldIs[TObject]) NotUnknown() WhereCondition[TObject] {
	return NewFieldCondition[TObject, bool](is.Field, IsNotNull[bool]())
}

func (is FieldIs[TObject, TAttribute]) Distinct(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.Field, IsDistinct[TAttribute](value))
}

func (is FieldIs[TObject, TAttribute]) NotDistinct(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.Field, IsNotDistinct[TAttribute](value))
}

func (is FieldIs[TObject, TAttribute]) In(values ...TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.Field, In(values))
}

func (is FieldIs[TObject, TAttribute]) NotIn(values ...TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.Field, NotIn(values))
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
func (is StringFieldIs[TObject]) Like(pattern string) WhereCondition[TObject] {
	return NewFieldCondition[TObject, string](is.Field, Like(pattern))
}

// Custom can be used to use other Operators, like database specific operators
func (is FieldIs[TObject, TAttribute]) Custom(op Operator[TAttribute]) WhereCondition[TObject] {
	return NewFieldCondition(is.Field, op)
}

// Dynamic transforms the FieldIs in a DynamicFieldIs to use dynamic operators
func (is FieldIs[TObject, TAttribute]) Dynamic() DynamicFieldIs[TObject, TAttribute] {
	return DynamicFieldIs[TObject, TAttribute]{
		field: is.Field,
	}
}

// Unsafe transforms the FieldIs in an UnsafeFieldIs to use unsafe operators
func (is FieldIs[TObject, TAttribute]) Unsafe() UnsafeFieldIs[TObject, TAttribute] {
	return UnsafeFieldIs[TObject, TAttribute]{
		field: is.Field,
	}
}

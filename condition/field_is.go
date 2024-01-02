package condition

import (
	"github.com/FrancoLiberali/cql/model"
)

type FieldIs[TObject model.Model, TAttribute any] struct {
	field Field[TObject, TAttribute]
}

type BoolFieldIs[TObject model.Model] struct {
	field Field[TObject, bool]
}

type StringFieldIs[TObject model.Model] struct {
	FieldIs[TObject, string]
}

// EqualTo
// NotDistinct must be used in cases where value can be NULL
func (is FieldIs[TObject, TAttribute]) Eq(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.field, Eq[TAttribute](value))
}

// NotEqualTo
// Distinct must be used in cases where value can be NULL
func (is FieldIs[TObject, TAttribute]) NotEq(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.field, NotEq[TAttribute](value))
}

// LessThan
func (is FieldIs[TObject, TAttribute]) Lt(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.field, Lt[TAttribute](value))
}

// LessThanOrEqualTo
func (is FieldIs[TObject, TAttribute]) LtOrEq(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.field, LtOrEq[TAttribute](value))
}

// GreaterThan
func (is FieldIs[TObject, TAttribute]) Gt(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.field, Gt[TAttribute](value))
}

// GreaterThanOrEqualTo
func (is FieldIs[TObject, TAttribute]) GtOrEq(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.field, GtOrEq[TAttribute](value))
}

// Equivalent to v1 < value < v2
func (is FieldIs[TObject, TAttribute]) Between(v1, v2 TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.field, Between[TAttribute](v1, v2))
}

// Equivalent to NOT (v1 < value < v2)
func (is FieldIs[TObject, TAttribute]) NotBetween(v1, v2 TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.field, NotBetween[TAttribute](v1, v2))
}

func (is FieldIs[TObject, TAttribute]) Null() WhereCondition[TObject] {
	return NewFieldCondition(is.field, IsNull[TAttribute]())
}

func (is FieldIs[TObject, TAttribute]) NotNull() WhereCondition[TObject] {
	return NewFieldCondition(is.field, IsNotNull[TAttribute]())
}

func (is BoolFieldIs[TObject]) True() WhereCondition[TObject] {
	return NewFieldCondition[TObject, bool](is.field, Eq[bool](true))
}

func (is BoolFieldIs[TObject]) NotTrue() WhereCondition[TObject] {
	return NewFieldCondition[TObject, bool](is.field, IsDistinct[bool](true))
}

func (is BoolFieldIs[TObject]) False() WhereCondition[TObject] {
	return NewFieldCondition[TObject, bool](is.field, Eq[bool](false))
}

func (is BoolFieldIs[TObject]) NotFalse() WhereCondition[TObject] {
	return NewFieldCondition[TObject, bool](is.field, IsDistinct[bool](false))
}

func (is BoolFieldIs[TObject]) Unknown() WhereCondition[TObject] {
	return NewFieldCondition[TObject, bool](is.field, IsNull[bool]())
}

func (is BoolFieldIs[TObject]) NotUnknown() WhereCondition[TObject] {
	return NewFieldCondition[TObject, bool](is.field, IsNotNull[bool]())
}

func (is FieldIs[TObject, TAttribute]) Distinct(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.field, IsDistinct[TAttribute](value))
}

func (is FieldIs[TObject, TAttribute]) NotDistinct(value TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.field, IsNotDistinct[TAttribute](value))
}

func (is FieldIs[TObject, TAttribute]) In(values ...TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.field, In(values))
}

func (is FieldIs[TObject, TAttribute]) NotIn(values ...TAttribute) WhereCondition[TObject] {
	return NewFieldCondition(is.field, NotIn(values))
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
	return NewFieldCondition[TObject, string](is.field, Like(pattern))
}

// Custom can be used to use other Operators, like database specific operators
func (is FieldIs[TObject, TAttribute]) Custom(op Operator[TAttribute]) WhereCondition[TObject] {
	return NewFieldCondition(is.field, op)
}

package condition

import (
	"strings"

	"github.com/FrancoLiberali/cql/model"
)

type IValue interface {
	getField() IField
	toSQL(query *GormQuery, table Table) (string, []any)
}

type ValueOfType[T any] interface {
	IValue
	getType() T
}

type FieldValue[TModel model.Model, TAttribute any] struct {
	field  Field[TModel, TAttribute]
	sql    string
	values []any
}

const columnSQL = "COLUMN_SQL"

func NewFieldValue[TModel model.Model, TAttribute any](field Field[TModel, TAttribute]) *FieldValue[TModel, TAttribute] {
	return &FieldValue[TModel, TAttribute]{
		field:  field,
		values: []any{},
		sql:    columnSQL,
	}
}

func (value FieldValue[TModel, TAttribute]) getField() IField {
	return value.field
}

func (value FieldValue[TModel, TAttribute]) toSQL(query *GormQuery, table Table) (string, []any) {
	return strings.Replace(
		value.sql,
		columnSQL,
		value.field.columnSQL(query, table),
		1,
	), value.values
}

//nolint:unused // necessary for ValueOfType[T any]
func (value FieldValue[TModel, TAttribute]) getType() TAttribute {
	return *new(TAttribute)
}

type NumericFieldValue[TModel model.Model, TAttribute any] struct {
	FieldValue[TModel, TAttribute]
}

type numeric struct{}

//nolint:unused // necessary for ValueOfType[Numeric]
func (value NumericFieldValue[TModel, TAttribute]) getType() numeric {
	return numeric{}
}

// Plus sums other to value
//
// Warning: in PostgreSQL the value received by parameter could be casted to integer
func (value *NumericFieldValue[TModel, TAttribute]) Plus(other float64) *NumericFieldValue[TModel, TAttribute] {
	return value.addOperation(other, "+")
}

// Minus subtracts other from the value
//
// Warning: in PostgreSQL the value received by parameter could be casted to integer
func (value *NumericFieldValue[TModel, TAttribute]) Minus(other float64) *NumericFieldValue[TModel, TAttribute] {
	return value.addOperation(other, "-")
}

// Times multiplies value by other
//
// Warning: in PostgreSQL the value received by parameter could be casted to integer
func (value *NumericFieldValue[TModel, TAttribute]) Times(other float64) *NumericFieldValue[TModel, TAttribute] {
	return value.addOperation(other, "*")
}

// Divided divides value by other
//
// Warning: in PostgreSQL the value received by parameter could be casted to integer
func (value *NumericFieldValue[TModel, TAttribute]) Divided(other float64) *NumericFieldValue[TModel, TAttribute] {
	return value.addOperation(other, "/")
}

func (value *NumericFieldValue[TModel, TAttribute]) addOperation(other float64, operator string) *NumericFieldValue[TModel, TAttribute] {
	value.sql = "(" + value.sql + " " + operator + " ?)"
	value.values = append(value.values, other)

	return value
}

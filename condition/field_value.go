package condition

import (
	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

type IValue interface {
	getField() IField
	toSQL(query *GormQuery, table Table) (string, []any, error)
}

type ValueOfType[T any] interface {
	IValue
	getType() T
}

type FieldValue[TModel model.Model, TAttribute any] struct {
	field     Field[TModel, TAttribute]
	values    []any
	functions []sql.FunctionByDialector
}

func NewFieldValue[TModel model.Model, TAttribute any](field Field[TModel, TAttribute]) *FieldValue[TModel, TAttribute] {
	return &FieldValue[TModel, TAttribute]{
		field:     field,
		values:    []any{},
		functions: []sql.FunctionByDialector{},
	}
}

func (value FieldValue[TModel, TAttribute]) getField() IField {
	return value.field
}

func (value *FieldValue[TModel, TAttribute]) addFunction(other any, function sql.FunctionByDialector) {
	value.functions = append(value.functions, function)
	value.values = append(value.values, other)
}

func (value FieldValue[TModel, TAttribute]) toSQL(query *GormQuery, table Table) (string, []any, error) {
	finalSQL := value.field.columnSQL(query, table)

	for _, functionByDialector := range value.functions {
		function, isPresent := functionByDialector.Get(query.Dialector())
		if !isPresent {
			return "", nil, functionError(ErrUnsupportedByDatabase, function)
		}

		finalSQL = function.ApplyTo(finalSQL)
	}

	return finalSQL, value.values, nil
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
func (value *NumericFieldValue[TModel, TAttribute]) Plus(other float64) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(other, sql.Plus)
	return value
}

// Minus subtracts other from the value
func (value *NumericFieldValue[TModel, TAttribute]) Minus(other float64) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(other, sql.Minus)
	return value
}

// Times multiplies value by other
func (value *NumericFieldValue[TModel, TAttribute]) Times(other float64) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(other, sql.Times)
	return value
}

// Divided divides value by other
func (value *NumericFieldValue[TModel, TAttribute]) Divided(other float64) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(other, sql.Divided)
	return value
}

// Modulo returns the remainder of the entire division
func (value *NumericFieldValue[TModel, TAttribute]) Modulo(other float64) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(other, sql.Modulo)
	return value
}

// Power elevates value to other
func (value *NumericFieldValue[TModel, TAttribute]) Power(other float64) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(other, sql.Power)
	return value
}

type StringFieldValue[TModel model.Model] struct {
	FieldValue[TModel, string]
}

// Concat concatenates other to value
func (value *StringFieldValue[TModel]) Concat(other string) *StringFieldValue[TModel] {
	value.addFunction(other, sql.Concat)
	return value
}

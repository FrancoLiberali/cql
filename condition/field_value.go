package condition

import (
	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

type IValue interface {
	// TODO agregar seccion a la doc de como trabajar con custom types: agregar GetValue y ToSQL
	ToSQL(query *CQLQuery) (string, []any, error)
}

type ValueOfType[T any] interface {
	IValue

	GetValue() T
}

type functionAndValues struct {
	function sql.FunctionByDialector
	values   int
}

type FieldValue[TModel model.Model, TAttribute any] struct {
	field     Field[TModel, TAttribute]
	values    []any
	functions []functionAndValues
}

func NewFieldValue[TModel model.Model, TAttribute any](field Field[TModel, TAttribute]) *FieldValue[TModel, TAttribute] {
	return &FieldValue[TModel, TAttribute]{
		field:     field,
		values:    []any{},
		functions: []functionAndValues{},
	}
}

func (value *FieldValue[TModel, TAttribute]) addFunction(function sql.FunctionByDialector, others ...any) {
	value.functions = append(value.functions, functionAndValues{function: function, values: len(others)})
	value.values = append(value.values, others...)
}

func (value FieldValue[TModel, TAttribute]) ToSQL(query *CQLQuery) (string, []any, error) {
	table, err := getModelTable(query, value.field)
	if err != nil {
		return "", nil, err
	}

	finalSQL := value.field.columnSQL(query, table)

	for _, functionAndValues := range value.functions {
		function, isPresent := functionAndValues.function.Get(query.Dialector())
		if !isPresent {
			return "", nil, functionError(ErrUnsupportedByDatabase, functionAndValues.function)
		}

		finalSQL = function.ApplyTo(finalSQL, functionAndValues.values)
	}

	return finalSQL, value.values, nil
}

func (value FieldValue[TModel, TAttribute]) GetValue() TAttribute {
	return *new(TAttribute)
}

type NumericFieldValue[TModel model.Model, TAttribute any] struct {
	FieldValue[TModel, TAttribute]
}

func (value NumericFieldValue[TModel, TAttribute]) GetValue() float64 {
	return 0
}

// Plus sums other to value
func (value *NumericFieldValue[TModel, TAttribute]) Plus(other float64) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.Plus, other)
	return value
}

// Minus subtracts other from the value
func (value *NumericFieldValue[TModel, TAttribute]) Minus(other float64) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.Minus, other)
	return value
}

// Times multiplies value by other
func (value *NumericFieldValue[TModel, TAttribute]) Times(other float64) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.Times, other)
	return value
}

// Divided divides value by other
func (value *NumericFieldValue[TModel, TAttribute]) Divided(other float64) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.Divided, other)
	return value
}

// Modulo returns the remainder of the entire division
func (value *NumericFieldValue[TModel, TAttribute]) Modulo(other int) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.Modulo, other)
	return value
}

// Power elevates value to other
//
// Warning: in sqlite DSQLITE_ENABLE_MATH_FUNCTIONS needs to be enabled or the error "no such function: POWER" will be returned
func (value *NumericFieldValue[TModel, TAttribute]) Power(other float64) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.Power, other)
	return value
}

// SquareRoot calculates the square root of the value
//
// Warning: in sqlite DSQLITE_ENABLE_MATH_FUNCTIONS needs to be enabled or the error "no such function: SQRT" will be returned
func (value *NumericFieldValue[TModel, TAttribute]) SquareRoot() *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.SquareRoot)
	return value
}

// Absolute calculates the absolute value of the value
func (value *NumericFieldValue[TModel, TAttribute]) Absolute() *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.Absolute)
	return value
}

// And calculates the bitwise AND between value and other
func (value *NumericFieldValue[TModel, TAttribute]) And(other int) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.BitAnd, other)
	return value
}

// Or calculates the bitwise OR between value and other
func (value *NumericFieldValue[TModel, TAttribute]) Or(other int) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.BitOr, other)
	return value
}

// Xor calculates the bitwise XOR (exclusive OR) between value and other
//
// Available for: postgres, mysql, sqlserver
func (value *NumericFieldValue[TModel, TAttribute]) Xor(other int) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.BitXor, other)
	return value
}

// Not calculates the bitwise NOT of value
func (value *NumericFieldValue[TModel, TAttribute]) Not() *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.BitNot)
	return value
}

// ShiftLeft shifts value amount bits to the left
func (value *NumericFieldValue[TModel, TAttribute]) ShiftLeft(amount int) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.BitShiftLeft, amount)
	return value
}

// ShiftRight shifts value amount bits to the right
func (value *NumericFieldValue[TModel, TAttribute]) ShiftRight(amount int) *NumericFieldValue[TModel, TAttribute] {
	value.addFunction(sql.BitShiftRight, amount)
	return value
}

type StringFieldValue[TModel model.Model] struct {
	FieldValue[TModel, string]
}

// Concat concatenates other to value
func (value *StringFieldValue[TModel]) Concat(other string) *StringFieldValue[TModel] {
	value.addFunction(sql.Concat, other)
	return value
}

package condition

import (
	"reflect"

	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

type IField interface {
	columnName(query *CQLQuery, table Table) string
	fieldName() string
	columnSQL(query *CQLQuery, table Table) string
	getModelType() reflect.Type
	getAppearance() (uint, bool)
}

type functionAndValues struct {
	function sql.FunctionByDialector
	values   int
}

type Field[TModel model.Model, TAttribute any] struct {
	column             string
	name               string
	columnPrefix       string
	appearance         uint
	appearanceSelected bool
	values             []any
	functions          []functionAndValues
}

// Is allows creating conditions that include the field and a static value
func (field Field[TModel, TAttribute]) Is() FieldIs[TModel, TAttribute] {
	return FieldIs[TModel, TAttribute]{field: field}
}

// Should not be used.
//
// IsUnsafe allows creating conditions that include the field and are not verified in compilation time.
func (field Field[TModel, TAttribute]) IsUnsafe() UnsafeFieldIs[TModel, TAttribute] {
	return UnsafeFieldIs[TModel, TAttribute]{field: field}
}

// Appearance allows to choose which number of appearance use
// when field's model is joined more than once.
func (field Field[TModel, TAttribute]) Appearance(number uint) Field[TModel, TAttribute] {
	newField := NewField[TModel, TAttribute](
		field.name, field.column, field.columnPrefix,
	)

	newField.appearanceSelected = true
	newField.appearance = number

	return newField
}

func (field Field[TModel, TAttribute]) getAppearance() (uint, bool) {
	if !field.appearanceSelected {
		return 0, false
	}

	return field.appearance, true
}

// Aggregate allows applying aggregation functions to the field inside a group by
func (field Field[TModel, TAttribute]) Aggregate() FieldAggregation[TAttribute] {
	return FieldAggregation[TAttribute]{field: field}
}

func (field Field[TModel, TAttribute]) getModelType() reflect.Type {
	return reflect.TypeOf(*new(TModel))
}

// Returns the name of the field identified
func (field Field[TModel, TAttribute]) fieldName() string {
	return field.name
}

// Returns the name of the column in which the field is saved in the table
func (field Field[TModel, TAttribute]) columnName(query *CQLQuery, table Table) string {
	columnName := field.column
	if columnName == "" {
		columnName = query.ColumnName(table, field.name)
	}

	// add column prefix and table name once we know the column name
	return field.columnPrefix + columnName
}

// Returns the SQL to get the value of the field in the table
func (field Field[TModel, TAttribute]) columnSQL(query *CQLQuery, table Table) string {
	return table.Alias + "." + field.columnName(query, table)
}

func (field *Field[TModel, TAttribute]) addFunction(function sql.FunctionByDialector, others ...any) {
	field.functions = append(field.functions, functionAndValues{function: function, values: len(others)})
	field.values = append(field.values, others...)
}

func (field Field[TModel, TAttribute]) ToSQL(query *CQLQuery) (string, []any, error) {
	table, err := getModelTable(query, field)
	if err != nil {
		return "", nil, err
	}

	finalSQL := field.columnSQL(query, table)

	for _, functionAndValues := range field.functions {
		function, isPresent := functionAndValues.function.Get(query.Dialector())
		if !isPresent {
			return "", nil, functionError(ErrUnsupportedByDatabase, functionAndValues.function)
		}

		finalSQL = function.ApplyTo(finalSQL, functionAndValues.values)
	}

	return finalSQL, field.values, nil
}

func (field Field[TModel, TAttribute]) GetValue() TAttribute {
	return *new(TAttribute)
}

func NewField[TModel model.Model, TAttribute any](name, column, columnPrefix string) Field[TModel, TAttribute] {
	return Field[TModel, TAttribute]{
		name:         name,
		column:       column,
		columnPrefix: columnPrefix,
	}
}

type UpdatableField[TModel model.Model, TAttribute any] struct {
	Field[TModel, TAttribute]
}

func (field UpdatableField[TModel, TAttribute]) Set() FieldSet[TModel, TAttribute] {
	return FieldSet[TModel, TAttribute]{field: field}
}

// Appearance allows to choose which number of appearance use
// when field's model is joined more than once.
func (field UpdatableField[TModel, TAttribute]) Appearance(number uint) UpdatableField[TModel, TAttribute] {
	return UpdatableField[TModel, TAttribute]{Field: field.Field.Appearance(number)}
}

func NewUpdatableField[TModel model.Model, TAttribute any](name, column, columnPrefix string) UpdatableField[TModel, TAttribute] {
	return UpdatableField[TModel, TAttribute]{
		Field: NewField[TModel, TAttribute](name, column, columnPrefix),
	}
}

type NullableField[TModel model.Model, TAttribute any] struct {
	UpdatableField[TModel, TAttribute]
}

func (field NullableField[TModel, TAttribute]) Set() NullableFieldSet[TModel, TAttribute] {
	return NullableFieldSet[TModel, TAttribute]{FieldSet[TModel, TAttribute]{field: field.UpdatableField}}
}

// Appearance allows to choose which number of appearance use
// when field's model is joined more than once.
func (field NullableField[TModel, TAttribute]) Appearance(number uint) NullableField[TModel, TAttribute] {
	return NullableField[TModel, TAttribute]{
		UpdatableField: UpdatableField[TModel, TAttribute]{Field: field.Field.Appearance(number)},
	}
}

func NewNullableField[TModel model.Model, TAttribute any](name, column, columnPrefix string) NullableField[TModel, TAttribute] {
	return NullableField[TModel, TAttribute]{
		UpdatableField: NewUpdatableField[TModel, TAttribute](name, column, columnPrefix),
	}
}

type BoolField[TModel model.Model] struct {
	UpdatableField[TModel, bool]
}

func (boolField BoolField[TModel]) Is() BoolFieldIs[TModel] {
	return BoolFieldIs[TModel]{
		field: boolField.Field,
	}
}

// Aggregate allows applying aggregation functions to the field inside a group by
func (boolField BoolField[TModel]) Aggregate() BoolFieldAggregation {
	return BoolFieldAggregation{FieldAggregation: boolField.Field.Aggregate()}
}

func NewBoolField[TModel model.Model](name, column, columnPrefix string) BoolField[TModel] {
	return BoolField[TModel]{
		UpdatableField: NewUpdatableField[TModel, bool](name, column, columnPrefix),
	}
}

type NullableBoolField[TModel model.Model] struct {
	NullableField[TModel, bool]
}

func (boolField NullableBoolField[TModel]) Is() BoolFieldIs[TModel] {
	return BoolFieldIs[TModel]{
		field: boolField.Field,
	}
}

func NewNullableBoolField[TModel model.Model](name, column, columnPrefix string) NullableBoolField[TModel] {
	return NullableBoolField[TModel]{
		NullableField: NewNullableField[TModel, bool](name, column, columnPrefix),
	}
}

type StringField[TModel model.Model] struct {
	UpdatableField[TModel, string]
}

func (stringField StringField[TModel]) Is() StringFieldIs[TModel] {
	return StringFieldIs[TModel]{
		FieldIs: FieldIs[TModel, string]{
			field: stringField.Field,
		},
	}
}

// Concat concatenates other to value
func (stringField StringField[TModel]) Concat(other string) StringField[TModel] {
	stringField.addFunction(sql.Concat, other)
	return stringField
}

// Appearance allows to choose which number of appearance use
// when field's model is joined more than once.
func (stringField StringField[TModel]) Appearance(number uint) StringField[TModel] {
	return StringField[TModel]{
		UpdatableField: UpdatableField[TModel, string]{Field: stringField.Field.Appearance(number)},
	}
}

func NewStringField[TModel model.Model](name, column, columnPrefix string) StringField[TModel] {
	return StringField[TModel]{
		UpdatableField: NewUpdatableField[TModel, string](name, column, columnPrefix),
	}
}

type NullableStringField[TModel model.Model] struct {
	NullableField[TModel, string]
}

func (stringField NullableStringField[TModel]) Is() StringFieldIs[TModel] {
	return StringFieldIs[TModel]{
		FieldIs: FieldIs[TModel, string]{
			field: stringField.Field,
		},
	}
}

// Appearance allows to choose which number of appearance use
// when field's model is joined more than once.
func (stringField NullableStringField[TModel]) Appearance(number uint) NullableStringField[TModel] {
	return NullableStringField[TModel]{
		NullableField: NullableField[TModel, string]{
			UpdatableField: UpdatableField[TModel, string]{Field: stringField.Field.Appearance(number)},
		},
	}
}

func NewNullableStringField[TModel model.Model](name, column, columnPrefix string) NullableStringField[TModel] {
	return NullableStringField[TModel]{
		NullableField: NewNullableField[TModel, string](name, column, columnPrefix),
	}
}

type NumericField[TModel model.Model, TAttribute int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64] struct {
	UpdatableField[TModel, TAttribute]
}

func NewNumericField[
	TModel model.Model,
	TAttribute int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64,
](name, column, columnPrefix string) NumericField[TModel, TAttribute] {
	return NumericField[TModel, TAttribute]{
		UpdatableField: NewUpdatableField[TModel, TAttribute](name, column, columnPrefix),
	}
}

func (numericField NumericField[TModel, TAttribute]) Is() NumericFieldIs[TModel] {
	return NumericFieldIs[TModel]{
		FieldIs: FieldIs[TModel, float64]{
			field: numericField.Field,
		},
	}
}

func (numericField NumericField[TModel, TAttribute]) Set() NumericFieldSet[TModel, TAttribute] {
	return NumericFieldSet[TModel, TAttribute]{field: numericField}
}

// Appearance allows to choose which number of appearance use
// when field's model is joined more than once.
func (numericField NumericField[TModel, TAttribute]) Appearance(number uint) NumericField[TModel, TAttribute] {
	return NumericField[TModel, TAttribute]{
		UpdatableField: UpdatableField[TModel, TAttribute]{Field: numericField.Field.Appearance(number)},
	}
}

// Aggregate allows applying aggregation functions to the field inside a group by
func (numericField NumericField[TModel, TAttribute]) Aggregate() NumericFieldAggregation {
	return NumericFieldAggregation{FieldAggregation: FieldAggregation[float64]{field: numericField}}
}

func (numericField NumericField[TModel, TAttribute]) GetValue() float64 {
	return 0
}

// Plus sums other to value
func (numericField NumericField[TModel, TAttribute]) Plus(other float64) NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.Plus, other)
	return numericField
}

// Minus subtracts other from the value
func (numericField NumericField[TModel, TAttribute]) Minus(other float64) NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.Minus, other)
	return numericField
}

// Times multiplies value by other
func (numericField NumericField[TModel, TAttribute]) Times(other float64) NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.Times, other)
	return numericField
}

// Divided divides value by other
func (numericField NumericField[TModel, TAttribute]) Divided(other float64) NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.Divided, other)
	return numericField
}

// Modulo returns the remainder of the entire division
func (numericField NumericField[TModel, TAttribute]) Modulo(other int) NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.Modulo, other)
	return numericField
}

// Power elevates value to other
//
// Warning: in sqlite DSQLITE_ENABLE_MATH_FUNCTIONS needs to be enabled or the error "no such function: POWER" will be returned
func (numericField NumericField[TModel, TAttribute]) Power(other float64) NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.Power, other)
	return numericField
}

// SquareRoot calculates the square root of the value
//
// Warning: in sqlite DSQLITE_ENABLE_MATH_FUNCTIONS needs to be enabled or the error "no such function: SQRT" will be returned
func (numericField NumericField[TModel, TAttribute]) SquareRoot() NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.SquareRoot)
	return numericField
}

// Absolute calculates the absolute value of the value
func (numericField NumericField[TModel, TAttribute]) Absolute() NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.Absolute)
	return numericField
}

// And calculates the bitwise AND between value and other
func (numericField NumericField[TModel, TAttribute]) And(other int) NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.BitAnd, other)
	return numericField
}

// Or calculates the bitwise OR between value and other
func (numericField NumericField[TModel, TAttribute]) Or(other int) NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.BitOr, other)
	return numericField
}

// Xor calculates the bitwise XOR (exclusive OR) between value and other
//
// Available for: postgres, mysql, sqlserver
func (numericField NumericField[TModel, TAttribute]) Xor(other int) NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.BitXor, other)
	return numericField
}

// Not calculates the bitwise NOT of value
func (numericField NumericField[TModel, TAttribute]) Not() NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.BitNot)
	return numericField
}

// ShiftLeft shifts value amount bits to the left
func (numericField NumericField[TModel, TAttribute]) ShiftLeft(amount int) NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.BitShiftLeft, amount)
	return numericField
}

// ShiftRight shifts value amount bits to the right
func (numericField NumericField[TModel, TAttribute]) ShiftRight(amount int) NumericField[TModel, TAttribute] {
	numericField.addFunction(sql.BitShiftRight, amount)
	return numericField
}

type NullableNumericField[
	TModel model.Model,
	TAttribute int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64,
] struct {
	NumericField[TModel, TAttribute]
}

func (field NullableNumericField[TModel, TAttribute]) Set() NullableFieldSet[TModel, TAttribute] {
	return NullableFieldSet[TModel, TAttribute]{FieldSet[TModel, TAttribute]{field: field.UpdatableField}}
}

// Appearance allows to choose which number of appearance use
// when field's model is joined more than once.
func (field NullableNumericField[TModel, TAttribute]) Appearance(number uint) NullableNumericField[TModel, TAttribute] {
	return NullableNumericField[TModel, TAttribute]{
		NumericField: NumericField[TModel, TAttribute]{
			UpdatableField: UpdatableField[TModel, TAttribute]{Field: field.Field.Appearance(number)},
		},
	}
}

func NewNullableNumericField[
	TModel model.Model,
	TAttribute int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64,
](name, column, columnPrefix string) NullableNumericField[TModel, TAttribute] {
	return NullableNumericField[TModel, TAttribute]{
		NumericField: NewNumericField[TModel, TAttribute](name, column, columnPrefix),
	}
}

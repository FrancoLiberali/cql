package condition

import (
	"reflect"
	"strings"

	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

type IField interface {
	ToSQL(query *CQLQuery) (string, []any, error)
	ToSQLForTable(query *CQLQuery, table Table) (string, []any, error)
	columnName(query *CQLQuery, table Table) string
	fieldName() string
	columnSQL(query *CQLQuery, table Table) string
	getModelType() reflect.Type
	getAppearance() (uint, bool)
	// getScannerErased returns *Scanner[TModel] as any (TModel is not
	// visible at the IField interface level). Nil if no scanner is wired
	// — i.e. when the field wasn't constructed by cql-gen output, or the
	// generated file predates scanner emission.
	getScannerErased() any
}

type FieldOfModel[T model.Model] interface {
	IField

	getModel() T
}

type functionAndValues struct {
	function sql.FunctionByDialector
	values   []IValue
}

type Field[TModel model.Model, TAttribute any] struct {
	column             string
	name               string
	columnPrefix       string
	appearance         uint
	appearanceSelected bool
	functions          []functionAndValues
	// scanner is the per-model fast-scan dispatcher, wired by cql-gen.
	// Shared by pointer across every Field on the same conditions struct,
	// so propagation through Condition values only costs a pointer copy.
	// nil when the field wasn't built by scanner-aware generated code, in
	// which case Query[T] falls back to the gorm scan path.
	scanner *Scanner[TModel]
}

func (field Field[TModel, TAttribute]) getModel() TModel {
	return *new(TModel)
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
	newField := Field[TModel, TAttribute]{
		name:               field.name,
		column:             field.column,
		columnPrefix:       field.columnPrefix,
		appearanceSelected: true,
		appearance:         number,
		scanner:            field.scanner,
	}

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
	if name, ok := field.schemaColumnName(table); ok {
		return name
	}

	columnName := field.column
	if columnName == "" {
		columnName = query.ColumnName(table, field.name)
	}

	// add column prefix and table name once we know the column name
	return field.columnPrefix + columnName
}

// schemaColumnName resolves the column name directly from gorm's parsed
// schema when available, returning (name, true) on a hit. This is a fast
// path: gorm's DBName is precomputed with any embedded-prefix applied, so it
// skips the per-call NamingStrategy allocations (strings.Split +
// Replacer.Replace). columnPrefix is already baked into DBName by gorm, so we
// must NOT prepend it again.
func (field Field[TModel, TAttribute]) schemaColumnName(table Table) (string, bool) {
	if field.column != "" || table.Schema == nil {
		return "", false
	}

	if field.columnPrefix == "" {
		// No prefix — top-level field OR field promoted from an anonymous
		// embed (e.g. User.ID coming from an anonymously-embedded
		// model.UIntModelWithTimestamps). Gorm distinguishes these via
		// EmbeddedBindNames: only NAMED embeds (`X SomeType `gorm:"embedded"``)
		// prepend the outer name; anonymous embeds keep EmbeddedBindNames at
		// length 1. So EmbeddedBindNames == 1 is the right guard — it accepts
		// promoted embedded fields whose DBName is unprefixed (matches what our
		// fallback would compute) and rejects prefixed embeds where DBName has
		// a prefix we must not miss.
		if f := table.Schema.LookUpField(field.name); f != nil && len(f.EmbeddedBindNames) == 1 {
			return f.DBName, true
		}

		return "", false
	}

	// Prefixed — the field lives inside an embedded struct. Iterate to find
	// the schema Field with matching Go name whose DBName starts with the CQL
	// columnPrefix. Embedded structs typically have few fields, so the scan is
	// cheap.
	for _, f := range table.Schema.Fields {
		if f.Name == field.name && strings.HasPrefix(f.DBName, field.columnPrefix) {
			return f.DBName, true
		}
	}

	return "", false
}

// Returns the SQL to get the value of the field in the table
func (field Field[TModel, TAttribute]) columnSQL(query *CQLQuery, table Table) string {
	return table.Alias + "." + field.columnName(query, table)
}

func (field Field[TModel, TAttribute]) addFunction(function sql.FunctionByDialector, others ...IValue) Field[TModel, TAttribute] {
	field.functions = append(field.functions, functionAndValues{function: function, values: others})

	return field
}

func (field Field[TModel, TAttribute]) ToSQLForTable(query *CQLQuery, table Table) (string, []any, error) {
	finalSQL := field.columnSQL(query, table)

	finalValues := []any{}

	for _, functionAndValues := range field.functions {
		function, isPresent := functionAndValues.function.Get(query.Dialector())
		if !isPresent {
			return "", nil, functionError(ErrUnsupportedByDatabase, functionAndValues.function)
		}

		valuesSQLs := []string{}

		for _, value := range functionAndValues.values {
			valueSQL, valueValue, err := value.ToSQL(query)
			if err != nil {
				return "", nil, err
			}

			if valueSQL != "" {
				valuesSQLs = append(valuesSQLs, valueSQL)
			} else {
				valuesSQLs = append(valuesSQLs, "?")
			}

			finalValues = append(finalValues, valueValue...)
		}

		finalSQL = function.ApplyTo(finalSQL, valuesSQLs)
	}

	return finalSQL, finalValues, nil
}

func (field Field[TModel, TAttribute]) ToSQL(query *CQLQuery) (string, []any, error) {
	table, err := getModelTable(query, field)
	if err != nil {
		return "", nil, err
	}

	return field.ToSQLForTable(query, table)
}

func (field Field[TModel, TAttribute]) GetValue() TAttribute {
	return *new(TAttribute)
}

// getScannerErased returns the per-model Scanner pointer as any so it can
// be carried through the non-generic IField interface. Returns a true nil
// (not a typed nil) when no scanner is wired, so caller's type assertion
// returns ok=false cleanly.
func (field Field[TModel, TAttribute]) getScannerErased() any {
	if field.scanner == nil {
		return nil
	}

	return field.scanner
}

// NewField builds a Field. The optional scanner is shared by pointer across
// every field on a model's conditions struct; cql-gen output passes it on
// every constructor call. Hand-written callers may omit it (the fast scan
// path simply falls back to gorm).
func NewField[TModel model.Model, TAttribute any](name, column, columnPrefix string, scanner ...*Scanner[TModel]) Field[TModel, TAttribute] {
	f := Field[TModel, TAttribute]{
		name:         name,
		column:       column,
		columnPrefix: columnPrefix,
	}

	if len(scanner) > 0 {
		f.scanner = scanner[0]
	}

	return f
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

func NewUpdatableField[TModel model.Model, TAttribute any](name, column, columnPrefix string, scanner ...*Scanner[TModel]) UpdatableField[TModel, TAttribute] {
	return UpdatableField[TModel, TAttribute]{
		Field: NewField[TModel, TAttribute](name, column, columnPrefix, scanner...),
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

func NewNullableField[TModel model.Model, TAttribute any](name, column, columnPrefix string, scanner ...*Scanner[TModel]) NullableField[TModel, TAttribute] {
	return NullableField[TModel, TAttribute]{
		UpdatableField: NewUpdatableField[TModel, TAttribute](name, column, columnPrefix, scanner...),
	}
}

type BoolField[TModel model.Model] struct {
	UpdatableField[TModel, bool]
}

func (boolField BoolField[TModel]) Is() BoolFieldIs[TModel] {
	return newBoolFieldIs(boolField.Field)
}

// Aggregate allows applying aggregation functions to the field inside a group by
func (boolField BoolField[TModel]) Aggregate() BoolFieldAggregation {
	return BoolFieldAggregation{FieldAggregation: boolField.Field.Aggregate()}
}

// Appearance allows to choose which number of appearance use
// when field's model is joined more than once.
func (boolField BoolField[TModel]) Appearance(number uint) BoolField[TModel] {
	return BoolField[TModel]{
		UpdatableField: boolField.UpdatableField.Appearance(number),
	}
}

func NewBoolField[TModel model.Model](name, column, columnPrefix string, scanner ...*Scanner[TModel]) BoolField[TModel] {
	return BoolField[TModel]{
		UpdatableField: NewUpdatableField[TModel, bool](name, column, columnPrefix, scanner...),
	}
}

type NullableBoolField[TModel model.Model] struct {
	BoolField[TModel]
}

func (boolField NullableBoolField[TModel]) Set() NullableFieldSet[TModel, bool] {
	return NullableFieldSet[TModel, bool]{FieldSet[TModel, bool]{field: boolField.UpdatableField}}
}

// Appearance allows to choose which number of appearance use
// when field's model is joined more than once.
func (boolField NullableBoolField[TModel]) Appearance(number uint) NullableBoolField[TModel] {
	return NullableBoolField[TModel]{
		BoolField: boolField.BoolField.Appearance(number),
	}
}

func NewNullableBoolField[TModel model.Model](name, column, columnPrefix string, scanner ...*Scanner[TModel]) NullableBoolField[TModel] {
	return NullableBoolField[TModel]{
		BoolField: NewBoolField[TModel](name, column, columnPrefix, scanner...),
	}
}

type NotUpdatableStringField[TModel model.Model] struct {
	Field[TModel, string]
}

func (stringField NotUpdatableStringField[TModel]) Is() StringFieldIs[TModel] {
	return newStringFieldIs(stringField.Field)
}

func (stringField NotUpdatableStringField[TModel]) addFunction(function sql.FunctionByDialector, others ...IValue) NotUpdatableStringField[TModel] {
	stringField.Field = stringField.Field.addFunction(function, others...)

	return stringField
}

// Concat concatenates other to value
func (stringField NotUpdatableStringField[TModel]) Concat(other ValueOfType[string]) NotUpdatableStringField[TModel] {
	return stringField.addFunction(sql.Concat, other)
}

type StringField[TModel model.Model] struct {
	UpdatableField[TModel, string]
}

func (stringField StringField[TModel]) Is() StringFieldIs[TModel] {
	return newStringFieldIs(stringField.Field)
}

func (stringField StringField[TModel]) toNotUpdatable() NotUpdatableStringField[TModel] {
	return NotUpdatableStringField[TModel]{Field: stringField.Field}
}

// Concat concatenates other to value
func (stringField StringField[TModel]) Concat(other ValueOfType[string]) NotUpdatableStringField[TModel] {
	return stringField.toNotUpdatable().Concat(other)
}

// Appearance allows to choose which number of appearance use
// when field's model is joined more than once.
func (stringField StringField[TModel]) Appearance(number uint) StringField[TModel] {
	return StringField[TModel]{
		UpdatableField: stringField.UpdatableField.Appearance(number),
	}
}

func NewStringField[TModel model.Model](name, column, columnPrefix string, scanner ...*Scanner[TModel]) StringField[TModel] {
	return StringField[TModel]{
		UpdatableField: NewUpdatableField[TModel, string](name, column, columnPrefix, scanner...),
	}
}

type NullableStringField[TModel model.Model] struct {
	StringField[TModel]
}

func (stringField NullableStringField[TModel]) Set() NullableFieldSet[TModel, string] {
	return NullableFieldSet[TModel, string]{FieldSet[TModel, string]{field: stringField.UpdatableField}}
}

// Appearance allows to choose which number of appearance use
// when field's model is joined more than once.
func (stringField NullableStringField[TModel]) Appearance(number uint) NullableStringField[TModel] {
	return NullableStringField[TModel]{
		StringField: stringField.StringField.Appearance(number),
	}
}

func NewNullableStringField[TModel model.Model](name, column, columnPrefix string, scanner ...*Scanner[TModel]) NullableStringField[TModel] {
	return NullableStringField[TModel]{
		StringField: NewStringField[TModel](name, column, columnPrefix, scanner...),
	}
}

type NotUpdatableNumericField[TModel model.Model, TAttribute Numeric] struct {
	Field[TModel, TAttribute]
}

func (numericField NotUpdatableNumericField[TModel, TAttribute]) GetValue() float64 {
	return 0
}

func (numericField NotUpdatableNumericField[TModel, TAttribute]) GetNumericValue() TAttribute {
	return 0
}

func (numericField NotUpdatableNumericField[TModel, TAttribute]) Is() NumericFieldIs[TModel] {
	return newNumericFieldIs(numericField.Field)
}

// Aggregate allows applying aggregation functions to the field inside a group by
func (numericField NotUpdatableNumericField[TModel, TAttribute]) Aggregate() NumericFieldAggregation {
	return NumericFieldAggregation{FieldAggregation: FieldAggregation[float64]{field: numericField}}
}

func (numericField NotUpdatableNumericField[TModel, TAttribute]) addFunction(
	function sql.FunctionByDialector, others ...IValue,
) NotUpdatableNumericField[TModel, TAttribute] {
	numericField.Field = numericField.Field.addFunction(function, others...)

	return numericField
}

// Plus sums other to value
func (numericField NotUpdatableNumericField[TModel, TAttribute]) Plus(other ValueOfType[float64]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.Plus, other)
}

// Minus subtracts other from the value
func (numericField NotUpdatableNumericField[TModel, TAttribute]) Minus(other ValueOfType[float64]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.Minus, other)
}

// Times multiplies value by other
func (numericField NotUpdatableNumericField[TModel, TAttribute]) Times(other ValueOfType[float64]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.Times, other)
}

// Divided divides value by other
func (numericField NotUpdatableNumericField[TModel, TAttribute]) Divided(other ValueOfType[float64]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.Divided, other)
}

// Modulo returns the remainder of the entire division
func (numericField NotUpdatableNumericField[TModel, TAttribute]) Modulo(other ValueOfType[float64]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.Modulo, other)
}

// Power elevates value to other
//
// Warning: in sqlite DSQLITE_ENABLE_MATH_FUNCTIONS needs to be enabled or the error "no such function: POWER" will be returned
func (numericField NotUpdatableNumericField[TModel, TAttribute]) Power(other ValueOfType[float64]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.Power, other)
}

// SquareRoot calculates the square root of the value
//
// Warning: in sqlite DSQLITE_ENABLE_MATH_FUNCTIONS needs to be enabled or the error "no such function: SQRT" will be returned
func (numericField NotUpdatableNumericField[TModel, TAttribute]) SquareRoot() NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.SquareRoot)
}

// Absolute calculates the absolute value of the value
func (numericField NotUpdatableNumericField[TModel, TAttribute]) Absolute() NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.Absolute)
}

// And calculates the bitwise AND between value and other
func (numericField NotUpdatableNumericField[TModel, TAttribute]) And(other NumericOfType[TAttribute]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.BitAnd, other)
}

// Or calculates the bitwise OR between value and other
func (numericField NotUpdatableNumericField[TModel, TAttribute]) Or(other NumericOfType[TAttribute]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.BitOr, other)
}

// Xor calculates the bitwise XOR (exclusive OR) between value and other
//
// Available for: postgres, mysql, sqlserver
func (numericField NotUpdatableNumericField[TModel, TAttribute]) Xor(other NumericOfType[TAttribute]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.BitXor, other)
}

// Not calculates the bitwise NOT of value
func (numericField NotUpdatableNumericField[TModel, TAttribute]) Not() NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.BitNot)
}

// ShiftLeft shifts value amount bits to the left
func (numericField NotUpdatableNumericField[TModel, TAttribute]) ShiftLeft(amount NumericOfType[int]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.BitShiftLeft, amount)
}

// ShiftRight shifts value amount bits to the right
func (numericField NotUpdatableNumericField[TModel, TAttribute]) ShiftRight(amount NumericOfType[int]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.addFunction(sql.BitShiftRight, amount)
}

type NumericField[TModel model.Model, TAttribute Numeric] struct {
	UpdatableField[TModel, TAttribute]
}

func NewNumericField[
	TModel model.Model,
	TAttribute Numeric,
](name, column, columnPrefix string, scanner ...*Scanner[TModel]) NumericField[TModel, TAttribute] {
	return NumericField[TModel, TAttribute]{
		UpdatableField: NewUpdatableField[TModel, TAttribute](name, column, columnPrefix, scanner...),
	}
}

func (numericField NumericField[TModel, TAttribute]) Is() NumericFieldIs[TModel] {
	return newNumericFieldIs(numericField.Field)
}

func (numericField NumericField[TModel, TAttribute]) Set() NumericFieldSet[TModel, TAttribute] {
	return NumericFieldSet[TModel, TAttribute]{field: numericField}
}

// Appearance allows to choose which number of appearance use
// when field's model is joined more than once.
func (numericField NumericField[TModel, TAttribute]) Appearance(number uint) NumericField[TModel, TAttribute] {
	return NumericField[TModel, TAttribute]{
		UpdatableField: numericField.UpdatableField.Appearance(number),
	}
}

// Aggregate allows applying aggregation functions to the field inside a group by
func (numericField NumericField[TModel, TAttribute]) Aggregate() NumericFieldAggregation {
	return NumericFieldAggregation{FieldAggregation: FieldAggregation[float64]{field: numericField}}
}

func (numericField NumericField[TModel, TAttribute]) GetValue() float64 {
	return 0
}

func (numericField NumericField[TModel, TAttribute]) GetNumericValue() TAttribute {
	return 0
}

func (numericField NumericField[TModel, TAttribute]) toNotUpdatable() NotUpdatableNumericField[TModel, TAttribute] {
	return NotUpdatableNumericField[TModel, TAttribute]{Field: numericField.Field}
}

// Plus sums other to value
func (numericField NumericField[TModel, TAttribute]) Plus(other ValueOfType[float64]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().Plus(other)
}

// Minus subtracts other from the value
func (numericField NumericField[TModel, TAttribute]) Minus(other ValueOfType[float64]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().Minus(other)
}

// Times multiplies value by other
func (numericField NumericField[TModel, TAttribute]) Times(other ValueOfType[float64]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().Times(other)
}

// Divided divides value by other
func (numericField NumericField[TModel, TAttribute]) Divided(other ValueOfType[float64]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().Divided(other)
}

// Modulo returns the remainder of the entire division
func (numericField NumericField[TModel, TAttribute]) Modulo(other ValueOfType[float64]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().Modulo(other)
}

// Power elevates value to other
//
// Warning: in sqlite DSQLITE_ENABLE_MATH_FUNCTIONS needs to be enabled or the error "no such function: POWER" will be returned
func (numericField NumericField[TModel, TAttribute]) Power(other ValueOfType[float64]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().Power(other)
}

// SquareRoot calculates the square root of the value
//
// Warning: in sqlite DSQLITE_ENABLE_MATH_FUNCTIONS needs to be enabled or the error "no such function: SQRT" will be returned
func (numericField NumericField[TModel, TAttribute]) SquareRoot() NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().SquareRoot()
}

// Absolute calculates the absolute value of the value
func (numericField NumericField[TModel, TAttribute]) Absolute() NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().Absolute()
}

// And calculates the bitwise AND between value and other
func (numericField NumericField[TModel, TAttribute]) And(other NumericOfType[TAttribute]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().And(other)
}

// Or calculates the bitwise OR between value and other
func (numericField NumericField[TModel, TAttribute]) Or(other NumericOfType[TAttribute]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().Or(other)
}

// Xor calculates the bitwise XOR (exclusive OR) between value and other
//
// Available for: postgres, mysql, sqlserver
func (numericField NumericField[TModel, TAttribute]) Xor(other NumericOfType[TAttribute]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().Xor(other)
}

// Not calculates the bitwise NOT of value
func (numericField NumericField[TModel, TAttribute]) Not() NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().Not()
}

// ShiftLeft shifts value amount bits to the left
func (numericField NumericField[TModel, TAttribute]) ShiftLeft(amount NumericOfType[int]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().ShiftLeft(amount)
}

// ShiftRight shifts value amount bits to the right
func (numericField NumericField[TModel, TAttribute]) ShiftRight(amount NumericOfType[int]) NotUpdatableNumericField[TModel, TAttribute] {
	return numericField.toNotUpdatable().ShiftRight(amount)
}

type NullableNumericField[
	TModel model.Model,
	TAttribute Numeric,
] struct {
	NumericField[TModel, TAttribute]
}

func (field NullableNumericField[TModel, TAttribute]) Set() NullableNumericFieldSet[TModel, TAttribute] {
	return newNullableNumericFieldSet(field.NumericField)
}

// Appearance allows to choose which number of appearance use
// when field's model is joined more than once.
func (field NullableNumericField[TModel, TAttribute]) Appearance(number uint) NullableNumericField[TModel, TAttribute] {
	return NullableNumericField[TModel, TAttribute]{
		NumericField: field.NumericField.Appearance(number),
	}
}

func NewNullableNumericField[
	TModel model.Model,
	TAttribute Numeric,
](name, column, columnPrefix string, scanner ...*Scanner[TModel]) NullableNumericField[TModel, TAttribute] {
	return NullableNumericField[TModel, TAttribute]{
		NumericField: NewNumericField[TModel, TAttribute](name, column, columnPrefix, scanner...),
	}
}

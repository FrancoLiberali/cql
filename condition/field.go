package condition

import (
	"reflect"

	"github.com/FrancoLiberali/cql/model"
)

type IField interface {
	columnName(query *GormQuery, table Table) string
	fieldName() string
	columnSQL(query *GormQuery, table Table) string
	getModelType() reflect.Type
	getAppearance() int
}

type Field[TModel model.Model, TAttribute any] struct {
	column             string
	name               string
	columnPrefix       string
	appearance         uint
	appearanceSelected bool
}

// Is allows creating conditions that include the field and a static value
func (field Field[TModel, TAttribute]) Is() FieldIs[TModel, TAttribute] {
	return FieldIs[TModel, TAttribute]{field: field}
}

// IsDynamic allows creating conditions that include the field and other fields
func (field Field[TModel, TAttribute]) IsDynamic() DynamicFieldIs[TModel, TAttribute] {
	return DynamicFieldIs[TModel, TAttribute]{field: field}
}

// Should not be used.
//
// IsUnsafe allows creating conditions that include the field and are not verified in compilation time.
func (field Field[TModel, TAttribute]) IsUnsafe() UnsafeFieldIs[TModel, TAttribute] {
	return UnsafeFieldIs[TModel, TAttribute]{field: field}
}

// Value allows using the value of the field inside dynamic conditions.
func (field Field[TModel, TAttribute]) Value() *FieldValue[TModel, TAttribute] {
	return NewFieldValue(field)
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

func (field Field[TModel, TAttribute]) getAppearance() int {
	if !field.appearanceSelected {
		return undefinedAppearance
	}

	return int(field.appearance)
}

func (field Field[TModel, TAttribute]) getModelType() reflect.Type {
	return reflect.TypeOf(*new(TModel))
}

// Returns the name of the field identified
func (field Field[TModel, TAttribute]) fieldName() string {
	return field.name
}

// Returns the name of the column in which the field is saved in the table
func (field Field[TModel, TAttribute]) columnName(query *GormQuery, table Table) string {
	columnName := field.column
	if columnName == "" {
		columnName = query.ColumnName(table, field.name)
	}

	// add column prefix and table name once we know the column name
	return field.columnPrefix + columnName
}

// Returns the SQL to get the value of the field in the table
func (field Field[TModel, TAttribute]) columnSQL(query *GormQuery, table Table) string {
	return table.Alias + "." + field.columnName(query, table)
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

// Value allows using the value of the field inside dynamic conditions.
func (stringField StringField[TModel]) Value() *StringFieldValue[TModel] {
	return &StringFieldValue[TModel]{FieldValue: *stringField.UpdatableField.Value()}
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

// Value allows using the value of the field inside dynamic conditions.
func (stringField NullableStringField[TModel]) Value() *StringFieldValue[TModel] {
	return &StringFieldValue[TModel]{FieldValue: *stringField.UpdatableField.Value()}
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

func (numericField NumericField[TModel, TAttribute]) IsDynamic() NumericDynamicFieldIs[TModel, TAttribute] {
	return NumericDynamicFieldIs[TModel, TAttribute]{field: numericField.Field}
}

// Value allows using the value of the field inside dynamic conditions.
func (numericField NumericField[TModel, TAttribute]) Value() *NumericFieldValue[TModel, TAttribute] {
	return &NumericFieldValue[TModel, TAttribute]{FieldValue: *numericField.UpdatableField.Value()}
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
func (numericField NullableNumericField[TModel, TAttribute]) Appearance(number uint) NullableNumericField[TModel, TAttribute] {
	return NullableNumericField[TModel, TAttribute]{
		NumericField: NumericField[TModel, TAttribute]{
			UpdatableField: UpdatableField[TModel, TAttribute]{Field: numericField.Field.Appearance(number)},
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

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
}

type FieldOfType[T any] interface {
	IField
	getType() T
}

type Field[TModel model.Model, TAttribute any] struct {
	column       string
	name         string
	columnPrefix string
}

func (field Field[TModel, TAttribute]) getType() TAttribute {
	return *new(TAttribute)
}

func (field Field[TModel, TAttribute]) Is() FieldIs[TModel, TAttribute] {
	return FieldIs[TModel, TAttribute]{field: field}
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
	return FieldSet[TModel, TAttribute]{Field: field}
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
	return NullableFieldSet[TModel, TAttribute]{FieldSet[TModel, TAttribute]{Field: field.UpdatableField}}
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

func NewNullableStringField[TModel model.Model](name, column, columnPrefix string) NullableStringField[TModel] {
	return NullableStringField[TModel]{
		NullableField: NewNullableField[TModel, string](name, column, columnPrefix),
	}
}

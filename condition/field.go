package condition

import (
	"reflect"

	"github.com/FrancoLiberali/cql/model"
)

type IField interface {
	ColumnName(query *GormQuery, table Table) string
	FieldName() string
	ColumnSQL(query *GormQuery, table Table) string
	GetModelType() reflect.Type
}

type FieldOfType[T any] interface {
	IField
	GetType() T
}

type Field[TModel model.Model, TAttribute any] struct {
	Column       string
	Name         string
	ColumnPrefix string
}

func (field Field[TModel, TAttribute]) GetType() TAttribute {
	return *new(TAttribute)
}

func (field Field[TModel, TAttribute]) Is() FieldIs[TModel, TAttribute] {
	return FieldIs[TModel, TAttribute]{Field: field}
}

func (field Field[TModel, TAttribute]) GetModelType() reflect.Type {
	return reflect.TypeOf(*new(TModel))
}

// Returns the name of the field identified
func (field Field[TModel, TAttribute]) FieldName() string {
	return field.Name
}

// Returns the name of the column in which the field is saved in the table
func (field Field[TModel, TAttribute]) ColumnName(query *GormQuery, table Table) string {
	columnName := field.Column
	if columnName == "" {
		columnName = query.ColumnName(table, field.Name)
	}

	// add column prefix and table name once we know the column name
	return field.ColumnPrefix + columnName
}

// Returns the SQL to get the value of the field in the table
func (field Field[TModel, TAttribute]) ColumnSQL(query *GormQuery, table Table) string {
	return table.Alias + "." + field.ColumnName(query, table)
}

type UpdatableField[TModel model.Model, TAttribute any] struct {
	Field[TModel, TAttribute]
}

func (field UpdatableField[TModel, TAttribute]) Set() FieldSet[TModel, TAttribute] {
	return FieldSet[TModel, TAttribute]{Field: field}
}

type NullableField[TModel model.Model, TAttribute any] struct {
	UpdatableField[TModel, TAttribute]
}

func (field NullableField[TModel, TAttribute]) Set() NullableFieldSet[TModel, TAttribute] {
	return NullableFieldSet[TModel, TAttribute]{FieldSet[TModel, TAttribute]{Field: field.UpdatableField}}
}

type BoolField[TModel model.Model] struct {
	UpdatableField[TModel, bool]
}

func (boolField BoolField[TModel]) Is() BoolFieldIs[TModel] {
	return BoolFieldIs[TModel]{
		Field: boolField.Field,
	}
}

type NullableBoolField[TModel model.Model] struct {
	NullableField[TModel, bool]
}

func (boolField NullableBoolField[TModel]) Is() BoolFieldIs[TModel] {
	return BoolFieldIs[TModel]{
		Field: boolField.Field,
	}
}

type StringField[TModel model.Model] struct {
	UpdatableField[TModel, string]
}

func (stringField StringField[TModel]) Is() StringFieldIs[TModel] {
	return StringFieldIs[TModel]{
		FieldIs: FieldIs[TModel, string]{
			Field: stringField.Field,
		},
	}
}

type NullableStringField[TModel model.Model] struct {
	NullableField[TModel, string]
}

func (stringField NullableStringField[TModel]) Is() StringFieldIs[TModel] {
	return StringFieldIs[TModel]{
		FieldIs: FieldIs[TModel, string]{
			Field: stringField.Field,
		},
	}
}

package orm

import (
	"reflect"

	"github.com/ditrit/badaas/orm/model"
)

type IField interface {
	// TODO ver si de todos estos podria sacar alguno
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

func (field Field[TModel, TAttribute]) Set() FieldSet[TModel, TAttribute] {
	return FieldSet[TModel, TAttribute]{Field: field}
}

func (field Field[TModel, TAttribute]) GetModelType() reflect.Type {
	// TODO ver si esto sigue siendo necesario
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

type BoolField[TModel model.Model] struct {
	Field[TModel, bool]
}

func (boolField BoolField[TModel]) Is() BoolFieldIs[TModel] {
	return BoolFieldIs[TModel]{
		FieldIs: FieldIs[TModel, bool]{
			Field: boolField.Field,
		},
	}
}

type StringField[TModel model.Model] struct {
	Field[TModel, string]
}

func (stringField StringField[TModel]) Is() StringFieldIs[TModel] {
	return StringFieldIs[TModel]{
		FieldIs: FieldIs[TModel, string]{
			Field: stringField.Field,
		},
	}
}

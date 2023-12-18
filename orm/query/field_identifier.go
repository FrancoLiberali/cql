package query

import (
	"reflect"
)

type IFieldIdentifier interface {
	ColumnName(query *GormQuery, table Table) string
	ColumnSQL(query *GormQuery, table Table) string
	GetModelType() reflect.Type
}

type FieldIdentifier[T any] struct {
	Column       string
	Field        string
	ColumnPrefix string
	ModelType    reflect.Type
}

func (fieldID FieldIdentifier[T]) GetModelType() reflect.Type {
	return fieldID.ModelType
}

// Returns the name of the column in which the field is saved in the table
func (fieldID FieldIdentifier[T]) ColumnName(query *GormQuery, table Table) string {
	columnName := fieldID.Column
	if columnName == "" {
		columnName = query.ColumnName(table, fieldID.Field)
	}

	// add column prefix and table name once we know the column name
	return fieldID.ColumnPrefix + columnName
}

// Returns the SQL to get the value of the field in the table
func (fieldID FieldIdentifier[T]) ColumnSQL(query *GormQuery, table Table) string {
	return table.Alias + "." + fieldID.ColumnName(query, table)
}

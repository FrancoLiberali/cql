// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	model "github.com/FrancoLiberali/cql/model"
	columndefinition "github.com/FrancoLiberali/cql/cql-cli/cmd/gen/conditions/tests/columndefinition"
	"time"
)

type columnDefinitionConditions struct {
	ID        condition.Field[columndefinition.ColumnDefinition, model.UUID]
	CreatedAt condition.Field[columndefinition.ColumnDefinition, time.Time]
	UpdatedAt condition.Field[columndefinition.ColumnDefinition, time.Time]
	DeletedAt condition.Field[columndefinition.ColumnDefinition, time.Time]
	String    condition.StringField[columndefinition.ColumnDefinition]
}

var ColumnDefinition = columnDefinitionConditions{
	CreatedAt: condition.Field[columndefinition.ColumnDefinition, time.Time]{Name: "CreatedAt"},
	DeletedAt: condition.Field[columndefinition.ColumnDefinition, time.Time]{Name: "DeletedAt"},
	ID:        condition.Field[columndefinition.ColumnDefinition, model.UUID]{Name: "ID"},
	String: condition.StringField[columndefinition.ColumnDefinition]{Field: condition.Field[columndefinition.ColumnDefinition, string]{
		Column: "string_something_else",
		Name:   "String",
	}},
	UpdatedAt: condition.Field[columndefinition.ColumnDefinition, time.Time]{Name: "UpdatedAt"},
}

// Preload allows preloading the ColumnDefinition when doing a query
func (columnDefinitionConditions columnDefinitionConditions) Preload() condition.Condition[columndefinition.ColumnDefinition] {
	return condition.NewPreloadCondition[columndefinition.ColumnDefinition](columnDefinitionConditions.ID, columnDefinitionConditions.CreatedAt, columnDefinitionConditions.UpdatedAt, columnDefinitionConditions.DeletedAt, columnDefinitionConditions.String)
}
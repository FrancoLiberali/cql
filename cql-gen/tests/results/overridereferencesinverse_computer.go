// Code generated by cql-gen v0.1.0, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	overridereferencesinverse "github.com/FrancoLiberali/cql/cql-gen/cmd/gen/conditions/tests/overridereferencesinverse"
	model "github.com/FrancoLiberali/cql/model"
	"time"
)

func (computerConditions computerConditions) Processor(conditions ...condition.Condition[overridereferencesinverse.Processor]) condition.JoinCondition[overridereferencesinverse.Computer] {
	return condition.NewJoinCondition[overridereferencesinverse.Computer, overridereferencesinverse.Processor](conditions, "Processor", "Name", computerConditions.preload(), "ComputerName", Processor.preload())
}

type computerConditions struct {
	ID        condition.Field[overridereferencesinverse.Computer, model.UUID]
	CreatedAt condition.Field[overridereferencesinverse.Computer, time.Time]
	UpdatedAt condition.Field[overridereferencesinverse.Computer, time.Time]
	DeletedAt condition.Field[overridereferencesinverse.Computer, time.Time]
	Name      condition.StringField[overridereferencesinverse.Computer]
}

var Computer = computerConditions{
	CreatedAt: condition.NewField[overridereferencesinverse.Computer, time.Time]("CreatedAt", "", ""),
	DeletedAt: condition.NewField[overridereferencesinverse.Computer, time.Time]("DeletedAt", "", ""),
	ID:        condition.NewField[overridereferencesinverse.Computer, model.UUID]("ID", "", ""),
	Name:      condition.NewStringField[overridereferencesinverse.Computer]("Name", "", ""),
	UpdatedAt: condition.NewField[overridereferencesinverse.Computer, time.Time]("UpdatedAt", "", ""),
}

// Preload allows preloading the Computer when doing a query
func (computerConditions computerConditions) preload() condition.Condition[overridereferencesinverse.Computer] {
	return condition.NewPreloadCondition[overridereferencesinverse.Computer](computerConditions.ID, computerConditions.CreatedAt, computerConditions.UpdatedAt, computerConditions.DeletedAt, computerConditions.Name)
}

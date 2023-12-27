// Code generated by cql-gen v0.0.5, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	overridereferencesinverse "github.com/FrancoLiberali/cql/cql-gen/cmd/gen/conditions/tests/overridereferencesinverse"
	model "github.com/FrancoLiberali/cql/model"
	"time"
)

type processorConditions struct {
	ID           condition.Field[overridereferencesinverse.Processor, model.UUID]
	CreatedAt    condition.Field[overridereferencesinverse.Processor, time.Time]
	UpdatedAt    condition.Field[overridereferencesinverse.Processor, time.Time]
	DeletedAt    condition.Field[overridereferencesinverse.Processor, time.Time]
	ComputerName condition.StringField[overridereferencesinverse.Processor]
}

var Processor = processorConditions{
	ComputerName: condition.StringField[overridereferencesinverse.Processor]{UpdatableField: condition.UpdatableField[overridereferencesinverse.Processor, string]{Field: condition.Field[overridereferencesinverse.Processor, string]{Name: "ComputerName"}}},
	CreatedAt:    condition.Field[overridereferencesinverse.Processor, time.Time]{Name: "CreatedAt"},
	DeletedAt:    condition.Field[overridereferencesinverse.Processor, time.Time]{Name: "DeletedAt"},
	ID:           condition.Field[overridereferencesinverse.Processor, model.UUID]{Name: "ID"},
	UpdatedAt:    condition.Field[overridereferencesinverse.Processor, time.Time]{Name: "UpdatedAt"},
}

// Preload allows preloading the Processor when doing a query
func (processorConditions processorConditions) preload() condition.Condition[overridereferencesinverse.Processor] {
	return condition.NewPreloadCondition[overridereferencesinverse.Processor](processorConditions.ID, processorConditions.CreatedAt, processorConditions.UpdatedAt, processorConditions.DeletedAt, processorConditions.ComputerName)
}

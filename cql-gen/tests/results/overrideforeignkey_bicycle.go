// Code generated by cql-gen v0.0.5, DO NOT EDIT.
package conditions

import (
	overrideforeignkey "github.com/FrancoLiberali/cql-gen/cmd/gen/conditions/tests/overrideforeignkey"
	condition "github.com/FrancoLiberali/cql/condition"
	model "github.com/FrancoLiberali/cql/model"
	"time"
)

func (bicycleConditions bicycleConditions) Owner(conditions ...condition.Condition[overrideforeignkey.Person]) condition.JoinCondition[overrideforeignkey.Bicycle] {
	return condition.NewJoinCondition[overrideforeignkey.Bicycle, overrideforeignkey.Person](conditions, "Owner", "OwnerSomethingID", bicycleConditions.preload(), "ID", Person.preload())
}

type bicycleConditions struct {
	ID               condition.Field[overrideforeignkey.Bicycle, model.UUID]
	CreatedAt        condition.Field[overrideforeignkey.Bicycle, time.Time]
	UpdatedAt        condition.Field[overrideforeignkey.Bicycle, time.Time]
	DeletedAt        condition.Field[overrideforeignkey.Bicycle, time.Time]
	OwnerSomethingID condition.StringField[overrideforeignkey.Bicycle]
}

var Bicycle = bicycleConditions{
	CreatedAt:        condition.Field[overrideforeignkey.Bicycle, time.Time]{Name: "CreatedAt"},
	DeletedAt:        condition.Field[overrideforeignkey.Bicycle, time.Time]{Name: "DeletedAt"},
	ID:               condition.Field[overrideforeignkey.Bicycle, model.UUID]{Name: "ID"},
	OwnerSomethingID: condition.StringField[overrideforeignkey.Bicycle]{UpdatableField: condition.UpdatableField[overrideforeignkey.Bicycle, string]{Field: condition.Field[overrideforeignkey.Bicycle, string]{Name: "OwnerSomethingID"}}},
	UpdatedAt:        condition.Field[overrideforeignkey.Bicycle, time.Time]{Name: "UpdatedAt"},
}

// Preload allows preloading the Bicycle when doing a query
func (bicycleConditions bicycleConditions) preload() condition.Condition[overrideforeignkey.Bicycle] {
	return condition.NewPreloadCondition[overrideforeignkey.Bicycle](bicycleConditions.ID, bicycleConditions.CreatedAt, bicycleConditions.UpdatedAt, bicycleConditions.DeletedAt, bicycleConditions.OwnerSomethingID)
}

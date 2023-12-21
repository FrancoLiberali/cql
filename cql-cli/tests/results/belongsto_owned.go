// Code generated by cql-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	belongsto "github.com/FrancoLiberali/cql/cql-cli/cmd/gen/conditions/tests/belongsto"
	model "github.com/FrancoLiberali/cql/model"
	"time"
)

func (ownedConditions ownedConditions) Owner(conditions ...condition.Condition[belongsto.Owner]) condition.JoinCondition[belongsto.Owned] {
	return condition.NewJoinCondition[belongsto.Owned, belongsto.Owner](conditions, "Owner", "OwnerID", ownedConditions.Preload(), "ID")
}
func (ownedConditions ownedConditions) PreloadOwner() condition.JoinCondition[belongsto.Owned] {
	return ownedConditions.Owner(Owner.Preload())
}

type ownedConditions struct {
	ID        condition.Field[belongsto.Owned, model.UUID]
	CreatedAt condition.Field[belongsto.Owned, time.Time]
	UpdatedAt condition.Field[belongsto.Owned, time.Time]
	DeletedAt condition.Field[belongsto.Owned, time.Time]
	OwnerID   condition.UpdatableField[belongsto.Owned, model.UUID]
}

var Owned = ownedConditions{
	CreatedAt: condition.Field[belongsto.Owned, time.Time]{Name: "CreatedAt"},
	DeletedAt: condition.Field[belongsto.Owned, time.Time]{Name: "DeletedAt"},
	ID:        condition.Field[belongsto.Owned, model.UUID]{Name: "ID"},
	OwnerID:   condition.UpdatableField[belongsto.Owned, model.UUID]{Field: condition.Field[belongsto.Owned, model.UUID]{Name: "OwnerID"}},
	UpdatedAt: condition.Field[belongsto.Owned, time.Time]{Name: "UpdatedAt"},
}

// Preload allows preloading the Owned when doing a query
func (ownedConditions ownedConditions) Preload() condition.Condition[belongsto.Owned] {
	return condition.NewPreloadCondition[belongsto.Owned](ownedConditions.ID, ownedConditions.CreatedAt, ownedConditions.UpdatedAt, ownedConditions.DeletedAt, ownedConditions.OwnerID)
}

// PreloadRelations allows preloading all the Owned's relation when doing a query
func (ownedConditions ownedConditions) PreloadRelations() []condition.Condition[belongsto.Owned] {
	return []condition.Condition[belongsto.Owned]{ownedConditions.PreloadOwner()}
}

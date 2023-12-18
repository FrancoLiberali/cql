// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	model "github.com/FrancoLiberali/cql/model"
	models "github.com/FrancoLiberali/cql/test/models"
	"time"
)

func (parent2Conditions parent2Conditions) ParentParent(conditions ...condition.Condition[models.ParentParent]) condition.JoinCondition[models.Parent2] {
	return condition.NewJoinCondition[models.Parent2, models.ParentParent](conditions, "ParentParent", "ParentParentID", parent2Conditions.Preload(), "ID")
}
func (parent2Conditions parent2Conditions) PreloadParentParent() condition.JoinCondition[models.Parent2] {
	return parent2Conditions.ParentParent(ParentParent.Preload())
}

type parent2Conditions struct {
	ID             condition.Field[models.Parent2, model.UUID]
	CreatedAt      condition.Field[models.Parent2, time.Time]
	UpdatedAt      condition.Field[models.Parent2, time.Time]
	DeletedAt      condition.Field[models.Parent2, time.Time]
	ParentParentID condition.Field[models.Parent2, model.UUID]
}

var Parent2 = parent2Conditions{
	CreatedAt:      condition.Field[models.Parent2, time.Time]{Name: "CreatedAt"},
	DeletedAt:      condition.Field[models.Parent2, time.Time]{Name: "DeletedAt"},
	ID:             condition.Field[models.Parent2, model.UUID]{Name: "ID"},
	ParentParentID: condition.Field[models.Parent2, model.UUID]{Name: "ParentParentID"},
	UpdatedAt:      condition.Field[models.Parent2, time.Time]{Name: "UpdatedAt"},
}

// Preload allows preloading the Parent2 when doing a query
func (parent2Conditions parent2Conditions) Preload() condition.Condition[models.Parent2] {
	return condition.NewPreloadCondition[models.Parent2](parent2Conditions.ID, parent2Conditions.CreatedAt, parent2Conditions.UpdatedAt, parent2Conditions.DeletedAt, parent2Conditions.ParentParentID)
}

// PreloadRelations allows preloading all the Parent2's relation when doing a query
func (parent2Conditions parent2Conditions) PreloadRelations() []condition.Condition[models.Parent2] {
	return []condition.Condition[models.Parent2]{parent2Conditions.PreloadParentParent()}
}
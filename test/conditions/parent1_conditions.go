// Code generated by cql-gen v0.0.6, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	model "github.com/FrancoLiberali/cql/model"
	models "github.com/FrancoLiberali/cql/test/models"
	"time"
)

func (parent1Conditions parent1Conditions) ParentParent(conditions ...condition.Condition[models.ParentParent]) condition.JoinCondition[models.Parent1] {
	return condition.NewJoinCondition[models.Parent1, models.ParentParent](conditions, "ParentParent", "ParentParentID", parent1Conditions.preload(), "ID", ParentParent.preload())
}

type parent1Conditions struct {
	ID             condition.Field[models.Parent1, model.UUID]
	CreatedAt      condition.Field[models.Parent1, time.Time]
	UpdatedAt      condition.Field[models.Parent1, time.Time]
	DeletedAt      condition.Field[models.Parent1, time.Time]
	ParentParentID condition.UpdatableField[models.Parent1, model.UUID]
}

var Parent1 = parent1Conditions{
	CreatedAt:      condition.Field[models.Parent1, time.Time]{Name: "CreatedAt"},
	DeletedAt:      condition.Field[models.Parent1, time.Time]{Name: "DeletedAt"},
	ID:             condition.Field[models.Parent1, model.UUID]{Name: "ID"},
	ParentParentID: condition.UpdatableField[models.Parent1, model.UUID]{Field: condition.Field[models.Parent1, model.UUID]{Name: "ParentParentID"}},
	UpdatedAt:      condition.Field[models.Parent1, time.Time]{Name: "UpdatedAt"},
}

// Preload allows preloading the Parent1 when doing a query
func (parent1Conditions parent1Conditions) preload() condition.Condition[models.Parent1] {
	return condition.NewPreloadCondition[models.Parent1](parent1Conditions.ID, parent1Conditions.CreatedAt, parent1Conditions.UpdatedAt, parent1Conditions.DeletedAt, parent1Conditions.ParentParentID)
}

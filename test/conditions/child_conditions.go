// Code generated by cql-gen v0.0.6, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	model "github.com/FrancoLiberali/cql/model"
	models "github.com/FrancoLiberali/cql/test/models"
	"time"
)

func (childConditions childConditions) Parent1(conditions ...condition.Condition[models.Parent1]) condition.JoinCondition[models.Child] {
	return condition.NewJoinCondition[models.Child, models.Parent1](conditions, "Parent1", "Parent1ID", childConditions.preload(), "ID", Parent1.preload())
}
func (childConditions childConditions) Parent2(conditions ...condition.Condition[models.Parent2]) condition.JoinCondition[models.Child] {
	return condition.NewJoinCondition[models.Child, models.Parent2](conditions, "Parent2", "Parent2ID", childConditions.preload(), "ID", Parent2.preload())
}

type childConditions struct {
	ID        condition.Field[models.Child, model.UUID]
	CreatedAt condition.Field[models.Child, time.Time]
	UpdatedAt condition.Field[models.Child, time.Time]
	DeletedAt condition.Field[models.Child, time.Time]
	Name      condition.StringField[models.Child]
	Number    condition.UpdatableField[models.Child, int]
	Parent1ID condition.UpdatableField[models.Child, model.UUID]
	Parent2ID condition.UpdatableField[models.Child, model.UUID]
}

var Child = childConditions{
	CreatedAt: condition.Field[models.Child, time.Time]{Name: "CreatedAt"},
	DeletedAt: condition.Field[models.Child, time.Time]{Name: "DeletedAt"},
	ID:        condition.Field[models.Child, model.UUID]{Name: "ID"},
	Name:      condition.StringField[models.Child]{UpdatableField: condition.UpdatableField[models.Child, string]{Field: condition.Field[models.Child, string]{Name: "Name"}}},
	Number:    condition.UpdatableField[models.Child, int]{Field: condition.Field[models.Child, int]{Name: "Number"}},
	Parent1ID: condition.UpdatableField[models.Child, model.UUID]{Field: condition.Field[models.Child, model.UUID]{Name: "Parent1ID"}},
	Parent2ID: condition.UpdatableField[models.Child, model.UUID]{Field: condition.Field[models.Child, model.UUID]{Name: "Parent2ID"}},
	UpdatedAt: condition.Field[models.Child, time.Time]{Name: "UpdatedAt"},
}

// Preload allows preloading the Child when doing a query
func (childConditions childConditions) preload() condition.Condition[models.Child] {
	return condition.NewPreloadCondition[models.Child](childConditions.ID, childConditions.CreatedAt, childConditions.UpdatedAt, childConditions.DeletedAt, childConditions.Name, childConditions.Number, childConditions.Parent1ID, childConditions.Parent2ID)
}

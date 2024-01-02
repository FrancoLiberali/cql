// Code generated by cql-gen v0.0.6, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	model "github.com/FrancoLiberali/cql/model"
	models "github.com/FrancoLiberali/cql/test/models"
	"time"
)

func (phoneConditions phoneConditions) Brand(conditions ...condition.Condition[models.Brand]) condition.JoinCondition[models.Phone] {
	return condition.NewJoinCondition[models.Phone, models.Brand](conditions, "Brand", "BrandID", phoneConditions.preload(), "ID", Brand.preload())
}

type phoneConditions struct {
	ID        condition.Field[models.Phone, model.UIntID]
	CreatedAt condition.Field[models.Phone, time.Time]
	UpdatedAt condition.Field[models.Phone, time.Time]
	DeletedAt condition.Field[models.Phone, time.Time]
	Name      condition.StringField[models.Phone]
	BrandID   condition.UpdatableField[models.Phone, uint]
}

var Phone = phoneConditions{
	BrandID:   condition.UpdatableField[models.Phone, uint]{Field: condition.Field[models.Phone, uint]{Name: "BrandID"}},
	CreatedAt: condition.Field[models.Phone, time.Time]{Name: "CreatedAt"},
	DeletedAt: condition.Field[models.Phone, time.Time]{Name: "DeletedAt"},
	ID:        condition.Field[models.Phone, model.UIntID]{Name: "ID"},
	Name:      condition.StringField[models.Phone]{UpdatableField: condition.UpdatableField[models.Phone, string]{Field: condition.Field[models.Phone, string]{Name: "Name"}}},
	UpdatedAt: condition.Field[models.Phone, time.Time]{Name: "UpdatedAt"},
}

// Preload allows preloading the Phone when doing a query
func (phoneConditions phoneConditions) preload() condition.Condition[models.Phone] {
	return condition.NewPreloadCondition[models.Phone](phoneConditions.ID, phoneConditions.CreatedAt, phoneConditions.UpdatedAt, phoneConditions.DeletedAt, phoneConditions.Name, phoneConditions.BrandID)
}

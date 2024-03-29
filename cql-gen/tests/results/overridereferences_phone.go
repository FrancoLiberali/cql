// Code generated by cql-gen v0.1.0, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	overridereferences "github.com/FrancoLiberali/cql/cql-gen/cmd/gen/conditions/tests/overridereferences"
	model "github.com/FrancoLiberali/cql/model"
	"time"
)

func (phoneConditions phoneConditions) Brand(conditions ...condition.Condition[overridereferences.Brand]) condition.JoinCondition[overridereferences.Phone] {
	return condition.NewJoinCondition[overridereferences.Phone, overridereferences.Brand](conditions, "Brand", "BrandName", phoneConditions.preload(), "Name", Brand.preload())
}

type phoneConditions struct {
	ID        condition.Field[overridereferences.Phone, model.UUID]
	CreatedAt condition.Field[overridereferences.Phone, time.Time]
	UpdatedAt condition.Field[overridereferences.Phone, time.Time]
	DeletedAt condition.Field[overridereferences.Phone, time.Time]
	BrandName condition.StringField[overridereferences.Phone]
}

var Phone = phoneConditions{
	BrandName: condition.NewStringField[overridereferences.Phone]("BrandName", "", ""),
	CreatedAt: condition.NewField[overridereferences.Phone, time.Time]("CreatedAt", "", ""),
	DeletedAt: condition.NewField[overridereferences.Phone, time.Time]("DeletedAt", "", ""),
	ID:        condition.NewField[overridereferences.Phone, model.UUID]("ID", "", ""),
	UpdatedAt: condition.NewField[overridereferences.Phone, time.Time]("UpdatedAt", "", ""),
}

// Preload allows preloading the Phone when doing a query
func (phoneConditions phoneConditions) preload() condition.Condition[overridereferences.Phone] {
	return condition.NewPreloadCondition[overridereferences.Phone](phoneConditions.ID, phoneConditions.CreatedAt, phoneConditions.UpdatedAt, phoneConditions.DeletedAt, phoneConditions.BrandName)
}

// Code generated by cql-gen v0.0.6, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	model "github.com/FrancoLiberali/cql/model"
	models "github.com/FrancoLiberali/cql/test/models"
	"time"
)

func (sellerConditions sellerConditions) Company(conditions ...condition.Condition[models.Company]) condition.JoinCondition[models.Seller] {
	return condition.NewJoinCondition[models.Seller, models.Company](conditions, "Company", "CompanyID", sellerConditions.preload(), "ID", Company.preload())
}
func (sellerConditions sellerConditions) University(conditions ...condition.Condition[models.University]) condition.JoinCondition[models.Seller] {
	return condition.NewJoinCondition[models.Seller, models.University](conditions, "University", "UniversityID", sellerConditions.preload(), "ID", University.preload())
}

type sellerConditions struct {
	ID           condition.Field[models.Seller, model.UUID]
	CreatedAt    condition.Field[models.Seller, time.Time]
	UpdatedAt    condition.Field[models.Seller, time.Time]
	DeletedAt    condition.Field[models.Seller, time.Time]
	Name         condition.StringField[models.Seller]
	CompanyID    condition.NullableField[models.Seller, model.UUID]
	UniversityID condition.NullableField[models.Seller, model.UUID]
}

var Seller = sellerConditions{
	CompanyID:    condition.NewNullableField[models.Seller, model.UUID]("CompanyID", "", ""),
	CreatedAt:    condition.NewField[models.Seller, time.Time]("CreatedAt", "", ""),
	DeletedAt:    condition.NewField[models.Seller, time.Time]("DeletedAt", "", ""),
	ID:           condition.NewField[models.Seller, model.UUID]("ID", "", ""),
	Name:         condition.NewStringField[models.Seller]("Name", "", ""),
	UniversityID: condition.NewNullableField[models.Seller, model.UUID]("UniversityID", "", ""),
	UpdatedAt:    condition.NewField[models.Seller, time.Time]("UpdatedAt", "", ""),
}

// Preload allows preloading the Seller when doing a query
func (sellerConditions sellerConditions) preload() condition.Condition[models.Seller] {
	return condition.NewPreloadCondition[models.Seller](sellerConditions.ID, sellerConditions.CreatedAt, sellerConditions.UpdatedAt, sellerConditions.DeletedAt, sellerConditions.Name, sellerConditions.CompanyID, sellerConditions.UniversityID)
}

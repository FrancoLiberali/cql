// Code generated by cql-gen v0.0.6, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	model "github.com/FrancoLiberali/cql/model"
	models "github.com/FrancoLiberali/cql/test/models"
	"time"
)

type companyConditions struct {
	ID        condition.Field[models.Company, model.UUID]
	CreatedAt condition.Field[models.Company, time.Time]
	UpdatedAt condition.Field[models.Company, time.Time]
	DeletedAt condition.Field[models.Company, time.Time]
	Name      condition.StringField[models.Company]
	Sellers   condition.Collection[models.Company, models.Seller]
}

var Company = companyConditions{
	CreatedAt: condition.NewField[models.Company, time.Time]("CreatedAt", "", ""),
	DeletedAt: condition.NewField[models.Company, time.Time]("DeletedAt", "", ""),
	ID:        condition.NewField[models.Company, model.UUID]("ID", "", ""),
	Name:      condition.NewStringField[models.Company]("Name", "", ""),
	Sellers:   condition.NewCollection[models.Company, models.Seller]("Sellers", "ID", "CompanyID"),
	UpdatedAt: condition.NewField[models.Company, time.Time]("UpdatedAt", "", ""),
}

// Preload allows preloading the Company when doing a query
func (companyConditions companyConditions) preload() condition.Condition[models.Company] {
	return condition.NewPreloadCondition[models.Company](companyConditions.ID, companyConditions.CreatedAt, companyConditions.UpdatedAt, companyConditions.DeletedAt, companyConditions.Name)
}

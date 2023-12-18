// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	model "github.com/FrancoLiberali/cql/model"
	models "github.com/FrancoLiberali/cql/test/models"
	"time"
)

func (cityConditions cityConditions) Country(conditions ...condition.Condition[models.Country]) condition.JoinCondition[models.City] {
	return condition.NewJoinCondition[models.City, models.Country](conditions, "Country", "CountryID", cityConditions.Preload(), "ID")
}
func (cityConditions cityConditions) PreloadCountry() condition.JoinCondition[models.City] {
	return cityConditions.Country(Country.Preload())
}

type cityConditions struct {
	ID        condition.Field[models.City, model.UUID]
	CreatedAt condition.Field[models.City, time.Time]
	UpdatedAt condition.Field[models.City, time.Time]
	DeletedAt condition.Field[models.City, time.Time]
	Name      condition.StringField[models.City]
	CountryID condition.Field[models.City, model.UUID]
}

var City = cityConditions{
	CountryID: condition.Field[models.City, model.UUID]{Name: "CountryID"},
	CreatedAt: condition.Field[models.City, time.Time]{Name: "CreatedAt"},
	DeletedAt: condition.Field[models.City, time.Time]{Name: "DeletedAt"},
	ID:        condition.Field[models.City, model.UUID]{Name: "ID"},
	Name:      condition.StringField[models.City]{Field: condition.Field[models.City, string]{Name: "Name"}},
	UpdatedAt: condition.Field[models.City, time.Time]{Name: "UpdatedAt"},
}

// Preload allows preloading the City when doing a query
func (cityConditions cityConditions) Preload() condition.Condition[models.City] {
	return condition.NewPreloadCondition[models.City](cityConditions.ID, cityConditions.CreatedAt, cityConditions.UpdatedAt, cityConditions.DeletedAt, cityConditions.Name, cityConditions.CountryID)
}

// PreloadRelations allows preloading all the City's relation when doing a query
func (cityConditions cityConditions) PreloadRelations() []condition.Condition[models.City] {
	return []condition.Condition[models.City]{cityConditions.PreloadCountry()}
}
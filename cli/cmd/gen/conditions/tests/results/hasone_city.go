// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	hasone "github.com/ditrit/badaas-cli/cmd/gen/conditions/tests/hasone"
	condition "github.com/ditrit/badaas/orm/condition"
	model "github.com/ditrit/badaas/orm/model"
	operator "github.com/ditrit/badaas/orm/operator"
	query "github.com/ditrit/badaas/orm/query"
	"reflect"
	"time"
)

var cityType = reflect.TypeOf(*new(hasone.City))
var CityIdField = query.FieldIdentifier[model.UUID]{
	Field:     "ID",
	ModelType: cityType,
}

func CityId(operator operator.Operator[model.UUID]) condition.WhereCondition[hasone.City] {
	return condition.NewFieldCondition[hasone.City, model.UUID](CityIdField, operator)
}

var CityCreatedAtField = query.FieldIdentifier[time.Time]{
	Field:     "CreatedAt",
	ModelType: cityType,
}

func CityCreatedAt(operator operator.Operator[time.Time]) condition.WhereCondition[hasone.City] {
	return condition.NewFieldCondition[hasone.City, time.Time](CityCreatedAtField, operator)
}

var CityUpdatedAtField = query.FieldIdentifier[time.Time]{
	Field:     "UpdatedAt",
	ModelType: cityType,
}

func CityUpdatedAt(operator operator.Operator[time.Time]) condition.WhereCondition[hasone.City] {
	return condition.NewFieldCondition[hasone.City, time.Time](CityUpdatedAtField, operator)
}

var CityDeletedAtField = query.FieldIdentifier[time.Time]{
	Field:     "DeletedAt",
	ModelType: cityType,
}

func CityDeletedAt(operator operator.Operator[time.Time]) condition.WhereCondition[hasone.City] {
	return condition.NewFieldCondition[hasone.City, time.Time](CityDeletedAtField, operator)
}
func CityCountry(conditions ...condition.Condition[hasone.Country]) condition.JoinCondition[hasone.City] {
	return condition.NewJoinCondition[hasone.City, hasone.Country](conditions, "Country", "CountryID", CityPreloadAttributes, "ID")
}

var CityPreloadCountry = CityCountry(CountryPreloadAttributes)
var CityCountryIdField = query.FieldIdentifier[model.UUID]{
	Field:     "CountryID",
	ModelType: cityType,
}

func CityCountryId(operator operator.Operator[model.UUID]) condition.WhereCondition[hasone.City] {
	return condition.NewFieldCondition[hasone.City, model.UUID](CityCountryIdField, operator)
}

var CityPreloadAttributes = condition.NewPreloadCondition[hasone.City](CityIdField, CityCreatedAtField, CityUpdatedAtField, CityDeletedAtField, CityCountryIdField)
var CityPreloadRelations = []condition.Condition[hasone.City]{CityPreloadCountry}

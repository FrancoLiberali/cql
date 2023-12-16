// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	hasone "github.com/ditrit/badaas-orm/cli/cmd/gen/conditions/tests/hasone"
	orm "github.com/ditrit/badaas/orm"
	"time"
)

func CountryId(operator orm.Operator[orm.UUID]) orm.WhereCondition[hasone.Country] {
	return orm.FieldCondition[hasone.Country, orm.UUID]{
		FieldIdentifier: orm.IDFieldID,
		Operator:        operator,
	}
}
func CountryCreatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[hasone.Country] {
	return orm.FieldCondition[hasone.Country, time.Time]{
		FieldIdentifier: orm.CreatedAtFieldID,
		Operator:        operator,
	}
}
func CountryUpdatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[hasone.Country] {
	return orm.FieldCondition[hasone.Country, time.Time]{
		FieldIdentifier: orm.UpdatedAtFieldID,
		Operator:        operator,
	}
}
func CountryDeletedAt(operator orm.Operator[time.Time]) orm.WhereCondition[hasone.Country] {
	return orm.FieldCondition[hasone.Country, time.Time]{
		FieldIdentifier: orm.DeletedAtFieldID,
		Operator:        operator,
	}
}
func CountryCapital(conditions ...orm.Condition[hasone.City]) orm.IJoinCondition[hasone.Country] {
	return orm.JoinCondition[hasone.Country, hasone.City]{
		Conditions:         conditions,
		RelationField:      "Capital",
		T1Field:            "ID",
		T1PreloadCondition: CountryPreloadAttributes,
		T2Field:            "CountryID",
	}
}

var CountryPreloadCapital = CountryCapital(CityPreloadAttributes)
var CountryPreloadAttributes = orm.NewPreloadCondition[hasone.Country]()
var CountryPreloadRelations = []orm.Condition[hasone.Country]{CountryPreloadCapital}

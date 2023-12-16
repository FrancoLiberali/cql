// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	hasone "github.com/ditrit/badaas-orm/cli/cmd/gen/conditions/tests/hasone"
	orm "github.com/ditrit/badaas/orm"
	gorm "gorm.io/gorm"
	"time"
)

func CountryId(v orm.UUID) orm.WhereCondition[hasone.Country] {
	return orm.WhereCondition[hasone.Country]{
		Field: "ID",
		Value: v,
	}
}
func CountryCreatedAt(v time.Time) orm.WhereCondition[hasone.Country] {
	return orm.WhereCondition[hasone.Country]{
		Field: "CreatedAt",
		Value: v,
	}
}
func CountryUpdatedAt(v time.Time) orm.WhereCondition[hasone.Country] {
	return orm.WhereCondition[hasone.Country]{
		Field: "UpdatedAt",
		Value: v,
	}
}
func CountryDeletedAt(v gorm.DeletedAt) orm.WhereCondition[hasone.Country] {
	return orm.WhereCondition[hasone.Country]{
		Field: "DeletedAt",
		Value: v,
	}
}
func CountryCapital(conditions ...orm.Condition[hasone.City]) orm.Condition[hasone.Country] {
	return orm.JoinCondition[hasone.Country, hasone.City]{
		Conditions: conditions,
		T1Field:    "ID",
		T2Field:    "CountryID",
	}
}
func CityCountry(conditions ...orm.Condition[hasone.Country]) orm.Condition[hasone.City] {
	return orm.JoinCondition[hasone.City, hasone.Country]{
		Conditions: conditions,
		T1Field:    "CountryID",
		T2Field:    "ID",
	}
}

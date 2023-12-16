// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	overridereferences "github.com/ditrit/badaas-orm/cli/cmd/gen/conditions/tests/overridereferences"
	orm "github.com/ditrit/badaas/orm"
	"time"
)

func PhoneId(operator orm.Operator[orm.UUID]) orm.WhereCondition[overridereferences.Phone] {
	return orm.FieldCondition[overridereferences.Phone, orm.UUID]{
		Field:    "ID",
		Operator: operator,
	}
}
func PhoneCreatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[overridereferences.Phone] {
	return orm.FieldCondition[overridereferences.Phone, time.Time]{
		Field:    "CreatedAt",
		Operator: operator,
	}
}
func PhoneUpdatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[overridereferences.Phone] {
	return orm.FieldCondition[overridereferences.Phone, time.Time]{
		Field:    "UpdatedAt",
		Operator: operator,
	}
}
func PhoneDeletedAt(operator orm.Operator[time.Time]) orm.WhereCondition[overridereferences.Phone] {
	return orm.FieldCondition[overridereferences.Phone, time.Time]{
		Field:    "DeletedAt",
		Operator: operator,
	}
}
func PhoneBrand(conditions ...orm.Condition[overridereferences.Brand]) orm.Condition[overridereferences.Phone] {
	return orm.JoinCondition[overridereferences.Phone, overridereferences.Brand]{
		Conditions: conditions,
		T1Field:    "BrandName",
		T2Field:    "Name",
	}
}
func PhoneBrandName(operator orm.Operator[string]) orm.WhereCondition[overridereferences.Phone] {
	return orm.FieldCondition[overridereferences.Phone, string]{
		Field:    "BrandName",
		Operator: operator,
	}
}

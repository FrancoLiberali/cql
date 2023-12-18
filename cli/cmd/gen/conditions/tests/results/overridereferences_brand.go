// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	overridereferences "github.com/ditrit/badaas-cli/cmd/gen/conditions/tests/overridereferences"
	orm "github.com/ditrit/badaas/orm"
	condition "github.com/ditrit/badaas/orm/condition"
	model "github.com/ditrit/badaas/orm/model"
	query "github.com/ditrit/badaas/orm/query"
	"reflect"
	"time"
)

var brandType = reflect.TypeOf(*new(overridereferences.Brand))

func (brandConditions brandConditions) IdIs() orm.FieldIs[overridereferences.Brand, model.UUID] {
	return orm.FieldIs[overridereferences.Brand, model.UUID]{FieldID: brandConditions.ID}
}
func (brandConditions brandConditions) CreatedAtIs() orm.FieldIs[overridereferences.Brand, time.Time] {
	return orm.FieldIs[overridereferences.Brand, time.Time]{FieldID: brandConditions.CreatedAt}
}
func (brandConditions brandConditions) UpdatedAtIs() orm.FieldIs[overridereferences.Brand, time.Time] {
	return orm.FieldIs[overridereferences.Brand, time.Time]{FieldID: brandConditions.UpdatedAt}
}
func (brandConditions brandConditions) DeletedAtIs() orm.FieldIs[overridereferences.Brand, time.Time] {
	return orm.FieldIs[overridereferences.Brand, time.Time]{FieldID: brandConditions.DeletedAt}
}
func (brandConditions brandConditions) NameIs() orm.StringFieldIs[overridereferences.Brand] {
	return orm.StringFieldIs[overridereferences.Brand]{FieldIs: orm.FieldIs[overridereferences.Brand, string]{FieldID: brandConditions.Name}}
}

type brandConditions struct {
	ID        query.FieldIdentifier[model.UUID]
	CreatedAt query.FieldIdentifier[time.Time]
	UpdatedAt query.FieldIdentifier[time.Time]
	DeletedAt query.FieldIdentifier[time.Time]
	Name      query.FieldIdentifier[string]
}

var Brand = brandConditions{
	CreatedAt: query.FieldIdentifier[time.Time]{
		Field:     "CreatedAt",
		ModelType: brandType,
	},
	DeletedAt: query.FieldIdentifier[time.Time]{
		Field:     "DeletedAt",
		ModelType: brandType,
	},
	ID: query.FieldIdentifier[model.UUID]{
		Field:     "ID",
		ModelType: brandType,
	},
	Name: query.FieldIdentifier[string]{
		Field:     "Name",
		ModelType: brandType,
	},
	UpdatedAt: query.FieldIdentifier[time.Time]{
		Field:     "UpdatedAt",
		ModelType: brandType,
	},
}

// Preload allows preloading the Brand when doing a query
func (brandConditions brandConditions) Preload() condition.Condition[overridereferences.Brand] {
	return condition.NewPreloadCondition[overridereferences.Brand](brandConditions.ID, brandConditions.CreatedAt, brandConditions.UpdatedAt, brandConditions.DeletedAt, brandConditions.Name)
}
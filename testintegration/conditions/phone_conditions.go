// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	condition "github.com/ditrit/badaas/orm/condition"
	model "github.com/ditrit/badaas/orm/model"
	query "github.com/ditrit/badaas/orm/query"
	models "github.com/ditrit/badaas/testintegration/models"
	"reflect"
	"time"
)

var phoneType = reflect.TypeOf(*new(models.Phone))

func (phoneConditions phoneConditions) IdIs() orm.FieldIs[models.Phone, model.UIntID] {
	return orm.FieldIs[models.Phone, model.UIntID]{FieldID: phoneConditions.ID}
}
func (phoneConditions phoneConditions) CreatedAtIs() orm.FieldIs[models.Phone, time.Time] {
	return orm.FieldIs[models.Phone, time.Time]{FieldID: phoneConditions.CreatedAt}
}
func (phoneConditions phoneConditions) UpdatedAtIs() orm.FieldIs[models.Phone, time.Time] {
	return orm.FieldIs[models.Phone, time.Time]{FieldID: phoneConditions.UpdatedAt}
}
func (phoneConditions phoneConditions) DeletedAtIs() orm.FieldIs[models.Phone, time.Time] {
	return orm.FieldIs[models.Phone, time.Time]{FieldID: phoneConditions.DeletedAt}
}
func (phoneConditions phoneConditions) NameIs() orm.StringFieldIs[models.Phone] {
	return orm.StringFieldIs[models.Phone]{FieldIs: orm.FieldIs[models.Phone, string]{FieldID: phoneConditions.Name}}
}
func (phoneConditions phoneConditions) Brand(conditions ...condition.Condition[models.Brand]) condition.JoinCondition[models.Phone] {
	return condition.NewJoinCondition[models.Phone, models.Brand](conditions, "Brand", "BrandID", phoneConditions.Preload(), "ID")
}
func (phoneConditions phoneConditions) PreloadBrand() condition.JoinCondition[models.Phone] {
	return phoneConditions.Brand(Brand.Preload())
}
func (phoneConditions phoneConditions) BrandIdIs() orm.FieldIs[models.Phone, uint] {
	return orm.FieldIs[models.Phone, uint]{FieldID: phoneConditions.BrandID}
}

type phoneConditions struct {
	ID        query.Field[model.UIntID]
	CreatedAt query.Field[time.Time]
	UpdatedAt query.Field[time.Time]
	DeletedAt query.Field[time.Time]
	Name      query.Field[string]
	BrandID   query.Field[uint]
}

var Phone = phoneConditions{
	BrandID: query.Field[uint]{
		Field:     "BrandID",
		ModelType: phoneType,
	},
	CreatedAt: query.Field[time.Time]{
		Field:     "CreatedAt",
		ModelType: phoneType,
	},
	DeletedAt: query.Field[time.Time]{
		Field:     "DeletedAt",
		ModelType: phoneType,
	},
	ID: query.Field[model.UIntID]{
		Field:     "ID",
		ModelType: phoneType,
	},
	Name: query.Field[string]{
		Field:     "Name",
		ModelType: phoneType,
	},
	UpdatedAt: query.Field[time.Time]{
		Field:     "UpdatedAt",
		ModelType: phoneType,
	},
}

// Preload allows preloading the Phone when doing a query
func (phoneConditions phoneConditions) Preload() condition.Condition[models.Phone] {
	return condition.NewPreloadCondition[models.Phone](phoneConditions.ID, phoneConditions.CreatedAt, phoneConditions.UpdatedAt, phoneConditions.DeletedAt, phoneConditions.Name, phoneConditions.BrandID)
}

// PreloadRelations allows preloading all the Phone's relation when doing a query
func (phoneConditions phoneConditions) PreloadRelations() []condition.Condition[models.Phone] {
	return []condition.Condition[models.Phone]{phoneConditions.PreloadBrand()}
}

func (phoneConditions phoneConditions) NameSet() query.FieldSet[models.Phone, string] {
	return query.FieldSet[models.Phone, string]{FieldID: phoneConditions.Name}
}

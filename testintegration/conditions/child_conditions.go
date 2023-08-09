// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	condition "github.com/ditrit/badaas/orm/condition"
	model "github.com/ditrit/badaas/orm/model"
	operator "github.com/ditrit/badaas/orm/operator"
	query "github.com/ditrit/badaas/orm/query"
	models "github.com/ditrit/badaas/testintegration/models"
	"reflect"
	"time"
)

var childType = reflect.TypeOf(*new(models.Child))
var ChildIdField = query.FieldIdentifier[model.UUID]{
	Field:     "ID",
	ModelType: childType,
}

func ChildId(operator operator.Operator[model.UUID]) condition.WhereCondition[models.Child] {
	return condition.NewFieldCondition[models.Child, model.UUID](ChildIdField, operator)
}

var ChildCreatedAtField = query.FieldIdentifier[time.Time]{
	Field:     "CreatedAt",
	ModelType: childType,
}

func ChildCreatedAt(operator operator.Operator[time.Time]) condition.WhereCondition[models.Child] {
	return condition.NewFieldCondition[models.Child, time.Time](ChildCreatedAtField, operator)
}

var ChildUpdatedAtField = query.FieldIdentifier[time.Time]{
	Field:     "UpdatedAt",
	ModelType: childType,
}

func ChildUpdatedAt(operator operator.Operator[time.Time]) condition.WhereCondition[models.Child] {
	return condition.NewFieldCondition[models.Child, time.Time](ChildUpdatedAtField, operator)
}

var ChildDeletedAtField = query.FieldIdentifier[time.Time]{
	Field:     "DeletedAt",
	ModelType: childType,
}

func ChildDeletedAt(operator operator.Operator[time.Time]) condition.WhereCondition[models.Child] {
	return condition.NewFieldCondition[models.Child, time.Time](ChildDeletedAtField, operator)
}

var ChildNameField = query.FieldIdentifier[string]{
	Field:     "Name",
	ModelType: childType,
}

func ChildName(operator operator.Operator[string]) condition.WhereCondition[models.Child] {
	return condition.NewFieldCondition[models.Child, string](ChildNameField, operator)
}

var ChildNumberField = query.FieldIdentifier[int]{
	Field:     "Number",
	ModelType: childType,
}

func ChildNumber(operator operator.Operator[int]) condition.WhereCondition[models.Child] {
	return condition.NewFieldCondition[models.Child, int](ChildNumberField, operator)
}
func ChildParent1(conditions ...condition.Condition[models.Parent1]) condition.JoinCondition[models.Child] {
	return condition.NewJoinCondition[models.Child, models.Parent1](conditions, "Parent1", "Parent1ID", ChildPreloadAttributes, "ID")
}

var ChildPreloadParent1 = ChildParent1(Parent1PreloadAttributes)
var ChildParent1IdField = query.FieldIdentifier[model.UUID]{
	Field:     "Parent1ID",
	ModelType: childType,
}

func ChildParent1Id(operator operator.Operator[model.UUID]) condition.WhereCondition[models.Child] {
	return condition.NewFieldCondition[models.Child, model.UUID](ChildParent1IdField, operator)
}
func ChildParent2(conditions ...condition.Condition[models.Parent2]) condition.JoinCondition[models.Child] {
	return condition.NewJoinCondition[models.Child, models.Parent2](conditions, "Parent2", "Parent2ID", ChildPreloadAttributes, "ID")
}

var ChildPreloadParent2 = ChildParent2(Parent2PreloadAttributes)
var ChildParent2IdField = query.FieldIdentifier[model.UUID]{
	Field:     "Parent2ID",
	ModelType: childType,
}

func ChildParent2Id(operator operator.Operator[model.UUID]) condition.WhereCondition[models.Child] {
	return condition.NewFieldCondition[models.Child, model.UUID](ChildParent2IdField, operator)
}

var ChildPreloadAttributes = condition.NewPreloadCondition[models.Child](ChildIdField, ChildCreatedAtField, ChildUpdatedAtField, ChildDeletedAtField, ChildNameField, ChildNumberField, ChildParent1IdField, ChildParent2IdField)
var ChildPreloadRelations = []condition.Condition[models.Child]{ChildPreloadParent1, ChildPreloadParent2}

// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func Parent2Id(operator orm.Operator[orm.UUID]) orm.WhereCondition[models.Parent2] {
	return orm.FieldCondition[models.Parent2, orm.UUID]{
		FieldIdentifier: orm.IDFieldID,
		Operator:        operator,
	}
}
func Parent2CreatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Parent2] {
	return orm.FieldCondition[models.Parent2, time.Time]{
		FieldIdentifier: orm.CreatedAtFieldID,
		Operator:        operator,
	}
}
func Parent2UpdatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Parent2] {
	return orm.FieldCondition[models.Parent2, time.Time]{
		FieldIdentifier: orm.UpdatedAtFieldID,
		Operator:        operator,
	}
}
func Parent2DeletedAt(operator orm.Operator[gorm.DeletedAt]) orm.WhereCondition[models.Parent2] {
	return orm.FieldCondition[models.Parent2, gorm.DeletedAt]{
		FieldIdentifier: orm.DeletedAtFieldID,
		Operator:        operator,
	}
}
func Parent2ParentParent(conditions ...orm.Condition[models.ParentParent]) orm.IJoinCondition[models.Parent2] {
	return orm.JoinCondition[models.Parent2, models.ParentParent]{
		Conditions:         conditions,
		RelationField:      "ParentParent",
		T1Field:            "ParentParentID",
		T1PreloadCondition: Parent2PreloadAttributes,
		T2Field:            "ID",
	}
}

var Parent2PreloadParentParent = Parent2ParentParent(ParentParentPreloadAttributes)
var parent2ParentParentIdFieldID = orm.FieldIdentifier{Field: "ParentParentID"}

func Parent2ParentParentId(operator orm.Operator[orm.UUID]) orm.WhereCondition[models.Parent2] {
	return orm.FieldCondition[models.Parent2, orm.UUID]{
		FieldIdentifier: parent2ParentParentIdFieldID,
		Operator:        operator,
	}
}

var Parent2PreloadAttributes = orm.NewPreloadCondition[models.Parent2](parent2ParentParentIdFieldID)
var Parent2PreloadRelations = []orm.Condition[models.Parent2]{Parent2PreloadParentParent}

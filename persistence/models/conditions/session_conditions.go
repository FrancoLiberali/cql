// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	condition "github.com/ditrit/badaas/orm/condition"
	model "github.com/ditrit/badaas/orm/model"
	query "github.com/ditrit/badaas/orm/query"
	models "github.com/ditrit/badaas/persistence/models"
	"reflect"
	"time"
)

var sessionType = reflect.TypeOf(*new(models.Session))

func (sessionConditions sessionConditions) IdIs() orm.FieldIs[models.Session, model.UUID] {
	return orm.FieldIs[models.Session, model.UUID]{FieldID: sessionConditions.ID}
}
func (sessionConditions sessionConditions) CreatedAtIs() orm.FieldIs[models.Session, time.Time] {
	return orm.FieldIs[models.Session, time.Time]{FieldID: sessionConditions.CreatedAt}
}
func (sessionConditions sessionConditions) UpdatedAtIs() orm.FieldIs[models.Session, time.Time] {
	return orm.FieldIs[models.Session, time.Time]{FieldID: sessionConditions.UpdatedAt}
}
func (sessionConditions sessionConditions) DeletedAtIs() orm.FieldIs[models.Session, time.Time] {
	return orm.FieldIs[models.Session, time.Time]{FieldID: sessionConditions.DeletedAt}
}
func (sessionConditions sessionConditions) UserIdIs() orm.FieldIs[models.Session, model.UUID] {
	return orm.FieldIs[models.Session, model.UUID]{FieldID: sessionConditions.UserID}
}
func (sessionConditions sessionConditions) ExpiresAtIs() orm.FieldIs[models.Session, time.Time] {
	return orm.FieldIs[models.Session, time.Time]{FieldID: sessionConditions.ExpiresAt}
}

type sessionConditions struct {
	ID        query.FieldIdentifier[model.UUID]
	CreatedAt query.FieldIdentifier[time.Time]
	UpdatedAt query.FieldIdentifier[time.Time]
	DeletedAt query.FieldIdentifier[time.Time]
	UserID    query.FieldIdentifier[model.UUID]
	ExpiresAt query.FieldIdentifier[time.Time]
}

var Session = sessionConditions{
	CreatedAt: query.FieldIdentifier[time.Time]{
		Field:     "CreatedAt",
		ModelType: sessionType,
	},
	DeletedAt: query.FieldIdentifier[time.Time]{
		Field:     "DeletedAt",
		ModelType: sessionType,
	},
	ExpiresAt: query.FieldIdentifier[time.Time]{
		Field:     "ExpiresAt",
		ModelType: sessionType,
	},
	ID: query.FieldIdentifier[model.UUID]{
		Field:     "ID",
		ModelType: sessionType,
	},
	UpdatedAt: query.FieldIdentifier[time.Time]{
		Field:     "UpdatedAt",
		ModelType: sessionType,
	},
	UserID: query.FieldIdentifier[model.UUID]{
		Field:     "UserID",
		ModelType: sessionType,
	},
}

func (sessionConditions sessionConditions) Preload() condition.Condition[models.Session] {
	return condition.NewPreloadCondition[models.Session](sessionConditions.ID, sessionConditions.CreatedAt, sessionConditions.UpdatedAt, sessionConditions.DeletedAt, sessionConditions.UserID, sessionConditions.ExpiresAt)
}
package orm

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/condition"
	ormErrors "github.com/ditrit/badaas/orm/errors"
	"github.com/ditrit/badaas/orm/model"
	ormQuery "github.com/ditrit/badaas/orm/query"
)

type Query[T model.Model] struct {
	gormQuery *ormQuery.GormQuery
	err       error
}

// Ascending specify an ascending order when retrieving models from database
// joinNumber can be used to select the join in case the field is joined more than once
func (query *Query[T]) Ascending(field ormQuery.IFieldIdentifier, joinNumber ...uint) *Query[T] {
	return query.order(field, false, joinNumber)
}

// Descending specify a descending order when retrieving models from database
// joinNumber can be used to select the join in case the field is joined more than once
func (query *Query[T]) Descending(field ormQuery.IFieldIdentifier, joinNumber ...uint) *Query[T] {
	return query.order(field, true, joinNumber)
}

// Order specify order when retrieving models from database
// if descending is true, the ordering is in descending direction
func (query *Query[T]) order(field ormQuery.IFieldIdentifier, descending bool, joinNumberList []uint) *Query[T] {
	err := query.gormQuery.Order(field, descending, getJoinNumber(joinNumberList))
	if err != nil && query.err == nil {
		methodName := "Ascending"
		if descending {
			methodName = "Descending"
		}

		query.err = methodError(err, methodName)
	}

	return query
}

// from a list of uint, return the first or UndefinedJoinNumber in case the list is empty
func getJoinNumber(joinNumberList []uint) int {
	if len(joinNumberList) == 0 {
		return ormQuery.UndefinedJoinNumber
	}

	return int(joinNumberList[0])
}

// Limit specify the number of models to be retrieved
//
// Limit conditions can be cancelled by using `Limit(-1)`
func (query *Query[T]) Limit(limit int) *Query[T] {
	query.gormQuery.Limit(limit)

	return query
}

// Offset specify the number of models to skip before starting to return the results
//
// Offset conditions can be cancelled by using `Offset(-1)`
//
// Warning: in MySQL Offset can only be used if Limit is also used
func (query *Query[T]) Offset(offset int) *Query[T] {
	query.gormQuery.Offset(offset)

	return query
}

// First finds the first model ordered by primary key, matching given conditions
// or returns gorm.ErrRecordNotFound is if no model does it
func (query *Query[T]) First() (*T, error) {
	if query.err != nil {
		return nil, query.err
	}

	var model *T

	return model, query.gormQuery.First(&model)
}

// Take finds the first model returned by the database in no specified order, matching given conditions
// or returns gorm.ErrRecordNotFound is if no model does it
func (query *Query[T]) Take() (*T, error) {
	if query.err != nil {
		return nil, query.err
	}

	var model *T

	return model, query.gormQuery.Take(&model)
}

// Last finds the last model ordered by primary key, matching given conditions
// or returns gorm.ErrRecordNotFound is if no model does it
func (query *Query[T]) Last() (*T, error) {
	if query.err != nil {
		return nil, query.err
	}

	var model *T

	return model, query.gormQuery.Last(&model)
}

// FindOne finds the only one model that matches given conditions
// or returns error if 0 or more than 1 are found.
func (query *Query[T]) FindOne() (*T, error) {
	models, err := query.Find()
	if err != nil {
		return nil, err
	}

	switch {
	case len(models) == 1:
		return models[0], nil
	case len(models) == 0:
		return nil, ormErrors.ErrObjectNotFound
	default:
		return nil, ormErrors.ErrMoreThanOneObjectFound
	}
}

// Find finds all models matching given conditions
func (query *Query[T]) Find() ([]*T, error) {
	if query.err != nil {
		return nil, query.err
	}

	var models []*T

	return models, query.gormQuery.Find(&models)
}

// Create a Query to which the conditions are applied inside transaction tx
func NewQuery[T model.Model](tx *gorm.DB, conditions ...condition.Condition[T]) *Query[T] {
	gormQuery, err := condition.ApplyConditions[T](tx, conditions)

	return &Query[T]{
		gormQuery: gormQuery,
		err:       err,
	}
}

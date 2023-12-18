package orm

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/condition"
	"github.com/ditrit/badaas/orm/errors"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/query"
)

type Query[T model.Model] struct {
	gormQuery *query.GormQuery
	err       error
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
		return nil, errors.ErrObjectNotFound
	default:
		return nil, errors.ErrMoreThanOneObjectFound
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

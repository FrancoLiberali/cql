package condition

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/model"
)

type Delete[T model.Model] struct {
	OrderLimitReturning[T]
}

// Ascending specify an ascending order when updating models
//
// available for: mysql
func (deleteS *Delete[T]) Ascending(field IField) *Delete[T] {
	deleteS.OrderLimitReturning.Ascending(field)

	return deleteS
}

// Descending specify a descending order when updating models
//
// available for: mysql
func (deleteS *Delete[T]) Descending(field IField) *Delete[T] {
	deleteS.OrderLimitReturning.Descending(field)

	return deleteS
}

// Limit specify the number of models to be updated
//
// Limit conditions can be cancelled by using `Limit(-1)`
//
// available for: mysql
func (deleteS *Delete[T]) Limit(limit int) *Delete[T] {
	deleteS.OrderLimitReturning.Limit(limit)

	return deleteS
}

// available for: postgres, sqlite, sqlserver
//
// warning: in sqlite preloads are not allowed
func (deleteS *Delete[T]) Returning(dest *[]T) *Delete[T] {
	deleteS.OrderLimitReturning.Returning(dest)

	return deleteS
}

func (deleteS *Delete[T]) Exec() (int64, error) {
	if deleteS.query.err != nil {
		return 0, deleteS.query.err
	}

	return deleteS.query.cqlQuery.Delete()
}

// Create a Delete to which the conditions are applied inside transaction tx
func NewDelete[T model.Model](tx *gorm.DB, conditions []Condition[T]) *Delete[T] {
	var err error

	if len(conditions) == 0 {
		err = methodError(ErrEmptyConditions, "Delete")
	}

	query := NewQuery(tx, conditions...)
	if err != nil {
		query.err = err
	}

	return &Delete[T]{
		OrderLimitReturning: OrderLimitReturning[T]{
			query:         query,
			orderByCalled: false,
		},
	}
}

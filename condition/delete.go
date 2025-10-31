package condition

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/model"
)

type Delete[T model.Model] struct {
	OrderLimitReturning[T]

	secondaryQuery *Query[T]
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
// warning: in mysql preloads are not allowed
func (deleteS *Delete[T]) Returning(dest *[]T) *Delete[T] {
	gormDB := deleteS.secondaryQuery.cqlQuery.gormDB

	if len(gormDB.Statement.Selects) > 1 ||
		len(gormDB.Statement.Preloads) > 0 {
		deleteS.query.addError(
			methodError(
				preloadsInReturningNotAllowed(deleteS.secondaryQuery.cqlQuery.Dialector()),
				"Returning",
			),
		)

		return deleteS
	}

	deleteS.OrderLimitReturning.Returning(dest)

	return deleteS
}

func (deleteS *Delete[T]) Exec() (int64, error) {
	if deleteS.query.err != nil {
		return 0, deleteS.query.err
	}

	return deleteS.query.cqlQuery.Delete(
		deleteS.secondaryQuery.cqlQuery,
	)
}

// Create a Delete to which the conditions are applied inside transaction tx
func NewDelete[T model.Model](tx *gorm.DB, conditions []Condition[T]) *Delete[T] {
	var err error

	if len(conditions) == 0 {
		err = methodError(ErrEmptyConditions, "Delete")
	}

	primaryQuery := NewQuery[T](tx)
	if err != nil {
		primaryQuery.err = err
	}

	secondaryQuery := NewQuery(tx, conditions...)

	return &Delete[T]{
		OrderLimitReturning: OrderLimitReturning[T]{
			query:         primaryQuery,
			orderByCalled: false,
		},
		secondaryQuery: secondaryQuery,
	}
}

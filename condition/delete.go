package condition

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/model"
)

type Delete[T model.Model] struct {
	OrderLimitReturning[T]

	secondaryQuery       *Query[T]
	softDeleteColumnName string
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
	gormDB := deleteS.query.cqlQuery.gormDB

	if deleteS.secondaryQuery != nil {
		gormDB = deleteS.secondaryQuery.cqlQuery.gormDB
	}

	if len(gormDB.Statement.Selects) > 1 ||
		len(gormDB.Statement.Preloads) > 0 {
		deleteS.query.addError(
			methodError(
				ErrPreloadsInDeleteReturningNotAllowed,
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

	if deleteS.softDeleteColumnName != "" {
		return deleteS.query.cqlQuery.SoftDelete(deleteS.softDeleteColumnName)
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

	model := *new(T)

	var primaryQuery, secondaryQuery *Query[T]

	softDeleteColumnName := model.SoftDeleteColumnName()

	if softDeleteColumnName != "" {
		// as soft delete is implemented with UPDATE, conditions can be applied directly to primary query
		primaryQuery = NewQuery(tx, conditions...)
	} else {
		// for DELETE statements, a secondary query is necessary
		primaryQuery = NewQuery[T](tx)
		secondaryQuery = NewQuery(tx, conditions...)
	}

	if err != nil {
		primaryQuery.err = err
	}

	return &Delete[T]{
		OrderLimitReturning: OrderLimitReturning[T]{
			query:         primaryQuery,
			orderByCalled: false,
		},
		secondaryQuery:       secondaryQuery,
		softDeleteColumnName: softDeleteColumnName,
	}
}

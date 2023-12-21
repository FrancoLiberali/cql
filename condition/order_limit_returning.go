package condition

import (
	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

type OrderLimitReturning[T model.Model] struct {
	query         *Query[T]
	orderByCalled bool
}

// Ascending specify an ascending order when updating models
//
// joinNumber can be used to select the join in case the field is joined more than once
//
// available for: mysql
func (olr *OrderLimitReturning[T]) Ascending(field IField, joinNumber ...uint) {
	if olr.query.gormQuery.Dialector() != sql.MySQL {
		olr.query.addError(methodError(ErrUnsupportedByDatabase, "Ascending"))
	}

	olr.orderByCalled = true
	olr.query.order(field, false, joinNumber)
}

// Descending specify a descending order when updating models
//
// joinNumber can be used to select the join in case the field is joined more than once
//
// available for: mysql
func (olr *OrderLimitReturning[T]) Descending(field IField, joinNumber ...uint) {
	if olr.query.gormQuery.Dialector() != sql.MySQL {
		olr.query.addError(methodError(ErrUnsupportedByDatabase, "Descending"))
	}

	olr.orderByCalled = true
	olr.query.order(field, true, joinNumber)
}

// Limit specify the number of models to be updated
//
// Limit conditions can be cancelled by using `Limit(-1)`
//
// available for: mysql
func (olr *OrderLimitReturning[T]) Limit(limit int) {
	if olr.query.gormQuery.Dialector() != sql.MySQL {
		olr.query.addError(methodError(ErrUnsupportedByDatabase, "Limit"))
	}

	if !olr.orderByCalled {
		olr.query.addError(methodError(ErrOrderByMustBeCalled, "Limit"))
	}

	olr.query.gormQuery.Limit(limit)
}

// available for: postgres, sqlite, sqlserver
//
// warning: in sqlite preloads are not allowed
func (olr OrderLimitReturning[T]) Returning(dest *[]T) {
	err := olr.query.gormQuery.Returning(dest)
	if err != nil {
		olr.query.addError(methodError(err, "Returning"))
	}
}

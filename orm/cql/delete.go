package cql

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/model"
)

type Delete[T model.Model] struct {
	OrderLimitReturning[T]
}

// Ascending specify an ascending order when updating models
//
// joinNumber can be used to select the join in case the field is joined more than once
//
// available for: mysql
func (deleteS *Delete[T]) Ascending(field IField, joinNumber ...uint) *Delete[T] {
	deleteS.OrderLimitReturning.Ascending(field, joinNumber...)

	return deleteS
}

// Descending specify a descending order when updating models
//
// joinNumber can be used to select the join in case the field is joined more than once
//
// available for: mysql
func (deleteS *Delete[T]) Descending(field IField, joinNumber ...uint) *Delete[T] {
	deleteS.OrderLimitReturning.Descending(field, joinNumber...)

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

	return deleteS.query.gormQuery.Delete()
}

// Create a Delete to which the conditions are applied inside transaction tx
func NewDelete[T model.Model](tx *gorm.DB, conditions ...Condition[T]) *Delete[T] {
	return &Delete[T]{
		OrderLimitReturning: OrderLimitReturning[T]{
			query:         NewQuery(tx, conditions...),
			orderByCalled: false,
		},
	}
}

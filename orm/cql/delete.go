package cql

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/model"
)

type Delete[T model.Model] struct {
	query         *Query[T]
	orderByCalled bool
}

func (deleteS *Delete[T]) Exec() (int64, error) {
	return deleteS.query.gormQuery.Delete()
}

// Create a Delete to which the conditions are applied inside transaction tx
func NewDelete[T model.Model](tx *gorm.DB, conditions ...Condition[T]) *Delete[T] {
	return &Delete[T]{
		query:         NewQuery(tx, conditions...),
		orderByCalled: false,
	}
}

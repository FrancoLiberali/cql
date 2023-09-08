package orm

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/condition"
	"github.com/ditrit/badaas/orm/errors"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/query"
)

type Update[T model.Model] struct {
	query         *Query[T]
	orderByCalled bool
}

// available for: postgres, sqlite, sqlserver
//
// warning: in sqlite preloads are not allowed
func (update *Update[T]) Returning(dest *[]T) *Update[T] {
	// TODO hacer el update del logger para que no muestre internals de ditrit/gorm
	err := update.query.gormQuery.Returning(dest)
	if err != nil {
		update.query.addError(methodError(err, "Returning"))
	}

	return update
}

func (update *Update[T]) Set(sets ...*query.Set[T]) (int64, error) {
	setsAsInterface := []query.ISet{}
	for _, set := range sets {
		setsAsInterface = append(setsAsInterface, set)
	}

	return update.unsafeSet(setsAsInterface)
}

// available for: mysql
func (update *Update[T]) SetMultiple(sets ...query.ISet) (int64, error) {
	// TODO hacer lo mismo con todos los operadores
	if update.query.gormQuery.Dialector() != query.MySQL {
		update.query.addError(methodError(errors.ErrUnsupportedByDatabase, "SetMultiple"))
	}

	// TODO que pasa si esta vacio?
	return update.unsafeSet(sets)
}

func (update *Update[T]) unsafeSet(sets []query.ISet) (int64, error) {
	if update.query.err != nil {
		return 0, update.query.err
	}

	updated, err := update.query.gormQuery.Update(sets)
	if err != nil {
		return 0, methodError(err, "Set")
	}

	return updated, nil
}

// Ascending specify an ascending order when updating models
// joinNumber can be used to select the join in case the field is joined more than once
//
// available for: mysql
func (update *Update[T]) Ascending(field query.IFieldIdentifier, joinNumber ...uint) *Update[T] {
	if update.query.gormQuery.Dialector() != query.MySQL {
		update.query.addError(methodError(errors.ErrUnsupportedByDatabase, "Ascending"))
	}

	update.orderByCalled = true
	update.query.order(field, false, joinNumber)

	return update
}

// Descending specify a descending order when updating models
// joinNumber can be used to select the join in case the field is joined more than once
//
// available for: mysql
func (update *Update[T]) Descending(field query.IFieldIdentifier, joinNumber ...uint) *Update[T] {
	if update.query.gormQuery.Dialector() != query.MySQL {
		update.query.addError(methodError(errors.ErrUnsupportedByDatabase, "Descending"))
	}

	update.orderByCalled = true
	update.query.order(field, true, joinNumber)

	return update
}

// Limit specify the number of models to be updated
//
// Limit conditions can be cancelled by using `Limit(-1)`
//
// available for: mysql
func (update *Update[T]) Limit(limit int) *Update[T] {
	if update.query.gormQuery.Dialector() != query.MySQL {
		update.query.addError(methodError(errors.ErrUnsupportedByDatabase, "Limit"))
	}

	if !update.orderByCalled {
		update.query.addError(methodError(errors.ErrOrderByMustBeCalled, "Limit"))
	}

	update.query.gormQuery.Limit(limit)

	return update
}

// Create a Update to which the conditions are applied inside transaction tx
func NewUpdate[T model.Model](tx *gorm.DB, conditions ...condition.Condition[T]) *Update[T] {
	return &Update[T]{
		query:         NewQuery(tx, conditions...),
		orderByCalled: false,
	}
}

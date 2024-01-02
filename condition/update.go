package condition

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

type Update[T model.Model] struct {
	OrderLimitReturning[T]
}

// Set allows updating multiple attributes of the same table.
func (update *Update[T]) Set(sets ...*Set[T]) (int64, error) {
	setsAsInterface := []ISet{}
	for _, set := range sets {
		setsAsInterface = append(setsAsInterface, set)
	}

	return update.unsafeSet(setsAsInterface)
}

// SetMultiple allows updating multiple tables in the same query.
//
// available for: mysql
func (update *Update[T]) SetMultiple(sets ...ISet) (int64, error) {
	if update.query.gormQuery.Dialector() != sql.MySQL {
		update.query.addError(methodError(ErrUnsupportedByDatabase, "SetMultiple"))
	}

	return update.unsafeSet(sets)
}

func (update *Update[T]) unsafeSet(sets []ISet) (int64, error) {
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
//
// joinNumber can be used to select the join in case the field is joined more than once
//
// available for: mysql
func (update *Update[T]) Ascending(field IField, joinNumber ...uint) *Update[T] {
	update.OrderLimitReturning.Ascending(field, joinNumber...)

	return update
}

// Descending specify a descending order when updating models
//
// joinNumber can be used to select the join in case the field is joined more than once
//
// available for: mysql
func (update *Update[T]) Descending(field IField, joinNumber ...uint) *Update[T] {
	update.OrderLimitReturning.Descending(field, joinNumber...)

	return update
}

// Limit specify the number of models to be updated
//
// Limit conditions can be cancelled by using `Limit(-1)`
//
// available for: mysql
func (update *Update[T]) Limit(limit int) *Update[T] {
	update.OrderLimitReturning.Limit(limit)

	return update
}

// available for: postgres, sqlite, sqlserver
//
// warning: in sqlite preloads are not allowed
func (update *Update[T]) Returning(dest *[]T) *Update[T] {
	update.OrderLimitReturning.Returning(dest)

	return update
}

// Create a Update to which the conditions are applied inside transaction tx
func NewUpdate[T model.Model](tx *gorm.DB, conditions ...Condition[T]) *Update[T] {
	return &Update[T]{
		OrderLimitReturning: OrderLimitReturning[T]{
			query:         NewQuery(tx, conditions...),
			orderByCalled: false,
		},
	}
}

type ISet interface {
	getField() IField
	getValue() any
	getJoinNumber() int
}

type Set[T model.Model] struct {
	field      IField
	value      any
	joinNumber int
}

func (set Set[T]) getField() IField {
	return set.field
}

func (set Set[T]) getValue() any {
	return set.value
}

func (set Set[T]) getJoinNumber() int {
	return set.joinNumber
}

type FieldSet[TModel model.Model, TAttribute any] struct {
	field UpdatableField[TModel, TAttribute]
}

func (set FieldSet[TModel, TAttribute]) Eq(value TAttribute) *Set[TModel] {
	return &Set[TModel]{
		field: set.field,
		value: value,
	}
}

// joinNumber can be used to select the join in case the field is joined more than once
func (set FieldSet[TModel, TAttribute]) Dynamic(field FieldOfType[TAttribute], joinNumber ...uint) *Set[TModel] {
	return &Set[TModel]{
		field:      set.field,
		value:      field,
		joinNumber: GetJoinNumber(joinNumber),
	}
}

func (set FieldSet[TModel, TAttribute]) Unsafe(value any) *Set[TModel] {
	return &Set[TModel]{
		field: set.field,
		value: value,
	}
}

type NullableFieldSet[TModel model.Model, TAttribute any] struct {
	FieldSet[TModel, TAttribute]
}

func (set NullableFieldSet[TModel, TAttribute]) Null() *Set[TModel] {
	return &Set[TModel]{
		field: set.field,
		value: nil,
	}
}

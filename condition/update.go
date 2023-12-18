package condition

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/model"
)

type Update[T model.Model] struct {
	OrderLimitReturning[T]
}

func (update *Update[T]) Set(sets ...*Set[T]) (int64, error) {
	setsAsInterface := []ISet{}
	for _, set := range sets {
		setsAsInterface = append(setsAsInterface, set)
	}

	return update.unsafeSet(setsAsInterface)
}

// available for: mysql
func (update *Update[T]) SetMultiple(sets ...ISet) (int64, error) {
	// TODO hacer lo mismo con todos los operadores
	if update.query.gormQuery.Dialector() != MySQL {
		update.query.addError(methodError(ErrUnsupportedByDatabase, "SetMultiple"))
	}

	// TODO que pasa si esta vacio?
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
	// TODO hacer el update del logger para que no muestre internals de ditrit/gorm
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

// TODO mover todo esto
type ISet interface {
	Field() IField
	Value() any
	JoinNumber() int
}

// TODO ver donde pongo esto
type Set[T model.Model] struct {
	field      IField
	value      any
	joinNumber int
}

func (set Set[T]) Field() IField {
	return set.field
}

func (set Set[T]) Value() any {
	return set.value
}

func (set Set[T]) JoinNumber() int {
	return set.joinNumber
}

// TODO ver donde pongo esto
// TODO nombre muy parecido
type FieldSet[TModel model.Model, TAttribute any] struct {
	Field Field[TModel, TAttribute]
}

func (set FieldSet[TModel, TAttribute]) Eq(value TAttribute) *Set[TModel] {
	return &Set[TModel]{
		field: set.Field,
		value: value,
	}
}

// joinNumber can be used to select the join in case the field is joined more than once
func (set FieldSet[TModel, TAttribute]) Dynamic(field FieldOfType[TAttribute], joinNumber ...uint) *Set[TModel] {
	return &Set[TModel]{
		field:      set.Field,
		value:      field,
		joinNumber: GetJoinNumber(joinNumber),
	}
}

func (set FieldSet[TModel, TAttribute]) Unsafe(value any) *Set[TModel] {
	return &Set[TModel]{
		field: set.Field,
		value: value,
	}
}
